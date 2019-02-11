package handler

import (
	"github.com/marcusnoble/kube-event-logger/pkg/config"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func WatchPods(kubeClient kubernetes.Interface, config *config.Config) cache.SharedIndexInformer {
	if !config.Resource.Pod {
		return nil
	}
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return kubeClient.CoreV1().Pods(config.Namespace).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return kubeClient.CoreV1().Pods(config.Namespace).Watch(options)
			},
		},
		&api_v1.Pod{},
		0,
		cache.Indexers{},
	)

	return informer
}
