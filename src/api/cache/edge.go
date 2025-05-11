package cache

import (
	"errors"
	"time"

	"github.com/GoLabra/labrago/src/api/entgql/entity"
)

var (
	ErrEdgesNotFound = errors.New("edges not found in cache")
)

var Edge Cache[string, []entity.Edge]

func NewEdgeCache(ttl time.Duration) {
	Edge = NewCache[string, []entity.Edge](ttl)
}
