package repo

import (
	"app/ent"
	"context"
)

type ForPermission interface {
    Get(ctx context.Context, where *ent.ForPermissionWhereInput, orderBy *ent.ForPermissionOrder, skip *int, first *int, last *int) ([]*ent.ForPermission, error)
    Connection(ctx context.Context, where *ent.ForPermissionWhereInput, orderBy *ent.ForPermissionOrder, skip *int, first *int, last *int) (*ent.ForPermissionConnection, error)
    GetOne(ctx context.Context, where ent.ForPermissionWhereUniqueInput) (*ent.ForPermission, error)
    GetTx(ctx context.Context, tx *ent.Tx, where *ent.ForPermissionWhereInput, orderBy *ent.ForPermissionOrder, skip *int, first *int, last *int) ([]*ent.ForPermission, error)
    GetOneTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereUniqueInput) (*ent.ForPermission, error)
    Create(ctx context.Context, u ent.CreateForPermissionInput) (*ent.ForPermission, error)
    CreateTx(ctx context.Context, tx *ent.Tx, u ent.CreateForPermissionInput) (*ent.ForPermission, error)
    CreateMany(ctx context.Context, data []ent.CreateForPermissionInput) ([]*ent.ForPermission, error)
    CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateForPermissionInput) ([]*ent.ForPermission, error)
    Update(ctx context.Context, where ent.ForPermissionWhereUniqueInput, data ent.UpdateForPermissionInput) (*ent.ForPermission, error)
    UpdateTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereUniqueInput, data ent.UpdateForPermissionInput) (*ent.ForPermission, error)
    Upsert(ctx context.Context, data ent.CreateForPermissionInput) (*ent.ForPermission, error)
    UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateForPermissionInput) (*ent.ForPermission, error)
    UpsertMany(ctx context.Context, data []ent.CreateForPermissionInput) (int, error)
    UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateForPermissionInput) (int, error)
    UpdateMany(ctx context.Context, where ent.ForPermissionWhereInput, data ent.UpdateForPermissionInput) (int, error)
    UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereInput, data ent.UpdateForPermissionInput) (int, error)
    Delete(ctx context.Context, where ent.ForPermissionWhereUniqueInput) (*ent.ForPermission, error)
    DeleteTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereUniqueInput) (*ent.ForPermission, error)
    DeleteMany(ctx context.Context, where ent.ForPermissionWhereInput) (int, error)
    DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereInput) (int, error)
}