package system

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/tnnmigga/corev2/log"
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

// 等候所有需要等待的协程退出
func WaitGoExit() {
	cancel()
	c := make(chan struct{}, 1)
	timer := time.After(time.Minute)
	go func() {
		wg.Wait()
		c <- struct{}{}
	}()
	select {
	case <-c:
		return
	case <-timer:
		// PrintCurrentGo()
		log.Errorf("wait goroutine exit timeout")
	}
}

func WaitExitSignal() os.Signal {
	signal.Notify(sign, exitSignals...)
	return <-sign
}

// 运行时故障触发进程退出流程
func Exit() {
	select {
	case sign <- syscall.SIGQUIT:
	default:
	}
}
