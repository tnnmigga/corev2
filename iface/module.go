package iface

import (
	"context"
	"reflect"
)

type IModule interface {
	Name() string
	Handle(mType reflect.Type, h func(any))
	Response(mType reflect.Type, h func(IRequestCtx))
	Assign(any)
	Init() error
	Run() error
	Exit() error
	Done() bool
}

type IRequestCtx interface {
	context.Context
	Body() any
	Return(any, error)
}
