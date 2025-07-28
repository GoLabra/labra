package utils

import (
	"fmt"
	"slices"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/GoLabra/labra/src/api/cache"
	"github.com/GoLabra/labra/src/api/config"
	"github.com/GoLabra/labra/src/api/entgql/annotations"
	"github.com/GoLabra/labra/src/api/entgql/entity"
	"github.com/GoLabra/labra/src/api/strcase"
	"github.com/mitchellh/mapstructure"
	"github.com/samborkent/uuidv7"
)

func NewUUIDV7() string {
	return uuidv7.New().String()
}

func LoadSchema(config *config.Config) {
	graph, err := entc.LoadGraph(config.EntSchemaPath, &gen.Config{})
	if err != nil {
		panic(err)
	}

	nodes := slices.DeleteFunc(graph.Nodes, func(node *gen.Type) bool {
		return node.Annotations["Entity"] == nil || node.Annotations["Entity"].(map[string]any)["Owner"] != "User"
	})

	graph, err = entc.LoadGraph("../api/entgql/ent/schema", &gen.Config{})
	if err != nil {
		panic(err)
	}

	nodes = append(nodes, graph.Nodes...)

	for _, node := range nodes {
		var entityAnnotations annotations.Entity
		err := mapstructure.Decode(node.Annotations["Entity"], &entityAnnotations)
		if err != nil {
			panic(err)
		}

		entityName := strcase.NodeNameToGraphqlName(node.Name)
		cache.Entity.Set(entityName, entity.Entity{
			Name:             entityName,
			EntName:          node.Name,
			Caption:          entityAnnotations.Caption,
			Owner:            entityAnnotations.Owner,
			DisplayFieldName: entityAnnotations.DisplayField,
		})

		fields := []entity.Field{
			{
				Caption: "Id",
				Name:    "id",
				Type:    string(entity.FieldTypeID),
			},
		}
		for _, nodeField := range node.Fields {
			var fieldAnnotations annotations.Field
			err := mapstructure.Decode(nodeField.Annotations["Field"], &fieldAnnotations)
			if err != nil {
				panic(err)
			}

			required := !nodeField.Optional
			unique := nodeField.Unique

			field := entity.Field{
				Name:           strcase.ToLowerCamel(nodeField.Name),
				EntName:        nodeField.Name,
				Caption:        fieldAnnotations.Caption,
				Type:           string(fieldAnnotations.Type),
				Required:       &required,
				Unique:         &unique,
				Nillable:       nodeField.Nillable,
				UpdateDefault:  nodeField.UpdateDefault,
				AcceptedValues: fieldAnnotations.AcceptedValues,
			}

			if nodeField.Default {
				defaultValue := fmt.Sprint(nodeField.DefaultValue())
				if fieldAnnotations.DefaultValue != "" {
					defaultValue = fieldAnnotations.DefaultValue
				}
				field.DefaultValue = &defaultValue
			}

			if fieldAnnotations.Min != "" {
				field.Min = &fieldAnnotations.Min
			}

			if fieldAnnotations.Max != "" {
				field.Max = &fieldAnnotations.Max
			}

			if fieldAnnotations.Private {
				field.Private = &fieldAnnotations.Private
			}

			fields = append(fields, field)
		}
		cache.Field.Set(entityName, fields)

		edges := []entity.Edge{}
		for _, edge := range node.Edges {
			var edgeAnnotations annotations.Edge
			err := mapstructure.Decode(edge.Annotations["Edge"], &edgeAnnotations)
			if err != nil {
				panic(err)
			}

			required := !edge.Optional
			ref := ""
			if edge.Ref != nil && edge.IsInverse() {
				ref = edge.Ref.Name
			}

			edges = append(edges, entity.Edge{
				Name:         strcase.ToLowerCamel(edge.Name),
				EntName:      edge.Name,
				Caption:      edgeAnnotations.Caption,
				Required:     &required,
				Type:         edge.Type.Name,
				Ref:          ref,
				RelationType: edgeAnnotations.RelationType,
			})
		}
		cache.Edge.Set(entityName, edges)
	}
}
