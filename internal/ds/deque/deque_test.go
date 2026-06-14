package deque

import "testing"

func TestPushBackPopFront(t *testing.T) {
	d := New[int](4)
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)

	v, ok := d.PopFront()
	if !ok || v != 1 {
		t.Errorf("PopFront = (%d, %v), want (1, true)", v, ok)
	}
	v, ok = d.PopFront()
	if !ok || v != 2 {
		t.Errorf("PopFront = (%d, %v), want (2, true)", v, ok)
	}
	v, ok = d.PopFront()
	if !ok || v != 3 {
		t.Errorf("PopFront = (%d, %v), want (3, true)", v, ok)
	}
}

func TestPushFrontPopFront(t *testing.T) {
	d := New[int](4)
	d.PushFront(3)
	d.PushFront(2)
	d.PushFront(1)

	v, _ := d.PopFront()
	if v != 1 {
		t.Errorf("PopFront = %d, want 1", v)
	}
	v, _ = d.PopFront()
	if v != 2 {
		t.Errorf("PopFront = %d, want 2", v)
	}
	v, _ = d.PopFront()
	if v != 3 {
		t.Errorf("PopFront = %d, want 3", v)
	}
}

func TestPushBackPopBack(t *testing.T) {
	d := New[int](4)
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)

	v, _ := d.PopBack()
	if v != 3 {
		t.Errorf("PopBack = %d, want 3", v)
	}
	v, _ = d.PopBack()
	if v != 2 {
		t.Errorf("PopBack = %d, want 2", v)
	}
	v, _ = d.PopBack()
	if v != 1 {
		t.Errorf("PopBack = %d, want 1", v)
	}
}

func TestMixedOperations(t *testing.T) {
	d := New[int](4)
	d.PushBack(1)
	d.PushFront(0)
	d.PushBack(2)

	v, _ := d.PopFront()
	if v != 0 {
		t.Errorf("PopFront = %d, want 0", v)
	}
	v, _ = d.PopBack()
	if v != 2 {
		t.Errorf("PopBack = %d, want 2", v)
	}
	v, _ = d.PopFront()
	if v != 1 {
		t.Errorf("PopFront = %d, want 1", v)
	}
}

func TestEmptyDeque(t *testing.T) {
	d := New[int](4)

	if !d.IsEmpty() {
		t.Error("new deque should be empty")
	}
	if d.Len() != 0 {
		t.Errorf("Len = %d, want 0", d.Len())
	}

	_, ok := d.PopFront()
	if ok {
		t.Error("PopFront on empty should return false")
	}

	_, ok = d.PopBack()
	if ok {
		t.Error("PopBack on empty should return false")
	}

	_, ok = d.PeekFront()
	if ok {
		t.Error("PeekFront on empty should return false")
	}

	_, ok = d.PeekBack()
	if ok {
		t.Error("PeekBack on empty should return false")
	}
}

func TestPeek(t *testing.T) {
	d := New[int](4)
	d.PushBack(1)
	d.PushBack(2)

	v, ok := d.PeekFront()
	if !ok || v != 1 {
		t.Errorf("PeekFront = (%d, %v), want (1, true)", v, ok)
	}
	if d.Len() != 2 {
		t.Error("Peek should not change length")
	}

	v, ok = d.PeekBack()
	if !ok || v != 2 {
		t.Errorf("PeekBack = (%d, %v), want (2, true)", v, ok)
	}
	if d.Len() != 2 {
		t.Error("Peek should not change length")
	}
}

func TestWrapAround(t *testing.T) {
	d := New[int](4)
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)
	d.PushBack(4)

	v, _ := d.PopFront()
	if v != 1 {
		t.Errorf("PopFront = %d, want 1", v)
	}
	v, _ = d.PopFront()
	if v != 2 {
		t.Errorf("PopFront = %d, want 2", v)
	}

	d.PushBack(5)
	d.PushBack(6)

	if d.Len() != 4 {
		t.Errorf("Len = %d, want 4", d.Len())
	}

	v, _ = d.PopFront()
	if v != 3 {
		t.Errorf("PopFront = %d, want 3", v)
	}
	v, _ = d.PopFront()
	if v != 4 {
		t.Errorf("PopFront = %d, want 4", v)
	}
	v, _ = d.PopFront()
	if v != 5 {
		t.Errorf("PopFront = %d, want 5", v)
	}
	v, _ = d.PopFront()
	if v != 6 {
		t.Errorf("PopFront = %d, want 6", v)
	}
}

func TestGrow(t *testing.T) {
	d := New[int](2)
	for i := range 10 {
		d.PushBack(i)
	}
	if d.Len() != 10 {
		t.Errorf("Len = %d, want 10", d.Len())
	}
	for i := range 10 {
		v, ok := d.PopFront()
		if !ok || v != i {
			t.Errorf("PopFront = (%d, %v), want (%d, true)", v, ok, i)
		}
	}
}

func TestZeroCapacityDefaults(t *testing.T) {
	d := New[int](0)
	if cap(d.buf) != 8 {
		t.Errorf("zero capacity should default to 8, got %d", cap(d.buf))
	}
}

func BenchmarkPushPop(b *testing.B) {
	d := New[int](1024)
	b.Run("PushBack", func(b *testing.B) {
		for b.Loop() {
			d.PushBack(42)
		}
	})
	b.Run("PushFront", func(b *testing.B) {
		d2 := New[int](1024)
		for b.Loop() {
			d2.PushFront(42)
		}
	})
}
