package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/lucsky/cuid"
	"github.com/robfig/cron/v3"
)

const (
	statusPending = "pending"
	statusRunning = "running"
	statusSuccess = "success"
	statusFailed  = "failed"
	statusTimeout = "timeout"

	defaultReconcileInterval = 15 * time.Second
)

type Handler interface {
	Execute(ctx context.Context, client *ent.Client, schedule CronSchedule) (int64, error)
}

type Scheduler struct {
	client           *ent.Client
	executableSchema graphql.ExecutableSchema
	reconcileEvery   time.Duration

	cronRunner *cron.Cron
	handlers   map[string]Handler

	mu         sync.Mutex
	entryByID  map[string]cron.EntryID
	loadedSpec map[string]string
	running    map[string]bool

	reconcileCancel context.CancelFunc
}

func New(client *ent.Client, executableSchema graphql.ExecutableSchema) *Scheduler {
	s := &Scheduler{
		client:           client,
		executableSchema: executableSchema,
		reconcileEvery:   defaultReconcileInterval,
		cronRunner:       cron.New(),
		handlers:         make(map[string]Handler),
		entryByID:        make(map[string]cron.EntryID),
		loadedSpec:       make(map[string]string),
		running:          make(map[string]bool),
	}

	s.handlers["graphql"] = &graphqlHandler{schema: executableSchema}
	s.handlers["builtin:cleanup_cron_history"] = &cleanupCronHistoryHandler{}
	return s
}

func (s *Scheduler) Start(ctx context.Context) error {
	if s.client == nil {
		return errors.New("scheduler client is nil")
	}
	if err := ensureTables(constants.WithInternalOperation(ctx), s.client, s.client.DialectName()); err != nil {
		return err
	}
	if err := seedDefaultSchedule(constants.WithInternalOperation(ctx), s.client, s.client.DialectName()); err != nil {
		return err
	}
	if err := s.reconcile(constants.WithInternalOperation(ctx)); err != nil {
		return err
	}

	s.cronRunner.Start()

	rCtx, cancel := context.WithCancel(context.Background())
	s.reconcileCancel = cancel
	go s.reconcileLoop(rCtx)

	return nil
}

func (s *Scheduler) Stop(ctx context.Context) error {
	if s.reconcileCancel != nil {
		s.reconcileCancel()
	}

	done := s.cronRunner.Stop().Done()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (s *Scheduler) reconcileLoop(ctx context.Context) {
	ticker := time.NewTicker(s.reconcileEvery)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.reconcile(constants.WithInternalOperation(context.Background())); err != nil {
				log.Printf("scheduler reconcile failed: %v", err)
			}
		}
	}
}

func (s *Scheduler) reconcile(ctx context.Context) error {
	schedules, err := loadSchedules(ctx, s.client, s.client.DialectName())
	if err != nil {
		return err
	}

	seen := make(map[string]struct{}, len(schedules))
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sch := range schedules {
		seen[sch.ID] = struct{}{}
		spec := scheduleSignature(sch)
		existingSpec, ok := s.loadedSpec[sch.ID]
		if ok && existingSpec == spec {
			continue
		}

		if eid, has := s.entryByID[sch.ID]; has {
			s.cronRunner.Remove(eid)
		}

		schedule := sch
		id, addErr := s.cronRunner.AddFunc(schedule.Expression, func() {
			s.runSchedule(schedule)
		})
		if addErr != nil {
			log.Printf("scheduler: skip invalid expression for %s: %v", schedule.Name, addErr)
			continue
		}
		s.entryByID[sch.ID] = id
		s.loadedSpec[sch.ID] = spec
	}

	for id, eid := range s.entryByID {
		if _, ok := seen[id]; ok {
			continue
		}
		s.cronRunner.Remove(eid)
		delete(s.entryByID, id)
		delete(s.loadedSpec, id)
		delete(s.running, id)
	}

	return nil
}

func (s *Scheduler) runSchedule(schedule CronSchedule) {
	if !s.markRunning(schedule.ID) {
		return
	}
	defer s.unmarkRunning(schedule.ID)

	startedAt := time.Now().UTC()
	jobID := cuid.New()
	iCtx := constants.WithInternalOperation(context.Background())

	if err := insertJob(iCtx, s.client, s.client.DialectName(), CronJob{
		ID:         jobID,
		ScheduleID: schedule.ID,
		StartedAt:  startedAt,
		Status:     statusRunning,
	}); err != nil {
		log.Printf("scheduler: failed to insert running job for %s: %v", schedule.Name, err)
		return
	}

	timeout := time.Duration(schedule.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 300 * time.Second
	}
	runCtx, cancel := context.WithTimeout(iCtx, timeout)
	defer cancel()

	var (
		status          = statusSuccess
		errText         *string
		recordsAffected *int64
	)

	handler, ok := s.handlers[schedule.Handler]
	if !ok {
		status = statusFailed
		msg := fmt.Sprintf("handler not found: %s", schedule.Handler)
		errText = &msg
	} else {
		affected, execErr := handler.Execute(runCtx, s.client, schedule)
		if execErr != nil {
			if errors.Is(execErr, context.DeadlineExceeded) || errors.Is(runCtx.Err(), context.DeadlineExceeded) {
				status = statusTimeout
			} else {
				status = statusFailed
			}
			msg := execErr.Error()
			errText = &msg
		} else {
			recordsAffected = &affected
		}
	}

	completedAt := time.Now().UTC()
	durationMs := completedAt.Sub(startedAt).Milliseconds()
	if err := completeJob(iCtx, s.client, s.client.DialectName(), JobResult{
		ID:              jobID,
		Status:          status,
		Error:           errText,
		RecordsAffected: recordsAffected,
		CompletedAt:     completedAt,
		DurationMs:      durationMs,
	}); err != nil {
		log.Printf("scheduler: failed to complete job %s: %v", jobID, err)
	}
}

func (s *Scheduler) markRunning(scheduleID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running[scheduleID] {
		return false
	}
	s.running[scheduleID] = true
	return true
}

func (s *Scheduler) unmarkRunning(scheduleID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.running, scheduleID)
}

type CronSchedule struct {
	ID             string
	Name           string
	Expression     string
	Handler        string
	HandlerConfig  map[string]any
	Enabled        bool
	TimeoutSeconds int
	RetentionDays  int
}

type CronJob struct {
	ID         string
	ScheduleID string
	StartedAt  time.Time
	Status     string
}

type JobResult struct {
	ID              string
	Status          string
	Error           *string
	RecordsAffected *int64
	CompletedAt     time.Time
	DurationMs      int64
}

func decodeConfig(raw string) map[string]any {
	if raw == "" {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return map[string]any{}
	}
	if out == nil {
		return map[string]any{}
	}
	return out
}

func scheduleSignature(s CronSchedule) string {
	cfgRaw, _ := json.Marshal(s.HandlerConfig)
	return s.Expression + "|" + s.Handler + "|" + string(cfgRaw) + "|" + fmt.Sprintf("%d", s.TimeoutSeconds)
}
