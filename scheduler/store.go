package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lucsky/cuid"
)

type sqlExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func ensureTables(ctx context.Context, db sqlExecutor, dialect string) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS cron_schedules (
	id VARCHAR(64) PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE,
	expression VARCHAR(255) NOT NULL,
	handler VARCHAR(255) NOT NULL,
	handler_config TEXT,
	enabled BOOLEAN NOT NULL DEFAULT TRUE,
	description TEXT,
	timeout_seconds INTEGER NOT NULL DEFAULT 300,
	retention_days INTEGER NOT NULL DEFAULT 30,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
)`,
		`CREATE TABLE IF NOT EXISTS cron_jobs (
	id VARCHAR(64) PRIMARY KEY,
	started_at TIMESTAMP NOT NULL,
	completed_at TIMESTAMP NULL,
	status VARCHAR(32) NOT NULL,
	error TEXT NULL,
	records_affected BIGINT NULL,
	duration_ms BIGINT NULL,
	created_at TIMESTAMP NOT NULL,
	cron_schedule_cron_jobs VARCHAR(64) NOT NULL
)`,
		`CREATE INDEX IF NOT EXISTS idx_cron_jobs_schedule ON cron_jobs(cron_schedule_cron_jobs)`,
		`CREATE INDEX IF NOT EXISTS idx_cron_jobs_created_at ON cron_jobs(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_cron_jobs_status ON cron_jobs(status)`,
	}
	for _, stmt := range stmts {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

func loadSchedules(ctx context.Context, db sqlExecutor, dialect string) ([]CronSchedule, error) {
	query := `SELECT id, name, expression, handler, handler_config, enabled, timeout_seconds, retention_days FROM cron_schedules WHERE enabled = ` + bindVar(dialect, 1)
	rows, err := db.QueryContext(ctx, query, true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []CronSchedule
	for rows.Next() {
		var (
			s          CronSchedule
			handlerCfg sql.NullString
			enabled    bool
			timeout    int
			retention  int
		)
		if err := rows.Scan(&s.ID, &s.Name, &s.Expression, &s.Handler, &handlerCfg, &enabled, &timeout, &retention); err != nil {
			return nil, err
		}
		s.Enabled = enabled
		s.TimeoutSeconds = timeout
		s.RetentionDays = retention
		if handlerCfg.Valid {
			s.HandlerConfig = decodeConfig(handlerCfg.String)
		} else {
			s.HandlerConfig = map[string]any{}
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func insertJob(ctx context.Context, db sqlExecutor, dialect string, job CronJob) error {
	now := time.Now().UTC()
	query := `INSERT INTO cron_jobs (id, started_at, status, created_at, cron_schedule_cron_jobs) VALUES (` +
		bindVar(dialect, 1) + `, ` + bindVar(dialect, 2) + `, ` + bindVar(dialect, 3) + `, ` + bindVar(dialect, 4) + `, ` + bindVar(dialect, 5) + `)`
	_, err := db.ExecContext(ctx, query, job.ID, job.StartedAt, job.Status, now, job.ScheduleID)
	return err
}

func completeJob(ctx context.Context, db sqlExecutor, dialect string, r JobResult) error {
	query := `UPDATE cron_jobs SET status = ` + bindVar(dialect, 1) +
		`, error = ` + bindVar(dialect, 2) +
		`, records_affected = ` + bindVar(dialect, 3) +
		`, duration_ms = ` + bindVar(dialect, 4) +
		`, completed_at = ` + bindVar(dialect, 5) +
		` WHERE id = ` + bindVar(dialect, 6)
	_, err := db.ExecContext(ctx, query, r.Status, r.Error, r.RecordsAffected, r.DurationMs, r.CompletedAt, r.ID)
	return err
}

func seedDefaultSchedule(ctx context.Context, db sqlExecutor, dialect string) error {
	now := time.Now().UTC()
	cronHistoryCfgRaw, _ := json.Marshal(map[string]any{})
	expiredTokenCfgRaw, _ := json.Marshal(map[string]any{
		"operation": "mutation CleanupExpiredTokens { cleanupExpiredTokens { count } }",
		"variables": map[string]any{},
	})

	switch normalizeDialect(dialect) {
	case "postgres":
		historyQuery := `
INSERT INTO cron_schedules (id, name, expression, handler, handler_config, enabled, description, timeout_seconds, retention_days, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (name) DO NOTHING`
		if _, err := db.ExecContext(ctx, historyQuery, cuid.New(), "Cleanup cron job history", "@daily", "builtin:cleanup_cron_history", string(cronHistoryCfgRaw), true, "Cleanup old cron job history", 300, 30, now, now); err != nil {
			return err
		}

		tokenQuery := `
INSERT INTO cron_schedules (id, name, expression, handler, handler_config, enabled, description, timeout_seconds, retention_days, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (name) DO NOTHING`
		_, err := db.ExecContext(ctx, tokenQuery, cuid.New(), "Cleanup expired token blacklist", "@hourly", "graphql", string(expiredTokenCfgRaw), true, "Delete expired entries from token_blacklist", 300, 30, now, now)
		return err
	default:
		historyQuery := `
INSERT OR IGNORE INTO cron_schedules (id, name, expression, handler, handler_config, enabled, description, timeout_seconds, retention_days, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		if _, err := db.ExecContext(ctx, historyQuery, cuid.New(), "Cleanup cron job history", "@daily", "builtin:cleanup_cron_history", string(cronHistoryCfgRaw), true, "Cleanup old cron job history", 300, 30, now, now); err != nil {
			return err
		}

		tokenQuery := `
INSERT OR IGNORE INTO cron_schedules (id, name, expression, handler, handler_config, enabled, description, timeout_seconds, retention_days, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := db.ExecContext(ctx, tokenQuery, cuid.New(), "Cleanup expired token blacklist", "@hourly", "graphql", string(expiredTokenCfgRaw), true, "Delete expired entries from token_blacklist", 300, 30, now, now)
		return err
	}
}

func bindVar(dialect string, idx int) string {
	if normalizeDialect(dialect) == "postgres" {
		return fmt.Sprintf("$%d", idx)
	}
	return "?"
}

func normalizeDialect(dialect string) string {
	d := strings.ToLower(strings.TrimSpace(dialect))
	switch d {
	case "postgresql":
		return "postgres"
	case "sqlite3":
		return "sqlite"
	default:
		return d
	}
}
