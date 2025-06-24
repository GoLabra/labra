package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/GoLabra/labra/src/api/entgql/date"
	"github.com/GoLabra/labra/src/api/utils"

	"github.com/GoLabra/labra/src/api/entgql/annotations"
	"github.com/GoLabra/labra/src/api/entgql/entity"
)

// File holds the schema definition for the  File entity.
type File struct {
	ent.Schema
}

func (File) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.MultiOrder(),
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		annotations.Entity{
			Caption:      "File",
			Owner:        entity.EntityOwnerAdmin,
			DisplayField: "id",
		},
	}
}

// Fields of the  File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").DefaultFunc(utils.NewUUIDV7).Annotations(
			entgql.OrderField("id"),
			annotations.Field{
				Caption: "Id",
				Type:    entity.FieldTypeID,
			},
		),

		field.Time("created_at").
			Optional().
			Nillable().
			Nillable().
			Default(func(val string) time.Time {
				if val == "now()" {
					return time.Now()
				}
				v, err := time.Parse(date.DateTimeFormat, val)
				if err != nil {
					return time.Now()
				}
				return v
			}("now()")).
			SchemaType(map[string]string{
				dialect.MySQL:    "datetime",
				dialect.Postgres: "timestamp",
			}).
			Annotations(
				entgql.Type("DateTime"),
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
			Nillable().
			Default(func(val string) time.Time {
				if val == "now()" {
					return time.Now()
				}
				v, err := time.Parse(date.DateTimeFormat, val)
				if err != nil {
					return time.Now()
				}
				return v
			}("now()")).
			UpdateDefault(time.Now).
			SchemaType(map[string]string{
				dialect.MySQL:    "datetime",
				dialect.Postgres: "timestamp",
			}).
			Annotations(
				entgql.Type("DateTime"),
				entgql.OrderField("updatedAt"),
				annotations.Field{
					Caption:      "Updated At",
					Type:         entity.FieldTypeDateTime,
					DefaultValue: "now()",
				},
			),

		field.String("caption").
			Optional().
			Unique().
			SchemaType(map[string]string{
				dialect.MySQL:    "VARCHAR(255)",
				dialect.Postgres: "VARCHAR(255)",
			}).
			Annotations(
				entgql.OrderField("caption"),
				annotations.Field{
					Caption: "Caption",
					Type:    entity.FieldTypeShortText,
				},
			),

		field.String("name").
			Unique().
			SchemaType(map[string]string{
				dialect.MySQL:    "VARCHAR(255)",
				dialect.Postgres: "VARCHAR(255)",
			}).
			Annotations(
				entgql.OrderField("name"),
				annotations.Field{
					Caption: "Name",
					Type:    entity.FieldTypeShortText,
				},
			),

		field.String("storage_file_name").
			Unique().
			SchemaType(map[string]string{
				dialect.MySQL:    "VARCHAR(255)",
				dialect.Postgres: "VARCHAR(255)",
			}).
			Annotations(
				entgql.OrderField("storageFileName"),
				// entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				annotations.Field{
					Caption: "Storage File Name",
					Type:    entity.FieldTypeShortText,
				},
			),

		field.Int64("size").
			Annotations(
				entgql.OrderField("size"),
				// entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				annotations.Field{
					Caption: "Size",
					Type:    entity.FieldTypeInteger,
				},
			),

		field.String("content").
			SchemaType(map[string]string{
				dialect.MySQL:    "MEDIUMTEXT",
				dialect.Postgres: "TEXT",
			}).
			Annotations(
				entsql.Skip(),
				annotations.Field{
					Caption: "Content",
					Type:    entity.FieldTypeLongText,
				},
			),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("created_by", User.Type).
			Unique().
			Annotations(
				annotations.Edge{
					Caption:      "Created By",
					RelationType: entity.RelationTypeOne,
				},
			),

		edge.To("updated_by", User.Type).
			Unique().
			Annotations(
				annotations.Edge{
					Caption:      "Updated By",
					RelationType: entity.RelationTypeOne,
				},
			),
	}
}
