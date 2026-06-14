package graph

type Vertex struct {
	ID    string
	Edges map[string]int
	Data  any
}

type Graph struct {
	vertices map[string]*Vertex
}

func New() *Graph {
	return &Graph{
		vertices: make(map[string]*Vertex),
	}
}

func (g *Graph) AddVertex(id string) {
	if _, exists := g.vertices[id]; !exists {
		g.vertices[id] = &Vertex{
			ID:    id,
			Edges: make(map[string]int),
		}
	}
}

func (g *Graph) AddEdge(v1, v2 string) {
	g.AddVertex(v1)
	g.AddVertex(v2)
	g.vertices[v1].Edges[v2] = 1
	g.vertices[v2].Edges[v1] = 1
}

func (g *Graph) IncrementEdge(v1, v2 string) {
	if v1 == v2 {
		return
	}
	g.AddVertex(v1)
	g.AddVertex(v2)
	g.vertices[v1].Edges[v2]++
	g.vertices[v2].Edges[v1]++
}

func (g *Graph) GetNeighbors(id string) map[string]int {
	v, ok := g.vertices[id]
	if !ok {
		return nil
	}
	return v.Edges
}

func (g *Graph) BFS(start string) []string {
	if _, ok := g.vertices[start]; !ok {
		return nil
	}
	visited := make(map[string]bool)
	var order []string
	queue := []string{start}
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		order = append(order, current)

		for neighbor := range g.vertices[current].Edges {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	return order
}

func (g *Graph) GetRecommendations(id string, k int) []string {
	if _, ok := g.vertices[id]; !ok {
		return nil
	}

	type pair struct {
		id     string
		weight int
	}

	var candidates []pair
	for neighbor, weight := range g.vertices[id].Edges {
		if weight > 0 {
			candidates = append(candidates, pair{neighbor, weight})
		}
	}

	for i := 0; i < len(candidates); i++ {
		maxIdx := i
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].weight > candidates[maxIdx].weight {
				maxIdx = j
			}
		}
		candidates[i], candidates[maxIdx] = candidates[maxIdx], candidates[i]
	}

	if k > len(candidates) {
		k = len(candidates)
	}

	result := make([]string, k)
	for i := range k {
		result[i] = candidates[i].id
	}
	return result
}

func (g *Graph) HasVertex(id string) bool {
	_, ok := g.vertices[id]
	return ok
}

func (g *Graph) VertexCount() int {
	return len(g.vertices)
}

func (g *Graph) EdgeCount() int {
	count := 0
	for _, v := range g.vertices {
		count += len(v.Edges)
	}
	return count / 2
}
