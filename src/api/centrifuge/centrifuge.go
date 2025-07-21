// Package centrifuge wraps the gocent client used for pushing events to
// Centrifugo. Files are located in src/api/centrifuge.
package centrifuge

import (
	"context"
	"encoding/json"

	"github.com/centrifugal/gocent/v3"
)

const (
	CodeGenerationChannel              = "channel"
	ErrCentrifugeClientNotSetInContext = "centrifuge client is not set in context"
)

// MessageContent represents payloads publishable to Centrifugo channels.
type MessageContent interface {
	Marshal() ([]byte, error)
}

// Publish sends the given content to the specified Centrifugo channel.
func Publish(ctx context.Context, client *gocent.Client, channel string, content any) error {
	c, err := json.Marshal(content)

	if err != nil {
		return err
	}

	_, err = client.Publish(ctx, channel, c)
	return err
}
