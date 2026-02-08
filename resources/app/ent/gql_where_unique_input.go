package ent

import (
    "app/ent/cycle"
    "app/ent/forpermission"
    "app/ent/lifecyclenot"
    "app/ent/miau"
    "app/ent/user"
    "errors"
	"app/ent/predicate"
)

var ErrEmptyCycleWhereUniqueInput = errors.New("empty predicate CycleWhereUniqueInput")

// CycleWhereUniqueInput represents a where input for filtering Cycle queries.
type CycleWhereUniqueInput struct {
	Predicates []predicate.Cycle  `json:"-"`

	ID             *string  `json:"id,omitempty"`
}

func (i *CycleWhereUniqueInput) AddPredicates(predicates ...predicate.Cycle) {
	i.Predicates = append(i.Predicates, predicates...)
}

func (i *CycleWhereUniqueInput) Filter(q *CycleQuery) (*CycleQuery, error) {
	if i == nil {
		return q, nil
	}
	p, err := i.P()
	if err != nil {
		return nil, err
	}
	return q.Where(p), nil
}

func (i *CycleWhereUniqueInput) P() (predicate.Cycle, error) {
	var predicates []predicate.Cycle

	if i.ID != nil {
		predicates = append(predicates, cycle.IDEQ(*i.ID))
	}

	switch len(predicates) {
	case 0:
		return nil, ErrEmptyCycleWhereUniqueInput
	case 1:
		return predicates[0], nil
	default:
		return cycle.And(predicates...), nil
	}
}

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

var ErrEmptyLifeCycleNotWhereUniqueInput = errors.New("empty predicate LifeCycleNotWhereUniqueInput")

// LifeCycleNotWhereUniqueInput represents a where input for filtering LifeCycleNot queries.
type LifeCycleNotWhereUniqueInput struct {
	Predicates []predicate.LifeCycleNot  `json:"-"`

	ID             *string  `json:"id,omitempty"`
}

func (i *LifeCycleNotWhereUniqueInput) AddPredicates(predicates ...predicate.LifeCycleNot) {
	i.Predicates = append(i.Predicates, predicates...)
}

func (i *LifeCycleNotWhereUniqueInput) Filter(q *LifeCycleNotQuery) (*LifeCycleNotQuery, error) {
	if i == nil {
		return q, nil
	}
	p, err := i.P()
	if err != nil {
		return nil, err
	}
	return q.Where(p), nil
}

func (i *LifeCycleNotWhereUniqueInput) P() (predicate.LifeCycleNot, error) {
	var predicates []predicate.LifeCycleNot

	if i.ID != nil {
		predicates = append(predicates, lifecyclenot.IDEQ(*i.ID))
	}

	switch len(predicates) {
	case 0:
		return nil, ErrEmptyLifeCycleNotWhereUniqueInput
	case 1:
		return predicates[0], nil
	default:
		return lifecyclenot.And(predicates...), nil
	}
}

var ErrEmptyMiauWhereUniqueInput = errors.New("empty predicate MiauWhereUniqueInput")

// MiauWhereUniqueInput represents a where input for filtering Miau queries.
type MiauWhereUniqueInput struct {
	Predicates []predicate.Miau  `json:"-"`

	ID             *string  `json:"id,omitempty"`
}

func (i *MiauWhereUniqueInput) AddPredicates(predicates ...predicate.Miau) {
	i.Predicates = append(i.Predicates, predicates...)
}

func (i *MiauWhereUniqueInput) Filter(q *MiauQuery) (*MiauQuery, error) {
	if i == nil {
		return q, nil
	}
	p, err := i.P()
	if err != nil {
		return nil, err
	}
	return q.Where(p), nil
}

func (i *MiauWhereUniqueInput) P() (predicate.Miau, error) {
	var predicates []predicate.Miau

	if i.ID != nil {
		predicates = append(predicates, miau.IDEQ(*i.ID))
	}

	switch len(predicates) {
	case 0:
		return nil, ErrEmptyMiauWhereUniqueInput
	case 1:
		return predicates[0], nil
	default:
		return miau.And(predicates...), nil
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