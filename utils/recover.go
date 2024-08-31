package utils

import (
	"runtime/debug"

	"github.com/tnnmigga/corev2/log"
)

func RecoverPanic() {
	if r := recover(); r != nil {
		log.Errorf("%v: %s", r, debug.Stack())
	}
}

func ExecAndRecover(fn func()) {
	defer RecoverPanic()
	fn()
}
