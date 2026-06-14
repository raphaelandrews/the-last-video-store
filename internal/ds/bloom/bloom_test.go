package bloom

import (
	"testing"
)

func TestAddAndContains(t *testing.T) {
	bf := New(1024, 3)
	bf.Add([]byte("banned-user-1"))

	if !bf.Contains([]byte("banned-user-1")) {
		t.Error("Contains should return true after Add")
	}
}

func TestEmptyFilter(t *testing.T) {
	bf := New(1024, 3)
	if bf.Contains([]byte("anything")) {
		t.Error("empty filter should not contain anything")
	}
}

func TestMultipleAdds(t *testing.T) {
	bf := New(1024, 3)
	items := []string{"user-1", "user-2", "user-3", "user-4", "user-5"}

	for _, item := range items {
		bf.Add([]byte(item))
	}

	for _, item := range items {
		if !bf.Contains([]byte(item)) {
			t.Errorf("Contains(%s) should return true after Add", item)
		}
	}
}

func TestFalsePositiveRate(t *testing.T) {
	size := uint64(10000)
	hashCount := 3
	bf := New(size, hashCount)

	for i := range 1000 {
		bf.Add([]byte{byte(i >> 8), byte(i & 0xff)})
	}

	falsePositives := 0
	tests := 10000
	for i := 1001; i < 1001+tests; i++ {
		if bf.Contains([]byte{byte(i >> 8), byte(i & 0xff)}) {
			falsePositives++
		}
	}

	rate := float64(falsePositives) / float64(tests)
	t.Logf("False positive rate: %.4f (%d/%d)", rate, falsePositives, tests)

	if rate > 0.1 {
		t.Errorf("False positive rate %.4f exceeds 10%%", rate)
	}
}

func TestEmptyData(t *testing.T) {
	bf := New(1024, 3)
	bf.Add([]byte{})
	if !bf.Contains([]byte{}) {
		t.Error("empty data should be found after Add")
	}
}

func BenchmarkAdd(b *testing.B) {
	bf := New(100000, 3)
	for b.Loop() {
		bf.Add([]byte("test-item"))
	}
}

func BenchmarkContains(b *testing.B) {
	bf := New(100000, 3)
	for i := range 10000 {
		bf.Add([]byte{byte(i >> 8), byte(i & 0xff)})
	}
	b.ResetTimer()
	for b.Loop() {
		bf.Contains([]byte{0, 5})
	}
}
