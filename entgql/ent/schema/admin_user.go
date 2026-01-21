package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/GoLabra/labra/entgql/annotations"
	"github.com/GoLabra/labra/entgql/entity"
	"github.com/lucsky/cuid"
)

// AdminUser holds the schema definition for the AdminUser annotations.
type AdminUser struct {
	ent.Schema
}

func (AdminUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.MultiOrder(),
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		annotations.Entity{
			Caption:      "AdminUser",
			Owner:        entity.EntityOwnerAdmin,
			DisplayField: "name",
		},
	}
}

// Fields of the AdminUser.
func (AdminUser) Fields() []ent.Field {
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
			Optional().
			Annotations(
				entgql.OrderField("name"),
				annotations.Field{
					Caption: "Name",
					Type:    entity.FieldTypeShortText,
				},
			),

		field.String("email").
			Unique().
			NotEmpty().
			Annotations(
				entgql.OrderField("email"),
				annotations.Field{
					Caption: "Email",
					Type:    entity.FieldTypeEmail,
				},
			),

		field.String("password").
			NotEmpty().
			Annotations(
				annotations.Field{
					Caption: "Password",
					Type:    entity.FieldTypeShortText,
					Private: true,
				},
			),

		field.String("activation_key").
			Optional().
			Nillable().
			Annotations(
				annotations.Field{
					Caption: "Activation Key",
					Type:    entity.FieldTypeShortText,
					Private: true,
				},
			),

		field.String("first_name").
			NotEmpty().
			Annotations(
				entgql.OrderField("firstName"),
				annotations.Field{
					Caption: "First Name",
					Type:    entity.FieldTypeShortText,
				},
			),

		field.String("last_name").
			NotEmpty().
			Annotations(
				entgql.OrderField("lastName"),
				annotations.Field{
					Caption: "Last Name",
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

// Edges of the AdminUser.
func (AdminUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("admin_created_by", AdminUser.Type).
			Unique().
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				annotations.Edge{
					Caption:      "Admin Created By",
					RelationType: entity.RelationTypeOne,
				},
			).
			From("ref_admin_created_by").
			Annotations(
				entgql.Skip(entgql.SkipAll),
				annotations.Edge{
					Caption:      "Ref Admin Created By",
					RelationType: entity.RelationTypeOne,
				},
			),
		edge.To("admin_updated_by", AdminUser.Type).
			Unique().
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				annotations.Edge{
					Caption:      "Admin Updated By",
					RelationType: entity.RelationTypeOne,
				},
			).
			From("ref_admin_updated_by").
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				annotations.Edge{
					Caption:      "Ref Admin Updated By",
					RelationType: entity.RelationTypeOne,
				},
			),

		edge.To("roles", Role.Type).
			Annotations(
				annotations.Edge{
					Caption:      "Roles",
					RelationType: entity.RelationTypeM2M,
				},
			),

		edge.To("default_role", Role.Type).
			Unique().
			Annotations(
				annotations.Edge{
					Caption:      "Default Role",
					RelationType: entity.RelationTypeOne,
				},
			),
	}
}
