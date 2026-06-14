package trie

type Node struct {
	children map[rune]*Node
	isEnd    bool
	value    any
}

type Trie struct {
	root *Node
	size int
}

func New() *Trie {
	return &Trie{
		root: &Node{children: make(map[rune]*Node)},
	}
}

func (t *Trie) Insert(word string, value any) {
	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			node.children[ch] = &Node{children: make(map[rune]*Node)}
		}
		node = node.children[ch]
	}
	if !node.isEnd {
		node.isEnd = true
		node.value = value
		t.size++
	} else {
		node.value = value
	}
}

func (t *Trie) Search(word string) (any, bool) {
	node := t.findNode(word)
	if node == nil || !node.isEnd {
		return nil, false
	}
	return node.value, true
}

func (t *Trie) StartsWith(prefix string) bool {
	if prefix == "" {
		return false
	}
	return t.findNode(prefix) != nil
}

func (t *Trie) Autocomplete(prefix string) []any {
	node := t.findNode(prefix)
	if node == nil {
		return nil
	}
	var results []any
	t.collect(node, &results)
	return results
}

func (t *Trie) Delete(word string) bool {
	node := t.findNode(word)
	if node == nil || !node.isEnd {
		return false
	}
	node.isEnd = false
	node.value = nil
	t.size--
	return true
}

func (t *Trie) Len() int {
	return t.size
}

func (t *Trie) findNode(prefix string) *Node {
	node := t.root
	for _, ch := range prefix {
		if node.children[ch] == nil {
			return nil
		}
		node = node.children[ch]
	}
	return node
}

func (t *Trie) collect(node *Node, results *[]any) {
	if node.isEnd && node.value != nil {
		*results = append(*results, node.value)
	}
	for _, child := range node.children {
		t.collect(child, results)
	}
}
