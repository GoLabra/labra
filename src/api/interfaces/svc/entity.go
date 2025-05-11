package svc

import (
	"context"
	"github.com/GoLabra/labrago/src/api/entgql/entity"
)

type Entity interface {
	UpdateEntity(ctx context.Context, where entity.EntityWhereUniqueInput, data entity.UpdateEntityInput) (*entity.Entity, error)
	CreateEntity(ctx context.Context, data entity.CreateEntityInput) (*entity.Entity, error)
	DeleteEntity(ctx context.Context, where entity.EntityWhereUniqueInput) (*entity.Entity, error)
	Entities(ctx context.Context) ([]*entity.Entity, error)
	Entity(ctx context.Context, where *entity.EntityWhereUniqueInput) (*entity.Entity, error)
	Fields(ctx context.Context) ([]*entity.Field, error)
	EntityRelatedEntity(ctx context.Context, obj *entity.Edge) (*entity.Entity, error)
	EntityFields(ctx context.Context, obj *entity.Entity) ([]*entity.Field, error)
	EntityEdges(ctx context.Context, obj *entity.Entity) ([]*entity.Edge, error)
}
