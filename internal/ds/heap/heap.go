package heap

import "cmp"

type Heap[T any] struct {
	items []T
	less  func(a, b T) bool
}

func New[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{less: less}
}

func NewMin[T cmp.Ordered]() *Heap[T] {
	return &Heap[T]{
		less: func(a, b T) bool { return a < b },
	}
}

func (h *Heap[T]) Len() int {
	return len(h.items)
}

func (h *Heap[T]) IsEmpty() bool {
	return len(h.items) == 0
}

func (h *Heap[T]) Peek() (T, bool) {
	if h.IsEmpty() {
		var zero T
		return zero, false
	}
	return h.items[0], true
}

func (h *Heap[T]) Push(v T) {
	h.items = append(h.items, v)
	h.siftUp(len(h.items) - 1)
}

func (h *Heap[T]) Pop() (T, bool) {
	if h.IsEmpty() {
		var zero T
		return zero, false
	}
	n := len(h.items)
	v := h.items[0]
	h.items[0] = h.items[n-1]
	var zero T
	h.items[n-1] = zero
	h.items = h.items[:n-1]
	if len(h.items) > 0 {
		h.siftDown(0)
	}
	return v, true
}

func (h *Heap[T]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if !h.less(h.items[i], h.items[parent]) {
			break
		}
		h.items[i], h.items[parent] = h.items[parent], h.items[i]
		i = parent
	}
}

func (h *Heap[T]) siftDown(i int) {
	n := len(h.items)
	for {
		smallest := i
		left := 2*i + 1
		right := 2*i + 2

		if left < n && h.less(h.items[left], h.items[smallest]) {
			smallest = left
		}
		if right < n && h.less(h.items[right], h.items[smallest]) {
			smallest = right
		}
		if smallest == i {
			break
		}
		h.items[i], h.items[smallest] = h.items[smallest], h.items[i]
		i = smallest
	}
}
