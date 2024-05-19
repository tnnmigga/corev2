package proc

import (
	"github.com/tnnmigga/corev2/iface"
)

const (
	StatusAfterInit = iota
	StatusBeforeRun
	StatusAfterRun
	StatusBeforeStop
	StatusAfterStop
	StatusBeforeExit
)

var (
	root  process
	hooks [StatusBeforeExit + 1][]func() error
)

func Root() iface.IProcess {
	return root
}

type process struct{}

func RegisterHook(status int, h func() error) {
	if status < StatusAfterInit || status > StatusBeforeExit {
		panic("invalid status")
	}
	hooks[status] = append(hooks[status], h)
}
