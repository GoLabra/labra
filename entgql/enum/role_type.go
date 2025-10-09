package enum

import (
	"fmt"
	"io"
	"strconv"
)

type RoleType string

const (
	RoleTypeAdmin RoleType = "Admin"
	RoleTypeApi   RoleType = "Api"
)

// Values provides list valid values for Enum.
func (RoleType) Values() (kinds []string) {
	for _, s := range []RoleType{RoleTypeAdmin, RoleTypeApi} {
		kinds = append(kinds, string(s))
	}
	return
}

func (e RoleType) String() string {
	return string(e)
}

// RoleTypeValidator is a validator for the "type" field enum values. It is called by the builders before save.
func RoleTypeValidator(e RoleType) error {
	switch e {
	case RoleTypeAdmin, RoleTypeApi:
		return nil
	default:
		return fmt.Errorf("fieldmodel: invalid enum value for field creator field: %q", e)
	}
}

// MarshalGQL implements graphql.Marshaler interface.
func (e RoleType) MarshalGQL(w io.Writer) {
	if _, err := io.WriteString(w, strconv.Quote(e.String())); err == nil {
		return
	}
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (e *RoleType) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", val)
	}
	*e = RoleType(str)
	if err := RoleTypeValidator(*e); err != nil {
		return fmt.Errorf("%s is not a valid Type", str)
	}
	return nil
}