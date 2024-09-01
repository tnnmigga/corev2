package idgen

import (
	"sync"
	"time"

	"github.com/tnnmigga/corev2/conf"
)

var uuidgen UUIDGenerater

type UUIDGenerater struct {
	sync.Mutex
	timestamp uint64
	index     uint64
}

func (idgen *UUIDGenerater) NewID() uint64 {
	idgen.Lock()
	defer idgen.Unlock()
	ms := newMs()
	if ms != idgen.timestamp {
		idgen.timestamp = ms
		idgen.index = 0
	}
	idgen.index++
	if idgen.index > 0x3FF {
		panic("idgen uuid index over limit")
	}
	serverID := uint64(conf.ServerID)
	if serverID >= 0xFFF {
		panic("UUIDGenerater.NewID server-id must be smaller than 4096")
	}
	return idgen.timestamp<<40 | idgen.index | idgen.index<<10 | serverID
}

func newMs() uint64 {
	ms := time.Now().UnixMilli()
	return uint64(ms - 1700000000000)
}

// 雪花算法生成一个新的UUID
func NewUUID() uint64 {
	return uuidgen.NewID()
}
