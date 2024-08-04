package system

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	exitSignals     = []os.Signal{syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM}
	sign            = make(chan os.Signal, 1)
	wg              = sync.WaitGroup{}
	rootCtx, cancel = context.WithCancel(context.Background())
)

func RootCtx() context.Context {
	return rootCtx
}

func WaitAdd() {
	wg.Add(1)
}

func WaitDone() {
	wg.Done()
}

// 等候所有由Go开辟的协程退出
func WaitGoDone(maxWaitTime time.Duration) error {
	cancel()
	c := make(chan struct{}, 1)
	timer := time.After(maxWaitTime)
	go func() {
		wg.Wait()
		c <- struct{}{}
	}()
	select {
	case <-c:
		return nil
	case <-timer:
		// PrintCurrentGo()
		return fmt.Errorf("wait goroutine exit timeout")
	}
}

func WaitExitSignal() os.Signal {
	signal.Notify(sign, exitSignals...)
	return <-sign
}

// 运行时故障触发进程退出流程
func Abort() {
	select {
	case sign <- syscall.SIGQUIT:
	default:
	}
}
