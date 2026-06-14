package list

type Node[T any] struct {
	Value T
	Prev  *Node[T]
	Next  *Node[T]
}

type List[T any] struct {
	Head *Node[T]
	Tail *Node[T]
	Len  int
}

func New[T any]() *List[T] {
	return &List[T]{}
}

func (l *List[T]) PushBack(v T) *Node[T] {
	n := &Node[T]{Value: v}
	if l.Tail == nil {
		l.Head = n
		l.Tail = n
	} else {
		n.Prev = l.Tail
		l.Tail.Next = n
		l.Tail = n
	}
	l.Len++
	return n
}

func (l *List[T]) PushFront(v T) *Node[T] {
	n := &Node[T]{Value: v}
	if l.Head == nil {
		l.Head = n
		l.Tail = n
	} else {
		n.Next = l.Head
		l.Head.Prev = n
		l.Head = n
	}
	l.Len++
	return n
}

func (l *List[T]) Remove(n *Node[T]) T {
	if n == nil {
		var zero T
		return zero
	}

	if n.Prev != nil {
		n.Prev.Next = n.Next
	} else {
		l.Head = n.Next
	}

	if n.Next != nil {
		n.Next.Prev = n.Prev
	} else {
		l.Tail = n.Prev
	}

	l.Len--
	n.Prev = nil
	n.Next = nil
	return n.Value
}

func (l *List[T]) PopFront() (T, bool) {
	if l.Head == nil {
		var zero T
		return zero, false
	}
	v := l.Remove(l.Head)
	return v, true
}

func (l *List[T]) PopBack() (T, bool) {
	if l.Tail == nil {
		var zero T
		return zero, false
	}
	v := l.Remove(l.Tail)
	return v, true
}

func (l *List[T]) Find(pred func(T) bool) *Node[T] {
	for n := l.Head; n != nil; n = n.Next {
		if pred(n.Value) {
			return n
		}
	}
	return nil
}

func (l *List[T]) Slice() []T {
	out := make([]T, 0, l.Len)
	for n := l.Head; n != nil; n = n.Next {
		out = append(out, n.Value)
	}
	return out
}

func (l *List[T]) IsEmpty() bool {
	return l.Len == 0
}
