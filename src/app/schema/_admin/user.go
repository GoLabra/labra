package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/GoLabra/labra/src/api/entgql/annotations"
	"github.com/GoLabra/labra/src/api/entgql/entity"
	"github.com/lucsky/cuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.MultiOrder(),
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		annotations.Entity{
			Caption:      "User",
			Owner:        entity.EntityOwnerUser,
			DisplayField: "email",
		},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
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
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("created_by", User.Type).
			Unique().
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				annotations.Edge{
					Caption:      "Created By",
					RelationType: entity.RelationTypeOne,
				},
			).
			From("ref_created_by").
			Annotations(
				entgql.Skip(entgql.SkipAll),
				annotations.Edge{
					Caption:      "Ref Created By",
					RelationType: entity.RelationTypeOne,
				},
			),
		edge.To("updated_by", User.Type).
			Unique().
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				annotations.Edge{
					Caption:      "Updated By",
					RelationType: entity.RelationTypeOne,
				},
			).
			From("ref_updated_by").
			Annotations(
				entgql.Skip(entgql.SkipAll),
				annotations.Edge{
					Caption:      "Ref Updated By",
					RelationType: entity.RelationTypeOne,
				},
			),

		edge.To("admin_created_by", AdminUser.Type).
			Unique().
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				annotations.Edge{
					Caption:      "Admin Created By",
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
