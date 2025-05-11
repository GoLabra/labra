package svc

import (
	"context"
	"sort"

	"github.com/GoLabra/labrago/src/api/cache"
	"github.com/GoLabra/labrago/src/api/entgql/entity"
	"github.com/GoLabra/labrago/src/api/entgql/generator"
	"github.com/GoLabra/labrago/src/api/strcase"
)

type SchemaManager interface {
	WriteEntityToSchema(entityTemplateData generator.EntityTemplateData) error
	RemoveEntityFromSchema(fileName string) error
	BackupSchema() error
	Generate(ctx context.Context, entities []entity.Entity) error
	RevertSchema() error
}

type Entity struct {
	schemaManager SchemaManager
}

func NewEntity(schemaManager SchemaManager) *Entity {
	return &Entity{
		schemaManager: schemaManager,
	}
}

func (e Entity) UpdateEntity(ctx context.Context, where entity.EntityWhereUniqueInput, data entity.UpdateEntityInput) (*entity.Entity, error) {
	var entitiesToGenerate []entity.Entity
	entityToUpdate, err := generator.ApplyEntityUpdate(where, data)

	if err != nil {
		return nil, err
	}

	err = e.schemaManager.BackupSchema()

	if err != nil {
		return nil, err
	}

	err = e.schemaManager.WriteEntityToSchema(*entityToUpdate)

	if err != nil {
		e.schemaManager.RevertSchema()
		return nil, err
	}

	entitiesToGenerate = append(entitiesToGenerate, entityToUpdate.Entity)

	for _, relatedEntity := range entityToUpdate.RelatedEntities {
		entitiesToGenerate = append(entitiesToGenerate, relatedEntity.Entity)

		err = e.schemaManager.WriteEntityToSchema(*relatedEntity)
		if err != nil {
			e.schemaManager.RevertSchema()
			return nil, err
		}
	}

	err = e.schemaManager.Generate(ctx, entitiesToGenerate)
	if err != nil {
		return nil, err
	}

	return &entityToUpdate.Entity, nil
}

func (e Entity) CreateEntity(ctx context.Context, data entity.CreateEntityInput) (*entity.Entity, error) {
	entityToUpdate, err := generator.ApplyEntityCreate(data)
	if err != nil {
		return nil, err
	}

	entitiesToGenerate := []entity.Entity{entityToUpdate.Entity}

	err = e.schemaManager.BackupSchema()

	if err != nil {
		return nil, err
	}

	err = e.schemaManager.WriteEntityToSchema(*entityToUpdate)

	if err != nil {
		e.schemaManager.RevertSchema()
		return nil, err
	}

	for _, relatedEntity := range entityToUpdate.RelatedEntities {
		entitiesToGenerate = append(entitiesToGenerate, relatedEntity.Entity)

		err = e.schemaManager.WriteEntityToSchema(*relatedEntity)
		if err != nil {
			e.schemaManager.RevertSchema()
			return nil, err
		}
	}

	err = e.schemaManager.Generate(ctx, entitiesToGenerate)
	if err != nil {
		return nil, err
	}

	return &entityToUpdate.Entity, nil
}

func (e Entity) DeleteEntity(ctx context.Context, where entity.EntityWhereUniqueInput) (*entity.Entity, error) {
	entityToDelete, err := generator.ApplyEntityDelete(where)
	if err != nil {
		return nil, err
	}

	err = e.schemaManager.BackupSchema()

	if err != nil {
		return nil, err
	}

	entitiesToGenerate := []entity.Entity{entityToDelete.Entity}

	for _, relatedEntity := range entityToDelete.RelatedEntities {
		entitiesToGenerate = append(entitiesToGenerate, relatedEntity.Entity)

		err = e.schemaManager.WriteEntityToSchema(*relatedEntity)
		if err != nil {
			e.schemaManager.RevertSchema()
			return nil, err
		}
	}

	err = e.schemaManager.RemoveEntityFromSchema(strcase.ToSnake(entityToDelete.Entity.EntName))

	if err != nil {
		e.schemaManager.RevertSchema()
		return nil, err
	}

	err = e.schemaManager.Generate(ctx, entitiesToGenerate)
	if err != nil {
		return nil, err
	}

	return &entityToDelete.Entity, nil
}

func (Entity) Entities(ctx context.Context) ([]*entity.Entity, error) {
	entities := cache.Entity.GetAll()

	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Caption < entities[j].Caption
	})

	var result = make([]*entity.Entity, len(entities))

	for index := range entities {
		result[index] = &entities[index]
	}

	return result, nil
}

func (Entity) Entity(ctx context.Context, where *entity.EntityWhereUniqueInput) (*entity.Entity, error) {
	entity, ok := cache.Entity.Get(*where.Name)
	if !ok {
		return nil, nil
	}
	return &entity, nil
}

func (Entity) Fields(ctx context.Context) ([]*entity.Field, error) {
	fieldBatches := cache.Field.GetAll()
	fields := []entity.Field{}
	for _, fieldBatch := range fieldBatches {
		fields = append(fields, fieldBatch...)
	}

	var result = make([]*entity.Field, len(fields))

	for index := range fields {
		result[index] = &fields[index]
	}

	return result, nil
}

func (Entity) EntityRelatedEntity(ctx context.Context, obj *entity.Edge) (*entity.Entity, error) {
	relatedEntity, ok := cache.Entity.Get(strcase.ToLowerCamel(obj.Type))
	if !ok {
		return nil, nil
	}
	return &relatedEntity, nil
}

func (Entity) EntityFields(ctx context.Context, obj *entity.Entity) ([]*entity.Field, error) {
	fields, ok := cache.Field.Get(obj.Name)
	if !ok {
		return nil, nil
	}

	var result = make([]*entity.Field, len(fields))

	for index := range fields {
		result[index] = &fields[index]
	}

	return result, nil
}

func (Entity) EntityEdges(ctx context.Context, obj *entity.Entity) ([]*entity.Edge, error) {
	edges, ok := cache.Edge.Get(obj.Name)
	if !ok {
		return nil, nil
	}

	var result = make([]*entity.Edge, len(edges))

	for index := range edges {
		result[index] = &edges[index]
	}

	return result, nil
}
