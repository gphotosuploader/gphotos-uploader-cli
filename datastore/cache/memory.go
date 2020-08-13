package cache

import (
	"sync"
)

type memory struct {
	mu    sync.RWMutex
	items map[string]*item
}

func NewMemoryCache() Cache {
	return &memory{
		items: make(map[string]*item),
	}
}

// Get returns an item from the cache.
// If there is no item on the cache will return an ErrNotFound.
func (c *memory) Get(key string) (interface{}, error) {
	c.mu.RLock()
	i, exists := c.items[key]
	c.mu.RUnlock()
	if !exists {
		return nil, ErrNotFound
	}
	return i.Object, nil
}

// Put adds an item to the cache, replacing any existing item.
func (c *memory) Put(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &item{Object:value}
	return nil
}

// Invalidate deletes an item from the cache.
func (c *memory) Invalidate(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

// Exists returns if an item is present in the cache.
func (c *memory) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.items[key]
	return exists
}
