package message

import (
	"reflect"

	"github.com/tnnmigga/corev2/iface"
)

func RegisterHandle[T any](m iface.IReactor, h func(*T)) {
	mType := reflect.TypeOf(new(T))
	m.RegisterHandle(mType, func(a any) {
		h(a.(*T))
	})
}

func RegisterRPC[T any](m iface.IReactor, rpc func(req *T, resp func(any), err func(error))) {
	mType := reflect.TypeOf(new(T))
	m.RegisterRPC(mType, func(req iface.IRPC) {
		body := req.RPCBody()
		rpc(body.(*T), req.Return, req.Error)
	})
}
