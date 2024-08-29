package iface

type IEventBus interface {
	Subscribe(func(IEvent), ...string)
	Publish(IEvent)
}

type IEventSubscriber interface {
	Name() string
	Subscribing(string) bool
	Handle(IEvent)
}

type IEvent interface {
	GetType() string
}
