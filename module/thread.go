package module

import (
	"reflect"

	"github.com/tnnmigga/corev2/conc"
	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/log"
)

type eventLoopModule struct {
	basic
	mq chan any
}

// 单线程异步 事件循环
func NewEventLoop(name string, mqLen int) iface.IModule {
	m := &eventLoopModule{
		mq: make(chan any, mqLen),
		basic: basic{
			name:    name,
			handles: map[reflect.Type]func(any){},
			rpcs:    map[reflect.Type]func(iface.IRPCCtx){},
		},
	}
	conc.Go(func() {
		for req := range m.mq {
			m.dispatch(req)
		}
	})
	return m
}

func (m *eventLoopModule) Assign(msg any) {
	select {
	case m.mq <- msg:
	default:
		log.Errorf("modele %s mq full, lose %#v", m.name, msg)
	}
}

type concurrencyModule struct {
	basic
}

// 每个请求一个goroutine执行
func NewConcurrency(name string) iface.IModule {
	m := &concurrencyModule{
		basic: basic{
			name:    name,
			handles: map[reflect.Type]func(any){},
			rpcs:    map[reflect.Type]func(iface.IRPCCtx){},
		},
	}
	return m
}

func (m *concurrencyModule) Name() string {
	return m.name
}

func (m *concurrencyModule) Assign(msg any) {
	conc.Go(func() {
		m.dispatch(msg)
	})
}

type threadPoolModule struct {
	basic
	mq chan any
}

// 线程池模式
func NewThreadPool(name string, mqLen int, goNum int) iface.IModule {
	m := &threadPoolModule{
		mq: make(chan any, mqLen),
		basic: basic{
			name:    name,
			handles: map[reflect.Type]func(any){},
			rpcs:    map[reflect.Type]func(iface.IRPCCtx){},
		},
	}
	for i := 0; i < goNum; i++ {
		conc.Go(func() {
			for req := range m.mq {
				m.dispatch(req)
			}
		})
	}
	return m
}

func (m *threadPoolModule) Assign(msg any) {
	select {
	case m.mq <- msg:
	default:
		log.Errorf("modele %s mq full, lose %#v", m.name, msg)
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
