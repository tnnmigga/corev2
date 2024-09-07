package utils

import (
	"errors"
	"sync"
	"time"
)

type WaitGroupWithTimeout struct {
	sync.WaitGroup
}

func (wg *WaitGroupWithTimeout) WaitWithTimeout(timeout time.Duration) error {
	sign := make(chan struct{}, 1)
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	go func() {
		wg.Wait()
		sign <- struct{}{}
	}()
	select {
	case <-sign:
		return nil
	case <-timer.C:
		return errors.New("timeout")
	}
}

func NowNs() time.Duration {
	return time.Duration(time.Now().UnixNano())
}
