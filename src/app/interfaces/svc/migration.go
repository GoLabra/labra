package svc

import (
	"app/ent"
	"context"
)

type Migration interface {
    Get(ctx context.Context, where *ent.MigrationWhereInput, orderBy *ent.MigrationOrder, skip *int, first *int, last *int) ([]*ent.Migration, error)
    Connection(ctx context.Context, where *ent.MigrationWhereInput, orderBy *ent.MigrationOrder, skip *int, first *int, last *int) (*ent.MigrationConnection, error)
    GetOne(ctx context.Context, where ent.MigrationWhereUniqueInput) (*ent.Migration, error)
    GetTx(ctx context.Context, tx *ent.Tx, where *ent.MigrationWhereInput, orderBy *ent.MigrationOrder, skip *int, first *int, last *int) ([]*ent.Migration, error)
    GetOneTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereUniqueInput) (*ent.Migration, error)
    Create(ctx context.Context, u ent.CreateMigrationInput) (*ent.Migration, error)
    CreateTx(ctx context.Context, tx *ent.Tx, u ent.CreateMigrationInput) (*ent.Migration, error)
    CreateMany(ctx context.Context, data []ent.CreateMigrationInput) ([]*ent.Migration, error)
    CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateMigrationInput) ([]*ent.Migration, error)
    Update(ctx context.Context, where ent.MigrationWhereUniqueInput, data ent.UpdateMigrationInput) (*ent.Migration, error)
    UpdateTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereUniqueInput, data ent.UpdateMigrationInput) (*ent.Migration, error)
    UpdateMany(ctx context.Context, where ent.MigrationWhereInput, data ent.UpdateMigrationInput) (int, error)
    UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereInput, data ent.UpdateMigrationInput) (int, error)
    Upsert(ctx context.Context, data ent.CreateMigrationInput) (*ent.Migration, error)
    UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateMigrationInput) (*ent.Migration, error)
    UpsertMany(ctx context.Context, data []ent.CreateMigrationInput) (int, error)
    UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateMigrationInput) (int, error)
    Delete(ctx context.Context, where ent.MigrationWhereUniqueInput) (*ent.Migration, error)
    DeleteTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereUniqueInput) (*ent.Migration, error)
    DeleteMany(ctx context.Context, where ent.MigrationWhereInput) (int, error)
    DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereInput) (int, error)
}