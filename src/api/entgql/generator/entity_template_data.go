package generator

import (
	"errors"
	"fmt"

	"github.com/GoLabra/labrago/src/api/cache"
	"github.com/GoLabra/labrago/src/api/entgql/entity"
	"github.com/GoLabra/labrago/src/api/strcase"
)

type EntityTemplateData struct {
	Entity          entity.Entity
	Fields          []entity.Field
	Edges           []entity.Edge
	RelatedEntities map[string]*EntityTemplateData
}

func (e *EntityTemplateData) Validate() error {
	err := e.Entity.Validate()
	if err != nil {
		return err
	}

	for _, field := range e.Fields {
		err = field.Validate()
		if err != nil {
			return err
		}
	}

	for _, edge := range e.Edges {
		err = edge.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *EntityTemplateData) UpdateRelatedEntity(ownerEdge entity.Edge) error {
	var err error

	if ownerEdge.Ref != "" ||
		ownerEdge.RelationType == entity.RelationTypeOne ||
		ownerEdge.RelationType == entity.RelationTypeMany {
		return nil
	}

	var relatedEntity, ok = cache.Entity.GetByUniqueInput(entity.EntityWhereUniqueInput{EntName: &ownerEdge.Type})

	if !ok {
		return fmt.Errorf("update related entity: %w", cache.ErrEntityNotFound)
	}

	if relatedEntity.Owner == entity.EntityOwnerAdmin {
		return nil
	}

	relatedEntityTemplateData, ok := e.RelatedEntities[relatedEntity.EntName]

	if !ok {
		relatedEntityTemplateData, err = LoadEntity(entity.EntityWhereUniqueInput{Name: &relatedEntity.Name})

		if err != nil {
			return fmt.Errorf("update related entity: %w", err)
		}

		e.RelatedEntities[relatedEntity.EntName] = relatedEntityTemplateData
	}

	if ownerEdge.BelongsToCaption == nil {
		return fmt.Errorf("belongs to caption is required for edges with inverse relation")
	}

	relatedEntityTemplateData.Edges = append(relatedEntityTemplateData.Edges, entity.Edge{
		Name:         strcase.ToLowerCamel(*ownerEdge.BelongsToCaption),
		EntName:      strcase.ToSnake(*ownerEdge.BelongsToCaption),
		Caption:      *ownerEdge.BelongsToCaption,
		Required:     ownerEdge.Required,
		RelationType: InverseRelationType(ownerEdge.RelationType),
		Private:      ownerEdge.Private,
		Type:         e.Entity.EntName,
		Ref:          ownerEdge.EntName,
	})

	return nil
}

func (e *EntityTemplateData) AddEdgesCreate(createEdgesInput []*entity.CreateEdgeInput) error {
	for _, createEdgeInput := range createEdgesInput {
		var relatedEntity, ok = cache.Entity.GetByUniqueInput(createEdgeInput.RelatedEntity.Connect)

		if !ok {
			return fmt.Errorf("add edges create: %w", cache.ErrEntityNotFound)
		}

		edge := entity.Edge{
			Name:             strcase.ToLowerCamel(createEdgeInput.Caption),
			EntName:          strcase.ToSnake(createEdgeInput.Caption),
			Caption:          createEdgeInput.Caption,
			Required:         createEdgeInput.Required,
			RelationType:     createEdgeInput.RelationType,
			Private:          createEdgeInput.Private,
			BelongsToCaption: createEdgeInput.BelongsToCaption,
			Type:             relatedEntity.EntName,
		}

		e.Edges = append(e.Edges, edge)

		err := e.UpdateRelatedEntity(edge)

		if err != nil {
			return fmt.Errorf("add edges create: %w", err)
		}
	}

	return nil
}

func (e *EntityTemplateData) AddDefaultEdges() {
	e.Edges = append(e.Edges, entity.Edge{
		Name:         "CreatedBy",
		EntName:      "created_by",
		Caption:      "Created By",
		RelationType: entity.RelationTypeOne,
		Type:         "User",
		// Required: &trueVal, // TODO will be required
	})

	e.Edges = append(e.Edges, entity.Edge{
		Name:         "UpdatedBy",
		EntName:      "updated_by",
		Caption:      "Updated By",
		RelationType: entity.RelationTypeOne,
		Type:         "User",
		// Required: &trueVal, // TODO will be required
	})
}

func InverseRelationType(relationType entity.RelationType) entity.RelationType {
	switch relationType {
	case entity.RelationTypeM2O:
		return entity.RelationTypeO2M
	case entity.RelationTypeO2M:
		return entity.RelationTypeM2O
	default:
		return relationType
	}
}

func (e *EntityTemplateData) RemoveEdgesDelete(edgesDeleteInput []*entity.EdgeWhereUniqueInput) error {
	for _, edgeDeleteInput := range edgesDeleteInput {
		var edgeDeleted = false
		for index, edge := range e.Edges {
			if edge.Equals(*edgeDeleteInput) {
				edgeDeleted = true
				e.Edges = append(e.Edges[:index], e.Edges[index+1:]...)

				// TODO remove ref edge

				break
			}
		}
		if !edgeDeleted {
			return errors.New("remove edges delete: edge does not exist")
		}
	}
	return nil
}

func (e *EntityTemplateData) EditEdgesUpdate(edgesUpdateInput []*entity.UpdateOneEdgeInput) error {
	for _, edgeUpdateInput := range edgesUpdateInput {
		var edgeUpdated = false
		for i := range e.Edges {
			if e.Edges[i].Equals(edgeUpdateInput.Where) {
				(&e.Edges[i]).ApplyUpdateInput(edgeUpdateInput.Data)
				edgeUpdated = true
			}
		}
		if !edgeUpdated {
			return errors.New("edit edges update: edge does not exist")
		}
	}

	return nil
}

func (e *EntityTemplateData) AddDefaultFields() {
	var nowDefaultValue = "now()"

	e.Fields = append(e.Fields, entity.Field{
		Name:    "id",
		Caption: "Id",
		Type:    string(entity.FieldTypeID),
	})

	e.Fields = append(e.Fields, entity.Field{
		Name:          "createdAt",
		EntName:       "created_at",
		Caption:       "Created At",
		Type:          string(entity.FieldTypeDateTime),
		DefaultValue:  &nowDefaultValue,
		Nillable:      true,
		UpdateDefault: false,
	})

	e.Fields = append(e.Fields, entity.Field{
		Name:          "updatedAt",
		EntName:       "updated_at",
		Caption:       "Updated At",
		Type:          string(entity.FieldTypeDateTime),
		DefaultValue:  &nowDefaultValue,
		Nillable:      true,
		UpdateDefault: true,
	})
}

func (e *EntityTemplateData) AddFieldsCreate(fieldsCreateInput []*entity.CreateFieldInput) error {
	for _, fieldCreateInput := range fieldsCreateInput {
		e.Fields = append(e.Fields, entity.Field{
			Name:           strcase.ToLowerCamel(fieldCreateInput.Caption),
			EntName:        strcase.ToSnake(fieldCreateInput.Caption),
			Caption:        fieldCreateInput.Caption,
			Type:           fieldCreateInput.Type,
			Required:       fieldCreateInput.Required,
			Unique:         fieldCreateInput.Unique,
			DefaultValue:   fieldCreateInput.DefaultValue,
			Min:            fieldCreateInput.Min,
			Max:            fieldCreateInput.Max,
			Private:        fieldCreateInput.Private,
			AcceptedValues: fieldCreateInput.AcceptedValues,
		})
	}

	return nil
}

func (e *EntityTemplateData) RemoveFieldsDelete(fieldsDeleteInput []*entity.FieldWhereUniqueInput) error {
	for _, fieldDeleteInput := range fieldsDeleteInput {
		var fieldDeleted = false
		for index, field := range e.Fields {
			if field.Equals(*fieldDeleteInput) {
				fieldDeleted = true
				e.Fields = append(e.Fields[:index], e.Fields[index+1:]...)
				break
			}
		}
		if !fieldDeleted {
			return errors.New("remove fields delete: field does not exist")
		}
	}
	return nil
}

func (e *EntityTemplateData) EditFieldsUpdate(fieldsUpdateInput []*entity.UpdateOneFieldInput) error {
	for _, fieldUpdateInput := range fieldsUpdateInput {
		var fieldUpdated = false
		for i := range e.Fields {
			if e.Fields[i].Equals(fieldUpdateInput.Where) {
				(&e.Fields[i]).ApplyUpdateInput(fieldUpdateInput.Data)
				fieldUpdated = true
			}
		}
		if !fieldUpdated {
			return errors.New("edit fields update: field does not exist")
		}
	}

	return nil
}

func LoadEntity(where entity.EntityWhereUniqueInput) (*EntityTemplateData, error) {
	entity, ok := cache.Entity.GetByUniqueInput(where)

	if !ok {
		return nil, fmt.Errorf("load entity: get entity by unique input error: %w", cache.ErrEntityNotFound)
	}

	fields, ok := cache.Field.Get(entity.Name)

	if !ok {
		return nil, fmt.Errorf("load entity: get fields by entity name error: %w", cache.ErrFieldsNotFound)
	}

	edges, ok := cache.Edge.Get(entity.Name)

	if !ok {
		return nil, fmt.Errorf("load entity: get edges by entity name error: %w", cache.ErrEdgesNotFound)
	}

	return &EntityTemplateData{
		Entity:          entity,
		Fields:          fields,
		Edges:           edges,
		RelatedEntities: make(map[string]*EntityTemplateData),
	}, nil
}

func ApplyEntityUpdate(where entity.EntityWhereUniqueInput, data entity.UpdateEntityInput) (*EntityTemplateData, error) {
	var entityTemplateData, err = LoadEntity(where)

	if err != nil {
		return nil, fmt.Errorf("apply entity update: %w", err)
	}

	(&entityTemplateData.Entity).ApplyUpdateInput(data)

	if data.Fields != nil {
		if data.Fields.Create != nil {
			err = entityTemplateData.AddFieldsCreate(data.Fields.Create)
			if err != nil {
				return nil, fmt.Errorf("apply entity update: %w", err)
			}
		}

		if data.Fields.Delete != nil {
			err = entityTemplateData.RemoveFieldsDelete(data.Fields.Delete)
			if err != nil {
				return nil, fmt.Errorf("apply entity update: %w", err)
			}
		}

		if data.Fields.Update != nil {
			err = entityTemplateData.EditFieldsUpdate(data.Fields.Update)
			if err != nil {
				return nil, fmt.Errorf("apply entity update: %w", err)
			}
		}
	}

	if data.Edges != nil {
		if data.Edges.Create != nil {
			err = entityTemplateData.AddEdgesCreate(data.Edges.Create)
			if err != nil {
				return nil, fmt.Errorf("apply entity update: %w", err)
			}
		}

		if data.Edges.Delete != nil {
			err = entityTemplateData.RemoveEdgesDelete(data.Edges.Delete)
			if err != nil {
				return nil, fmt.Errorf("apply entity update: %w", err)
			}
		}

		if data.Edges.Update != nil {
			err = entityTemplateData.EditEdgesUpdate(data.Edges.Update)
			if err != nil {
				return nil, fmt.Errorf("apply entity update: %w", err)
			}
		}
	}

	if data.DisplayField != nil {
		var displayFieldSet = false
		for _, field := range entityTemplateData.Fields {
			if field.Equals(*data.DisplayField) {
				entityTemplateData.Entity.DisplayFieldName = field.Name
				displayFieldSet = true
			}
		}
		if !displayFieldSet {
			return nil, fmt.Errorf("cannot find display field")
		}
	}

	return entityTemplateData, nil
}

func ApplyEntityCreate(data entity.CreateEntityInput) (*EntityTemplateData, error) {
	var err error
	var entityTemplateData = &EntityTemplateData{
		RelatedEntities: make(map[string]*EntityTemplateData),
	}

	entityTemplateData.Entity = entity.Entity{
		Name:    strcase.ToLowerCamel(data.Caption),
		EntName: strcase.ToCamel(data.Caption),
		Caption: data.Caption,
		Owner:   "User",
	}

	entityTemplateData.AddDefaultFields()

	if data.Fields != nil && data.Fields.Create != nil {
		err = entityTemplateData.AddFieldsCreate(data.Fields.Create)
		if err != nil {
			return entityTemplateData, err
		}
	}

	entityTemplateData.AddDefaultEdges()

	if data.Edges != nil && data.Edges.Create != nil {
		err = entityTemplateData.AddEdgesCreate(data.Edges.Create)
		if err != nil {
			return entityTemplateData, err
		}
	}

	var displayFieldSet = false

	for _, field := range entityTemplateData.Fields {
		if field.Equals(data.DisplayField) {
			entityTemplateData.Entity.DisplayFieldName = field.Name
			displayFieldSet = true
		}
	}

	if !displayFieldSet {
		return nil, fmt.Errorf("cannot find display field")
	}

	return entityTemplateData, nil
}

func ApplyEntityDelete(data entity.EntityWhereUniqueInput) (*EntityTemplateData, error) {
	var entityTemplateData = &EntityTemplateData{
		RelatedEntities: make(map[string]*EntityTemplateData),
	}

	entityToDelete, ok := cache.Entity.GetByUniqueInput(data)

	if !ok {
		return nil, fmt.Errorf("apply entity delete: %w", cache.ErrEntityNotFound)
	}

	entityTemplateData.Entity = entityToDelete

	for _, relatedEntity := range cache.Entity.GetAll() {
		edges, _ := cache.Edge.Get(relatedEntity.Name)
		for _, edge := range edges {
			if edge.Type == entityToDelete.EntName {
				updatedRelatedEntity, err := ApplyEntityUpdate(
					entity.EntityWhereUniqueInput{Name: &relatedEntity.Name},
					entity.UpdateEntityInput{
						Edges: &entity.UpdateManyEdgesInput{
							Delete: []*entity.EdgeWhereUniqueInput{
								{
									Name: &edge.Name,
								},
							},
						},
					})
				if err != nil {
					return nil, err
				}
				entityTemplateData.RelatedEntities[relatedEntity.EntName] = updatedRelatedEntity
			}
		}
	}

	return entityTemplateData, nil
}
