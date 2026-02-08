package repo

import (
	"app/ent"
	"context"

	"github.com/GoLabra/labra/entgql/entity"
)

type Cycle interface {
    Get(ctx context.Context, where *ent.CycleWhereInput, orderBy *ent.CycleOrder, skip *int, first *int, last *int) ([]*ent.Cycle, error)
    Connection(ctx context.Context, where *ent.CycleWhereInput, orderBy *ent.CycleOrder, skip *int, first *int, last *int) (*ent.CycleConnection, error)
    GetOne(ctx context.Context, where ent.CycleWhereUniqueInput) (*ent.Cycle, error)
    GetTx(ctx context.Context, tx *ent.Tx, where *ent.CycleWhereInput, orderBy *ent.CycleOrder, skip *int, first *int, last *int) ([]*ent.Cycle, error)
    GetOneTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereUniqueInput) (*ent.Cycle, error)
    Create(ctx context.Context, u ent.CreateCycleInput) (*ent.Cycle, error)
    CreateTx(ctx context.Context, tx *ent.Tx, u ent.CreateCycleInput) (*ent.Cycle, error)
    CreateMany(ctx context.Context, data []ent.CreateCycleInput) ([]*ent.Cycle, error)
    CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateCycleInput) ([]*ent.Cycle, error)
    Update(ctx context.Context, where ent.CycleWhereUniqueInput, data ent.UpdateCycleInput) (*ent.Cycle, error)
    UpdateState(ctx context.Context, where ent.CycleWhereUniqueInput, state entity.EntityState) (*ent.Cycle, error)
    UpdateStateTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereUniqueInput, state entity.EntityState) (*ent.Cycle, error)
    UpdateTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereUniqueInput, data ent.UpdateCycleInput) (*ent.Cycle, error)
    Upsert(ctx context.Context, data ent.CreateCycleInput) (*ent.Cycle, error)
    UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateCycleInput) (*ent.Cycle, error)
    UpsertMany(ctx context.Context, data []ent.CreateCycleInput) (int, error)
    UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateCycleInput) (int, error)
    UpdateMany(ctx context.Context, where ent.CycleWhereInput, data ent.UpdateCycleInput) (int, error)
    UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereInput, data ent.UpdateCycleInput) (int, error)
    Delete(ctx context.Context, where ent.CycleWhereUniqueInput) (*ent.Cycle, error)
    DeleteTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereUniqueInput) (*ent.Cycle, error)
    DeleteMany(ctx context.Context, where ent.CycleWhereInput) (int, error)
    DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereInput) (int, error)
}