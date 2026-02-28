package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/ent"
)

type graphqlHandler struct {
	schema graphql.ExecutableSchema
}

func (h *graphqlHandler) Execute(ctx context.Context, _ *ent.Client, schedule CronSchedule) (int64, error) {
	operation, ok := schedule.HandlerConfig["operation"].(string)
	if !ok || operation == "" {
		return 0, errors.New("graphql handler requires handler_config.operation")
	}

	var variables map[string]any
	if v, ok := schedule.HandlerConfig["variables"].(map[string]any); ok {
		variables = v
	} else {
		variables = map[string]any{}
	}

	exec := executor.New(h.schema)
	raw := &graphql.RawParams{
		Query:     operation,
		Variables: variables,
	}

	iCtx := constants.WithInternalOperation(ctx)
	opCtx, errs := exec.CreateOperationContext(iCtx, raw)
	if len(errs) > 0 {
		return 0, fmt.Errorf("graphql operation context error: %v", errs.Error())
	}

	respHandler, innerCtx := exec.DispatchOperation(iCtx, opCtx)
	resp := respHandler(innerCtx)
	if resp == nil {
		return 0, errors.New("graphql response is nil")
	}
	if len(resp.Errors) > 0 {
		return 0, fmt.Errorf("graphql execution error: %v", resp.Errors.Error())
	}

	return extractCount(resp.Data), nil
}

type cleanupCronHistoryHandler struct{}

func (h *cleanupCronHistoryHandler) Execute(ctx context.Context, client *ent.Client, _ CronSchedule) (int64, error) {
	schedules, err := loadAllSchedules(ctx, client, client.DialectName())
	if err != nil {
		return 0, err
	}

	var total int64
	for _, schedule := range schedules {
		if schedule.RetentionDays <= 0 {
			continue
		}
		cutoff := time.Now().UTC().Add(-time.Duration(schedule.RetentionDays) * 24 * time.Hour)
		query := `DELETE FROM cron_jobs WHERE cron_schedule_cron_jobs = ` + bindVar(client.DialectName(), 1) +
			` AND status != ` + bindVar(client.DialectName(), 2) +
			` AND created_at < ` + bindVar(client.DialectName(), 3)
		res, execErr := client.ExecContext(ctx, query, schedule.ID, statusRunning, cutoff)
		if execErr != nil {
			return total, execErr
		}
		affected, rowsErr := res.RowsAffected()
		if rowsErr == nil {
			total += affected
		}
	}
	return total, nil
}

func loadAllSchedules(ctx context.Context, db sqlExecutor, dialect string) ([]CronSchedule, error) {
	query := `SELECT id, name, expression, handler, handler_config, enabled, timeout_seconds, retention_days FROM cron_schedules`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []CronSchedule
	for rows.Next() {
		var (
			s      CronSchedule
			cfgRaw sql.NullString
		)
		if err := rows.Scan(&s.ID, &s.Name, &s.Expression, &s.Handler, &cfgRaw, &s.Enabled, &s.TimeoutSeconds, &s.RetentionDays); err != nil {
			return nil, err
		}
		if cfgRaw.Valid {
			s.HandlerConfig = decodeConfig(cfgRaw.String)
		} else {
			s.HandlerConfig = map[string]any{}
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func extractCount(data []byte) int64 {
	if len(data) == 0 {
		return 0
	}
	var root map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		return 0
	}
	for _, v := range root {
		if m, ok := v.(map[string]any); ok {
			if countVal, has := m["count"]; has {
				switch c := countVal.(type) {
				case float64:
					return int64(c)
				case int64:
					return c
				case int:
					return int64(c)
				}
			}
		}
	}
	return 0
}
