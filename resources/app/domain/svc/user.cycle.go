package svc

import (
	"context"
	"app/ent"
	"app/domain/repo"

	"github.com/GoLabra/labra/entgql/entity"
)

type Cycle struct {
    repository *repo.Repository
}

func NewCycle(r *repo.Repository) *Cycle {
	return &Cycle{repository: r}
}


func (s *Cycle) Get(ctx context.Context, where *ent.CycleWhereInput, orderBy *ent.CycleOrder, skip *int, first *int, last *int) ([]*ent.Cycle, error) {
	return s.repository.Cycle.Get(ctx, where, orderBy, skip, first, last)
}

func (s *Cycle) Connection(ctx context.Context, where *ent.CycleWhereInput, orderBy *ent.CycleOrder, skip *int, first *int, last *int) (*ent.CycleConnection, error) {
	return s.repository.Cycle.Connection(ctx, where, orderBy, skip, first, last)
}

func (s *Cycle) GetOne(ctx context.Context, where ent.CycleWhereUniqueInput) (*ent.Cycle, error) {
	return s.repository.Cycle.GetOne(ctx, where)
}

func (s *Cycle) GetTx(ctx context.Context, tx *ent.Tx, where *ent.CycleWhereInput, orderBy *ent.CycleOrder, skip *int, first *int, last *int) ([]*ent.Cycle, error) {
	return s.repository.Cycle.GetTx(ctx, tx, where, orderBy, skip, first, last)
}

func (s *Cycle) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereUniqueInput) (*ent.Cycle, error) {
	return s.repository.Cycle.GetOneTx(ctx, tx, where)
}
func (s *Cycle) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateCycleInput) (*ent.Cycle, error) {
	createdInput, err := s.repository.Cycle.CreateTx(ctx, tx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *Cycle) Create(ctx context.Context, data ent.CreateCycleInput) (*ent.Cycle, error) {
	createdInput, err := s.repository.Cycle.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *Cycle) CreateMany(ctx context.Context, data []ent.CreateCycleInput) ([]*ent.Cycle, error) {
	return s.repository.Cycle.CreateMany(ctx, data)
}

func (s *Cycle) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateCycleInput) ([]*ent.Cycle, error) {
	return s.repository.Cycle.CreateManyTx(ctx, tx, data)
}

func (s *Cycle) Update(ctx context.Context, where ent.CycleWhereUniqueInput, data ent.UpdateCycleInput) (*ent.Cycle, error) {
	return s.repository.Cycle.Update(ctx, where, data)
}
func (s *Cycle) UpdateState(ctx context.Context, where ent.CycleWhereUniqueInput, state entity.EntityState) (*ent.Cycle, error) {
	return s.repository.Cycle.UpdateState(ctx, where, state)
}

func (s *Cycle) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereUniqueInput, data ent.UpdateCycleInput) (*ent.Cycle, error) {
	return s.repository.Cycle.UpdateTx(ctx, tx, where, data)
}

func (s *Cycle) UpdateMany(ctx context.Context, where ent.CycleWhereInput, data ent.UpdateCycleInput) (int, error) {
	return s.repository.Cycle.UpdateMany(ctx, where, data)
}

func (s *Cycle) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereInput, data ent.UpdateCycleInput) (int, error) {
	return s.repository.Cycle.UpdateManyTx(ctx, tx, where, data)
}

func (s *Cycle) Upsert(ctx context.Context, data ent.CreateCycleInput) (*ent.Cycle, error) {
	return s.repository.Cycle.Upsert(ctx, data)
}

func (s *Cycle) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateCycleInput) (*ent.Cycle, error) {
	return s.repository.Cycle.UpsertTx(ctx, tx, data)
}

func (s *Cycle) UpsertMany(ctx context.Context, data []ent.CreateCycleInput) (int, error) {
	return s.repository.Cycle.UpsertMany(ctx, data)
}

func (s *Cycle) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateCycleInput) (int, error) {
	return s.repository.Cycle.UpsertManyTx(ctx, tx, data)
}

func (s *Cycle) Delete(ctx context.Context, where ent.CycleWhereUniqueInput) (*ent.Cycle, error) {
	return s.repository.Cycle.Delete(ctx, where)
}

func (s *Cycle) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereUniqueInput) (*ent.Cycle, error) {
	return s.repository.Cycle.DeleteTx(ctx, tx, where)
}

func (s *Cycle) DeleteMany(ctx context.Context, where ent.CycleWhereInput) (int, error) {
	return s.repository.Cycle.DeleteMany(ctx, where)
}

func (s *Cycle) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereInput) (int, error) {
	return s.repository.Cycle.DeleteManyTx(ctx, tx, where)
}