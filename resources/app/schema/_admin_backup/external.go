package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/GoLabra/labra/entgql/annotations"
	"github.com/GoLabra/labra/entgql/entity"
)

type AdminUser struct {
	ent.Schema
}

func (AdminUser) IsExternal() bool {
	return true
}

func (AdminUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// entsql.Skip(),
		entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
	}
}

func (AdminUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Annotations(
				entgql.Type("ID"),
			),
	}
}

type File struct {
	ent.Schema
}

func (File) IsExternal() bool {
	return true
}

func (File) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// entsql.Skip(),
		// entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
	}
}

func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Annotations(
				entgql.Type("ID"),
			),
	}
}

func (File) Edges() []ent.Edge {
	return additionalFileEdges
}

type Role struct {
	ent.Schema
}

func (Role) IsExternal() bool {
	return true
}

func (Role) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
	}
}

func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Annotations(
				entgql.Type("ID"),
			),
		field.String("name").
			NotEmpty().
			Unique().
			Annotations(
				entgql.OrderField("name"),
				annotations.Field{
					Caption: "Name",
					Type:    entity.FieldTypeShortText,
				},
			),
	}
}

func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user_roles", User.Type).
			Ref("roles").
			Annotations(
				annotations.Edge{
					Caption:      "Role Users",
					RelationType: entity.RelationTypeM2M,
				},
			),
	}
}
