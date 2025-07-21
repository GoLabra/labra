// Package cache defines caches for field metadata used by the API. The file is
// located at src/api/cache.
package cache

import (
	"errors"
	"time"

	"github.com/GoLabra/labra/src/api/entgql/entity"
)

// ErrFieldsNotFound is returned when field metadata is missing from the cache.
var (
	ErrFieldsNotFound = errors.New("fields not found in cache")
)

// Field holds cached slices of field definitions keyed by entity name.
var Field Cache[string, []entity.Field]

// NewFieldCache initializes the global Field cache with a TTL.
func NewFieldCache(ttl time.Duration) {
	Field = NewCache[string, []entity.Field](ttl)
}
