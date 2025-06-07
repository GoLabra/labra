package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/lucsky/cuid"
    "entgo.io/ent/dialect"
    "github.com/GoLabra/labra/src/api/entgql/date"
    "time"
    
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
			Caption: "File",
			Owner:   entity.EntityOwnerAdmin,
            DisplayField: "id",
		},
	}
}

// Fields of the  File.
func (File) Fields() []ent.Field {
	return []ent.Field{
            
		field.String("id").DefaultFunc(cuid.New).Annotations(
			entgql.OrderField("id"),
			annotations.Field{
				Caption:   "Id",
				Type: entity.FieldTypeID,
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
                    Caption:   "Created At",
                    Type: entity.FieldTypeDateTime,
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
                    Caption:   "Updated At",
                    Type: entity.FieldTypeDateTime,
                    DefaultValue: "now()",
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
					Caption:   "Name",
					Type: entity.FieldTypeShortText,
				},
			),
            
		field.String("content").
			SchemaType(map[string]string{
				dialect.MySQL:    "MEDIUMTEXT",
				dialect.Postgres: "TEXT",
			}).
			Annotations(
				entgql.OrderField("content"),
				annotations.Field{
					Caption:   "Content",
					Type: entity.FieldTypeLongText,
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