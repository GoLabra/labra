package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/lucsky/cuid"
    "entgo.io/ent/dialect"
    "github.com/GoLabra/labra/entgql/date"
    "time"
    
	"github.com/GoLabra/labra/entgql/annotations"
	"github.com/GoLabra/labra/entgql/entity"
)

// Miau holds the schema definition for the  Miau entity.
type Miau struct {
	ent.Schema
}

func (Miau) Annotations() []schema.Annotation {
	annotationsList := []schema.Annotation{
		entgql.MultiOrder(),
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		annotations.Entity{
			Caption: "Miau",
			Owner:   entity.EntityOwnerUser,
            DisplayField: "name",
		},
	}
	return annotationsList
}

// Fields of the  Miau.
func (Miau) Fields() []ent.Field {
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
			Optional().
            Nillable().
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
	}
}

// Edges of the Miau.
func (Miau) Edges() []ent.Edge {
	return []ent.Edge{
        
        edge.To("created_by", User.Type).
            Unique().
            Annotations(
                entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				annotations.Edge{
					Caption:      "Created By",
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
	}
}