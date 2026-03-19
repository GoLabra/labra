package repo

import (
	"app/ent"
	"context"
	"time"
)

type UserRefreshToken interface {
	Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*ent.UserRefreshToken, error)
	MarkUsedAndRevoked(ctx context.Context, id int, at time.Time) error
}
