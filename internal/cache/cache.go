package cache

import (
    "sync"
)

type Cache interface {
    Get(key string) (float64, bool)
    Set(key string, value float64)
}

type memoryCache struct {
    mu    sync.RWMutex
    items map[string]float64
}

func NewMemoryCache() Cache {
    return &memoryCache{
        items: make(map[string]float64),
    }
}

func (c *memoryCache) Get(key string) (float64, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, exists := c.items[key]
    return value, exists
}

func (c *memoryCache) Set(key string, value float64) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = value
}