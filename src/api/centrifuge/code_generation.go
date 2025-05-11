package centrifuge

import (
	"context"
	"time"

	"github.com/centrifugal/gocent/v3"
)

type AppStatusEvent string

const (
	AppStatusGenerationStarted   AppStatusEvent = "CODE_GENERATION_STARTED"
	AppStatusGenerationCompleted AppStatusEvent = "CODE_GENERATION_COMPLETED"
	AppStatusGenerationReverted  AppStatusEvent = "CODE_GENERATION_REVERTED"
	AppStatusGenerationFailed    AppStatusEvent = "CODE_GENERATION_FAILED"
	AppStatusShutdown            AppStatusEvent = "SHUTDOWN"
	AppStatusUp                  AppStatusEvent = "UP"
)

type CodeGenerationMessage struct {
	Event     AppStatusEvent `json:"event"`
	Timestamp time.Time      `json:"timestamp"`
	Entities  any            `json:"entities,omitempty"`
}

func PublishAppStatusMessage(ctx context.Context, client *gocent.Client, event AppStatusEvent, entities any) error {
	content := CodeGenerationMessage{
		Event:     event,
		Timestamp: time.Now(),
		Entities:  entities,
	}

	return Publish(ctx, client, CodeGenerationChannel, content)
}
