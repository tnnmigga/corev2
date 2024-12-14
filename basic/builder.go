package basic

import "github.com/tnnmigga/corev2/iface"

var ModuleFactory = map[string]func() iface.IModule{}

func RegModule[T interface{ Name() string }](f func() *T) {
	var m T
	name := m.Name()
	if len(name) == 0 {
		panic("module unnamed")
	}
	ModuleFactory[name] = func() iface.IModule {
		return any(f()).(iface.IModule)
	}
}

func Create(names ...string) []iface.IModule {
	ms := make([]iface.IModule, 0, len(names))
	for _, name := range names {
		f := ModuleFactory[name]
		m := f()
		ms = append(ms, m)
	}
	return ms
}
