package schema

import (
	"fmt"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/GoLabra/labra/src/api/entgql/date"
	"github.com/lucsky/cuid"

	"github.com/GoLabra/labra/src/api/entgql/annotations"
	"github.com/GoLabra/labra/src/api/entgql/entity"
)

// Permission holds the schema definition for the  Permission entity.
type Permission struct {
	ent.Schema
}

func (Permission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
		entgql.MultiOrder(),
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		annotations.Entity{
			Caption:      "Permission",
			Owner:        entity.EntityOwnerAdmin,
			DisplayField: "id",
		},
	}
}

// Fields of the  Permission.
func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").DefaultFunc(cuid.New).Annotations(
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
				annotations.Field{
					Caption:      "Updated At",
					Type:         entity.FieldTypeDateTime,
					DefaultValue: "now()",
				},
			),

		field.String("entity").
			SchemaType(map[string]string{
				dialect.MySQL:    "VARCHAR(255)",
				dialect.Postgres: "VARCHAR(255)",
			}).
			Annotations(
				entgql.OrderField("entity"),
				annotations.Field{
					Caption: "Entity",
					Type:    entity.FieldTypeShortText,
				},
			),

		field.String("operation").
			Default("").
			Validate(func(val string) error {
				for _, acceptedValue := range []string{"Create", "Update", "Delete", "Read", ""} {
					if val == acceptedValue {
						return nil
					}
				}

				return fmt.Errorf("value \"%s\" is not an accepted value", val)
			}).
			Annotations(
				entgql.OrderField("operation"),
				annotations.Field{
					Caption:        "Operation",
					Type:           entity.FieldTypeSingleChoice,
					AcceptedValues: []string{"Create", "Update", "Delete", "Read", ""},
				},
			),
	}
}

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
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

		edge.To("role", Role.Type).
			Unique().
			Annotations(
				annotations.Edge{
					Caption:      "Role",
					RelationType: entity.RelationTypeO2M,
				},
			),
	}
}
