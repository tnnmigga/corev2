package module

import (
	"reflect"

	"github.com/tnnmigga/corev2/conc"
	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/logger"
)

type module struct {
	name    string
	mq      chan any
	handles map[reflect.Type]func(any)
	rpcs    map[reflect.Type](func(iface.IRPCCtx))
}

func New(name string, workerNum int, mqLen int) iface.IModule {
	m := &module{
		name:    name,
		mq:      make(chan any, mqLen),
		handles: map[reflect.Type]func(any){},
		rpcs:    map[reflect.Type]func(iface.IRPCCtx){},
	}
	for i := 0; i < workerNum; i++ {
		conc.Go(func() {
			for req := range m.mq {
				m.dispatch(req)
			}
		})
	}
	return m
}

func (m *module) MQ() chan any {
	return m.mq
}

func (m *module) Name() string {
	return m.name
}

func (m *module) Assign(msg any) {
	select {
	case m.mq <- msg:
	default:
		logger.Errorf("modele %s mq full, lose %#v", m.name, msg)
	}
}

// func WaitMsgHandle(timeout ...time.Duration) {
// 	if len(timeout) == 0 {
// 		timeout = append(timeout, time.Minute)
// 	}
// 	const interval = 100 * time.Millisecond
// 	count := int(timeout[0] / interval)
// 	for i := 0; i < count; i++ {
// 		flag := true
// 		for _, m := range modules {
// 			if len(m.mq) != 0 {
// 				flag = false
// 				break
// 			}
// 		}
// 		if flag {
// 			return
// 		}
// 		time.Sleep(interval)
// 	}
// 	for _, m := range modules {
// 		if len(m.mq) != 0 {
// 			zlog.Errorf("module mq remain %d", len(m.mq))
// 		}
// 	}
// }
