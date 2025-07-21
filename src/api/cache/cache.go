// Package cache provides generic in-memory caching utilities for the Labra API.
// Located at src/api/cache, it stores frequently accessed entities to reduce
// database lookups.
package cache

import (
	"sync"
	"time"
)

// Cache stores key-value pairs with an optional TTL. It is thread safe and used
// by various services to hold frequently accessed objects.
type Cache[K comparable, V any] struct {
	ttl    time.Duration
	mu     sync.RWMutex
	Values map[K]Item[V]
}

// Item represents a cached value with its expiration time.
type Item[V any] struct {
	expirationTime time.Time
	Value          V
}

// NewCache initializes a Cache with the provided TTL.
func NewCache[K comparable, V any](ttl time.Duration) Cache[K, V] {
	return Cache[K, V]{
		ttl:    ttl,
		mu:     sync.RWMutex{},
		Values: make(map[K]Item[V]),
	}
}

// Get returns a cached value and a flag indicating whether it exists.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.Values[key]

	return value.Value, ok
}

// GetAll retrieves all values in the cache in no particular order.
func (c *Cache[K, V]) GetAll() []V {
	c.mu.RLock()
	defer c.mu.RUnlock()
	values := []V{}

	for _, value := range c.Values {
		values = append(values, value.Value)
	}

	return values
}

// Set adds or replaces a value for the given key and sets its expiration time.
func (c *Cache[K, V]) Set(key K, val V) {
	var expirationTime time.Time

	if c.ttl > 0 {
		expirationTime = time.Now().Add(c.ttl)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.Values[key] = Item[V]{
		expirationTime: expirationTime,
		Value:          val,
	}
}

// Delete removes the cached value for the provided key.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.Values, key)
}
