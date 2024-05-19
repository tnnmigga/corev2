package message

import (
	"reflect"

	"github.com/tnnmigga/corev2/codec"
	"github.com/tnnmigga/corev2/iface"
)

var recvers = map[reflect.Type][]iface.IReactor{}

func bind(mType reflect.Type, m iface.IReactor) {
	recvers[mType] = append(recvers[mType], m)
}

func RegisterHandle[T any](m iface.IReactor, h func(*T)) {
	codec.Register[T]()
	mType := reflect.TypeOf(new(T))
	bind(mType, m)
	m.RegisterHandle(mType, func(a any) {
		h(a.(*T))
	})
}

func RegisterRPC[T any](m iface.IReactor, rpc func(req *T, resp func(any), err func(error))) {
	codec.Register[T]()
	mType := reflect.TypeOf(new(T))
	bind(mType, m)
	m.RegisterRPC(mType, func(req iface.IRPC) {
		body := req.RPCBody()
		rpc(body.(*T), req.Return, req.Error)
	})
}
