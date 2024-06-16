package process

import (
	"github.com/tnnmigga/corev2/event"
	"github.com/tnnmigga/corev2/iface"
)

type ProcessEvent string

func (e ProcessEvent) GetType() string {
	return string(e)
}

const (
	EventAfterInit  ProcessEvent = "process:afterInit"
	EventBeforeRun  ProcessEvent = "process:beforeRun"
	EventAfterRun   ProcessEvent = "process:afterRun"
	EventBeforeStop ProcessEvent = "process:afterStop"
	EventAfterStop  ProcessEvent = "process:afterStop"
)

func Create(mods ...iface.IModule) {
	event.Cast(EventAfterInit)
	
}

func Exit() {

}
