package iface

import "reflect"

type IProcess interface {
}

type IReactor interface {
	RegisterHandle(mType reflect.Type, h func(any))
	RegisterRPC(mType reflect.Type, rpc func(IRPC))
	Assign(any)
}

type IRPC interface {
	RPCBody() any
	Return(any)
	Error(error)
}

type IRPCCb interface {
	RPCCb()
}

type IAsyncCtx interface {
	AsyncCb()
}

type IMsgDest interface {
	String() string
}
