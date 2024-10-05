package basic

import (
	"fmt"
	"reflect"

	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/utils"
)

type handle struct {
	handleFns map[reflect.Type]func(any)
	respFns   map[reflect.Type](func(iface.IRequestCtx))
}

func (m *handle) Name() string {
	return "unknown"
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
		h := m.respFns[mType]
		h(req)
	default:
		mType := reflect.TypeOf(msg)
		h := m.handleFns[mType]
		h(msg)
	}
}

func (m *handle) Run() error {
	return nil
}
