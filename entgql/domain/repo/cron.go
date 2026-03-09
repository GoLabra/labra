package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/GoLabra/labra/entgql/ent"
	gqlgen "github.com/GoLabra/labra/entgql/generated"
	"github.com/lucsky/cuid"
)

type Cron struct {
	client *ent.Client
}

func NewCron(client *ent.Client) *Cron {
	return &Cron{client: client}
}

func (r *Cron) CreateSchedule(ctx context.Context, data gqlgen.CreateCronScheduleInput) (*gqlgen.CronSchedule, error) {
	handlerConfig := map[string]any{}
	if data.HandlerConfig != nil {
		handlerConfig = data.HandlerConfig
	}
	cfgRaw, err := json.Marshal(handlerConfig)
	if err != nil {
		return nil, fmt.Errorf("marshal handlerConfig: %w", err)
	}

	now := time.Now().UTC()
	id := cuid.New()
	query := `INSERT INTO cron_schedules
		(id, name, expression, handler, handler_config, enabled, description, timeout_seconds, retention_days, created_at, updated_at)
		VALUES (` + bindVar(r.client.DialectName(), 1) + `, ` + bindVar(r.client.DialectName(), 2) + `, ` + bindVar(r.client.DialectName(), 3) + `, ` + bindVar(r.client.DialectName(), 4) + `, ` + bindVar(r.client.DialectName(), 5) + `, ` + bindVar(r.client.DialectName(), 6) + `, ` + bindVar(r.client.DialectName(), 7) + `, ` + bindVar(r.client.DialectName(), 8) + `, ` + bindVar(r.client.DialectName(), 9) + `, ` + bindVar(r.client.DialectName(), 10) + `, ` + bindVar(r.client.DialectName(), 11) + `)`
	if _, err := r.client.ExecContext(ctx, query, id, data.Name, data.Expression, data.Handler, string(cfgRaw), data.Enabled, data.Description, data.TimeoutSeconds, data.RetentionDays, now, now); err != nil {
		return nil, err
	}

	return r.loadScheduleByID(ctx, id)
}

func (r *Cron) UpdateSchedule(ctx context.Context, where gqlgen.CronScheduleWhereUniqueInput, data gqlgen.UpdateCronScheduleInput) (*gqlgen.CronSchedule, error) {
	existing, err := r.GetSchedule(ctx, where)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	sets := make([]string, 0, 10)
	args := make([]any, 0, 10)
	argIdx := 1
	addSet := func(field string, value any) {
		sets = append(sets, field+" = "+bindVar(r.client.DialectName(), argIdx))
		args = append(args, value)
		argIdx++
	}

	if data.Name != nil {
		addSet("name", *data.Name)
	}
	if data.Expression != nil {
		addSet("expression", *data.Expression)
	}
	if data.Handler != nil {
		addSet("handler", *data.Handler)
	}
	if data.TimeoutSeconds != nil {
		addSet("timeout_seconds", *data.TimeoutSeconds)
	}
	if data.RetentionDays != nil {
		addSet("retention_days", *data.RetentionDays)
	}
	if data.Enabled != nil {
		addSet("enabled", *data.Enabled)
	}
	if data.Description != nil {
		addSet("description", *data.Description)
	} else if data.ClearDescription != nil && *data.ClearDescription {
		addSet("description", nil)
	}
	if data.HandlerConfig != nil {
		raw, err := json.Marshal(data.HandlerConfig)
		if err != nil {
			return nil, fmt.Errorf("marshal handlerConfig: %w", err)
		}
		addSet("handler_config", string(raw))
	} else if data.ClearHandlerConfig != nil && *data.ClearHandlerConfig {
		addSet("handler_config", "{}")
	}

	addSet("updated_at", time.Now().UTC())
	whereCol, whereVal, err := cronWhere(where)
	if err != nil {
		return nil, err
	}

	if len(sets) == 0 {
		return existing, nil
	}
	query := `UPDATE cron_schedules SET ` + strings.Join(sets, ", ") + ` WHERE ` + whereCol + ` = ` + bindVar(r.client.DialectName(), argIdx)
	args = append(args, whereVal)
	if _, err := r.client.ExecContext(ctx, query, args...); err != nil {
		return nil, err
	}

	return r.loadScheduleByID(ctx, existing.ID)
}

func (r *Cron) DeleteSchedule(ctx context.Context, where gqlgen.CronScheduleWhereUniqueInput) (*gqlgen.CronSchedule, error) {
	existing, err := r.GetSchedule(ctx, where)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	whereCol, whereVal, err := cronWhere(where)
	if err != nil {
		return nil, err
	}
	query := `DELETE FROM cron_schedules WHERE ` + whereCol + ` = ` + bindVar(r.client.DialectName(), 1)
	if _, err := r.client.ExecContext(ctx, query, whereVal); err != nil {
		return nil, err
	}

	return existing, nil
}

