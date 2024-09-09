package corev2

import (
	"time"

	"github.com/tnnmigga/corev2/conf"
	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/log"
	"github.com/tnnmigga/corev2/message"
	"github.com/tnnmigga/corev2/system"
)

type App struct {
	modules []iface.IModule
}

func DefaultApp() *App {
	return &App{}
}

func (app *App) Append(mods ...iface.IModule) {
	app.modules = append(app.modules, mods...)
}

func (app *App) Launch() {
	err := message.Start()
	if err != nil {
		panic(err)
	}
	for _, m := range app.modules {
		err = m.Run()
		if err != nil {
			log.Panicf("Launch Run module %s error %v", m.Name(), err)
		}
	}
	log.Infof("server %d launch success", conf.ServerID)
}

func (app *App) Shutdown() {
	message.Stop()
	app.waitMsgHandle()
	for i := len(app.modules) - 1; i >= 0; i-- {
		err := app.modules[i].Exit()
		if err != nil {
			log.Errorf("module exit error %v", err)
		}
	}
	system.WaitGoExit()
	log.Infof("server %d shutdown success", conf.ServerID)
}

func (app *App) waitMsgHandle() {
	timeout := time.Minute
	const interval = 100 * time.Millisecond
	count := int(timeout / interval)
	for i := 0; i < count; i++ {
		done := true
		for _, m := range app.modules {
			if !m.Done() {
				done = false
				break
			}
		}
		if done {
			return
		}
		time.Sleep(interval)
	}
	for _, m := range app.modules {
		if !m.Done() {
			log.Errorf("module %v done timeout", m.Name())
		}
	}
}
