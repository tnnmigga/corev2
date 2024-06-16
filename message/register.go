package message

import (
	"reflect"

	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/message/codec"
)

var recvers = map[reflect.Type][]iface.IModule{}

func bind(mType reflect.Type, m iface.IModule) {
	recvers[mType] = append(recvers[mType], m)
}

func RegisterHandler[T any](m iface.IModule, h func(*T)) {
	codec.Register[T]()
	mType := reflect.TypeOf(new(T))
	bind(mType, m)
	m.RegisterHandler(mType, func(a any) {
		h(a.(*T))
	})
}

func RegisterRPC[T any](m iface.IModule, rpc func(req *T, resp func(any), err func(error))) {
	codec.Register[T]()
	mType := reflect.TypeOf(new(T))
	bind(mType, m)
	m.RegisterRPC(mType, func(req iface.IRPCCtx) {
		body := req.RPCBody()
		rpc(body.(*T), req.Return, req.Error)
	})
}
