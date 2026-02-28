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

// CronSchedule holds the schema definition for the CronSchedule entity.
type CronSchedule struct {
	ent.Schema
}

func (CronSchedule) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.MultiOrder(),
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		annotations.Entity{
			Caption:      "Cron Schedule",
			Owner:        entity.EntityOwnerAdmin,
			DisplayField: "name",
		},
	}
}

func (CronSchedule) Fields() []ent.Field {
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
		field.String("name").
			NotEmpty().
			Unique().
			Annotations(
				entgql.OrderField("name"),
				annotations.Field{
					Caption: "Name",
					Type:    entity.FieldTypeShortText,
				},
			),
		field.String("expression").
			NotEmpty().
			Annotations(
				entgql.OrderField("expression"),
				annotations.Field{
					Caption: "Expression",
					Type:    entity.FieldTypeShortText,
				},
			),
		field.String("handler").
			NotEmpty().
			Annotations(
				entgql.OrderField("handler"),
				annotations.Field{
					Caption: "Handler",
					Type:    entity.FieldTypeShortText,
				},
			),
		field.JSON("handler_config", map[string]interface{}{}).
			Optional().
			Default(map[string]interface{}{}).
			Annotations(
				annotations.Field{
					Caption: "Handler Config",
					Type:    entity.FieldTypeJson,
				},
			),
		field.Bool("enabled").
			Default(true).
			Annotations(
				entgql.OrderField("enabled"),
				annotations.Field{
					Caption:      "Enabled",
					Type:         entity.FieldTypeBoolean,
					DefaultValue: "true",
				},
			),
		field.String("description").
			Optional().
			Nillable().
			Annotations(
				entgql.OrderField("description"),
				annotations.Field{
					Caption: "Description",
					Type:    entity.FieldTypeLongText,
				},
			),
		field.Int("timeout_seconds").
			Default(300).
			Positive().
			Annotations(
				entgql.OrderField("timeoutSeconds"),
				annotations.Field{
					Caption:      "Timeout Seconds",
					Type:         entity.FieldTypeInteger,
					DefaultValue: "300",
				},
			),
		field.Int("retention_days").
			Default(30).
			Positive().
			Annotations(
				entgql.OrderField("retentionDays"),
				annotations.Field{
					Caption:      "Retention Days",
					Type:         entity.FieldTypeInteger,
					DefaultValue: "30",
				},
			),
		field.Time("created_at").
			Optional().
			Nillable().
			Immutable().
			Default(time.Now).
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
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
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				entgql.OrderField("updatedAt"),
				annotations.Field{
					Caption:      "Updated At",
					Type:         entity.FieldTypeDateTime,
					DefaultValue: "now()",
				},
			),
	}
}

func (CronSchedule) Edges() []ent.Edge {
	return []ent.Edge{
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
		edge.To("cron_jobs", CronJob.Type).
			Annotations(
				annotations.Edge{
					Caption:      "Cron Jobs",
					RelationType: entity.RelationTypeO2M,
				},
			),
	}
}
