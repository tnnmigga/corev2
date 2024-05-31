package iface

type IEventSubscriber interface {
	Name() string
	Subscribing() []string
	Handler(any)
}

type IEvent interface {
	GetType()
}
