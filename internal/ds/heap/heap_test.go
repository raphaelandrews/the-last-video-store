package heap

import (
	"testing"
)

func TestMinHeapOrdered(t *testing.T) {
	h := NewMin[int]()
	for _, v := range []int{5, 3, 7, 1, 4, 2, 6} {
		h.Push(v)
	}

	expected := []int{1, 2, 3, 4, 5, 6, 7}
	for _, want := range expected {
		v, ok := h.Pop()
		if !ok {
			t.Fatal("unexpected empty heap")
		}
		if v != want {
			t.Errorf("Pop = %d, want %d", v, want)
		}
	}
}

func TestMinHeapCustomLess(t *testing.T) {
	type Item struct {
		val       string
		timestamp int64
	}

	h := New(func(a, b Item) bool { return a.timestamp < b.timestamp })
	h.Push(Item{"oldest", 100})
	h.Push(Item{"newest", 300})
	h.Push(Item{"middle", 200})

	v, _ := h.Pop()
	if v.timestamp != 100 {
		t.Errorf("first pop = %d, want 100", v.timestamp)
	}
	v, _ = h.Pop()
	if v.timestamp != 200 {
		t.Errorf("second pop = %d, want 200", v.timestamp)
	}
	v, _ = h.Pop()
	if v.timestamp != 300 {
		t.Errorf("third pop = %d, want 300", v.timestamp)
	}
}

func TestEmptyHeap(t *testing.T) {
	h := NewMin[int]()

	if !h.IsEmpty() {
		t.Error("new heap should be empty")
	}
	if h.Len() != 0 {
		t.Errorf("Len = %d, want 0", h.Len())
	}

	v, ok := h.Pop()
	if ok {
		t.Errorf("Pop on empty = %d, want false", v)
	}

	v, ok = h.Peek()
	if ok {
		t.Errorf("Peek on empty = %d, want false", v)
	}
}

func TestPeek(t *testing.T) {
	h := NewMin[int]()
	h.Push(42)
	h.Push(10)

	v, ok := h.Peek()
	if !ok || v != 10 {
		t.Errorf("Peek = (%d, %v), want (10, true)", v, ok)
	}
	if h.Len() != 2 {
		t.Error("Peek should not change length")
	}
}

func TestSingleElement(t *testing.T) {
	h := NewMin[int]()
	h.Push(99)

	v, ok := h.Pop()
	if !ok || v != 99 {
		t.Errorf("Pop = (%d, %v), want (99, true)", v, ok)
	}
	if !h.IsEmpty() {
		t.Error("heap should be empty after popping only element")
	}
}

func TestDuplicateValues(t *testing.T) {
	h := NewMin[int]()
	h.Push(5)
	h.Push(5)
	h.Push(5)

	v, _ := h.Pop()
	if v != 5 {
		t.Errorf("first pop = %d, want 5", v)
	}
	v, _ = h.Pop()
	if v != 5 {
		t.Errorf("second pop = %d, want 5", v)
	}
	v, _ = h.Pop()
	if v != 5 {
		t.Errorf("third pop = %d, want 5", v)
	}
}

func BenchmarkPush(b *testing.B) {
	h := NewMin[int]()
	for b.Loop() {
		h.Push(42)
	}
}

func BenchmarkPop(b *testing.B) {
	h := NewMin[int]()
	for i := range 10000 {
		h.Push(i)
	}
	b.ResetTimer()
	for b.Loop() {
		if h.IsEmpty() {
			for i := range 10000 {
				h.Push(i)
			}
		}
		h.Pop()
	}
}
