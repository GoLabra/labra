package repo

import (
	"context"

	"github.com/GoLabra/labra/src/api/entgql/ent"
)

type File interface {
	Get(ctx context.Context, where *ent.FileWhereInput, orderBy *ent.FileOrder, skip *int, first *int, last *int) ([]*ent.File, error)
	Connection(ctx context.Context, where *ent.FileWhereInput, orderBy *ent.FileOrder, skip *int, first *int, last *int) (*ent.FileConnection, error)
	GetOne(ctx context.Context, where ent.FileWhereUniqueInput) (*ent.File, error)
	GetTx(ctx context.Context, tx *ent.Tx, where *ent.FileWhereInput, orderBy *ent.FileOrder, skip *int, first *int, last *int) ([]*ent.File, error)
	GetOneTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereUniqueInput) (*ent.File, error)
	Create(ctx context.Context, u ent.CreateFileInput) (*ent.File, error)
	CreateTx(ctx context.Context, tx *ent.Tx, u ent.CreateFileInput) (*ent.File, error)
	CreateMany(ctx context.Context, data []ent.CreateFileInput) ([]*ent.File, error)
	CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateFileInput) ([]*ent.File, error)
	Update(ctx context.Context, where ent.FileWhereUniqueInput, data ent.UpdateFileInput) (*ent.File, error)
	UpdateTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereUniqueInput, data ent.UpdateFileInput) (*ent.File, error)
	Upsert(ctx context.Context, data ent.CreateFileInput) (*ent.File, error)
	UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateFileInput) (*ent.File, error)
	UpsertMany(ctx context.Context, data []ent.CreateFileInput) (int, error)
	UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateFileInput) (int, error)
	UpdateMany(ctx context.Context, where ent.FileWhereInput, data ent.UpdateFileInput) (int, error)
	UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereInput, data ent.UpdateFileInput) (int, error)
	Delete(ctx context.Context, where ent.FileWhereUniqueInput) (*ent.File, error)
	DeleteTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereUniqueInput) (*ent.File, error)
	DeleteMany(ctx context.Context, where ent.FileWhereInput) (int, error)
	DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereInput) (int, error)
}
