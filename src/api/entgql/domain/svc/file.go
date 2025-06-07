package svc

import (
	"context"

	"github.com/GoLabra/labra/src/api/entgql/domain/repo"
	"github.com/GoLabra/labra/src/api/entgql/ent"
)

type File struct {
	repository *repo.Repository
}

func NewFile(r *repo.Repository) *File {
	return &File{repository: r}
}

func (s *File) Get(ctx context.Context, where *ent.FileWhereInput, orderBy *ent.FileOrder, skip *int, first *int, last *int) ([]*ent.File, error) {
	return s.repository.File.Get(ctx, where, orderBy, skip, first, last)
}

func (s *File) Connection(ctx context.Context, where *ent.FileWhereInput, orderBy *ent.FileOrder, skip *int, first *int, last *int) (*ent.FileConnection, error) {
	return s.repository.File.Connection(ctx, where, orderBy, skip, first, last)
}

func (s *File) GetOne(ctx context.Context, where ent.FileWhereUniqueInput) (*ent.File, error) {
	return s.repository.File.GetOne(ctx, where)
}

func (s *File) GetTx(ctx context.Context, tx *ent.Tx, where *ent.FileWhereInput, orderBy *ent.FileOrder, skip *int, first *int, last *int) ([]*ent.File, error) {
	return s.repository.File.GetTx(ctx, tx, where, orderBy, skip, first, last)
}

func (s *File) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereUniqueInput) (*ent.File, error) {
	return s.repository.File.GetOneTx(ctx, tx, where)
}
func (s *File) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateFileInput) (*ent.File, error) {
	createdInput, err := s.repository.File.CreateTx(ctx, tx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *File) Create(ctx context.Context, data ent.CreateFileInput) (*ent.File, error) {
	createdInput, err := s.repository.File.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *File) CreateMany(ctx context.Context, data []ent.CreateFileInput) ([]*ent.File, error) {
	return s.repository.File.CreateMany(ctx, data)
}

func (s *File) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateFileInput) ([]*ent.File, error) {
	return s.repository.File.CreateManyTx(ctx, tx, data)
}

func (s *File) Update(ctx context.Context, where ent.FileWhereUniqueInput, data ent.UpdateFileInput) (*ent.File, error) {
	return s.repository.File.Update(ctx, where, data)
}

func (s *File) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereUniqueInput, data ent.UpdateFileInput) (*ent.File, error) {
	return s.repository.File.UpdateTx(ctx, tx, where, data)
}

func (s *File) UpdateMany(ctx context.Context, where ent.FileWhereInput, data ent.UpdateFileInput) (int, error) {
	return s.repository.File.UpdateMany(ctx, where, data)
}

func (s *File) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereInput, data ent.UpdateFileInput) (int, error) {
	return s.repository.File.UpdateManyTx(ctx, tx, where, data)
}

func (s *File) Upsert(ctx context.Context, data ent.CreateFileInput) (*ent.File, error) {
	return s.repository.File.Upsert(ctx, data)
}

func (s *File) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateFileInput) (*ent.File, error) {
	return s.repository.File.UpsertTx(ctx, tx, data)
}

func (s *File) UpsertMany(ctx context.Context, data []ent.CreateFileInput) (int, error) {
	return s.repository.File.UpsertMany(ctx, data)
}

func (s *File) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateFileInput) (int, error) {
	return s.repository.File.UpsertManyTx(ctx, tx, data)
}

func (s *File) Delete(ctx context.Context, where ent.FileWhereUniqueInput) (*ent.File, error) {
	return s.repository.File.Delete(ctx, where)
}

func (s *File) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereUniqueInput) (*ent.File, error) {
	return s.repository.File.DeleteTx(ctx, tx, where)
}

func (s *File) DeleteMany(ctx context.Context, where ent.FileWhereInput) (int, error) {
	return s.repository.File.DeleteMany(ctx, where)
}

func (s *File) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereInput) (int, error) {
	return s.repository.File.DeleteManyTx(ctx, tx, where)
}
