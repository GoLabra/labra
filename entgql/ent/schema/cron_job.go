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

// CronJob holds the schema definition for the CronJob entity.
type CronJob struct {
	ent.Schema
}

func (CronJob) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.MultiOrder(),
		entgql.RelayConnection(),
		annotations.Entity{
			Caption:      "Cron Job",
			Owner:        entity.EntityOwnerAdmin,
			DisplayField: "id",
		},
	}
}

func (CronJob) Fields() []ent.Field {
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
		field.Time("started_at").
			Default(time.Now).
			Annotations(
				entgql.OrderField("startedAt"),
				annotations.Field{
					Caption:      "Started At",
					Type:         entity.FieldTypeDateTime,
					DefaultValue: "now()",
				},
			),
		field.Time("completed_at").
			Optional().
			Nillable().
			Annotations(
				entgql.OrderField("completedAt"),
				annotations.Field{
					Caption: "Completed At",
					Type:    entity.FieldTypeDateTime,
				},
			),
		field.Enum("status").
			Values("pending", "running", "success", "failed", "timeout").
			Default("pending").
			Annotations(
				entgql.OrderField("status"),
				annotations.Field{
					Caption:        "Status",
					Type:           entity.FieldTypeSingleChoice,
					DefaultValue:   "pending",
					AcceptedValues: []string{"pending", "running", "success", "failed", "timeout"},
				},
			),
		field.String("error").
			Optional().
			Nillable().
			Annotations(
				entgql.OrderField("error"),
				annotations.Field{
					Caption: "Error",
					Type:    entity.FieldTypeLongText,
				},
			),
		field.Int64("records_affected").
			Optional().
			Nillable().
			Annotations(
				entgql.OrderField("recordsAffected"),
				annotations.Field{
					Caption: "Records Affected",
					Type:    entity.FieldTypeInteger,
				},
			),
		field.Int64("duration_ms").
			Optional().
			Nillable().
			Annotations(
				entgql.OrderField("durationMs"),
				annotations.Field{
					Caption: "Duration Ms",
					Type:    entity.FieldTypeInteger,
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
	}
}

func (CronJob) Edges() []ent.Edge {
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
		edge.From("cron_schedule", CronSchedule.Type).
			Ref("cron_jobs").
			Unique().
			Required().
			Annotations(
				annotations.Edge{
					Caption:      "Cron Schedule",
					RelationType: entity.RelationTypeM2O,
				},
			),
	}
}
