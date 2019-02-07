package logger

import (
	"encoding/json"
	"fmt"

	"github.com/marcusnoble/kube-event-logger/pkg/event"
)

type kubeEvent struct {
	event.Event
	Message string
}

type logMessage struct {
	KubeEvent kubeEvent
}

type Logger interface {
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(oldObj, newObj interface{})
}

type Default struct {
}

func (d *Default) ObjectCreated(obj event.Event) {
	fmt.Println(prepareLogMessage(obj))
}

func (d *Default) ObjectDeleted(obj event.Event) {
	fmt.Println(prepareLogMessage(obj))
}

func (d *Default) ObjectUpdated(obj event.Event) {
	fmt.Println(prepareLogMessage(obj))
}

func prepareLogMessage(e event.Event) string {
	msg := logMessage{
		KubeEvent: kubeEvent{e, e.Message()},
	}
	b, _ := json.Marshal(msg)
	return string(b)
}
