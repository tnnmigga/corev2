package iface

import (
	"context"
	"reflect"
)

type IModule interface {
	Name() string
	Handle(mType reflect.Type, h func(any))
	Response(mType reflect.Type, h func(IReqCtx))
	Assign(any)
	Run() error
	Exit() error
	Done() bool
}

type IReqCtx interface {
	context.Context
	ReqBody() any
	Return(any, error)
}
