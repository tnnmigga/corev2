package module

import (
	"fmt"
	"reflect"

	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/utils/stack"
	"github.com/tnnmigga/corev2/zlog"
)

func (m *module) RegisterHandler(mType reflect.Type, h func(any)) {
	if _, ok := m.handles[mType]; ok {
		panic(fmt.Errorf("duplicate registration %s", mType.String()))
	}
	m.handles[mType] = h
}

func (m *module) RegisterRPC(mType reflect.Type, rpc func(iface.IRPCCtx)) {
	if _, ok := m.rpcs[mType]; ok {
		panic(fmt.Errorf("duplicate registration %s", mType.String()))
	}
	m.rpcs[mType] = rpc
}

func (m *module) dispatch(msg any) {
	defer stack.RecoverPanic()
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
		zlog.Errorf("module %s rpc not found %s", m.name, mType.String())
	default:
		mType := reflect.TypeOf(msg)
		h, ok := m.handles[mType]
		if ok {
			h(msg)
			return
		}
		zlog.Errorf("module %s handle not found %s", m.name, mType.String())
	}
}
