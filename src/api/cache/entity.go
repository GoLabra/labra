// Package cache stores entities, fields and edges for quick retrieval. This
// file contains cache logic specific to entities.
package cache

import (
	"fmt"
	"time"

	"github.com/GoLabra/labra/src/api/entgql/entity"
)

// ErrEntityNotFound indicates that an entity was not found in the cache.
var (
	ErrEntityNotFound = fmt.Errorf("entity not found in cache")
)

// EntityCache wraps Cache to provide lookups by various unique inputs.
type EntityCache struct {
	Cache[string, entity.Entity]
}

// GetByUniqueInput searches for an entity based on its unique fields.
func (c *EntityCache) GetByUniqueInput(where entity.EntityWhereUniqueInput) (entity.Entity, bool) {
	if where.Name != nil {
		return c.Get(*where.Name)
	}

	for _, entity := range c.Values {
		if where.Caption != nil && entity.Value.Caption == *where.Caption {
			return entity.Value, true
		}
		if where.EntName != nil && entity.Value.EntName == *where.EntName {
			return entity.Value, true
		}
	}

	return entity.Entity{}, false
}

// Entity is the global cache of entities keyed by name.
var Entity EntityCache

// NewEntityCache initializes the global Entity cache with the provided TTL.
func NewEntityCache(ttl time.Duration) {
	Entity = EntityCache{
		Cache: NewCache[string, entity.Entity](ttl),
	}
}
