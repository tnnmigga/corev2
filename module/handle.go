package module

import (
	"fmt"
	"reflect"

	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/utils"
	"github.com/tnnmigga/corev2/zlog"
)

func (m *module) RegisterHandler(mType reflect.Type, h func(any)) {
	if _, ok := m.handles[mType]; ok {
		panic(fmt.Errorf("duplicate registration %s", mType.String()))
	}
	m.handles[mType] = h
}

func (m *module) RegisterRPC(mType reflect.Type, rpc func(iface.IRPC)) {
	if _, ok := m.rpcs[mType]; ok {
		panic(fmt.Errorf("duplicate registration %s", mType.String()))
	}
	m.rpcs[mType] = rpc
}

func (m *module) dispatch(req any) {
	defer utils.RecoverPanic()
	switch data := req.(type) {
	case iface.IAsyncCtx:
		data.AsyncCb()
	case iface.IRPCCb:
		data.RPCCb()
	case iface.IRPC:
		body := data.RPCBody()
		mType := reflect.TypeOf(body)
		rpc, ok := m.rpcs[mType]
		if ok {
			rpc(data)
			return
		}
		zlog.Errorf("module %s rpc not found %s", m.name, mType.String())
	default:
		mType := reflect.TypeOf(req)
		h, ok := m.handles[mType]
		if ok {
			h(req)
			return
		}
		zlog.Errorf("module %s handle not found %s", m.name, mType.String())
	}
}
