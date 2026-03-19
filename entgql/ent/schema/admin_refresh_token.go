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

// AdminRefreshToken stores server-side refresh tokens for admin sessions.
type AdminRefreshToken struct {
	ent.Schema
}

func (AdminRefreshToken) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(entgql.SkipAll),
		annotations.Entity{
			Caption:      "AdminRefreshToken",
			Owner:        entity.EntityOwnerAdmin,
			DisplayField: "id",
		},
	}
}

func (AdminRefreshToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(cuid.New).
			Annotations(
				annotations.Field{
					Caption: "ID",
					Type:    entity.FieldTypeID,
				},
			),
		field.String("token_hash").
			Unique().
			NotEmpty().
			Sensitive().
			Annotations(
				annotations.Field{
					Caption: "Token Hash",
					Type:    entity.FieldTypeShortText,
					Private: true,
				},
			),
		field.Time("expires_at").
			Annotations(
				annotations.Field{
					Caption: "Expires At",
					Type:    entity.FieldTypeDateTime,
				},
			),
		field.Time("used_at").
			Optional().
			Nillable().
			Annotations(
				annotations.Field{
					Caption: "Used At",
					Type:    entity.FieldTypeDateTime,
				},
			),
		field.Time("revoked_at").
			Optional().
			Nillable().
			Annotations(
				annotations.Field{
					Caption: "Revoked At",
					Type:    entity.FieldTypeDateTime,
				},
			),
		field.Time("created_at").
			Optional().
			Nillable().
			Immutable().
			Default(time.Now).
			Annotations(
				annotations.Field{
					Caption: "Created At",
					Type:    entity.FieldTypeDateTime,
				},
			),
		field.Time("updated_at").
			Optional().
			Nillable().
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(
				annotations.Field{
					Caption: "Updated At",
					Type:    entity.FieldTypeDateTime,
				},
			),
	}
}

func (AdminRefreshToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("admin_user", AdminUser.Type).
			Ref("refresh_tokens").
			Unique().
			Required().
			Annotations(
				annotations.Edge{
					Caption:      "Admin User",
					RelationType: entity.RelationTypeM2O,
				},
			),
	}
}
