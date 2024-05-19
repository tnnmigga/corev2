package conc

// 信号量 PV操作
type Semaphore interface {
	P()
	V()
}

type semaphore struct {
	count chan struct{}
}

func NewSemaphore(n int) Semaphore {
	return &semaphore{count: make(chan struct{}, n)}
}

func (s *semaphore) P() {
	s.count <- struct{}{}
}

func (s *semaphore) V() {
	<-s.count
}
