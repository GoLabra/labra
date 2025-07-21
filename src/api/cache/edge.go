// Package cache provides caches for edges, fields, and entities in the API.
// This file defines an edge-specific cache stored in src/api/cache.
package cache

import (
	"errors"
	"time"

	"github.com/GoLabra/labra/src/api/entgql/entity"
)

// ErrEdgesNotFound is returned when requested edges are absent from the cache.
var (
	ErrEdgesNotFound = errors.New("edges not found in cache")
)

// Edge holds cached edges keyed by entity name.
var Edge Cache[string, []entity.Edge]

// NewEdgeCache initializes the global Edge cache with the given TTL.
func NewEdgeCache(ttl time.Duration) {
	Edge = NewCache[string, []entity.Edge](ttl)
}
