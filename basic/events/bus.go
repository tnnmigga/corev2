package events

import (
	"sync"

	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/utils"
)

var (
	subscribers = map[string][]iface.IEventSubscriber{}
	rw          = sync.RWMutex{}
)

type subscriber[T iface.IEvent] struct {
	name string
	m    iface.IModule
	cb   func(*T)
}

func (sub subscriber[T]) Name() string {
	return sub.name
}

func (sub *subscriber[T]) Cb(e any) {
	sub.m.Assign(func() {
		sub.cb(e.(*T))
	})
}

func Subscribe[T iface.IEvent](m iface.IModule, cb func(*T), events ...string) {
	Unsubscribe(cb)
	if len(events) == 0 {
		return
	}
	rw.Lock()
	defer rw.Unlock()
	name := utils.FuncName(cb)
	sub := &subscriber[T]{
		name: name,
		m:    m,
		cb:   cb,
	}
	for _, e := range events {
		subscribers[e] = append(subscribers[e], sub)
	}
}

func Unsubscribe[T iface.IEvent](cb func(*T)) {
	name := utils.FuncName(cb)
	rw.Lock()
	defer rw.Unlock()
	for event := range subscribers {
		for idx, sub := range subscribers[event] {
			if sub.Name() == name {
				subscribers[event] = append(subscribers[event][:idx], subscribers[event][idx+1:]...)
			}
		}
	}
}

func Publish(events ...iface.IEvent) {
	rw.RLock()
	defer rw.RUnlock()
	for _, e := range events {
		eventName := e.Name()
		for _, sub := range subscribers[eventName] {
			sub.Cb(e)
		}
	}
}
