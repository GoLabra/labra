package ent

import (
	"errors"

	"github.com/GoLabra/labra/src/api/entgql/ent/file"
	"github.com/GoLabra/labra/src/api/entgql/ent/permission"
	"github.com/GoLabra/labra/src/api/entgql/ent/predicate"
	"github.com/GoLabra/labra/src/api/entgql/ent/role"
	"github.com/GoLabra/labra/src/api/entgql/ent/user"
)

var ErrEmptyFileWhereUniqueInput = errors.New("empty predicate FileWhereUniqueInput")

// RoleWhereInput represents a where input for filtering Role queries.
type FileWhereUniqueInput struct {
	Predicates []predicate.File `json:"-"`

	ID              *string `json:"id,omitempty"`
	Caption         *string `json:"caption,omitempty"`
	Name            *string `json:"name,omitempty"`
	StorageFileName *string `json:"storageFileName,omitempty"`
}

func (i *FileWhereUniqueInput) AddPredicates(predicates ...predicate.File) {
	i.Predicates = append(i.Predicates, predicates...)
}

func (i *FileWhereUniqueInput) Filter(q *FileQuery) (*FileQuery, error) {
	if i == nil {
		return q, nil
	}
	p, err := i.P()
	if err != nil {
		return nil, err
	}
	return q.Where(p), nil
}

func (i *FileWhereUniqueInput) P() (predicate.File, error) {
	var predicates []predicate.File

	if i.ID != nil {
		predicates = append(predicates, file.IDEQ(*i.ID))
	}
	if i.Caption != nil {
		predicates = append(predicates, file.CaptionEQ(*i.Caption))
	}
	if i.Name != nil {
		predicates = append(predicates, file.NameEQ(*i.Name))
	}
	if i.StorageFileName != nil {
		predicates = append(predicates, file.StorageFileNameEQ(*i.StorageFileName))
	}

	switch len(predicates) {
	case 0:
		return nil, ErrEmptyFileWhereUniqueInput
	case 1:
		return predicates[0], nil
	default:
		return file.And(predicates...), nil
	}
}

var ErrEmptyPermissionWhereUniqueInput = errors.New("empty predicate PermissionWhereUniqueInput")

// RoleWhereInput represents a where input for filtering Role queries.
type PermissionWhereUniqueInput struct {
	Predicates []predicate.Permission `json:"-"`

	ID *string `json:"id,omitempty"`
}

func (i *PermissionWhereUniqueInput) AddPredicates(predicates ...predicate.Permission) {
	i.Predicates = append(i.Predicates, predicates...)
}

func (i *PermissionWhereUniqueInput) Filter(q *PermissionQuery) (*PermissionQuery, error) {
	if i == nil {
		return q, nil
	}
	p, err := i.P()
	if err != nil {
		return nil, err
	}
	return q.Where(p), nil
}

func (i *PermissionWhereUniqueInput) P() (predicate.Permission, error) {
	var predicates []predicate.Permission

	if i.ID != nil {
		predicates = append(predicates, permission.IDEQ(*i.ID))
	}

	switch len(predicates) {
	case 0:
		return nil, ErrEmptyPermissionWhereUniqueInput
	case 1:
		return predicates[0], nil
	default:
		return permission.And(predicates...), nil
	}
}

var ErrEmptyRoleWhereUniqueInput = errors.New("empty predicate RoleWhereUniqueInput")

// RoleWhereInput represents a where input for filtering Role queries.
type RoleWhereUniqueInput struct {
	Predicates []predicate.Role `json:"-"`

	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
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
	if i.Name != nil {
		predicates = append(predicates, role.NameEQ(*i.Name))
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
