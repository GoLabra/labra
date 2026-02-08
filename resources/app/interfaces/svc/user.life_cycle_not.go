package svc

import (
	"app/ent"
	"context"

	"github.com/GoLabra/labra/entgql/entity"
)

type LifeCycleNot interface {
    Get(ctx context.Context, where *ent.LifeCycleNotWhereInput, orderBy *ent.LifeCycleNotOrder, skip *int, first *int, last *int) ([]*ent.LifeCycleNot, error)
    Connection(ctx context.Context, where *ent.LifeCycleNotWhereInput, orderBy *ent.LifeCycleNotOrder, skip *int, first *int, last *int) (*ent.LifeCycleNotConnection, error)
    GetOne(ctx context.Context, where ent.LifeCycleNotWhereUniqueInput) (*ent.LifeCycleNot, error)
    GetTx(ctx context.Context, tx *ent.Tx, where *ent.LifeCycleNotWhereInput, orderBy *ent.LifeCycleNotOrder, skip *int, first *int, last *int) ([]*ent.LifeCycleNot, error)
    GetOneTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereUniqueInput) (*ent.LifeCycleNot, error)
    Create(ctx context.Context, u ent.CreateLifeCycleNotInput) (*ent.LifeCycleNot, error)
    CreateTx(ctx context.Context, tx *ent.Tx, u ent.CreateLifeCycleNotInput) (*ent.LifeCycleNot, error)
    CreateMany(ctx context.Context, data []ent.CreateLifeCycleNotInput) ([]*ent.LifeCycleNot, error)
    CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateLifeCycleNotInput) ([]*ent.LifeCycleNot, error)
    Update(ctx context.Context, where ent.LifeCycleNotWhereUniqueInput, data ent.UpdateLifeCycleNotInput) (*ent.LifeCycleNot, error)
    UpdateState(ctx context.Context, where ent.LifeCycleNotWhereUniqueInput, state entity.EntityState) (*ent.LifeCycleNot, error)
    UpdateTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereUniqueInput, data ent.UpdateLifeCycleNotInput) (*ent.LifeCycleNot, error)
    UpdateMany(ctx context.Context, where ent.LifeCycleNotWhereInput, data ent.UpdateLifeCycleNotInput) (int, error)
    UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereInput, data ent.UpdateLifeCycleNotInput) (int, error)
    Upsert(ctx context.Context, data ent.CreateLifeCycleNotInput) (*ent.LifeCycleNot, error)
    UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateLifeCycleNotInput) (*ent.LifeCycleNot, error)
    UpsertMany(ctx context.Context, data []ent.CreateLifeCycleNotInput) (int, error)
    UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateLifeCycleNotInput) (int, error)
    Delete(ctx context.Context, where ent.LifeCycleNotWhereUniqueInput) (*ent.LifeCycleNot, error)
    DeleteTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereUniqueInput) (*ent.LifeCycleNot, error)
    DeleteMany(ctx context.Context, where ent.LifeCycleNotWhereInput) (int, error)
    DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereInput) (int, error)
}