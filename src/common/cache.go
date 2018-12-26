package common

import "sync"

// Cache is simple in-memory cache implementation.
type Cache struct {
	mtx  sync.RWMutex
	data map[interface{}]interface{}
	size int
}

// NewCache creates new in-memory cache.
func NewCache(cacheSize int) *Cache {
	return &Cache{
		data: make(map[interface{}]interface{}),
		size: cacheSize,
	}
}

// Add adds item to in-memory cache.
func (c *Cache) Add(key interface{}, value interface{}) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.size > c.len() {
		c.data[key] = value
		return true
	}
	return false
}

// Get gets item from in-memory cache.
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.get(key)
}

// Remove removes item from in-memory cache.
func (c *Cache) Remove(key interface{}) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if _, ok := c.get(key); !ok {
		return false
	}
	delete(c.data, key)
	return true
}

// Purge clears a in-memory cache.
func (c *Cache) Purge() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.data = make(map[interface{}]interface{})
}

// Len gets the number of items in a in-memory cache.
func (c *Cache) Len() int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.len()
}

func (c *Cache) get(key interface{}) (interface{}, bool) {
	if val, ok := c.data[key]; ok {
		return val, true
	}
	return nil, false
}

func (c *Cache) len() int {
	return len(c.data)
}
