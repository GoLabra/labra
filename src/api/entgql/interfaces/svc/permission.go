package svc

import (
	"context"

	"github.com/GoLabra/labra/src/api/entgql/ent"
)

type Permission interface {
	Get(ctx context.Context, where *ent.PermissionWhereInput, orderBy *ent.PermissionOrder, skip *int, first *int, last *int) ([]*ent.Permission, error)
	Connection(ctx context.Context, where *ent.PermissionWhereInput, orderBy *ent.PermissionOrder, skip *int, first *int, last *int) (*ent.PermissionConnection, error)
	GetOne(ctx context.Context, where ent.PermissionWhereUniqueInput) (*ent.Permission, error)
	GetTx(ctx context.Context, tx *ent.Tx, where *ent.PermissionWhereInput, orderBy *ent.PermissionOrder, skip *int, first *int, last *int) ([]*ent.Permission, error)
	GetOneTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereUniqueInput) (*ent.Permission, error)
	Create(ctx context.Context, u ent.CreatePermissionInput) (*ent.Permission, error)
	CreateTx(ctx context.Context, tx *ent.Tx, u ent.CreatePermissionInput) (*ent.Permission, error)
	CreateMany(ctx context.Context, data []ent.CreatePermissionInput) ([]*ent.Permission, error)
	CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreatePermissionInput) ([]*ent.Permission, error)
	Update(ctx context.Context, where ent.PermissionWhereUniqueInput, data ent.UpdatePermissionInput) (*ent.Permission, error)
	UpdateTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereUniqueInput, data ent.UpdatePermissionInput) (*ent.Permission, error)
	UpdateMany(ctx context.Context, where ent.PermissionWhereInput, data ent.UpdatePermissionInput) (int, error)
	UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereInput, data ent.UpdatePermissionInput) (int, error)
	Upsert(ctx context.Context, data ent.CreatePermissionInput) (*ent.Permission, error)
	UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreatePermissionInput) (*ent.Permission, error)
	UpsertMany(ctx context.Context, data []ent.CreatePermissionInput) (int, error)
	UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreatePermissionInput) (int, error)
	Delete(ctx context.Context, where ent.PermissionWhereUniqueInput) (*ent.Permission, error)
	DeleteTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereUniqueInput) (*ent.Permission, error)
	DeleteMany(ctx context.Context, where ent.PermissionWhereInput) (int, error)
	DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereInput) (int, error)
}
