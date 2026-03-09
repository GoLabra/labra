package svc

import (
	"context"
	"errors"

	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/domain/repo"
	gqlgen "github.com/GoLabra/labra/entgql/generated"
)

type Cron struct {
	repository *repo.Repository
}

func NewCron(r *repo.Repository) *Cron {
	return &Cron{repository: r}
}

func (s *Cron) CreateSchedule(ctx context.Context, data gqlgen.CreateCronScheduleInput) (*gqlgen.CronSchedule, error) {
	enabled := true
	if data.Enabled != nil {
		enabled = *data.Enabled
	}
	timeoutSeconds := 300
	if data.TimeoutSeconds != nil {
		timeoutSeconds = *data.TimeoutSeconds
	}
	retentionDays := 30
	if data.RetentionDays != nil {
		retentionDays = *data.RetentionDays
	}
	if timeoutSeconds <= 0 {
		return nil, errors.New("timeoutSeconds must be greater than 0")
	}
	if retentionDays <= 0 {
		return nil, errors.New("retentionDays must be greater than 0")
	}
	data.Enabled = &enabled
	data.TimeoutSeconds = &timeoutSeconds
	data.RetentionDays = &retentionDays

	return s.repository.Cron.CreateSchedule(constants.WithInternalOperation(ctx), data)
}

func (s *Cron) UpdateSchedule(ctx context.Context, where gqlgen.CronScheduleWhereUniqueInput, data gqlgen.UpdateCronScheduleInput) (*gqlgen.CronSchedule, error) {
	if data.TimeoutSeconds != nil && *data.TimeoutSeconds <= 0 {
		return nil, errors.New("timeoutSeconds must be greater than 0")
	}
	if data.RetentionDays != nil && *data.RetentionDays <= 0 {
		return nil, errors.New("retentionDays must be greater than 0")
	}
	return s.repository.Cron.UpdateSchedule(constants.WithInternalOperation(ctx), where, data)
}

func (s *Cron) DeleteSchedule(ctx context.Context, where gqlgen.CronScheduleWhereUniqueInput) (*gqlgen.CronSchedule, error) {
	return s.repository.Cron.DeleteSchedule(constants.WithInternalOperation(ctx), where)
}

func (s *Cron) GetSchedule(ctx context.Context, where gqlgen.CronScheduleWhereUniqueInput) (*gqlgen.CronSchedule, error) {
	return s.repository.Cron.GetSchedule(constants.WithInternalOperation(ctx), where)
}

func (s *Cron) ListSchedules(ctx context.Context, enabled *bool, limit *int, offset *int) ([]*gqlgen.CronSchedule, error) {
	l := 50
	if limit != nil {
		l = *limit
	}
	if l <= 0 {
		l = 50
	}
	o := 0
	if offset != nil && *offset > 0 {
		o = *offset
	}
	return s.repository.Cron.ListSchedules(constants.WithInternalOperation(ctx), enabled, l, o)
}

func (s *Cron) ListJobs(ctx context.Context, scheduleID *string, limit *int, offset *int) ([]*gqlgen.CronJob, error) {
	l := 100
	if limit != nil {
		l = *limit
	}
	if l <= 0 {
		l = 100
	}
	o := 0
	if offset != nil && *offset > 0 {
		o = *offset
	}
	return s.repository.Cron.ListJobs(constants.WithInternalOperation(ctx), scheduleID, l, o)
}
