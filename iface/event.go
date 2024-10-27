package iface

type IEventSubscriber interface {
	Name() string
	Cb(IEvent)
}

type IEvent interface {
	Name() string
}
