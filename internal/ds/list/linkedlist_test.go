package list

import (
	"testing"
)

func TestPushBack(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	if l.Len != 3 {
		t.Fatalf("len = %d, want 3", l.Len)
	}
	if l.Head.Value != 1 {
		t.Errorf("head = %d, want 1", l.Head.Value)
	}
	if l.Tail.Value != 3 {
		t.Errorf("tail = %d, want 3", l.Tail.Value)
	}
}

func TestPushFront(t *testing.T) {
	l := New[int]()
	l.PushFront(3)
	l.PushFront(2)
	l.PushFront(1)

	if l.Len != 3 {
		t.Fatalf("len = %d, want 3", l.Len)
	}
	if l.Head.Value != 1 {
		t.Errorf("head = %d, want 1", l.Head.Value)
	}
	if l.Tail.Value != 3 {
		t.Errorf("tail = %d, want 3", l.Tail.Value)
	}
}

func TestRemove(t *testing.T) {
	l := New[int]()
	n1 := l.PushBack(1)
	n2 := l.PushBack(2)
	n3 := l.PushBack(3)

	v := l.Remove(n2)
	if v != 2 {
		t.Errorf("removed = %d, want 2", v)
	}
	if l.Len != 2 {
		t.Errorf("len = %d, want 2", l.Len)
	}
	if l.Head.Value != 1 {
		t.Errorf("head = %d, want 1", l.Head.Value)
	}
	if l.Tail.Value != 3 {
		t.Errorf("tail = %d, want 3", l.Tail.Value)
	}
	if n1.Next != n3 {
		t.Error("n1.Next should point to n3 after removal")
	}
	if n3.Prev != n1 {
		t.Error("n3.Prev should point to n1 after removal")
	}
}

func TestRemoveHead(t *testing.T) {
	l := New[int]()
	n1 := l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	l.Remove(n1)
	if l.Len != 2 {
		t.Errorf("len = %d, want 2", l.Len)
	}
	if l.Head.Value != 2 {
		t.Errorf("head = %d, want 2", l.Head.Value)
	}
}

func TestRemoveTail(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)
	n3 := l.PushBack(3)

	l.Remove(n3)
	if l.Len != 2 {
		t.Errorf("len = %d, want 2", l.Len)
	}
	if l.Tail.Value != 2 {
		t.Errorf("tail = %d, want 2", l.Tail.Value)
	}
}

func TestPopFront(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)

	v, ok := l.PopFront()
	if !ok {
		t.Fatal("PopFront should succeed")
	}
	if v != 1 {
		t.Errorf("popped = %d, want 1", v)
	}
	if l.Len != 1 {
		t.Errorf("len = %d, want 1", l.Len)
	}
}

func TestPopFrontEmpty(t *testing.T) {
	l := New[int]()
	v, ok := l.PopFront()
	if ok {
		t.Error("PopFront on empty list should return false")
	}
	if v != 0 {
		t.Errorf("zero value = %d, want 0", v)
	}
}

func TestPopBack(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)

	v, ok := l.PopBack()
	if !ok {
		t.Fatal("PopBack should succeed")
	}
	if v != 2 {
		t.Errorf("popped = %d, want 2", v)
	}
	if l.Len != 1 {
		t.Errorf("len = %d, want 1", l.Len)
	}
}

func TestPopBackEmpty(t *testing.T) {
	l := New[int]()
	v, ok := l.PopBack()
	if ok {
		t.Error("PopBack on empty list should return false")
	}
	if v != 0 {
		t.Errorf("zero value = %d, want 0", v)
	}
}

func TestFind(t *testing.T) {
	l := New[int]()
	l.PushBack(10)
	l.PushBack(20)
	l.PushBack(30)

	n := l.Find(func(v int) bool { return v == 20 })
	if n == nil {
		t.Fatal("Find should find 20")
	}
	if n.Value != 20 {
		t.Errorf("found = %d, want 20", n.Value)
	}

	n = l.Find(func(v int) bool { return v == 99 })
	if n != nil {
		t.Error("Find should return nil for missing value")
	}
}

func TestSlice(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	s := l.Slice()
	if len(s) != 3 {
		t.Fatalf("slice len = %d, want 3", len(s))
	}
	expected := []int{1, 2, 3}
	for i, v := range expected {
		if s[i] != v {
			t.Errorf("s[%d] = %d, want %d", i, s[i], v)
		}
	}
}

func TestSliceEmpty(t *testing.T) {
	l := New[int]()
	s := l.Slice()
	if len(s) != 0 {
		t.Errorf("empty slice len = %d, want 0", len(s))
	}
}

func TestIsEmpty(t *testing.T) {
	l := New[int]()
	if !l.IsEmpty() {
		t.Error("new list should be empty")
	}
	l.PushBack(1)
	if l.IsEmpty() {
		t.Error("list with item should not be empty")
	}
}

func TestNilRemove(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	v := l.Remove(nil)
	if v != 0 {
		t.Error("removing nil should return zero value")
	}
	if l.Len != 1 {
		t.Error("removing nil should not change length")
	}
}

func TestSingleElementList(t *testing.T) {
	l := New[int]()
	n := l.PushBack(42)

	if l.Head != n || l.Tail != n {
		t.Error("in single-element list, head and tail should be the same node")
	}
	if n.Prev != nil || n.Next != nil {
		t.Error("single node should have nil prev and next")
	}

	l.Remove(n)
	if l.Len != 0 || l.Head != nil || l.Tail != nil {
		t.Error("after removing only element, list should be empty")
	}
}

func BenchmarkPushBack(b *testing.B) {
	l := New[int]()
	for b.Loop() {
		l.PushBack(1)
	}
}

func BenchmarkFind(b *testing.B) {
	l := New[int]()
	for i := range 1000 {
		l.PushBack(i)
	}
	b.ResetTimer()
	for b.Loop() {
		l.Find(func(v int) bool { return v == 500 })
	}
}
