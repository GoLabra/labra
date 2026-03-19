package svc

import (
	"app/ent"
	"context"
	"time"
)

type UserRefreshToken interface {
	Persist(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	Use(ctx context.Context, tokenHash string) (*ent.User, *ent.Role, error)
}
