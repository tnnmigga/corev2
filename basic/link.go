package basic

import (
	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/log"
)

var Link = moduleLink{modules: map[string]iface.IModule{}}

type moduleLink struct {
	modules map[string]iface.IModule
}

func (ms *moduleLink) regModule(m iface.IModule) {
	name := m.Name()
	if _, ok := ms.modules[name]; ok {
		log.Panicf("module %s exist", name)
	}
	ms.modules[name] = m
}

func (ms *moduleLink) Get(name string) iface.IModule {
	return ms.modules[name]
}

func (ms *moduleLink) Send(name string, msg any) {
	if m, ok := ms.modules[name]; ok {
		m.Assign(msg)
		return
	}
	log.Panicf("module %s not found", name)
}
