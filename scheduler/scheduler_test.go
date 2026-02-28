package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/GoLabra/labra/entgql/ent"
	_ "github.com/mattn/go-sqlite3"
)

type handlerFunc func(context.Context, *ent.Client, CronSchedule) (int64, error)

func (f handlerFunc) Execute(ctx context.Context, client *ent.Client, schedule CronSchedule) (int64, error) {
	return f(ctx, client, schedule)
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := ensureTables(context.Background(), db, "sqlite"); err != nil {
		t.Fatalf("ensure tables: %v", err)
	}
	return db
}

func TestSeedDefaultScheduleIsIdempotent(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	if err := seedDefaultSchedule(ctx, db, "sqlite"); err != nil {
		t.Fatalf("first seed: %v", err)
	}
	if err := seedDefaultSchedule(ctx, db, "sqlite"); err != nil {
		t.Fatalf("second seed: %v", err)
	}

	var count int
	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM cron_schedules WHERE name = ?`, "Cleanup cron job history").Scan(&count)
	if err != nil {
		t.Fatalf("count seed rows: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 seeded schedule, got %d", count)
	}

	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM cron_schedules WHERE name = ?`, "Cleanup expired token blacklist").Scan(&count)
	if err != nil {
		t.Fatalf("count token cleanup seed rows: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 seeded token cleanup schedule, got %d", count)
	}

	var (
		handler   string
		expr      string
		cfgRaw    string
		isEnabled bool
	)
	err = db.QueryRowContext(ctx, `SELECT handler, expression, handler_config, enabled FROM cron_schedules WHERE name = ?`, "Cleanup expired token blacklist").
		Scan(&handler, &expr, &cfgRaw, &isEnabled)
	if err != nil {
		t.Fatalf("read seeded token cleanup schedule: %v", err)
	}
	if handler != "graphql" {
		t.Fatalf("expected graphql handler, got %s", handler)
	}
	if expr != "@hourly" {
		t.Fatalf("expected @hourly expression, got %s", expr)
	}
	if !isEnabled {
		t.Fatalf("expected token cleanup schedule to be enabled")
	}
	if !strings.Contains(cfgRaw, "cleanupExpiredTokens") {
		t.Fatalf("expected handler_config operation to include cleanupExpiredTokens")
	}
}

func TestLoadSchedulesReturnsEnabledOnly(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	enabledCfg, _ := json.Marshal(map[string]any{
		"operation": "mutation { noop }",
		"variables": map[string]any{"olderThanDays": 7},
	})
	_, err := db.ExecContext(ctx, `
INSERT INTO cron_schedules (id, name, expression, handler, handler_config, enabled, timeout_seconds, retention_days, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?),
       (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`,
		"s1", "enabled schedule", "@daily", "graphql", string(enabledCfg), true, 30, 10, time.Now().UTC(), time.Now().UTC(),
		"s2", "disabled schedule", "@daily", "graphql", string(enabledCfg), false, 30, 10, time.Now().UTC(), time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("insert schedules: %v", err)
	}

	schedules, err := loadSchedules(ctx, db, "sqlite")
	if err != nil {
		t.Fatalf("load schedules: %v", err)
	}
	if len(schedules) != 1 {
		t.Fatalf("expected 1 enabled schedule, got %d", len(schedules))
	}
	if schedules[0].ID != "s1" {
		t.Fatalf("expected s1, got %s", schedules[0].ID)
	}
	if schedules[0].HandlerConfig["operation"] == nil {
		t.Fatalf("expected handler_config to be decoded")
	}
}

