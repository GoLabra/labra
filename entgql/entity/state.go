package entity

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type EntityState string

const (
	EntityStateDraft     EntityState = "DRAFT"
	EntityStatePublished EntityState = "PUBLISHED"
	EntityStateArchived  EntityState = "ARCHIVED"
)

var AllEntityState = []EntityState{
	EntityStateDraft,
	EntityStatePublished,
	EntityStateArchived,
}

func (e EntityState) IsValid() bool {
	switch e {
	case EntityStateDraft, EntityStatePublished, EntityStateArchived:
		return true
	}
	return false
}

func (e EntityState) String() string {
	return string(e)
}

func (e *EntityState) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EntityState(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EntityState", str)
	}
	return nil
}

func (e EntityState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *EntityState) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e EntityState) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}

// LifecycleStateGetter is implemented by generated ent models that have entity state (lifecycle) enabled.
// Used by hooks to check lifecycle state without reflection.
type LifecycleStateGetter interface {
	GetLifecycleState() EntityState
}
