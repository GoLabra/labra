package repo

import (
	"app/ent"
	"context"
)

type Miau interface {
    Get(ctx context.Context, where *ent.MiauWhereInput, orderBy *ent.MiauOrder, skip *int, first *int, last *int) ([]*ent.Miau, error)
    Connection(ctx context.Context, where *ent.MiauWhereInput, orderBy *ent.MiauOrder, skip *int, first *int, last *int) (*ent.MiauConnection, error)
    GetOne(ctx context.Context, where ent.MiauWhereUniqueInput) (*ent.Miau, error)
    GetTx(ctx context.Context, tx *ent.Tx, where *ent.MiauWhereInput, orderBy *ent.MiauOrder, skip *int, first *int, last *int) ([]*ent.Miau, error)
    GetOneTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereUniqueInput) (*ent.Miau, error)
    Create(ctx context.Context, u ent.CreateMiauInput) (*ent.Miau, error)
    CreateTx(ctx context.Context, tx *ent.Tx, u ent.CreateMiauInput) (*ent.Miau, error)
    CreateMany(ctx context.Context, data []ent.CreateMiauInput) ([]*ent.Miau, error)
    CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateMiauInput) ([]*ent.Miau, error)
    Update(ctx context.Context, where ent.MiauWhereUniqueInput, data ent.UpdateMiauInput) (*ent.Miau, error)
    UpdateTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereUniqueInput, data ent.UpdateMiauInput) (*ent.Miau, error)
    Upsert(ctx context.Context, data ent.CreateMiauInput) (*ent.Miau, error)
    UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateMiauInput) (*ent.Miau, error)
    UpsertMany(ctx context.Context, data []ent.CreateMiauInput) (int, error)
    UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateMiauInput) (int, error)
    UpdateMany(ctx context.Context, where ent.MiauWhereInput, data ent.UpdateMiauInput) (int, error)
    UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereInput, data ent.UpdateMiauInput) (int, error)
    Delete(ctx context.Context, where ent.MiauWhereUniqueInput) (*ent.Miau, error)
    DeleteTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereUniqueInput) (*ent.Miau, error)
    DeleteMany(ctx context.Context, where ent.MiauWhereInput) (int, error)
    DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereInput) (int, error)
}