func TestInsertAndCompleteJobLifecycle(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	now := time.Now().UTC()

	_, err := db.ExecContext(ctx, `
INSERT INTO cron_schedules (id, name, expression, handler, handler_config, enabled, timeout_seconds, retention_days, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, "s1", "schedule", "@daily", "graphql", `{}`, true, 300, 30, now, now)
	if err != nil {
		t.Fatalf("insert schedule: %v", err)
	}

	if err := insertJob(ctx, db, "sqlite", CronJob{
		ID:         "job1",
		ScheduleID: "s1",
		StartedAt:  now,
		Status:     statusRunning,
	}); err != nil {
		t.Fatalf("insert job: %v", err)
	}

	records := int64(13)
	completedAt := now.Add(2 * time.Second)
	if err := completeJob(ctx, db, "sqlite", JobResult{
		ID:              "job1",
		Status:          statusSuccess,
		RecordsAffected: &records,
		CompletedAt:     completedAt,
		DurationMs:      2000,
	}); err != nil {
		t.Fatalf("complete job: %v", err)
	}

	var (
		status   string
		duration int64
		affected int64
	)
	if err := db.QueryRowContext(ctx, `SELECT status, duration_ms, records_affected FROM cron_jobs WHERE id = ?`, "job1").
		Scan(&status, &duration, &affected); err != nil {
		t.Fatalf("read completed job: %v", err)
	}
	if status != statusSuccess {
		t.Fatalf("expected status %s, got %s", statusSuccess, status)
	}
	if duration != 2000 {
		t.Fatalf("expected duration 2000, got %d", duration)
	}
	if affected != 13 {
		t.Fatalf("expected affected 13, got %d", affected)
	}
}

func TestExtractCount(t *testing.T) {
	got := extractCount([]byte(`{"cleanupExpiredTokens":{"count":7}}`))
	if got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}

	got = extractCount([]byte(`{"cleanupExpiredTokens":{"ok":true}}`))
	if got != 0 {
		t.Fatalf("expected 0 when count missing, got %d", got)
	}
}

func TestCleanupCronHistoryHandlerExecute(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	now := time.Now().UTC()

	_, err := db.ExecContext(ctx, `
INSERT INTO cron_schedules (id, name, expression, handler, handler_config, enabled, timeout_seconds, retention_days, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, "s1", "schedule 1", "@daily", "builtin:cleanup_cron_history", `{}`, true, 300, 1, now, now)
	if err != nil {
		t.Fatalf("insert schedule: %v", err)
	}

	old := now.Add(-72 * time.Hour)
	recent := now.Add(-2 * time.Hour)
	_, err = db.ExecContext(ctx, `
INSERT INTO cron_jobs (id, started_at, completed_at, status, created_at, cron_schedule_cron_jobs)
VALUES (?, ?, ?, ?, ?, ?),
       (?, ?, ?, ?, ?, ?),
       (?, ?, ?, ?, ?, ?)
`,
		"j-old-success", old, old, statusSuccess, old, "s1",
		"j-old-running", old, nil, statusRunning, old, "s1",
		"j-recent-success", recent, recent, statusSuccess, recent, "s1",
	)
	if err != nil {
		t.Fatalf("insert jobs: %v", err)
	}

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))
	defer func() { _ = client.Close() }()

	handler := &cleanupCronHistoryHandler{}
	affected, err := handler.Execute(ctx, client, CronSchedule{})
	if err != nil {
		t.Fatalf("cleanup execute: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected row, got %d", affected)
	}

	var remaining int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM cron_jobs`).Scan(&remaining); err != nil {
		t.Fatalf("count remaining jobs: %v", err)
	}
	if remaining != 2 {
		t.Fatalf("expected 2 jobs to remain, got %d", remaining)
	}
}

func TestRunScheduleTimeoutMarksJobAsTimeout(t *testing.T) {
	db := newTestDB(t)
	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))
	defer func() { _ = client.Close() }()

	s := New(client, nil)
	s.handlers["test:timeout"] = handlerFunc(func(ctx context.Context, _ *ent.Client, _ CronSchedule) (int64, error) {
		<-ctx.Done()
		return 0, ctx.Err()
	})

	schedule := CronSchedule{
		ID:             "timeout-schedule",
		Name:           "timeout schedule",
		Expression:     "@every 1m",
		Handler:        "test:timeout",
		TimeoutSeconds: 1,
	}

	s.runSchedule(schedule)

	var (
		status string
		errMsg sql.NullString
	)
	err := db.QueryRowContext(context.Background(),
		`SELECT status, error FROM cron_jobs WHERE cron_schedule_cron_jobs = ? ORDER BY created_at DESC LIMIT 1`,
		schedule.ID,
	).Scan(&status, &errMsg)
	if err != nil {
		t.Fatalf("read timeout job: %v", err)
	}
	if status != statusTimeout {
		t.Fatalf("expected status %s, got %s", statusTimeout, status)
	}
	if !errMsg.Valid || errMsg.String == "" {
		t.Fatalf("expected timeout error message to be set")
	}
}

func TestRunScheduleOverlapGuardSkipsConcurrentRun(t *testing.T) {
	db := newTestDB(t)
	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))
	defer func() { _ = client.Close() }()

	started := make(chan struct{})
	release := make(chan struct{})

	s := New(client, nil)
	s.handlers["test:block"] = handlerFunc(func(ctx context.Context, _ *ent.Client, _ CronSchedule) (int64, error) {
		select {
		case <-started:
		default:
			close(started)
		}
		select {
		case <-release:
			return 1, nil
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	})

	schedule := CronSchedule{
		ID:             "overlap-schedule",
		Name:           "overlap schedule",
		Expression:     "@every 1m",
		Handler:        "test:block",
		TimeoutSeconds: 10,
	}

	done := make(chan struct{})
	go func() {
		s.runSchedule(schedule)
		close(done)
	}()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatalf("first run did not start in time")
	}

	s.runSchedule(schedule)

	close(release)
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("first run did not finish in time")
	}

	var count int
	err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM cron_jobs WHERE cron_schedule_cron_jobs = ?`, schedule.ID).Scan(&count)
	if err != nil {
		t.Fatalf("count overlap jobs: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 job due to overlap guard, got %d", count)
	}
}

