package svc

import (
	"context"

	"github.com/GoLabra/labra/entgql/ent"
)

type AdminUser interface {
	Get(ctx context.Context, where *ent.AdminUserWhereInput, orderBy *ent.AdminUserOrder, skip *int, first *int, last *int) ([]*ent.AdminUser, error)
	Connection(ctx context.Context, where *ent.AdminUserWhereInput, orderBy *ent.AdminUserOrder, skip *int, first *int, last *int) (*ent.AdminUserConnection, error)
	GetOne(ctx context.Context, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error)
	GetTx(ctx context.Context, tx *ent.Tx, where *ent.AdminUserWhereInput, orderBy *ent.AdminUserOrder, skip *int, first *int, last *int) ([]*ent.AdminUser, error)
	GetOneTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error)
	Create(ctx context.Context, u ent.CreateAdminUserInput) (*ent.AdminUser, error)
	CreateTx(ctx context.Context, tx *ent.Tx, u ent.CreateAdminUserInput) (*ent.AdminUser, error)
	CreateMany(ctx context.Context, data []ent.CreateAdminUserInput) ([]*ent.AdminUser, error)
	CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateAdminUserInput) ([]*ent.AdminUser, error)
	Update(ctx context.Context, where ent.AdminUserWhereUniqueInput, data ent.UpdateAdminUserInput) (*ent.AdminUser, error)
	UpdateTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereUniqueInput, data ent.UpdateAdminUserInput) (*ent.AdminUser, error)
	UpdateMany(ctx context.Context, where ent.AdminUserWhereInput, data ent.UpdateAdminUserInput) (int, error)
	UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereInput, data ent.UpdateAdminUserInput) (int, error)
	Upsert(ctx context.Context, data ent.CreateAdminUserInput) (*ent.AdminUser, error)
	UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateAdminUserInput) (*ent.AdminUser, error)
	UpsertMany(ctx context.Context, data []ent.CreateAdminUserInput) (int, error)
	UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateAdminUserInput) (int, error)
	UpdatePassword(ctx context.Context, where ent.AdminUserWhereUniqueInput, password string) (*ent.AdminUser, error)
	Delete(ctx context.Context, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error)
	DeleteTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error)
	DeleteMany(ctx context.Context, where ent.AdminUserWhereInput) (int, error)
	DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereInput) (int, error)
}
