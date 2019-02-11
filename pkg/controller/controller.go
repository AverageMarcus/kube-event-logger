package controller

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/marcusnoble/kube-event-logger/pkg/event"
	"github.com/marcusnoble/kube-event-logger/pkg/handler"
	"github.com/marcusnoble/kube-event-logger/pkg/logger"
	"github.com/marcusnoble/kube-event-logger/pkg/utils"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	kubeClient kubernetes.Interface
	queue      workqueue.RateLimitingInterface
	informer   cache.SharedIndexInformer
}

const maxRetries = 5

var serverStartTime time.Time

func Start() {
	var kubeClient kubernetes.Interface
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfigPath := os.Getenv("KUBECONFIG")
		if kubeconfigPath == "" {
			kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
	kubeClient, err = kubernetes.NewForConfig(config)

	if err != nil {
		panic(err)
	}

	informers := []cache.SharedIndexInformer{}

	informers = append(informers, handler.WatchPods(kubeClient))
	informers = append(informers, handler.WatchDeployments(kubeClient))
	informers = append(informers, handler.WatchReplicaSets(kubeClient))
	informers = append(informers, handler.WatchJobs(kubeClient))
	informers = append(informers, handler.WatchNamespaces(kubeClient))
	informers = append(informers, handler.WatchIngresses(kubeClient))
	informers = append(informers, handler.WatchPersistentVolumes(kubeClient))
	informers = append(informers, handler.WatchReplicationControllers(kubeClient))
	informers = append(informers, handler.WatchServices(kubeClient))
	informers = append(informers, handler.WatchSecrets(kubeClient))
	informers = append(informers, handler.WatchConfigMaps(kubeClient))

	for _, informer := range informers {
		queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				queue.Add(event.New(obj, "created"))
			},
			UpdateFunc: func(old, new interface{}) {
				queue.Add(event.New(new, "updated"))
			},
			DeleteFunc: func(obj interface{}) {
				queue.Add(event.New(obj, "deleted"))
			},
		})
		c := &Controller{
			kubeClient: kubeClient,
			informer:   informer,
			queue:      queue,
		}

		stopCh := make(chan struct{})
		defer close(stopCh)
		go c.Run(stopCh)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	defer c.queue.ShutDown()

	serverStartTime = time.Now().Local()

	go c.informer.Run(stopCh)

	wait.Until(c.runWorker, time.Second, stopCh)
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	newEvent, quit := c.queue.Get()

	if quit {
		return false
	}
	defer c.queue.Done(newEvent)
	err := c.processItem(newEvent.(event.Event))
	if err == nil {
		// No error, reset the ratelimit counters
		c.queue.Forget(newEvent)
	} else if c.queue.NumRequeues(newEvent) < maxRetries {
		c.queue.AddRateLimited(newEvent)
	} else {
		c.queue.Forget(newEvent)
	}

	return true
}

func (c *Controller) processItem(newEvent event.Event) error {
	obj, _, err := c.informer.GetIndexer().GetByKey(newEvent.Key)
	if err != nil {
		return fmt.Errorf("Error fetching object with key %s from store: %v\n", newEvent.Name, err)
	}
	objectMeta := utils.GetObjectMetaData(obj)

	log := new(logger.Default)

	switch newEvent.Action {
	case "created":
		if objectMeta.CreationTimestamp.Sub(serverStartTime).Seconds() > 0 {
			log.ObjectCreated(newEvent)
		}
	case "updated":
		log.ObjectUpdated(newEvent)
	case "deleted":
		log.ObjectDeleted(newEvent)
	}
	return nil
}
