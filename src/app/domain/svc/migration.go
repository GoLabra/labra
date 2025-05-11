package svc

import (
	"context"
	"app/ent"
    "app/domain/repo"
)

type Migration struct {
    repository *repo.Repository
}

func NewMigration(r *repo.Repository) *Migration {
	return &Migration{repository: r}
}


func (s *Migration) Get(ctx context.Context, where *ent.MigrationWhereInput, orderBy *ent.MigrationOrder, skip *int, first *int, last *int) ([]*ent.Migration, error) {
	return s.repository.Migration.Get(ctx, where, orderBy, skip, first, last)
}

func (s *Migration) Connection(ctx context.Context, where *ent.MigrationWhereInput, orderBy *ent.MigrationOrder, skip *int, first *int, last *int) (*ent.MigrationConnection, error) {
	return s.repository.Migration.Connection(ctx, where, orderBy, skip, first, last)
}

func (s *Migration) GetOne(ctx context.Context, where ent.MigrationWhereUniqueInput) (*ent.Migration, error) {
	return s.repository.Migration.GetOne(ctx, where)
}

func (s *Migration) GetTx(ctx context.Context, tx *ent.Tx, where *ent.MigrationWhereInput, orderBy *ent.MigrationOrder, skip *int, first *int, last *int) ([]*ent.Migration, error) {
	return s.repository.Migration.GetTx(ctx, tx, where, orderBy, skip, first, last)
}

func (s *Migration) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereUniqueInput) (*ent.Migration, error) {
	return s.repository.Migration.GetOneTx(ctx, tx, where)
}

func (s *Migration) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateMigrationInput) (*ent.Migration, error) {
	createdInput, err := s.repository.Migration.CreateTx(ctx, tx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *Migration) Create(ctx context.Context, data ent.CreateMigrationInput) (*ent.Migration, error) {
	createdInput, err := s.repository.Migration.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *Migration) CreateMany(ctx context.Context, data []ent.CreateMigrationInput) ([]*ent.Migration, error) {
	return s.repository.Migration.CreateMany(ctx, data)
}

func (s *Migration) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateMigrationInput) ([]*ent.Migration, error) {
	return s.repository.Migration.CreateManyTx(ctx, tx, data)
}

func (s *Migration) Update(ctx context.Context, where ent.MigrationWhereUniqueInput, data ent.UpdateMigrationInput) (*ent.Migration, error) {
	return s.repository.Migration.Update(ctx, where, data)
}

func (s *Migration) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereUniqueInput, data ent.UpdateMigrationInput) (*ent.Migration, error) {
	return s.repository.Migration.UpdateTx(ctx, tx, where, data)
}

func (s *Migration) UpdateMany(ctx context.Context, where ent.MigrationWhereInput, data ent.UpdateMigrationInput) (int, error) {
	return s.repository.Migration.UpdateMany(ctx, where, data)
}

func (s *Migration) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereInput, data ent.UpdateMigrationInput) (int, error) {
	return s.repository.Migration.UpdateManyTx(ctx, tx, where, data)
}

func (s *Migration) Upsert(ctx context.Context, data ent.CreateMigrationInput) (*ent.Migration, error) {
	return s.repository.Migration.Upsert(ctx, data)
}

func (s *Migration) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateMigrationInput) (*ent.Migration, error) {
	return s.repository.Migration.UpsertTx(ctx, tx, data)
}

func (s *Migration) UpsertMany(ctx context.Context, data []ent.CreateMigrationInput) (int, error) {
	return s.repository.Migration.UpsertMany(ctx, data)
}

func (s *Migration) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateMigrationInput) (int, error) {
	return s.repository.Migration.UpsertManyTx(ctx, tx, data)
}

func (s *Migration) Delete(ctx context.Context, where ent.MigrationWhereUniqueInput) (*ent.Migration, error) {
	return s.repository.Migration.Delete(ctx, where)
}

func (s *Migration) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereUniqueInput) (*ent.Migration, error) {
	return s.repository.Migration.DeleteTx(ctx, tx, where)
}

func (s *Migration) DeleteMany(ctx context.Context, where ent.MigrationWhereInput) (int, error) {
	return s.repository.Migration.DeleteMany(ctx, where)
}

func (s *Migration) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereInput) (int, error) {
	return s.repository.Migration.DeleteManyTx(ctx, tx, where)
}