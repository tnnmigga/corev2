package event

import (
	fmt "fmt"
	"strconv"

	"github.com/tnnmigga/corev2/iface"
)

var (
	subscribers = map[string][]iface.IEventSubscriber{}
)

func RegisterHandler[T any](m iface.IModule, eventType string, h func(*Event)) {
	subscribers[eventType] = append(subscribers[eventType], eventCbWrapper{})
}

func RegisterSubscriber(sub iface.IEventSubscriber) {

}

type eventCbWrapper struct {
	name        string
	subscribing []string
}

func (w eventCbWrapper) Name() string {
	return w.name
}

func (w eventCbWrapper) Subscribing() []string {
	return w.subscribing
}

func (e *Event) Int(name string) int {
	if n, err := strconv.Atoi(e.Str(name)); err == nil {
		return n
	}
	panic(fmt.Errorf("event param %s not a number ", name))
}

func (e *Event) Str(name string) (arg string) {
	if e.Args != nil {
		if v, ok := e.Args[name]; ok {
			return v
		}
	}
	panic(fmt.Errorf("event param %s not found", name))
}
