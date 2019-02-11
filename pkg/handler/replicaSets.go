package handler

import (
	"github.com/marcusnoble/kube-event-logger/pkg/config"
	apps_v1 "k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func WatchReplicaSets(kubeClient kubernetes.Interface, config *config.Config) cache.SharedIndexInformer {
	if !config.Resource.ReplicaSet {
		return nil
	}
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return kubeClient.AppsV1().ReplicaSets(config.Namespace).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return kubeClient.AppsV1().ReplicaSets(config.Namespace).Watch(options)
			},
		},
		&apps_v1.ReplicaSet{},
		0,
		cache.Indexers{},
	)

	return informer
}
