package pv

import "golang.org/x/exp/constraints"

// 信号量 PV操作
type PV interface {
	P()
	V()
}

type h struct {
	count chan struct{}
}

func New[N constraints.Integer](n int) PV {
	return &h{count: make(chan struct{}, n)}
}

func (s *h) P() {
	s.count <- struct{}{}
}

func (s *h) V() {
	<-s.count
}
