package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
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
	}
}
