package ent

import (
	"app/ent/predicate"
	"app/ent/user"
	"app/ent/userrefreshtoken"
	"errors"
)

var ErrEmptyUserWhereUniqueInput = errors.New("empty predicate UserWhereUniqueInput")

// UserWhereUniqueInput represents a where input for filtering User queries.
type UserWhereUniqueInput struct {
	Predicates []predicate.User `json:"-"`

	ID    *string `json:"id,omitempty"`
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

var ErrEmptyUserRefreshTokenWhereUniqueInput = errors.New("empty predicate UserRefreshTokenWhereUniqueInput")

// UserRefreshTokenWhereUniqueInput represents a where input for filtering UserRefreshToken queries.
type UserRefreshTokenWhereUniqueInput struct {
	Predicates []predicate.UserRefreshToken `json:"-"`

	ID        *int    `json:"id,omitempty"`
	Tokenhash *string `json:"tokenHash,omitempty"`
}

func (i *UserRefreshTokenWhereUniqueInput) AddPredicates(predicates ...predicate.UserRefreshToken) {
	i.Predicates = append(i.Predicates, predicates...)
}

func (i *UserRefreshTokenWhereUniqueInput) Filter(q *UserRefreshTokenQuery) (*UserRefreshTokenQuery, error) {
	if i == nil {
		return q, nil
	}
	p, err := i.P()
	if err != nil {
		return nil, err
	}
	return q.Where(p), nil
}

func (i *UserRefreshTokenWhereUniqueInput) P() (predicate.UserRefreshToken, error) {
	var predicates []predicate.UserRefreshToken

	if i.ID != nil {
		predicates = append(predicates, userrefreshtoken.IDEQ(*i.ID))
	}
	if i.Tokenhash != nil {
		predicates = append(predicates, userrefreshtoken.TokenHashEQ(*i.Tokenhash))
	}

	switch len(predicates) {
	case 0:
		return nil, ErrEmptyUserRefreshTokenWhereUniqueInput
	case 1:
		return predicates[0], nil
	default:
		return userrefreshtoken.And(predicates...), nil
	}
}
