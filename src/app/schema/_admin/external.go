// Package schema defines ent schemas for admin entities. File path:
// src/app/schema/_admin/external.go.
package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// User defines the minimal fields required for external user references.
type User struct {
	ent.Schema
}

func (User) IsExternal() bool {
	return true
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// entsql.Skip(),
		entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
	}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Annotations(
				entgql.Type("ID"),
			),
	}
}

// File defines the minimal fields required for external file references.
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

// Edges exposes additional edges injected by admin code generation.
func (File) Edges() []ent.Edge {
	return additionalFileEdges
}
