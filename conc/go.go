package conc

import (
	"sync"

	"github.com/tnnmigga/corev2/system"
	"github.com/tnnmigga/corev2/utils"
)

var (
	wkg = newWorkerGroup()
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

// 开启一个安全的协程
func Go(fn func()) {
	system.WaitAdd()
	go func() {
		defer utils.RecoverPanic()
		defer system.WaitDone()
		fn()
	}()
}

// 规则同Go, 但是可以通过name参数对协程进行分组
// 同分组下的任务会等候上一个执行完毕后再执行
func GoWithGroup(name string, fn func()) {
	wkg.run(name, fn)
}
