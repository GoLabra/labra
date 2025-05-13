package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/GoLabra/labra/src/api/entgql/annotations"
	"github.com/GoLabra/labra/src/api/entgql/entity"
	"github.com/lucsky/cuid"
)

// Role holds the schema definition for the Role annotations.
type Role struct {
	ent.Schema
}

func (Role) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.MultiOrder(),
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		annotations.Entity{
			Caption:      "Role",
			Owner:        entity.EntityOwnerAdmin,
			DisplayField: "name",
		},
	}
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
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
			Unique().
			Annotations(
				entgql.OrderField("name"),
				annotations.Field{
					Caption: "Name",
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

// Edges of the Role.
func (Role) Edges() []ent.Edge {
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

		edge.From("user_roles", User.Type).
			Ref("roles").
			Annotations(
				annotations.Edge{
					Caption:      "Role Users",
					RelationType: entity.RelationTypeM2M,
				},
			),

		edge.From("permissions", Permission.Type).
			Ref("role").
			Annotations(
				annotations.Edge{
					Caption:      "Permissions",
					RelationType: entity.RelationTypeM2O,
				},
			),
	}
}
