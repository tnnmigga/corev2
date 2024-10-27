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
	h   = heap.New[uint64, time.Duration, *timer]()
	mtx = sync.Mutex{}
)

type timer struct {
	m  iface.IModule
	id uint64
	do any
	at time.Duration
}

func (t *timer) Key() uint64 {
	return t.id
}

func (t *timer) Value() time.Duration {
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
		item.m.Assign(item.do)
	}
}

func After(m iface.IModule, do any, delay time.Duration) uint64 {
	t := &timer{id: idgen.NewUUID(), do: do, at: utils.NowNs()}
	mtx.Lock()
	h.Push(t)
	mtx.Unlock()
	return t.id
}

func Cancel(id uint64) bool {
	mtx.Lock()
	item := h.Remove(id)
	mtx.Unlock()
	return item.id == id
}

func CancelAll() []*timer {
	mtx.Lock()
	defer mtx.Unlock()
	items := h.Items
	h.Items = nil
	return items
}
