package repo

import (
	"context"
	"time"

	"github.com/GoLabra/labra/tokenrevocation"
)

type TokenRevocation struct {
	client  tokenrevocation.SQLExecutor
	dialect string
}

func NewTokenRevocation(client interface {
	tokenrevocation.SQLExecutor
	DialectName() string
}) *TokenRevocation {
	return &TokenRevocation{
		client:  client,
		dialect: client.DialectName(),
	}
}

func (r *TokenRevocation) RevokeSubjectTokens(ctx context.Context, subjectEmail, subjectType string, revokedAt time.Time) error {
	return tokenrevocation.RevokeSubjectTokens(ctx, r.client, r.dialect, subjectEmail, subjectType, revokedAt)
}
