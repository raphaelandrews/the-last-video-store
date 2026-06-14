package trie

import (
	"testing"
)

func TestInsertAndSearch(t *testing.T) {
	tr := New()
	tr.Insert("matrix", "movie-1")
	tr.Insert("matilda", "movie-2")

	v, ok := tr.Search("matrix")
	if !ok || v != "movie-1" {
		t.Errorf("Search(matrix) = (%v, %v), want (movie-1, true)", v, ok)
	}

	v, ok = tr.Search("matilda")
	if !ok || v != "movie-2" {
		t.Errorf("Search(matilda) = (%v, %v), want (movie-2, true)", v, ok)
	}
}

func TestSearchMissing(t *testing.T) {
	tr := New()
	tr.Insert("matrix", "movie-1")

	_, ok := tr.Search("mat")
	if ok {
		t.Error("Search(mat) should return false for partial match")
	}

	_, ok = tr.Search("matrixx")
	if ok {
		t.Error("Search(matrixx) should return false")
	}

	_, ok = tr.Search("pulp")
	if ok {
		t.Error("Search(pulp) should return false")
	}
}

func TestStartsWith(t *testing.T) {
	tr := New()
	tr.Insert("matrix", "movie-1")
	tr.Insert("matilda", "movie-2")

	if !tr.StartsWith("mat") {
		t.Error("StartsWith(mat) should be true")
	}
	if !tr.StartsWith("matrix") {
		t.Error("StartsWith(matrix) should be true")
	}
	if tr.StartsWith("max") {
		t.Error("StartsWith(max) should be false")
	}
	if tr.StartsWith("") {
		t.Error("StartsWith(empty) should return false")
	}
}

func TestAutocomplete(t *testing.T) {
	tr := New()
	tr.Insert("matrix", "movie-1")
	tr.Insert("matilda", "movie-2")
	tr.Insert("match point", "movie-3")
	tr.Insert("pulp fiction", "movie-4")

	results := tr.Autocomplete("mat")
	if len(results) != 3 {
		t.Errorf("Autocomplete(mat) len = %d, want 3", len(results))
	}

	results = tr.Autocomplete("pul")
	if len(results) != 1 {
		t.Errorf("Autocomplete(pul) len = %d, want 1", len(results))
	}

	results = tr.Autocomplete("zzz")
	if results != nil {
		t.Error("Autocomplete(zzz) should return nil")
	}
}

func TestDelete(t *testing.T) {
	tr := New()
	tr.Insert("matrix", "movie-1")
	tr.Insert("matilda", "movie-2")

	if !tr.Delete("matrix") {
		t.Error("Delete(matrix) should succeed")
	}

	_, ok := tr.Search("matrix")
	if ok {
		t.Error("Search(matrix) should fail after delete")
	}

	_, ok = tr.Search("matilda")
	if !ok {
		t.Error("Search(matilda) should still succeed after deleting matrix")
	}

	if !tr.StartsWith("mat") {
		t.Error("StartsWith(mat) should still be true (matilda remains)")
	}

	if tr.Delete("matrix") {
		t.Error("Delete(matrix) twice should return false")
	}
	if tr.Delete("nonexistent") {
		t.Error("Delete(nonexistent) should return false")
	}
}

func TestCaseSensitivity(t *testing.T) {
	tr := New()
	tr.Insert("Matrix", "m1")
	tr.Insert("matrix", "m2")

	v, _ := tr.Search("Matrix")
	if v != "m1" {
		t.Errorf("Search(Matrix) = %v, want m1", v)
	}
	v, _ = tr.Search("matrix")
	if v != "m2" {
		t.Errorf("Search(matrix) = %v, want m2", v)
	}
}

func TestUpdateValue(t *testing.T) {
	tr := New()
	tr.Insert("matrix", "v1")
	tr.Insert("matrix", "v2")

	if tr.Len() != 1 {
		t.Errorf("Len = %d, want 1 after updating same key", tr.Len())
	}

	v, _ := tr.Search("matrix")
	if v != "v2" {
		t.Errorf("updated value = %v, want v2", v)
	}
}

func BenchmarkInsert(b *testing.B) {
	tr := New()
	words := make([]string, 5000)
	for i := range 5000 {
		words[i] = string(rune('a'+i%26)) + string(rune('a'+(i/26)%26)) + string(rune('a'+(i/676)%26))
	}
	b.ResetTimer()
	b.Run("Insert5000", func(b *testing.B) {
		for b.Loop() {
			for _, w := range words {
				tr.Insert(w, w)
			}
		}
	})
}

func BenchmarkAutocomplete(b *testing.B) {
	tr := New()
	for i := range 5000 {
		w := string(rune('a'+i%26)) + string(rune('a'+(i/26)%26)) + string(rune('a'+(i/676)%26))
		tr.Insert(w, w)
	}
	b.ResetTimer()
	for b.Loop() {
		tr.Autocomplete("aa")
	}
}
