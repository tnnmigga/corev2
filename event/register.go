package event

import (
	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/utils/stack"
)

var (
	subscribers = map[string][]iface.IEventSubscriber{}
)

func RegisterSubscriber(subscriber iface.IEventSubscriber) {
	subs := subscriber.Subscribing()
	for _, sub := range subs {
		subscribers[sub] = append(subscribers[sub], subscriber)
	}
}

func RegisterHandler[T ~string](m iface.IModule, h func(iface.IEvent), eventType ...T) {
	subscribing := make([]string, len(eventType))
	for i, s := range eventType {
		subscribing[i] = string(s)
	}
	wrapper := eventCbWrapper{
		IModule:     m,
		name:        stack.FuncName(h),
		subscribing: subscribing,
		handler:     h,
	}
	RegisterSubscriber(wrapper)
}

func Cast(e iface.IEvent) {
	eType := e.GetType()
	for _, sub := range subscribers[eType] {
		func() {
			defer stack.RecoverPanic()
			sub.Handler(e)
		}()
	}
}
