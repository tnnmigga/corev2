// Package mpsc provides an efficient implementation of a multi-producer, single-consumer lock-free queue.
//
// The Push function is safe to call from multiple goroutines. The Pop and Empty APIs must only be
// called from a single, consumer goroutine.
//
package algorithm

// This implementation is based on http://www.1024cores.net/home/lock-free-algorithms/queues/non-intrusive-mpsc-node-based-queue

import (
	"sync/atomic"
	"unsafe"
)

type queueNode[T any] struct {
	next *queueNode[T]
	val  T
}

type MpscQueue[T any] struct {
	head, tail *queueNode[T]
	_nil       T
}

func NewMpscQueue[T any]() *MpscQueue[T] {
	q := &MpscQueue[T]{}
	stub := &queueNode[T]{}
	q.head = stub
	q.tail = stub
	return q
}

// Push adds x to the back of the queue.
//
// Push can be safely called from multiple goroutines
func (q *MpscQueue[T]) Push(x T) {
	n := new(queueNode[T])
	n.val = x
	// current producer acquires head node
	prev := (*queueNode[T])(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(n)))

	// release node to consumer
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prev.next)), unsafe.Pointer(n))
}

// Pop removes the item from the front of the queue or nil if the queue is empty
//
// Pop must be called from a single, consumer goroutine
func (q *MpscQueue[T]) Pop() T {
	tail := q.tail
	next := (*queueNode[T])(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next)))) // acquire
	if next != nil {
		q.tail = next
		v := next.val
		return v
	}
	return q._nil
}

// Empty returns true if the queue is empty
//
// Empty must be called from a single, consumer goroutine
func (q *MpscQueue[T]) Empty() bool {
	tail := q.tail
	next := (*queueNode[T])(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next))))
	return next == nil
}
