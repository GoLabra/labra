package svc

import (
	"context"

	gqlgen "github.com/GoLabra/labra/entgql/generated"
)

type Cron interface {
	CreateSchedule(ctx context.Context, data gqlgen.CreateCronScheduleInput) (*gqlgen.CronSchedule, error)
	UpdateSchedule(ctx context.Context, where gqlgen.CronScheduleWhereUniqueInput, data gqlgen.UpdateCronScheduleInput) (*gqlgen.CronSchedule, error)
	DeleteSchedule(ctx context.Context, where gqlgen.CronScheduleWhereUniqueInput) (*gqlgen.CronSchedule, error)
	GetSchedule(ctx context.Context, where gqlgen.CronScheduleWhereUniqueInput) (*gqlgen.CronSchedule, error)
	ListSchedules(ctx context.Context, enabled *bool, limit *int, offset *int) ([]*gqlgen.CronSchedule, error)
	ListJobs(ctx context.Context, scheduleID *string, limit *int, offset *int) ([]*gqlgen.CronJob, error)
}
