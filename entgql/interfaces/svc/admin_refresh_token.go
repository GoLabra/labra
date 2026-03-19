package svc

import (
	"context"
	"time"

	"github.com/GoLabra/labra/entgql/ent"
)

type AdminRefreshToken interface {
	Persist(ctx context.Context, adminUserID, tokenHash string, expiresAt time.Time) error
	Use(ctx context.Context, tokenHash string) (*ent.AdminUser, *ent.Role, error)
}
