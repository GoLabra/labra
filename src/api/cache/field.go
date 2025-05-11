package cache

import (
	"errors"
	"github.com/GoLabra/labrago/src/api/entgql/entity"
	"time"
)

var (
	ErrFieldsNotFound = errors.New("fields not found in cache")
)

var Field Cache[string, []entity.Field]

func NewFieldCache(ttl time.Duration) {
	Field = NewCache[string, []entity.Field](ttl)
}
