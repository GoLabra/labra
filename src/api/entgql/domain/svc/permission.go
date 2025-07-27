package svc

import (
	"context"

	"github.com/GoLabra/labra/src/api/entgql/domain/repo"
	"github.com/GoLabra/labra/src/api/entgql/ent"
)

type Permission struct {
	repository *repo.Repository
}

func NewPermission(r *repo.Repository) *Permission {
	return &Permission{repository: r}
}

func (s *Permission) Get(ctx context.Context, where *ent.PermissionWhereInput, orderBy *ent.PermissionOrder, skip *int, first *int, last *int) ([]*ent.Permission, error) {
	return s.repository.Permission.Get(ctx, where, orderBy, skip, first, last)
}

func (s *Permission) Connection(ctx context.Context, where *ent.PermissionWhereInput, orderBy *ent.PermissionOrder, skip *int, first *int, last *int) (*ent.PermissionConnection, error) {
	return s.repository.Permission.Connection(ctx, where, orderBy, skip, first, last)
}

func (s *Permission) GetOne(ctx context.Context, where ent.PermissionWhereUniqueInput) (*ent.Permission, error) {
	return s.repository.Permission.GetOne(ctx, where)
}

func (s *Permission) GetTx(ctx context.Context, tx *ent.Tx, where *ent.PermissionWhereInput, orderBy *ent.PermissionOrder, skip *int, first *int, last *int) ([]*ent.Permission, error) {
	return s.repository.Permission.GetTx(ctx, tx, where, orderBy, skip, first, last)
}

func (s *Permission) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereUniqueInput) (*ent.Permission, error) {
	return s.repository.Permission.GetOneTx(ctx, tx, where)
}

func (s *Permission) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreatePermissionInput) (*ent.Permission, error) {
	createdInput, err := s.repository.Permission.CreateTx(ctx, tx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *Permission) Create(ctx context.Context, data ent.CreatePermissionInput) (*ent.Permission, error) {
	createdInput, err := s.repository.Permission.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *Permission) CreateMany(ctx context.Context, data []ent.CreatePermissionInput) ([]*ent.Permission, error) {
	return s.repository.Permission.CreateMany(ctx, data)
}

func (s *Permission) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreatePermissionInput) ([]*ent.Permission, error) {
	return s.repository.Permission.CreateManyTx(ctx, tx, data)
}

func (s *Permission) Update(ctx context.Context, where ent.PermissionWhereUniqueInput, data ent.UpdatePermissionInput) (*ent.Permission, error) {
	return s.repository.Permission.Update(ctx, where, data)
}

func (s *Permission) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereUniqueInput, data ent.UpdatePermissionInput) (*ent.Permission, error) {
	return s.repository.Permission.UpdateTx(ctx, tx, where, data)
}

func (s *Permission) UpdateMany(ctx context.Context, where ent.PermissionWhereInput, data ent.UpdatePermissionInput) (int, error) {
	return s.repository.Permission.UpdateMany(ctx, where, data)
}

func (s *Permission) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereInput, data ent.UpdatePermissionInput) (int, error) {
	return s.repository.Permission.UpdateManyTx(ctx, tx, where, data)
}

func (s *Permission) Upsert(ctx context.Context, data ent.CreatePermissionInput) (*ent.Permission, error) {
	return s.repository.Permission.Upsert(ctx, data)
}

func (s *Permission) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreatePermissionInput) (*ent.Permission, error) {
	return s.repository.Permission.UpsertTx(ctx, tx, data)
}

func (s *Permission) UpsertMany(ctx context.Context, data []ent.CreatePermissionInput) (int, error) {
	return s.repository.Permission.UpsertMany(ctx, data)
}

func (s *Permission) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreatePermissionInput) (int, error) {
	return s.repository.Permission.UpsertManyTx(ctx, tx, data)
}

func (s *Permission) Delete(ctx context.Context, where ent.PermissionWhereUniqueInput) (*ent.Permission, error) {
	return s.repository.Permission.Delete(ctx, where)
}

func (s *Permission) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereUniqueInput) (*ent.Permission, error) {
	return s.repository.Permission.DeleteTx(ctx, tx, where)
}

func (s *Permission) DeleteMany(ctx context.Context, where ent.PermissionWhereInput) (int, error) {
	return s.repository.Permission.DeleteMany(ctx, where)
}

func (s *Permission) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereInput) (int, error) {
	return s.repository.Permission.DeleteManyTx(ctx, tx, where)
}
