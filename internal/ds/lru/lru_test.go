package lru

import (
	"testing"
)

func TestPutAndGet(t *testing.T) {
	c := New[string, int](4)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)

	v, ok := c.Get("a")
	if !ok || v != 1 {
		t.Errorf("Get(a) = (%d, %v), want (1, true)", v, ok)
	}
	v, ok = c.Get("b")
	if !ok || v != 2 {
		t.Errorf("Get(b) = (%d, %v), want (2, true)", v, ok)
	}
}

func TestGetMissing(t *testing.T) {
	c := New[string, int](4)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("Get on missing key should return false")
	}
}

func TestEviction(t *testing.T) {
	c := New[string, int](3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	c.Put("d", 4)

	_, ok := c.Get("a")
	if ok {
		t.Error("a should have been evicted (LRU)")
	}

	v, ok := c.Get("b")
	if !ok || v != 2 {
		t.Errorf("b should still be in cache, got (%d, %v)", v, ok)
	}
	v, ok = c.Get("c")
	if !ok || v != 3 {
		t.Errorf("c should still be in cache, got (%d, %v)", v, ok)
	}
	v, ok = c.Get("d")
	if !ok || v != 4 {
		t.Errorf("d should be in cache, got (%d, %v)", v, ok)
	}
}

func TestUpdateMovesToFront(t *testing.T) {
	c := New[string, int](3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)

	c.Get("a")
	c.Put("d", 4)

	_, ok := c.Get("b")
	if ok {
		t.Error("b should be evicted after accessing a")
	}
	v, ok := c.Get("a")
	if !ok || v != 1 {
		t.Errorf("a should still be in cache after access, got (%d, %v)", v, ok)
	}
}

func TestRemove(t *testing.T) {
	c := New[string, int](4)
	c.Put("a", 1)
	c.Put("b", 2)

	if !c.Remove("a") {
		t.Error("Remove existing key should return true")
	}
	_, ok := c.Get("a")
	if ok {
		t.Error("Get after remove should return false")
	}
	if c.Remove("a") {
		t.Error("Remove twice should return false")
	}
	if c.Len() != 1 {
		t.Errorf("Len = %d, want 1", c.Len())
	}
}

func TestContains(t *testing.T) {
	c := New[string, int](4)
	c.Put("a", 1)

	if !c.Contains("a") {
		t.Error("Contains should return true for existing key")
	}
	if c.Contains("b") {
		t.Error("Contains should return false for missing key")
	}
}

func TestLen(t *testing.T) {
	c := New[string, int](3)
	if c.Len() != 0 {
		t.Errorf("initial Len = %d, want 0", c.Len())
	}
	c.Put("a", 1)
	c.Put("b", 2)
	if c.Len() != 2 {
		t.Errorf("Len = %d, want 2", c.Len())
	}
}

func BenchmarkPut(b *testing.B) {
	c := New[int, int](1000)
	for b.Loop() {
		c.Put(b.N, b.N)
	}
}

func BenchmarkGet(b *testing.B) {
	c := New[int, int](1000)
	for i := range 1000 {
		c.Put(i, i)
	}
	b.ResetTimer()
	for b.Loop() {
		c.Get(b.N % 1000)
	}
}
