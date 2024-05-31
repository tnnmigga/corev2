package process

import (
	"github.com/tnnmigga/corev2/iface"
)

const (
	EventAfterInit = iota
	EventBeforeRun
	EventAfterRun
	EventBeforeStop
	EventAfterStop
)

func Create(mods ...iface.IModule) {

}

func Exit() {

}
