package deque

type Deque[T any] struct {
	buf  []T
	head int
	tail int
	size int
}

func New[T any](capacity int) *Deque[T] {
	if capacity < 1 {
		capacity = 8
	}
	return &Deque[T]{
		buf: make([]T, capacity),
	}
}

func (d *Deque[T]) Len() int {
	return d.size
}

func (d *Deque[T]) IsEmpty() bool {
	return d.size == 0
}

func (d *Deque[T]) PushBack(v T) {
	if d.size == len(d.buf) {
		d.grow()
	}
	d.buf[d.tail] = v
	d.tail = (d.tail + 1) % len(d.buf)
	d.size++
}

func (d *Deque[T]) PushFront(v T) {
	if d.size == len(d.buf) {
		d.grow()
	}
	d.head = (d.head - 1 + len(d.buf)) % len(d.buf)
	d.buf[d.head] = v
	d.size++
}

func (d *Deque[T]) PopFront() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	v := d.buf[d.head]
	var zero T
	d.buf[d.head] = zero
	d.head = (d.head + 1) % len(d.buf)
	d.size--
	return v, true
}

func (d *Deque[T]) PopBack() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	d.tail = (d.tail - 1 + len(d.buf)) % len(d.buf)
	v := d.buf[d.tail]
	var zero T
	d.buf[d.tail] = zero
	d.size--
	return v, true
}

func (d *Deque[T]) PeekFront() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	return d.buf[d.head], true
}

func (d *Deque[T]) PeekBack() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	idx := (d.tail - 1 + len(d.buf)) % len(d.buf)
	return d.buf[idx], true
}

func (d *Deque[T]) grow() {
	newBuf := make([]T, len(d.buf)*2)
	if d.head < d.tail {
		copy(newBuf, d.buf[d.head:d.tail])
	} else {
		n := copy(newBuf, d.buf[d.head:])
		copy(newBuf[n:], d.buf[:d.tail])
	}
	d.buf = newBuf
	d.head = 0
	d.tail = d.size
}
