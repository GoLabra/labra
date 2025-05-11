package ent

import (
    "errors"

	"app/ent/predicate"
    "app/ent/migration"
    "app/ent/role"
    "app/ent/user"
)

var ErrEmptyMigrationWhereUniqueInput = errors.New("empty predicate MigrationWhereUniqueInput")

// RoleWhereInput represents a where input for filtering Role queries.
type MigrationWhereUniqueInput struct {
	Predicates []predicate.Migration  `json:"-"`

	ID             *string  `json:"id,omitempty"`
}

func (i *MigrationWhereUniqueInput) AddPredicates(predicates ...predicate.Migration) {
	i.Predicates = append(i.Predicates, predicates...)
}

func (i *MigrationWhereUniqueInput) Filter(q *MigrationQuery) (*MigrationQuery, error) {
	if i == nil {
		return q, nil
	}
	p, err := i.P()
	if err != nil {
		return nil, err
	}
	return q.Where(p), nil
}

func (i *MigrationWhereUniqueInput) P() (predicate.Migration, error) {
	var predicates []predicate.Migration

	if i.ID != nil {
		predicates = append(predicates, migration.IDEQ(*i.ID))
	}

	switch len(predicates) {
	case 0:
		return nil, ErrEmptyMigrationWhereUniqueInput
	case 1:
		return predicates[0], nil
	default:
		return migration.And(predicates...), nil
	}
}

var ErrEmptyRoleWhereUniqueInput = errors.New("empty predicate RoleWhereUniqueInput")

// RoleWhereInput represents a where input for filtering Role queries.
type RoleWhereUniqueInput struct {
	Predicates []predicate.Role  `json:"-"`

	ID             *string  `json:"id,omitempty"`
}

func (i *RoleWhereUniqueInput) AddPredicates(predicates ...predicate.Role) {
	i.Predicates = append(i.Predicates, predicates...)
}

func (i *RoleWhereUniqueInput) Filter(q *RoleQuery) (*RoleQuery, error) {
	if i == nil {
		return q, nil
	}
	p, err := i.P()
	if err != nil {
		return nil, err
	}
	return q.Where(p), nil
}

func (i *RoleWhereUniqueInput) P() (predicate.Role, error) {
	var predicates []predicate.Role

	if i.ID != nil {
		predicates = append(predicates, role.IDEQ(*i.ID))
	}

	switch len(predicates) {
	case 0:
		return nil, ErrEmptyRoleWhereUniqueInput
	case 1:
		return predicates[0], nil
	default:
		return role.And(predicates...), nil
	}
}

var ErrEmptyUserWhereUniqueInput = errors.New("empty predicate UserWhereUniqueInput")

// RoleWhereInput represents a where input for filtering Role queries.
type UserWhereUniqueInput struct {
	Predicates []predicate.User  `json:"-"`

	ID             *string  `json:"id,omitempty"`
    Email *string `json:"email,omitempty"`
}

func (i *UserWhereUniqueInput) AddPredicates(predicates ...predicate.User) {
	i.Predicates = append(i.Predicates, predicates...)
}

func (i *UserWhereUniqueInput) Filter(q *UserQuery) (*UserQuery, error) {
	if i == nil {
		return q, nil
	}
	p, err := i.P()
	if err != nil {
		return nil, err
	}
	return q.Where(p), nil
}

func (i *UserWhereUniqueInput) P() (predicate.User, error) {
	var predicates []predicate.User

	if i.ID != nil {
		predicates = append(predicates, user.IDEQ(*i.ID))
	}
	if i.Email != nil {
		predicates = append(predicates, user.EmailEQ(*i.Email))
	}

	switch len(predicates) {
	case 0:
		return nil, ErrEmptyUserWhereUniqueInput
	case 1:
		return predicates[0], nil
	default:
		return user.And(predicates...), nil
	}
}