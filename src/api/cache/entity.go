package cache

import (
	"fmt"
	"time"

	"github.com/GoLabra/labrago/src/api/entgql/entity"
)

var (
	ErrEntityNotFound = fmt.Errorf("entity not found in cache")
)

type EntityCache struct {
	Cache[string, entity.Entity]
}

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

var Entity EntityCache

func NewEntityCache(ttl time.Duration) {
	Entity = EntityCache{
		Cache: NewCache[string, entity.Entity](ttl),
	}
}
