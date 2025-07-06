package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

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

// func (User) Fields() []ent.Field {
// 	return adminEnt.User{}.Fields()
// }

// func (User) Edges() []ent.Edge {
// 	return adminEnt.User{}.Edges()
// }
