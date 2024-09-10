package basic

import (
	"fmt"
	"reflect"

	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/log"
	"github.com/tnnmigga/corev2/utils"
)

type handle struct {
	name    string
	handles map[reflect.Type]func(any)
	rpcs    map[reflect.Type](func(iface.IRPCCtx))
}

func (m *handle) Name() string {
	return m.name
}

func (m *handle) Handle(mType reflect.Type, h func(any)) {
	if _, ok := m.handles[mType]; ok {
		panic(fmt.Errorf("duplicate registration %s", mType.String()))
	}
	m.handles[mType] = h
}

func (m *handle) RegisterRPC(mType reflect.Type, rpc func(iface.IRPCCtx)) {
	if _, ok := m.rpcs[mType]; ok {
		panic(fmt.Errorf("duplicate registration %s", mType.String()))
	}
	m.rpcs[mType] = rpc
}

func (m *handle) dispatch(msg any) {
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
		log.Errorf("module %s rpc not found %s", m.name, mType.String())
	default:
		mType := reflect.TypeOf(msg)
		h, ok := m.handles[mType]
		if ok {
			h(msg)
			return
		}
		log.Errorf("module %s handle not found %s", m.name, mType.String())
	}
}

func (m *handle) Run() error {
	return nil
}
