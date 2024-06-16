package iface

type IEventSubscriber interface {
	IModule
	Name() string
	Subscribing() []string
	Handler(IEvent)
}

type IEvent interface {
	GetType() string
}
