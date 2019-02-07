package event

import (
	"fmt"

	"github.com/marcusnoble/kube-event-logger/pkg/utils"
	apps_v1beta1 "k8s.io/api/apps/v1beta1"
	batch_v1 "k8s.io/api/batch/v1"
	api_v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
)

type Event struct {
	Key       string
	Namespace string
	Kind      string
	Host      string
	Action    string
	Name      string
	Phase     api_v1.PodPhase
}

func New(obj interface{}, action string) Event {
	objectMeta := utils.GetObjectMetaData(obj)
	kubeEvent := Event{
		Namespace: objectMeta.Namespace,
		Action:    action,
		Name:      objectMeta.Name,
	}
	kubeEvent.Key, _ = cache.MetaNamespaceKeyFunc(obj)

	switch object := obj.(type) {
	case *ext_v1beta1.DaemonSet:
		kubeEvent.Kind = "daemon set"
	case *apps_v1beta1.Deployment:
		kubeEvent.Kind = "deployment"
	case *batch_v1.Job:
		kubeEvent.Kind = "job"
	case *api_v1.Namespace:
		kubeEvent.Kind = "namespace"
	case *ext_v1beta1.Ingress:
		kubeEvent.Kind = "ingress"
	case *api_v1.PersistentVolume:
		kubeEvent.Kind = "persistent volume"
	case *api_v1.Pod:
		kubeEvent.Kind = "pod"
		kubeEvent.Host = object.Spec.NodeName

		if action == "updated" {
			kubeEvent.Phase = object.Status.Phase
		}
	case *api_v1.ReplicationController:
		kubeEvent.Kind = "replication controller"
	case *ext_v1beta1.ReplicaSet:
		kubeEvent.Kind = "replica set"
	case *api_v1.Service:
		kubeEvent.Kind = "service"
	case *api_v1.Secret:
		kubeEvent.Kind = "secret"
	case *api_v1.ConfigMap:
		kubeEvent.Kind = "configmap"
	case Event:
		kubeEvent.Name = object.Name
		kubeEvent.Kind = object.Kind
		kubeEvent.Namespace = object.Namespace
	}

	return kubeEvent
}

func (e *Event) Message() (msg string) {
	if e.Kind == "namespace" {
		msg = fmt.Sprintf("A namespace `%s` has been `%s`", e.Name, e.Action)
	} else {
		msg = fmt.Sprintf("A `%s` in namespace `%s` has been `%s`:\n`%s`", e.Kind, e.Namespace, e.Action, e.Name)
	}

	return msg
}
