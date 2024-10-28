package idgen

import (
	"sync/atomic"
	"time"

	"github.com/tnnmigga/corev2/conf"
	"github.com/tnnmigga/corev2/log"
)

var (
	uuidLastAt  atomic.Uint64
	uuidLastIdx atomic.Uint64
)

// 雪花算法生成一个新的UUID
func NewUUID() uint64 {
	index := uuidLastIdx.Add(1)
	if index >= 128 {
		log.Errorf("NewUUID too fast")
		time.Sleep(time.Millisecond)
	}
	ms := uint64(time.Now().UnixMilli())
	if ms != uuidLastAt.Load() {
		uuidLastAt.Store(ms)
		uuidLastIdx.Store(0)
		index = 0
	}
	serverID := uint64(conf.ServerID)
	return ms<<19 | index<<12 | serverID
}
