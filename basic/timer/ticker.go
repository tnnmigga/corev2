package timer

import (
	"sync"
	"time"

	"github.com/tnnmigga/corev2/conc"
	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/system"
	"github.com/tnnmigga/corev2/utils"
	"github.com/tnnmigga/corev2/utils/heap"
	"github.com/tnnmigga/corev2/utils/idgen"
)

func init() {
	conc.Go(ticker)
}

var (
	h   = heap.New[uint64, time.Duration, *Timer]()
	mtx = sync.Mutex{}
)

type Timer struct {
	Recv iface.IModule
	ID   uint64
	Do   any
	at   time.Duration
}

func (t *Timer) Key() uint64 {
	return t.ID
}

func (t *Timer) Value() time.Duration {
	return t.at
}

func ticker() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			tryTrigger()
		case <-system.RootCtx().Done():
			return
		}
	}
}

func tryTrigger() {
	now := utils.NowNs()
	mtx.Lock()
	defer mtx.Unlock()
	for {
		item := h.Top()
		if item.at > now {
			return
		}
		h.Pop()
		item.Recv.Assign(item.Do)
	}
}

func After(m iface.IModule, do any, delay time.Duration) uint64 {
	t := &Timer{ID: idgen.NewUUID(), Do: do, at: utils.NowNs()}
	mtx.Lock()
	h.Push(t)
	mtx.Unlock()
	return t.ID
}

func Cancel(id uint64) bool {
	mtx.Lock()
	defer mtx.Unlock()
	item := h.Remove(id)
	return item.ID == id
}

func CancelByFilter(filter func(*Timer) bool) []*Timer {
	mtx.Lock()
	defer mtx.Unlock()
	var items []*Timer
	for _, item := range h.Items {
		if filter(item) {
			items = append(items, item)
		}
	}
	for _, item := range items {
		h.Remove(item.ID)
	}
	return items
}
