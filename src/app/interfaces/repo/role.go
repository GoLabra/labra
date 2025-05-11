package repo

import (
	"app/ent"
	"context"
)

type Role interface {
    Get(ctx context.Context, where *ent.RoleWhereInput, orderBy *ent.RoleOrder, skip *int, first *int, last *int) ([]*ent.Role, error)
    Connection(ctx context.Context, where *ent.RoleWhereInput, orderBy *ent.RoleOrder, skip *int, first *int, last *int) (*ent.RoleConnection, error)
    GetOne(ctx context.Context, where ent.RoleWhereUniqueInput) (*ent.Role, error)
    GetTx(ctx context.Context, tx *ent.Tx, where *ent.RoleWhereInput, orderBy *ent.RoleOrder, skip *int, first *int, last *int) ([]*ent.Role, error)
    GetOneTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereUniqueInput) (*ent.Role, error)
    Create(ctx context.Context, u ent.CreateRoleInput) (*ent.Role, error)
    CreateTx(ctx context.Context, tx *ent.Tx, u ent.CreateRoleInput) (*ent.Role, error)
    CreateMany(ctx context.Context, data []ent.CreateRoleInput) ([]*ent.Role, error)
    CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateRoleInput) ([]*ent.Role, error)
    Update(ctx context.Context, where ent.RoleWhereUniqueInput, data ent.UpdateRoleInput) (*ent.Role, error)
    UpdateTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereUniqueInput, data ent.UpdateRoleInput) (*ent.Role, error)
    Upsert(ctx context.Context, data ent.CreateRoleInput) (*ent.Role, error)
    UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateRoleInput) (*ent.Role, error)
    UpsertMany(ctx context.Context, data []ent.CreateRoleInput) (int, error)
    UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateRoleInput) (int, error)
    UpdateMany(ctx context.Context, where ent.RoleWhereInput, data ent.UpdateRoleInput) (int, error)
    UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereInput, data ent.UpdateRoleInput) (int, error)
    Delete(ctx context.Context, where ent.RoleWhereUniqueInput) (*ent.Role, error)
    DeleteTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereUniqueInput) (*ent.Role, error)
    DeleteMany(ctx context.Context, where ent.RoleWhereInput) (int, error)
    DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereInput) (int, error)
}