func TestRunScheduleMissingHandlerMarksFailed(t *testing.T) {
	db := newTestDB(t)
	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))
	defer func() { _ = client.Close() }()

	s := New(client, nil)
	schedule := CronSchedule{
		ID:             "missing-handler",
		Name:           "missing handler",
		Expression:     "@every 1m",
		Handler:        "does:not:exist",
		TimeoutSeconds: 2,
	}

	s.runSchedule(schedule)

	var (
		status string
		errMsg sql.NullString
	)
	err := db.QueryRowContext(context.Background(),
		`SELECT status, error FROM cron_jobs WHERE cron_schedule_cron_jobs = ? ORDER BY created_at DESC LIMIT 1`,
		schedule.ID,
	).Scan(&status, &errMsg)
	if err != nil {
		t.Fatalf("read failed job: %v", err)
	}
	if status != statusFailed {
		t.Fatalf("expected status %s, got %s", statusFailed, status)
	}
	if !errMsg.Valid || !strings.Contains(errMsg.String, "handler not found") {
		t.Fatalf("expected failure error message to be set")
	}
}

func TestReconcileAddsAndRemovesSchedulesByEnabledFlag(t *testing.T) {
	db := newTestDB(t)
	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))
	defer func() { _ = client.Close() }()

	s := New(client, nil)
	ctx := context.Background()
	now := time.Now().UTC()

	_, err := db.ExecContext(ctx, `
INSERT INTO cron_schedules (id, name, expression, handler, handler_config, enabled, timeout_seconds, retention_days, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, "reconcile-s1", "reconcile schedule", "@daily", "graphql", `{"operation":"query { __typename }"}`, true, 60, 7, now, now)
	if err != nil {
		t.Fatalf("insert schedule: %v", err)
	}

	if err := s.reconcile(ctx); err != nil {
		t.Fatalf("reconcile enabled: %v", err)
	}
	if len(s.entryByID) != 1 {
		t.Fatalf("expected 1 active entry, got %d", len(s.entryByID))
	}
	if _, ok := s.entryByID["reconcile-s1"]; !ok {
		t.Fatalf("expected entry for reconcile-s1")
	}

	_, err = db.ExecContext(ctx, `UPDATE cron_schedules SET enabled = ? WHERE id = ?`, false, "reconcile-s1")
	if err != nil {
		t.Fatalf("disable schedule: %v", err)
	}

	if err := s.reconcile(ctx); err != nil {
		t.Fatalf("reconcile disabled: %v", err)
	}
	if len(s.entryByID) != 0 {
		t.Fatalf("expected 0 active entries after disable, got %d", len(s.entryByID))
	}
}

func TestReconcileReschedulesOnExpressionChange(t *testing.T) {
	db := newTestDB(t)
	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))
	defer func() { _ = client.Close() }()

	s := New(client, nil)
	ctx := context.Background()
	now := time.Now().UTC()

	_, err := db.ExecContext(ctx, `
INSERT INTO cron_schedules (id, name, expression, handler, handler_config, enabled, timeout_seconds, retention_days, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, "reconcile-s2", "reconcile expr", "@daily", "graphql", `{"operation":"query { __typename }"}`, true, 60, 7, now, now)
	if err != nil {
		t.Fatalf("insert schedule: %v", err)
	}

	if err := s.reconcile(ctx); err != nil {
		t.Fatalf("initial reconcile: %v", err)
	}
	before := s.entryByID["reconcile-s2"]

	_, err = db.ExecContext(ctx, `UPDATE cron_schedules SET expression = ? WHERE id = ?`, "@every 2m", "reconcile-s2")
	if err != nil {
		t.Fatalf("update expression: %v", err)
	}

	if err := s.reconcile(ctx); err != nil {
		t.Fatalf("reconcile after expression change: %v", err)
	}
	after := s.entryByID["reconcile-s2"]
	if before == after {
		t.Fatalf("expected entry id to change after expression update")
	}
}

