package handler

import (
	"github.com/marcusnoble/kube-event-logger/pkg/config"
	ext_v1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func WatchIngresses(kubeClient kubernetes.Interface, config *config.Config) cache.SharedIndexInformer {
	if !config.Resource.Ingress {
		return nil
	}
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return kubeClient.ExtensionsV1beta1().Ingresses(config.Namespace).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return kubeClient.ExtensionsV1beta1().Ingresses(config.Namespace).Watch(options)
			},
		},
		&ext_v1.Ingress{},
		0,
		cache.Indexers{},
	)

	return informer
}
