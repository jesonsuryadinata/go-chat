package main

import (
	"sync"
)

type Cache struct {
	data map[string]interface{}
	mu   sync.Mutex // for thread safety
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.data[key]
	return val, ok
}