func TestReconcileReschedulesOnHandlerConfigChange(t *testing.T) {
	db := newTestDB(t)
	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))
	defer func() { _ = client.Close() }()

	s := New(client, nil)
	ctx := context.Background()
	now := time.Now().UTC()

	_, err := db.ExecContext(ctx, `
INSERT INTO cron_schedules (id, name, expression, handler, handler_config, enabled, timeout_seconds, retention_days, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, "reconcile-s3", "reconcile config", "@daily", "graphql", `{"operation":"query { __typename }","variables":{"x":1}}`, true, 60, 7, now, now)
	if err != nil {
		t.Fatalf("insert schedule: %v", err)
	}

	if err := s.reconcile(ctx); err != nil {
		t.Fatalf("initial reconcile: %v", err)
	}
	before := s.entryByID["reconcile-s3"]

	_, err = db.ExecContext(ctx, `UPDATE cron_schedules SET handler_config = ? WHERE id = ?`, `{"operation":"query { __typename }","variables":{"x":2}}`, "reconcile-s3")
	if err != nil {
		t.Fatalf("update handler_config: %v", err)
	}

	if err := s.reconcile(ctx); err != nil {
		t.Fatalf("reconcile after handler_config change: %v", err)
	}
	after := s.entryByID["reconcile-s3"]
	if before == after {
		t.Fatalf("expected entry id to change after handler_config update")
	}
}

func TestSchedulerStartStopSeedsAndLoadsSchedules(t *testing.T) {
	db := newTestDB(t)
	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))
	defer func() { _ = client.Close() }()

	ctx := context.Background()
	now := time.Now().UTC()
	_, err := db.ExecContext(ctx, `
INSERT INTO cron_schedules (id, name, expression, handler, handler_config, enabled, timeout_seconds, retention_days, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, "manual-s1", "manual schedule", "@every 10m", "builtin:cleanup_cron_history", `{}`, true, 60, 7, now, now)
	if err != nil {
		t.Fatalf("insert manual schedule: %v", err)
	}

	s := New(client, nil)
	if err := s.Start(ctx); err != nil {
		t.Fatalf("start scheduler: %v", err)
	}
	stopped := false
	t.Cleanup(func() {
		if stopped {
			return
		}
		stopCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = s.Stop(stopCtx)
	})

	// Give reconcile startup a brief moment to register entries.
	time.Sleep(50 * time.Millisecond)

	s.mu.Lock()
	entries := len(s.entryByID)
	_, hasManual := s.entryByID["manual-s1"]
	s.mu.Unlock()
	if !hasManual {
		t.Fatalf("expected manual schedule to be registered")
	}
	if entries < 3 {
		t.Fatalf("expected at least 3 entries (manual + 2 default seeds), got %d", entries)
	}

	var seeded int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM cron_schedules WHERE name = ?`, "Cleanup cron job history").Scan(&seeded); err != nil {
		t.Fatalf("count seeded schedules: %v", err)
	}
	if seeded != 1 {
		t.Fatalf("expected seeded cleanup schedule count to be 1, got %d", seeded)
	}

	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM cron_schedules WHERE name = ?`, "Cleanup expired token blacklist").Scan(&seeded); err != nil {
		t.Fatalf("count seeded token cleanup schedule: %v", err)
	}
	if seeded != 1 {
		t.Fatalf("expected seeded token cleanup schedule count to be 1, got %d", seeded)
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := s.Stop(stopCtx); err != nil {
		t.Fatalf("stop scheduler: %v", err)
	}
	stopped = true
}
