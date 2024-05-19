package algorithm

import "sync"

type Counter[T comparable] struct {
	count map[T]int
	m     sync.Mutex
}

func NewCounter[T comparable]() *Counter[T] {
	return &Counter[T]{count: map[T]int{}}
}

func (c *Counter[T]) Change(key T, n int) {
	c.m.Lock()
	defer c.m.Unlock()
	c.count[key] += n
}

func (c *Counter[T]) Count(key T) int {
	c.m.Lock()
	defer c.m.Unlock()
	return c.count[key]
}

func (c *Counter[T]) Keys() []T {
	c.m.Lock()
	defer c.m.Unlock()
	keys := make([]T, 0, len(c.count))
	for k := range c.count {
		keys = append(keys, k)
	}
	return keys
}

func (c *Counter[T]) Len() int {
	c.m.Lock()
	defer c.m.Unlock()
	return len(c.count)
}

func (c *Counter[T]) Range(f func(T, int)) {
	c.m.Lock()
	defer c.m.Unlock()
	for k, v := range c.count {
		f(k, v)
	}
}
