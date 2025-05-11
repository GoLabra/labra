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

type MessageContent interface {
	Marshal() ([]byte, error)
}

func Publish(ctx context.Context, client *gocent.Client, channel string, content any) error {
	c, err := json.Marshal(content)

	if err != nil {
		return err
	}

	_, err = client.Publish(ctx, channel, c)
	return err
}
