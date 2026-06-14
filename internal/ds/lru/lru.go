package lru

import (
	"github.com/thelastvideostore/internal/ds/list"
)

type entry[K comparable, V any] struct {
	key   K
	value V
}

type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]*list.Node[entry[K, V]]
	order    *list.List[entry[K, V]]
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity < 1 {
		capacity = 16
	}
	return &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Node[entry[K, V]]),
		order:    list.New[entry[K, V]](),
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	node, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.order.Remove(node)
	c.items[key] = c.order.PushBack(node.Value)
	return node.Value.value, true
}

func (c *Cache[K, V]) Put(key K, value V) {
	if node, ok := c.items[key]; ok {
		c.order.Remove(node)
	} else if c.order.Len >= c.capacity {
		oldest, ok := c.order.PopFront()
		if ok {
			delete(c.items, oldest.key)
		}
	}
	e := entry[K, V]{key: key, value: value}
	c.items[key] = c.order.PushBack(e)
}

func (c *Cache[K, V]) Remove(key K) bool {
	node, ok := c.items[key]
	if !ok {
		return false
	}
	c.order.Remove(node)
	delete(c.items, key)
	return true
}

func (c *Cache[K, V]) Len() int {
	return c.order.Len
}

func (c *Cache[K, V]) Contains(key K) bool {
	_, ok := c.items[key]
	return ok
}
