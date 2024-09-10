package conc

import (
	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/utils"
)

func Async[T any](m iface.IModule, f func() (T, error), cb func(T, error), groupKey ...string) {
	fn := func() {
		defer utils.RecoverPanic()
		result, err := f()
		m.Assign(func() {
			cb(result, err)
		})
	}
	if len(groupKey) > 0 {
		GoWithGroup(groupKey[0], fn)
	} else {
		Go(fn)
	}
}
