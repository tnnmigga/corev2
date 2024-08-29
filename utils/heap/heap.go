package heap

import (
	"golang.org/x/exp/constraints"
)

type heapItem[K comparable, V constraints.Ordered] interface {
	Key() K
	Value() V
}

func New[K comparable, V constraints.Ordered]() *Heap[K, V, heapItem[K, V]] {
	return &Heap[K, V, heapItem[K, V]]{}
}

// 小顶堆
type Heap[K comparable, V constraints.Ordered, T heapItem[K, V]] struct {
	Items []T
}

func (h *Heap[K, V, T]) Top() (top T) {
	if h.Len() != 0 {
		top = h.Items[0]
	}
	return top
}

func (h *Heap[K, V, T]) Push(x T) {
	h.Items = append(h.Items, x)
	h.up(h.Len() - 1)
}

func (h *Heap[K, V, T]) Pop() (item T) {
	if h.Len() == 0 {
		return
	}
	item = h.Items[0]
	n := h.Len() - 1
	h.swap(0, n)
	h.down(0, n)
	h.Items = h.Items[:n]
	return item
}

func (h *Heap[K, V, T]) Remove(key K) (item T) {
	i := h.Find(key)
	return h.RemoveByIndex(i)
}

func (h *Heap[K, V, T]) RemoveByIndex(i int) (item T) {
	if i < 0 || i >= h.Len() {
		return
	}
	n := h.Len() - 1
	if n != i {
		h.swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	item = h.Items[n]
	h.Items = h.Items[:n]
	return item
}

func (h *Heap[K, V, T]) Len() int {
	return len(h.Items)
}

func (h *Heap[K, V, T]) Find(key K) int {
	for index, item := range h.Items {
		if item.Key() == key {
			return index
		}
	}
	return -1
}

func (h *Heap[K, V, T]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *Heap[K, V, T]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(j, i) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}

func (h *Heap[K, V, T]) swap(i, j int) {
	h.Items[i], h.Items[j] = h.Items[j], h.Items[i]
}

func (h *Heap[K, V, T]) less(i, j int) bool {
	return h.Items[i].Value() < h.Items[j].Value()
}
