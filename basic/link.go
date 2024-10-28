package basic

import (
	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/log"
)

var Link = moduleLink{modules: map[string]iface.IModule{}}

type moduleLink struct {
	modules map[string]iface.IModule
}

func (ml *moduleLink) regModule(m iface.IModule) {
	name := m.Name()
	if _, ok := ml.modules[name]; ok {
		log.Panicf("module %s exist", name)
	}
	ml.modules[name] = m
}

func (ml *moduleLink) Get(name string) iface.IModule {
	return ml.modules[name]
}

func (ml *moduleLink) Send(name string, msg any) {
	if m, ok := ml.modules[name]; ok {
		m.Assign(msg)
		return
	}
	log.Panicf("module %s not found", name)
}
