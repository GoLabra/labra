package cache

import (
	"errors"
	"time"

	"github.com/GoLabra/labra/src/api/entgql/entity"
)

var (
	ErrFieldsNotFound = errors.New("fields not found in cache")
)

var Field Cache[string, []entity.Field]

func NewFieldCache(ttl time.Duration) {
	Field = NewCache[string, []entity.Field](ttl)
}
