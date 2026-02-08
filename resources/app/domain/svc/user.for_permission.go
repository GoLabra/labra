package svc

import (
	"context"
	"app/ent"
	"app/domain/repo"
)

type ForPermission struct {
    repository *repo.Repository
}

func NewForPermission(r *repo.Repository) *ForPermission {
	return &ForPermission{repository: r}
}


func (s *ForPermission) Get(ctx context.Context, where *ent.ForPermissionWhereInput, orderBy *ent.ForPermissionOrder, skip *int, first *int, last *int) ([]*ent.ForPermission, error) {
	return s.repository.ForPermission.Get(ctx, where, orderBy, skip, first, last)
}

func (s *ForPermission) Connection(ctx context.Context, where *ent.ForPermissionWhereInput, orderBy *ent.ForPermissionOrder, skip *int, first *int, last *int) (*ent.ForPermissionConnection, error) {
	return s.repository.ForPermission.Connection(ctx, where, orderBy, skip, first, last)
}

func (s *ForPermission) GetOne(ctx context.Context, where ent.ForPermissionWhereUniqueInput) (*ent.ForPermission, error) {
	return s.repository.ForPermission.GetOne(ctx, where)
}

func (s *ForPermission) GetTx(ctx context.Context, tx *ent.Tx, where *ent.ForPermissionWhereInput, orderBy *ent.ForPermissionOrder, skip *int, first *int, last *int) ([]*ent.ForPermission, error) {
	return s.repository.ForPermission.GetTx(ctx, tx, where, orderBy, skip, first, last)
}

func (s *ForPermission) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereUniqueInput) (*ent.ForPermission, error) {
	return s.repository.ForPermission.GetOneTx(ctx, tx, where)
}
func (s *ForPermission) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateForPermissionInput) (*ent.ForPermission, error) {
	createdInput, err := s.repository.ForPermission.CreateTx(ctx, tx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *ForPermission) Create(ctx context.Context, data ent.CreateForPermissionInput) (*ent.ForPermission, error) {
	createdInput, err := s.repository.ForPermission.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *ForPermission) CreateMany(ctx context.Context, data []ent.CreateForPermissionInput) ([]*ent.ForPermission, error) {
	return s.repository.ForPermission.CreateMany(ctx, data)
}

func (s *ForPermission) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateForPermissionInput) ([]*ent.ForPermission, error) {
	return s.repository.ForPermission.CreateManyTx(ctx, tx, data)
}

func (s *ForPermission) Update(ctx context.Context, where ent.ForPermissionWhereUniqueInput, data ent.UpdateForPermissionInput) (*ent.ForPermission, error) {
	return s.repository.ForPermission.Update(ctx, where, data)
}

func (s *ForPermission) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereUniqueInput, data ent.UpdateForPermissionInput) (*ent.ForPermission, error) {
	return s.repository.ForPermission.UpdateTx(ctx, tx, where, data)
}

func (s *ForPermission) UpdateMany(ctx context.Context, where ent.ForPermissionWhereInput, data ent.UpdateForPermissionInput) (int, error) {
	return s.repository.ForPermission.UpdateMany(ctx, where, data)
}

func (s *ForPermission) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereInput, data ent.UpdateForPermissionInput) (int, error) {
	return s.repository.ForPermission.UpdateManyTx(ctx, tx, where, data)
}

func (s *ForPermission) Upsert(ctx context.Context, data ent.CreateForPermissionInput) (*ent.ForPermission, error) {
	return s.repository.ForPermission.Upsert(ctx, data)
}

func (s *ForPermission) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateForPermissionInput) (*ent.ForPermission, error) {
	return s.repository.ForPermission.UpsertTx(ctx, tx, data)
}

func (s *ForPermission) UpsertMany(ctx context.Context, data []ent.CreateForPermissionInput) (int, error) {
	return s.repository.ForPermission.UpsertMany(ctx, data)
}

func (s *ForPermission) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateForPermissionInput) (int, error) {
	return s.repository.ForPermission.UpsertManyTx(ctx, tx, data)
}

func (s *ForPermission) Delete(ctx context.Context, where ent.ForPermissionWhereUniqueInput) (*ent.ForPermission, error) {
	return s.repository.ForPermission.Delete(ctx, where)
}

func (s *ForPermission) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereUniqueInput) (*ent.ForPermission, error) {
	return s.repository.ForPermission.DeleteTx(ctx, tx, where)
}

func (s *ForPermission) DeleteMany(ctx context.Context, where ent.ForPermissionWhereInput) (int, error) {
	return s.repository.ForPermission.DeleteMany(ctx, where)
}

func (s *ForPermission) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereInput) (int, error) {
	return s.repository.ForPermission.DeleteManyTx(ctx, tx, where)
}