package svc

import (
	"app/ent"
	"context"
)

type User interface {
    Get(ctx context.Context, where *ent.UserWhereInput, orderBy *ent.UserOrder, skip *int, first *int, last *int) ([]*ent.User, error)
    Connection(ctx context.Context, where *ent.UserWhereInput, orderBy *ent.UserOrder, skip *int, first *int, last *int) (*ent.UserConnection, error)
    GetOne(ctx context.Context, where ent.UserWhereUniqueInput) (*ent.User, error)
    GetTx(ctx context.Context, tx *ent.Tx, where *ent.UserWhereInput, orderBy *ent.UserOrder, skip *int, first *int, last *int) ([]*ent.User, error)
    GetOneTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereUniqueInput) (*ent.User, error)
    Create(ctx context.Context, u ent.CreateUserInput) (*ent.User, error)
    CreateTx(ctx context.Context, tx *ent.Tx, u ent.CreateUserInput) (*ent.User, error)
    CreateMany(ctx context.Context, data []ent.CreateUserInput) ([]*ent.User, error)
    CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateUserInput) ([]*ent.User, error)
    Update(ctx context.Context, where ent.UserWhereUniqueInput, data ent.UpdateUserInput) (*ent.User, error)
    UpdateTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereUniqueInput, data ent.UpdateUserInput) (*ent.User, error)
    UpdateMany(ctx context.Context, where ent.UserWhereInput, data ent.UpdateUserInput) (int, error)
    UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereInput, data ent.UpdateUserInput) (int, error)
    Upsert(ctx context.Context, data ent.CreateUserInput) (*ent.User, error)
    UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateUserInput) (*ent.User, error)
    UpsertMany(ctx context.Context, data []ent.CreateUserInput) (int, error)
    UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateUserInput) (int, error)
    Delete(ctx context.Context, where ent.UserWhereUniqueInput) (*ent.User, error)
    DeleteTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereUniqueInput) (*ent.User, error)
    DeleteMany(ctx context.Context, where ent.UserWhereInput) (int, error)
    DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereInput) (int, error)
}