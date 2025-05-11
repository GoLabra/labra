package svc

import (
	"context"
	"app/ent"
    "app/domain/repo"
)

type Role struct {
    repository *repo.Repository
}

func NewRole(r *repo.Repository) *Role {
	return &Role{repository: r}
}


func (s *Role) Get(ctx context.Context, where *ent.RoleWhereInput, orderBy *ent.RoleOrder, skip *int, first *int, last *int) ([]*ent.Role, error) {
	return s.repository.Role.Get(ctx, where, orderBy, skip, first, last)
}

func (s *Role) Connection(ctx context.Context, where *ent.RoleWhereInput, orderBy *ent.RoleOrder, skip *int, first *int, last *int) (*ent.RoleConnection, error) {
	return s.repository.Role.Connection(ctx, where, orderBy, skip, first, last)
}

func (s *Role) GetOne(ctx context.Context, where ent.RoleWhereUniqueInput) (*ent.Role, error) {
	return s.repository.Role.GetOne(ctx, where)
}

func (s *Role) GetTx(ctx context.Context, tx *ent.Tx, where *ent.RoleWhereInput, orderBy *ent.RoleOrder, skip *int, first *int, last *int) ([]*ent.Role, error) {
	return s.repository.Role.GetTx(ctx, tx, where, orderBy, skip, first, last)
}

func (s *Role) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereUniqueInput) (*ent.Role, error) {
	return s.repository.Role.GetOneTx(ctx, tx, where)
}

func (s *Role) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateRoleInput) (*ent.Role, error) {
	createdInput, err := s.repository.Role.CreateTx(ctx, tx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *Role) Create(ctx context.Context, data ent.CreateRoleInput) (*ent.Role, error) {
	createdInput, err := s.repository.Role.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *Role) CreateMany(ctx context.Context, data []ent.CreateRoleInput) ([]*ent.Role, error) {
	return s.repository.Role.CreateMany(ctx, data)
}

func (s *Role) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateRoleInput) ([]*ent.Role, error) {
	return s.repository.Role.CreateManyTx(ctx, tx, data)
}

func (s *Role) Update(ctx context.Context, where ent.RoleWhereUniqueInput, data ent.UpdateRoleInput) (*ent.Role, error) {
	return s.repository.Role.Update(ctx, where, data)
}

func (s *Role) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereUniqueInput, data ent.UpdateRoleInput) (*ent.Role, error) {
	return s.repository.Role.UpdateTx(ctx, tx, where, data)
}

func (s *Role) UpdateMany(ctx context.Context, where ent.RoleWhereInput, data ent.UpdateRoleInput) (int, error) {
	return s.repository.Role.UpdateMany(ctx, where, data)
}

func (s *Role) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereInput, data ent.UpdateRoleInput) (int, error) {
	return s.repository.Role.UpdateManyTx(ctx, tx, where, data)
}

func (s *Role) Upsert(ctx context.Context, data ent.CreateRoleInput) (*ent.Role, error) {
	return s.repository.Role.Upsert(ctx, data)
}

func (s *Role) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateRoleInput) (*ent.Role, error) {
	return s.repository.Role.UpsertTx(ctx, tx, data)
}

func (s *Role) UpsertMany(ctx context.Context, data []ent.CreateRoleInput) (int, error) {
	return s.repository.Role.UpsertMany(ctx, data)
}

func (s *Role) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateRoleInput) (int, error) {
	return s.repository.Role.UpsertManyTx(ctx, tx, data)
}

func (s *Role) Delete(ctx context.Context, where ent.RoleWhereUniqueInput) (*ent.Role, error) {
	return s.repository.Role.Delete(ctx, where)
}

func (s *Role) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereUniqueInput) (*ent.Role, error) {
	return s.repository.Role.DeleteTx(ctx, tx, where)
}

func (s *Role) DeleteMany(ctx context.Context, where ent.RoleWhereInput) (int, error) {
	return s.repository.Role.DeleteMany(ctx, where)
}

func (s *Role) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereInput) (int, error) {
	return s.repository.Role.DeleteManyTx(ctx, tx, where)
}