package iface

import "reflect"

type IModule interface {
	Name() string
	Handle(mType reflect.Type, h func(any))
	RegisterRPC(mType reflect.Type, rpc func(IRPCCtx))
	Assign(any)
}

type IRPCCtx interface {
	RPCBody() any
	Return(any, error)
}
