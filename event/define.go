package event

import (
	"sync"

	"github.com/tnnmigga/corev2/iface"
)

type EventBus struct {
	subscribers map[string][]iface.IEventSubscriber
	rw          sync.RWMutex
}

func (bus *EventBus) Subscribe(cb func(iface.IEvent), events ...string) {
	bus.rw.Lock()
	defer bus.rw.Unlock()
	
}

type eventCbWrapper struct {
	iface.IModule
	name        string
	subscribing []string
	handler     func(iface.IEvent)
}

func (w eventCbWrapper) Name() string {
	return w.name
}

func (w eventCbWrapper) Subscribing() []string {
	return w.subscribing
}

func (w eventCbWrapper) Handler(e iface.IEvent) {
	fn := func() { w.handler(e) }
	w.Assign(fn)
}
