package module

import (
	"reflect"

	"github.com/tnnmigga/corev2/iface"
)

type multi struct {
	basic
}

func Multi(name string, mqLen int) iface.IModule {
	m := &multi{
		basic: basic{
			name:    name,
			handles: map[reflect.Type]func(any){},
			rpcs:    map[reflect.Type]func(iface.IRPCCtx){},
		},
	}
	return m
}

func (m *multi) Name() string {
	return m.name
}
