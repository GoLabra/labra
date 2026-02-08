package svc

import (
	"context"
	"app/ent"
	"app/domain/repo"
)

type Miau struct {
    repository *repo.Repository
}

func NewMiau(r *repo.Repository) *Miau {
	return &Miau{repository: r}
}


func (s *Miau) Get(ctx context.Context, where *ent.MiauWhereInput, orderBy *ent.MiauOrder, skip *int, first *int, last *int) ([]*ent.Miau, error) {
	return s.repository.Miau.Get(ctx, where, orderBy, skip, first, last)
}

func (s *Miau) Connection(ctx context.Context, where *ent.MiauWhereInput, orderBy *ent.MiauOrder, skip *int, first *int, last *int) (*ent.MiauConnection, error) {
	return s.repository.Miau.Connection(ctx, where, orderBy, skip, first, last)
}

func (s *Miau) GetOne(ctx context.Context, where ent.MiauWhereUniqueInput) (*ent.Miau, error) {
	return s.repository.Miau.GetOne(ctx, where)
}

func (s *Miau) GetTx(ctx context.Context, tx *ent.Tx, where *ent.MiauWhereInput, orderBy *ent.MiauOrder, skip *int, first *int, last *int) ([]*ent.Miau, error) {
	return s.repository.Miau.GetTx(ctx, tx, where, orderBy, skip, first, last)
}

func (s *Miau) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereUniqueInput) (*ent.Miau, error) {
	return s.repository.Miau.GetOneTx(ctx, tx, where)
}
func (s *Miau) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateMiauInput) (*ent.Miau, error) {
	createdInput, err := s.repository.Miau.CreateTx(ctx, tx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *Miau) Create(ctx context.Context, data ent.CreateMiauInput) (*ent.Miau, error) {
	createdInput, err := s.repository.Miau.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *Miau) CreateMany(ctx context.Context, data []ent.CreateMiauInput) ([]*ent.Miau, error) {
	return s.repository.Miau.CreateMany(ctx, data)
}

func (s *Miau) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateMiauInput) ([]*ent.Miau, error) {
	return s.repository.Miau.CreateManyTx(ctx, tx, data)
}

func (s *Miau) Update(ctx context.Context, where ent.MiauWhereUniqueInput, data ent.UpdateMiauInput) (*ent.Miau, error) {
	return s.repository.Miau.Update(ctx, where, data)
}

func (s *Miau) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereUniqueInput, data ent.UpdateMiauInput) (*ent.Miau, error) {
	return s.repository.Miau.UpdateTx(ctx, tx, where, data)
}

func (s *Miau) UpdateMany(ctx context.Context, where ent.MiauWhereInput, data ent.UpdateMiauInput) (int, error) {
	return s.repository.Miau.UpdateMany(ctx, where, data)
}

func (s *Miau) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereInput, data ent.UpdateMiauInput) (int, error) {
	return s.repository.Miau.UpdateManyTx(ctx, tx, where, data)
}

func (s *Miau) Upsert(ctx context.Context, data ent.CreateMiauInput) (*ent.Miau, error) {
	return s.repository.Miau.Upsert(ctx, data)
}

func (s *Miau) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateMiauInput) (*ent.Miau, error) {
	return s.repository.Miau.UpsertTx(ctx, tx, data)
}

func (s *Miau) UpsertMany(ctx context.Context, data []ent.CreateMiauInput) (int, error) {
	return s.repository.Miau.UpsertMany(ctx, data)
}

func (s *Miau) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateMiauInput) (int, error) {
	return s.repository.Miau.UpsertManyTx(ctx, tx, data)
}

func (s *Miau) Delete(ctx context.Context, where ent.MiauWhereUniqueInput) (*ent.Miau, error) {
	return s.repository.Miau.Delete(ctx, where)
}

func (s *Miau) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereUniqueInput) (*ent.Miau, error) {
	return s.repository.Miau.DeleteTx(ctx, tx, where)
}

func (s *Miau) DeleteMany(ctx context.Context, where ent.MiauWhereInput) (int, error) {
	return s.repository.Miau.DeleteMany(ctx, where)
}

func (s *Miau) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereInput) (int, error) {
	return s.repository.Miau.DeleteManyTx(ctx, tx, where)
}