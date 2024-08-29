package module

import (
	"fmt"
	"reflect"

	"github.com/tnnmigga/corev2/conc"
	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/logger"
	"github.com/tnnmigga/corev2/message"
	"github.com/tnnmigga/corev2/message/codec"
	"github.com/tnnmigga/corev2/utils"
)

type basic struct {
	name    string
	handles map[reflect.Type]func(any)
	rpcs    map[reflect.Type](func(iface.IRPCCtx))
}

func (m *basic) Name() string {
	return m.name
}

func (m *basic) Handle(mType reflect.Type, h func(any)) {
	if _, ok := m.handles[mType]; ok {
		panic(fmt.Errorf("duplicate registration %s", mType.String()))
	}
	m.handles[mType] = h
}

func (m *basic) RegisterRPC(mType reflect.Type, rpc func(iface.IRPCCtx)) {
	if _, ok := m.rpcs[mType]; ok {
		panic(fmt.Errorf("duplicate registration %s", mType.String()))
	}
	m.rpcs[mType] = rpc
}

func (m *basic) dispatch(msg any) {
	defer utils.RecoverPanic()
	switch req := msg.(type) {
	case func():
		req()
	case iface.IRPCCtx:
		body := req.RPCBody()
		mType := reflect.TypeOf(body)
		rpc, ok := m.rpcs[mType]
		if ok {
			rpc(req)
			return
		}
		logger.Errorf("module %s rpc not found %s", m.name, mType.String())
	default:
		mType := reflect.TypeOf(msg)
		h, ok := m.handles[mType]
		if ok {
			h(msg)
			return
		}
		logger.Errorf("module %s handle not found %s", m.name, mType.String())
	}
}

func Handle[T any](m iface.IModule, h func(*T)) {
	codec.Register[T]()
	mType := reflect.TypeOf(new(T))
	message.Subscribe[T](func(msg *T) {
		m.Assign(msg)
	})
	m.Handle(mType, func(a any) {
		h(a.(*T))
	})
}

func RegisterRPC[T any](m iface.IModule, rpc func(req *T, resp func(any), err func(error))) {
	codec.Register[T]()
	mType := reflect.TypeOf(new(T))
	message.Subscribe[T](func(msg *T) {
		m.Assign(msg)
	})
	m.RegisterRPC(mType, func(req iface.IRPCCtx) {
		body := req.RPCBody()
		rpc(body.(*T), req.Return, req.Error)
	})
}

func Async[T any](m iface.IModule, f func() (T, error), cb func(T, error)) {
	conc.Go(func() {
		defer utils.RecoverPanic()
		result, err := f()
		m.Assign(func() {
			cb(result, err)
		})
	})
}