func (r *Cron) GetSchedule(ctx context.Context, where gqlgen.CronScheduleWhereUniqueInput) (*gqlgen.CronSchedule, error) {
	col, val, err := cronWhere(where)
	if err != nil {
		return nil, err
	}
	query := `SELECT id, name, expression, handler, handler_config, enabled, description, timeout_seconds, retention_days, created_at, updated_at
		FROM cron_schedules WHERE ` + col + ` = ` + bindVar(r.client.DialectName(), 1)
	rows, err := r.client.QueryContext(ctx, query, val)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	s, err := scanCronSchedule(rows)
	if err != nil {
		return nil, err
	}
	return s, rows.Err()
}

func (r *Cron) ListSchedules(ctx context.Context, enabled *bool, limit int, offset int) ([]*gqlgen.CronSchedule, error) {
	query := `SELECT id, name, expression, handler, handler_config, enabled, description, timeout_seconds, retention_days, created_at, updated_at
		FROM cron_schedules`
	args := make([]any, 0, 3)
	if enabled != nil {
		query += ` WHERE enabled = ` + bindVar(r.client.DialectName(), 1)
		args = append(args, *enabled)
	}
	query += ` ORDER BY created_at DESC LIMIT ` + bindVar(r.client.DialectName(), len(args)+1) + ` OFFSET ` + bindVar(r.client.DialectName(), len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*gqlgen.CronSchedule
	for rows.Next() {
		s, err := scanCronSchedule(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *Cron) ListJobs(ctx context.Context, scheduleID *string, limit int, offset int) ([]*gqlgen.CronJob, error) {
	query := `SELECT id, started_at, completed_at, status, error, records_affected, duration_ms, created_at, cron_schedule_cron_jobs
		FROM cron_jobs`
	args := make([]any, 0, 3)
	if scheduleID != nil && *scheduleID != "" {
		query += ` WHERE cron_schedule_cron_jobs = ` + bindVar(r.client.DialectName(), 1)
		args = append(args, *scheduleID)
	}
	query += ` ORDER BY created_at DESC LIMIT ` + bindVar(r.client.DialectName(), len(args)+1) + ` OFFSET ` + bindVar(r.client.DialectName(), len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*gqlgen.CronJob
	for rows.Next() {
		job, err := scanCronJob(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, job)
	}
	return out, rows.Err()
}

func (r *Cron) loadScheduleByID(ctx context.Context, id string) (*gqlgen.CronSchedule, error) {
	return r.GetSchedule(ctx, gqlgen.CronScheduleWhereUniqueInput{ID: &id})
}

func cronWhere(where gqlgen.CronScheduleWhereUniqueInput) (string, string, error) {
	if where.ID != nil && *where.ID != "" {
		return "id", *where.ID, nil
	}
	if where.Name != nil && *where.Name != "" {
		return "name", *where.Name, nil
	}
	return "", "", errors.New("where requires id or name")
}

func scanCronSchedule(rows *sql.Rows) (*gqlgen.CronSchedule, error) {
	var (
		s         gqlgen.CronSchedule
		cfgRaw    sql.NullString
		desc      sql.NullString
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)
	if err := rows.Scan(&s.ID, &s.Name, &s.Expression, &s.Handler, &cfgRaw, &s.Enabled, &desc, &s.TimeoutSeconds, &s.RetentionDays, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	if cfgRaw.Valid && strings.TrimSpace(cfgRaw.String) != "" {
		var cfg map[string]any
		if err := json.Unmarshal([]byte(cfgRaw.String), &cfg); err != nil {
			cfg = map[string]any{}
		}
		s.HandlerConfig = cfg
	} else {
		s.HandlerConfig = map[string]any{}
	}
	if desc.Valid {
		s.Description = &desc.String
	}
	if createdAt.Valid {
		v := createdAt.Time.UTC()
		s.CreatedAt = &v
	}
	if updatedAt.Valid {
		v := updatedAt.Time.UTC()
		s.UpdatedAt = &v
	}
	return &s, nil
}

func scanCronJob(rows *sql.Rows) (*gqlgen.CronJob, error) {
	var (
		j           gqlgen.CronJob
		completedAt sql.NullTime
		errText     sql.NullString
		records     sql.NullInt64
		duration    sql.NullInt64
	)
	if err := rows.Scan(&j.ID, &j.StartedAt, &completedAt, &j.Status, &errText, &records, &duration, &j.CreatedAt, &j.CronScheduleID); err != nil {
		return nil, err
	}
	if completedAt.Valid {
		v := completedAt.Time.UTC()
		j.CompletedAt = &v
	}
	if errText.Valid {
		j.Error = &errText.String
	}
	if records.Valid {
		v := int(records.Int64)
		j.RecordsAffected = &v
	}
	if duration.Valid {
		v := int(duration.Int64)
		j.DurationMs = &v
	}
	return &j, nil
}

func bindVar(dialect string, idx int) string {
	d := strings.ToLower(strings.TrimSpace(dialect))
	if d == "postgres" || d == "postgresql" {
		return fmt.Sprintf("$%d", idx)
	}
	return "?"
}
