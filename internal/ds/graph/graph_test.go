package graph

import (
	"testing"
)

func TestAddVertex(t *testing.T) {
	g := New()
	g.AddVertex("The Matrix")

	if !g.HasVertex("The Matrix") {
		t.Error("HasVertex should return true after AddVertex")
	}
	if g.VertexCount() != 1 {
		t.Errorf("VertexCount = %d, want 1", g.VertexCount())
	}
}

func TestAddEdge(t *testing.T) {
	g := New()
	g.AddEdge("The Matrix", "Inception")

	if g.VertexCount() != 2 {
		t.Errorf("VertexCount = %d, want 2", g.VertexCount())
	}
	if g.EdgeCount() != 1 {
		t.Errorf("EdgeCount = %d, want 1", g.EdgeCount())
	}

	neighbors := g.GetNeighbors("The Matrix")
	if neighbors["Inception"] != 1 {
		t.Errorf("edge weight = %d, want 1", neighbors["Inception"])
	}

	neighbors = g.GetNeighbors("Inception")
	if neighbors["The Matrix"] != 1 {
		t.Error("edge should be undirected")
	}
}

func TestIncrementEdge(t *testing.T) {
	g := New()
	g.IncrementEdge("The Matrix", "Inception")
	g.IncrementEdge("The Matrix", "Inception")
	g.IncrementEdge("The Matrix", "Inception")

	neighbors := g.GetNeighbors("The Matrix")
	if neighbors["Inception"] != 3 {
		t.Errorf("edge weight = %d, want 3", neighbors["Inception"])
	}
}

func TestSelfEdge(t *testing.T) {
	g := New()
	g.IncrementEdge("The Matrix", "The Matrix")

	if g.EdgeCount() != 0 {
		t.Error("self-edges should be ignored")
	}
}

func TestGetNeighborsMissing(t *testing.T) {
	g := New()
	n := g.GetNeighbors("nonexistent")
	if n != nil {
		t.Error("GetNeighbors on missing vertex should return nil")
	}
}

func TestBFS(t *testing.T) {
	g := New()
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "D")

	order := g.BFS("A")
	if len(order) != 4 {
		t.Errorf("BFS order len = %d, want 4", len(order))
	}
	if order[0] != "A" {
		t.Errorf("BFS[0] = %s, want A", order[0])
	}
}

func TestBFSMissing(t *testing.T) {
	g := New()
	order := g.BFS("nonexistent")
	if order != nil {
		t.Error("BFS on missing vertex should return nil")
	}
}

func TestGetRecommendations(t *testing.T) {
	g := New()
	g.IncrementEdge("A", "B")
	g.IncrementEdge("A", "B")
	g.IncrementEdge("A", "B")
	g.IncrementEdge("A", "C")
	g.IncrementEdge("A", "C")
	g.IncrementEdge("A", "D")

	recs := g.GetRecommendations("A", 2)
	if len(recs) != 2 {
		t.Fatalf("GetRecommendations len = %d, want 2", len(recs))
	}
	if recs[0] != "B" {
		t.Errorf("first recommendation = %s, want B (weight 3)", recs[0])
	}
	if recs[1] != "C" {
		t.Errorf("second recommendation = %s, want C (weight 2)", recs[1])
	}
}

func TestGetRecommendationsMissing(t *testing.T) {
	g := New()
	recs := g.GetRecommendations("nonexistent", 5)
	if recs != nil {
		t.Error("GetRecommendations on missing vertex should return nil")
	}
}

func TestEmptyGraph(t *testing.T) {
	g := New()
	if g.VertexCount() != 0 {
		t.Errorf("VertexCount = %d, want 0", g.VertexCount())
	}
	if g.EdgeCount() != 0 {
		t.Errorf("EdgeCount = %d, want 0", g.EdgeCount())
	}
	if g.HasVertex("anything") {
		t.Error("empty graph should have no vertices")
	}
}

func BenchmarkBuildGraph(b *testing.B) {
	g := New()
	vertices := make([]string, 100)
	for i := range 100 {
		vertices[i] = string(rune('A'+i%26)) + string(rune('0'+i/26))
	}
	for i := range vertices {
		for j := i + 1; j < len(vertices); j++ {
			if (i+j)%3 == 0 {
				g.AddEdge(vertices[i], vertices[j])
			}
		}
	}
	b.ResetTimer()
	for b.Loop() {
		g.GetRecommendations(vertices[0], 5)
	}
}
