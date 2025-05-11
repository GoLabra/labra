package schema

import (
	"entgo.io/ent/schema/edge"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/GoLabra/labrago/src/api/entgql/annotations"
	"github.com/GoLabra/labrago/src/api/entgql/entity"
	"github.com/lucsky/cuid"
)

type Migration struct {
	ent.Schema
}

func (Migration) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.MultiOrder(),
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		annotations.Entity{
			Caption:      "Migration",
			Owner:        entity.EntityOwnerAdmin,
			DisplayField: "name",
		},
	}
}

func (Migration) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(cuid.New).
			Annotations(
				entgql.OrderField("id"),
				annotations.Field{
					Caption: "ID",
					Type:    entity.FieldTypeID,
				},
			),

		field.String("name").
			NotEmpty().
			Annotations(
				entgql.OrderField("name"),
				annotations.Field{
					Caption: "Migration Name",
					Type:    entity.FieldTypeShortText,
				},
			),

		field.Enum("type").
			Values("sql", "graphql", "go").
			Annotations(
				entgql.OrderField("type"),
				annotations.Field{
					Caption: "Type",
					Type:    entity.FieldTypeEnum,
				},
			),

		field.Enum("direction").
			Values("up", "down").
			Annotations(
				entgql.OrderField("direction"),
				annotations.Field{
					Caption: "Direction",
					Type:    entity.FieldTypeEnum,
				},
			),

		field.String("plugin").
			Optional().
			Annotations(
				entgql.OrderField("plugin"),
				annotations.Field{
					Caption: "Plugin",
					Type:    entity.FieldTypeShortText,
				},
			),

		field.Time("created_at").
			Optional().
			Nillable().
			Immutable().
			Default(time.Now).
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				entgql.OrderField("createdAt"),
				annotations.Field{
					Caption:      "Created At",
					Type:         entity.FieldTypeDateTime,
					DefaultValue: "now()",
				},
			),

		field.Time("updated_at").
			Optional().
			Nillable().
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				entgql.OrderField("updatedAt"),
				annotations.Field{
					Caption:      "Updated At",
					Type:         entity.FieldTypeDateTime,
					DefaultValue: "now()",
				},
			),
	}
}
func (Migration) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("created_by", User.Type).
			Unique().
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				annotations.Edge{
					Caption:      "Created By",
					RelationType: entity.RelationTypeO2M,
				},
			),
		edge.To("updated_by", User.Type).
			Unique().
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				annotations.Edge{
					Caption:      "Updated By",
					RelationType: entity.RelationTypeO2M,
				},
			),
	}
}
