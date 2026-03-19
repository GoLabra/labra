package repo

import (
	"context"
	"time"

	"github.com/GoLabra/labra/entgql/ent"
)

type AdminRefreshToken interface {
	Create(ctx context.Context, adminUserID, tokenHash string, expiresAt time.Time) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*ent.AdminRefreshToken, error)
	MarkUsedAndRevoked(ctx context.Context, id string, at time.Time) error
}
