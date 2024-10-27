package iface

type IEventSubscriber interface {
	Name() string
	Cb(any)
}

type IEvent interface {
	Name() string
}
