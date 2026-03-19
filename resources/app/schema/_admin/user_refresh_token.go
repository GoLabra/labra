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
)

type UserRefreshToken struct {
	ent.Schema
}

func (UserRefreshToken) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(entgql.SkipAll),
		annotations.Entity{
			Caption:      "User Refresh Token",
			Owner:        entity.EntityOwnerUser,
			DisplayField: "token_hash",
		},
	}
}

func (UserRefreshToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("token_hash").
			NotEmpty().
			Unique(),
		field.Time("expires_at"),
		field.Time("used_at").
			Optional().
			Nillable(),
		field.Time("revoked_at").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (UserRefreshToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Unique().
			Required(),
	}
}
