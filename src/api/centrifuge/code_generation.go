// Package centrifuge contains helpers for publishing code generation events to
// Centrifugo. File location: src/api/centrifuge.
package centrifuge

import (
	"context"
	"time"

	"github.com/centrifugal/gocent/v3"
)

// AppStatusEvent enumerates application lifecycle events published to Centrifugo.
type AppStatusEvent string

const (
	AppStatusGenerationStarted   AppStatusEvent = "CODE_GENERATION_STARTED"
	AppStatusGenerationCompleted AppStatusEvent = "CODE_GENERATION_COMPLETED"
	AppStatusGenerationReverted  AppStatusEvent = "CODE_GENERATION_REVERTED"
	AppStatusGenerationFailed    AppStatusEvent = "CODE_GENERATION_FAILED"
	AppStatusShutdown            AppStatusEvent = "SHUTDOWN"
	AppStatusUp                  AppStatusEvent = "UP"
)

// CodeGenerationMessage is the payload published for code generation events.
type CodeGenerationMessage struct {
	Event     AppStatusEvent `json:"event"`
	Timestamp time.Time      `json:"timestamp"`
	Entities  any            `json:"entities,omitempty"`
}

// PublishAppStatusMessage publishes a code generation status message to Centrifugo.
func PublishAppStatusMessage(ctx context.Context, client *gocent.Client, event AppStatusEvent, entities any) error {
	content := CodeGenerationMessage{
		Event:     event,
		Timestamp: time.Now(),
		Entities:  entities,
	}

	return Publish(ctx, client, CodeGenerationChannel, content)
}
