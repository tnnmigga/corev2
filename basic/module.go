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

const DefaultMQLen = 100000

type eventLoopModule struct {
	handle
	mq      chan any
	pending atomic.Int32
	wg      utils.WaitGroupWithTimeout
}

// 单线程异步 事件循环
func NewEventLoop(mqLen int) iface.IModule {
	m := &eventLoopModule{
		mq: make(chan any, mqLen),
		handle: handle{
			handleFns: map[reflect.Type]func(any){},
			respFns:   map[reflect.Type]func(iface.IRequestCtx){},
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
	Link.regModule(m)
	return m
}

func (m *eventLoopModule) Assign(msg any) {
	select {
	case m.mq <- msg:
	default:
		log.Errorf("modele mq full, msg %#v", msg)
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
	handle
	pending atomic.Int32
}

// 每个请求一个goroutine执行
func NewConcurrency() iface.IModule {
	m := &concurrencyModule{
		handle: handle{
			handleFns: map[reflect.Type]func(any){},
			respFns:   map[reflect.Type]func(iface.IRequestCtx){},
		},
	}
	return m
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

type goroutinePoolModule struct {
	handle
	mq      chan any
	wg      utils.WaitGroupWithTimeout
	pending atomic.Int32
}

// 线程池模式
func NewGoPool(mqLen int, goNum int) iface.IModule {
	m := &goroutinePoolModule{
		mq: make(chan any, mqLen),
		handle: handle{
			handleFns: map[reflect.Type]func(any){},
			respFns:   map[reflect.Type]func(iface.IRequestCtx){},
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
	Link.regModule(m)
	return m
}

func (m *goroutinePoolModule) Assign(msg any) {
	select {
	case m.mq <- msg:
	default:
		log.Errorf("modele mq full, lose %#v", msg)
	}
}

func (m *goroutinePoolModule) Done() bool {
	return len(m.mq) == 0 && m.pending.Load() == 0
}

func (m *goroutinePoolModule) Exit() error {
	close(m.mq)
	return m.wg.WaitWithTimeout(time.Minute)
}
