package repo

import (
	"context"
	"time"
)

type TokenRevocation interface {
	RevokeSubjectTokens(ctx context.Context, subjectEmail, subjectType string, revokedAt time.Time) error
}
