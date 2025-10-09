package ent

import (
    "app/ent/forpermission"
    "app/ent/user"
    "errors"
	"app/ent/predicate"
)

var ErrEmptyForPermissionWhereUniqueInput = errors.New("empty predicate ForPermissionWhereUniqueInput")

// ForPermissionWhereUniqueInput represents a where input for filtering ForPermission queries.
type ForPermissionWhereUniqueInput struct {
	Predicates []predicate.ForPermission  `json:"-"`

	ID             *string  `json:"id,omitempty"`
}

func (i *ForPermissionWhereUniqueInput) AddPredicates(predicates ...predicate.ForPermission) {
	i.Predicates = append(i.Predicates, predicates...)
}

func (i *ForPermissionWhereUniqueInput) Filter(q *ForPermissionQuery) (*ForPermissionQuery, error) {
	if i == nil {
		return q, nil
	}
	p, err := i.P()
	if err != nil {
		return nil, err
	}
	return q.Where(p), nil
}

func (i *ForPermissionWhereUniqueInput) P() (predicate.ForPermission, error) {
	var predicates []predicate.ForPermission

	if i.ID != nil {
		predicates = append(predicates, forpermission.IDEQ(*i.ID))
	}

	switch len(predicates) {
	case 0:
		return nil, ErrEmptyForPermissionWhereUniqueInput
	case 1:
		return predicates[0], nil
	default:
		return forpermission.And(predicates...), nil
	}
}

var ErrEmptyUserWhereUniqueInput = errors.New("empty predicate UserWhereUniqueInput")

// UserWhereUniqueInput represents a where input for filtering User queries.
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