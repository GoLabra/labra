package cache

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	ttl    time.Duration
	mu     sync.RWMutex
	Values map[K]Item[V]
}

type Item[V any] struct {
	expirationTime time.Time
	Value          V
}

func NewCache[K comparable, V any](ttl time.Duration) Cache[K, V] {
	return Cache[K, V]{
		ttl:    ttl,
		mu:     sync.RWMutex{},
		Values: make(map[K]Item[V]),
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.Values[key]

	return value.Value, ok
}

func (c *Cache[K, V]) GetAll() []V {
	c.mu.RLock()
	defer c.mu.RUnlock()
	values := []V{}

	for _, value := range c.Values {
		values = append(values, value.Value)
	}

	return values
}

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

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.Values, key)
}
