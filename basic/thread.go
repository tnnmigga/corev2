package basic

import (
	"reflect"
	"sync/atomic"
	"time"

	"github.com/tnnmigga/corev2/conc"
	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/log"
	"github.com/tnnmigga/corev2/utils"
)

type eventLoopModule struct {
	basic
	mq      chan any
	pending atomic.Int32
	wg      utils.WaitGroupWithTimeout
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
	m.wg.Add(1)
	conc.Go(func() {
		defer m.wg.Done()
		for req := range m.mq {
			m.pending.Add(1)
			m.dispatch(req)
			m.pending.Add(-1)
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

func (m *eventLoopModule) Done() bool {
	return len(m.mq) == 0 && m.pending.Load() == 0
}

func (m *eventLoopModule) Exit() error {
	close(m.mq)
	return m.wg.WaitWithTimeout(time.Minute)
}

type concurrencyModule struct {
	basic
	pending atomic.Int32
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
	m.pending.Add(1)
	conc.Go(func() {
		defer m.pending.Add(-1)
		m.dispatch(msg)
	})
}

func (m *concurrencyModule) Done() bool {
	return m.pending.Load() == 0
}

func (m *concurrencyModule) Exit() error {
	return nil
}

type threadPoolModule struct {
	basic
	mq      chan any
	wg      utils.WaitGroupWithTimeout
	pending atomic.Int32
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
		m.wg.Add(1)
		conc.Go(func() {
			defer m.wg.Done()
			for req := range m.mq {
				m.pending.Add(1)
				m.dispatch(req)
				m.pending.Add(-1)
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

func (m *threadPoolModule) Done() bool {
	return len(m.mq) == 0 && m.pending.Load() == 0
}

func (m *threadPoolModule) Exit() error {
	close(m.mq)
	return m.wg.WaitWithTimeout(time.Minute)
}
