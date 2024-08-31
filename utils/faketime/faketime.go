package faketime

import (
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/procodr/monkey"
	"github.com/tnnmigga/corev2/log"
)

var addTime = &atomic.Int64{}

func FakeNow() time.Time {
	t := &syscall.Timeval{}
	syscall.Gettimeofday(t)
	return time.Unix(t.Unix()).Add(time.Duration(addTime.Load()))
}

func init() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("faketime err %v", r)
		}
	}()
	monkey.Patch(time.Now, FakeNow)
}

func AddFakeTime(arg string) {
	n, _ := strconv.Atoi(arg[:len(arg)-1])
	add := time.Duration(n)
	switch arg[len(arg)-1] {
	case 'd':
		add *= time.Hour * 24
	case 'h':
		add *= time.Hour
	case 'm':
		add *= time.Minute
	}
	addTime.Add(int64(add))
}

func ResetFakeTime() {
	addTime.Store(0)
}
