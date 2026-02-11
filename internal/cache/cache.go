package cache

import (
	"container/list"
	"fmt"
	"sync"
)

type entry struct {
	key   string
	value interface{}
}

type Cache struct {
	mu       sync.RWMutex
	capacity int
	items    map[string]*list.Element
	lru      *list.List
}

func New(capacity int) *Cache {
	if capacity <= 0 {
		panic(fmt.Sprintf("cache capacity must be positive, got %d", capacity))
	}

	return &Cache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		lru:      list.New(),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, found := c.items[key]
	if !found {
		return nil, false
	}

	c.lru.MoveToFront(elem)
	return elem.Value.(*entry).value, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		c.lru.MoveToFront(elem)
		elem.Value.(*entry).value = value
		return
	}

	if c.lru.Len() >= c.capacity {
		c.evictOldest()
	}

	e := &entry{key: key, value: value}
	elem := c.lru.PushFront(e)
	c.items[key] = elem
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		c.removeElement(elem)
	}
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element)
	c.lru.Init()
}

func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lru.Len()
}

func (c *Cache) evictOldest() {
	oldest := c.lru.Back()
	if oldest != nil {
		c.removeElement(oldest)
	}
}

func (c *Cache) removeElement(elem *list.Element) {
	c.lru.Remove(elem)
	e := elem.Value.(*entry)
	delete(c.items, e.key)
}
