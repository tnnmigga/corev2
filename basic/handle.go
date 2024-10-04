package basic

import (
	"fmt"
	"reflect"

	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/log"
	"github.com/tnnmigga/corev2/utils"
)

type handle struct {
	name      string
	handleFns map[reflect.Type]func(any)
	respFns   map[reflect.Type](func(iface.IRequestCtx))
}

func (m *handle) Name() string {
	return m.name
}

func (m *handle) Handle(mType reflect.Type, h func(any)) {
	if _, ok := m.handleFns[mType]; ok {
		panic(fmt.Errorf("duplicate registration %s", mType.String()))
	}
	m.handleFns[mType] = h
}

func (m *handle) Response(mType reflect.Type, h func(iface.IRequestCtx)) {
	if _, ok := m.respFns[mType]; ok {
		panic(fmt.Errorf("duplicate registration %s", mType.String()))
	}
	m.respFns[mType] = h
}

func (m *handle) dispatch(msg any) {
	defer utils.RecoverPanic()
	switch req := msg.(type) {
	case func():
		req()
	case iface.IRequestCtx:
		body := req.Body()
		mType := reflect.TypeOf(body)
		h, ok := m.respFns[mType]
		if ok {
			h(req)
			return
		}
		log.Errorf("module %s request not found %s", m.name, mType.String())
	default:
		mType := reflect.TypeOf(msg)
		h, ok := m.handleFns[mType]
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
