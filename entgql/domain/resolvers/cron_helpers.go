package resolvers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/GoLabra/labra/constants"
	gqlgen "github.com/GoLabra/labra/entgql/generated"
)

type sqlExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	DialectName() string
}

func cronDB(ctx context.Context) (sqlExecutor, context.Context, error) {
	client, ok := ctx.Value(constants.AdminEntClientContextValue).(sqlExecutor)
	if !ok || client == nil {
		return nil, nil, errors.New("admin ent client not found in context")
	}
	return client, constants.WithInternalOperation(ctx), nil
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

func loadCronScheduleByWhere(ctx context.Context, db sqlExecutor, where gqlgen.CronScheduleWhereUniqueInput) (*gqlgen.CronSchedule, error) {
	col, val, err := cronWhere(where)
	if err != nil {
		return nil, err
	}
	query := `SELECT id, name, expression, handler, handler_config, enabled, description, timeout_seconds, retention_days, created_at, updated_at
		FROM cron_schedules WHERE ` + col + ` = ` + cronBindVar(db.DialectName(), 1)
	rows, err := db.QueryContext(ctx, query, val)
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

func loadCronScheduleByID(ctx context.Context, db sqlExecutor, id string) (*gqlgen.CronSchedule, error) {
	return loadCronScheduleByWhere(ctx, db, gqlgen.CronScheduleWhereUniqueInput{ID: &id})
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

func cronBindVar(dialect string, idx int) string {
	d := strings.ToLower(strings.TrimSpace(dialect))
	if d == "postgres" || d == "postgresql" {
		return fmt.Sprintf("$%d", idx)
	}
	return "?"
}
