package conc

import (
	"context"
	"sync"

	"github.com/tnnmigga/corev2/algorithm"
	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/proc"
	"github.com/tnnmigga/corev2/utils"
	"github.com/tnnmigga/corev2/zlog"
)

var (
	rootCtx, cancelGo = context.WithCancel(context.Background())
	wkg               = newWorkerGroup()
	running           = algorithm.NewCounter[string]()
)

func newWorkerGroup() *workerGroup {
	return &workerGroup{
		workerPool: sync.Pool{
			New: func() any {
				return &worker{
					pending: make(chan func(), 256),
				}
			},
		},
	}
}

type gocall interface {
	func(context.Context) | func()
}

type workerGroup struct {
	group      sync.Map
	workerPool sync.Pool
	mu         sync.Mutex
}

func (wkg *workerGroup) run(name string, fn func()) {
	wkg.mu.Lock()
	var w *worker
	value, ok := wkg.group.Load(name)
	if !ok {
		w = wkg.workerPool.Get().(*worker)
		w.name = name
		wkg.group.Store(name, w)
	} else {
		w = value.(*worker)
	}
	w.count++
	pending := w.count
	wkg.mu.Unlock()
	w.pending <- fn
	if pending == 1 {
		Go(w.work)
	}
}

type worker struct {
	name    string
	pending chan func()
	count   int32
}

func (w *worker) work() {
	for {
		select {
		case fn := <-w.pending:
			utils.ExecAndRecover(fn)
			w.count--
		default:
			wkg.mu.Lock()
			var empty bool
			if w.count == 0 {
				wkg.group.Delete(w.name)
				wkg.workerPool.Put(w)
				empty = true
			}
			wkg.mu.Unlock()
			if empty {
				return
			}
		}
	}
}

// 开启一个受到一定监督的协程
// 若fn参数包含context.Context类型, 当系统准备退出时, 此ctx会Done, 此时必须退出协程
// 系统会等候所有由Go开辟的协程退出后再退出
func Go[T gocall](fn T) {
	switch f := any(fn).(type) {
	case func(context.Context):
		proc.WaitAdd()
		go func() {
			// GoRunMark(utils.FuncName(fn))
			// defer GoDoneMark(utils.FuncName(fn))
			defer utils.RecoverPanic()
			defer proc.WaitDone()
			f(rootCtx)
		}()
	case func():
		proc.WaitAdd()
		go func() {
			// GoRunMark(utils.FuncName(fn))
			// defer GoDoneMark(utils.FuncName(fn))
			defer utils.RecoverPanic()
			defer proc.WaitDone()
			f()
		}()
	}
}

// 规则同Go, 但是可以通过name参数对协程进行分组
// 同分组下的任务会等候上一个执行完毕后再执行
func GoWithGroup(name string, fn func()) {
	wkg.run(name, fn)
}

func GoRunMark(key string) {
	running.Change(key, 1)
}

func GoDoneMark(key string) {
	running.Change(key, -1)
}

func PrintCurrentGo() {
	running.Range(func(s string, i int) {
		if i > 0 {
			zlog.Debugf("Go %s: %d", s, i)
		}
	})
}

func Async[T any](m iface.IReactor, f func() (T, error), cb func(T, error)) {
	Go(func() {
		defer utils.RecoverPanic()
		c := &asyncCtx[T]{}
		c.res, c.err = f()
		m.Assign(c)
	})
}

type asyncCtx[T any] struct {
	res T
	err error
	cb  func(T, error)
}

func (c *asyncCtx[T]) AsyncCb() {
	c.cb(c.res, c.err)
}
