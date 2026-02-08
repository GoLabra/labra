package svc

import (
	"context"
	"app/ent"
	"app/domain/repo"

	"github.com/GoLabra/labra/entgql/entity"
)

type LifeCycleNot struct {
    repository *repo.Repository
}

func NewLifeCycleNot(r *repo.Repository) *LifeCycleNot {
	return &LifeCycleNot{repository: r}
}


func (s *LifeCycleNot) Get(ctx context.Context, where *ent.LifeCycleNotWhereInput, orderBy *ent.LifeCycleNotOrder, skip *int, first *int, last *int) ([]*ent.LifeCycleNot, error) {
	return s.repository.LifeCycleNot.Get(ctx, where, orderBy, skip, first, last)
}

func (s *LifeCycleNot) Connection(ctx context.Context, where *ent.LifeCycleNotWhereInput, orderBy *ent.LifeCycleNotOrder, skip *int, first *int, last *int) (*ent.LifeCycleNotConnection, error) {
	return s.repository.LifeCycleNot.Connection(ctx, where, orderBy, skip, first, last)
}

func (s *LifeCycleNot) GetOne(ctx context.Context, where ent.LifeCycleNotWhereUniqueInput) (*ent.LifeCycleNot, error) {
	return s.repository.LifeCycleNot.GetOne(ctx, where)
}

func (s *LifeCycleNot) GetTx(ctx context.Context, tx *ent.Tx, where *ent.LifeCycleNotWhereInput, orderBy *ent.LifeCycleNotOrder, skip *int, first *int, last *int) ([]*ent.LifeCycleNot, error) {
	return s.repository.LifeCycleNot.GetTx(ctx, tx, where, orderBy, skip, first, last)
}

func (s *LifeCycleNot) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereUniqueInput) (*ent.LifeCycleNot, error) {
	return s.repository.LifeCycleNot.GetOneTx(ctx, tx, where)
}
func (s *LifeCycleNot) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateLifeCycleNotInput) (*ent.LifeCycleNot, error) {
	createdInput, err := s.repository.LifeCycleNot.CreateTx(ctx, tx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *LifeCycleNot) Create(ctx context.Context, data ent.CreateLifeCycleNotInput) (*ent.LifeCycleNot, error) {
	createdInput, err := s.repository.LifeCycleNot.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *LifeCycleNot) CreateMany(ctx context.Context, data []ent.CreateLifeCycleNotInput) ([]*ent.LifeCycleNot, error) {
	return s.repository.LifeCycleNot.CreateMany(ctx, data)
}

func (s *LifeCycleNot) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateLifeCycleNotInput) ([]*ent.LifeCycleNot, error) {
	return s.repository.LifeCycleNot.CreateManyTx(ctx, tx, data)
}

func (s *LifeCycleNot) Update(ctx context.Context, where ent.LifeCycleNotWhereUniqueInput, data ent.UpdateLifeCycleNotInput) (*ent.LifeCycleNot, error) {
	return s.repository.LifeCycleNot.Update(ctx, where, data)
}
func (s *LifeCycleNot) UpdateState(ctx context.Context, where ent.LifeCycleNotWhereUniqueInput, state entity.EntityState) (*ent.LifeCycleNot, error) {
	return s.repository.LifeCycleNot.UpdateState(ctx, where, state)
}

func (s *LifeCycleNot) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereUniqueInput, data ent.UpdateLifeCycleNotInput) (*ent.LifeCycleNot, error) {
	return s.repository.LifeCycleNot.UpdateTx(ctx, tx, where, data)
}

func (s *LifeCycleNot) UpdateMany(ctx context.Context, where ent.LifeCycleNotWhereInput, data ent.UpdateLifeCycleNotInput) (int, error) {
	return s.repository.LifeCycleNot.UpdateMany(ctx, where, data)
}

func (s *LifeCycleNot) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereInput, data ent.UpdateLifeCycleNotInput) (int, error) {
	return s.repository.LifeCycleNot.UpdateManyTx(ctx, tx, where, data)
}

func (s *LifeCycleNot) Upsert(ctx context.Context, data ent.CreateLifeCycleNotInput) (*ent.LifeCycleNot, error) {
	return s.repository.LifeCycleNot.Upsert(ctx, data)
}

func (s *LifeCycleNot) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateLifeCycleNotInput) (*ent.LifeCycleNot, error) {
	return s.repository.LifeCycleNot.UpsertTx(ctx, tx, data)
}

func (s *LifeCycleNot) UpsertMany(ctx context.Context, data []ent.CreateLifeCycleNotInput) (int, error) {
	return s.repository.LifeCycleNot.UpsertMany(ctx, data)
}

func (s *LifeCycleNot) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateLifeCycleNotInput) (int, error) {
	return s.repository.LifeCycleNot.UpsertManyTx(ctx, tx, data)
}

func (s *LifeCycleNot) Delete(ctx context.Context, where ent.LifeCycleNotWhereUniqueInput) (*ent.LifeCycleNot, error) {
	return s.repository.LifeCycleNot.Delete(ctx, where)
}

func (s *LifeCycleNot) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereUniqueInput) (*ent.LifeCycleNot, error) {
	return s.repository.LifeCycleNot.DeleteTx(ctx, tx, where)
}

func (s *LifeCycleNot) DeleteMany(ctx context.Context, where ent.LifeCycleNotWhereInput) (int, error) {
	return s.repository.LifeCycleNot.DeleteMany(ctx, where)
}

func (s *LifeCycleNot) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereInput) (int, error) {
	return s.repository.LifeCycleNot.DeleteManyTx(ctx, tx, where)
}