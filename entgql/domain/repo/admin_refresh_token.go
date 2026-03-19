package repo

import (
	"context"
	"time"

	"github.com/GoLabra/labra/entgql/ent"
	"github.com/GoLabra/labra/entgql/ent/adminrefreshtoken"
)

type AdminRefreshToken struct {
	client *ent.Client
}

func NewAdminRefreshToken(c *ent.Client) *AdminRefreshToken {
	return &AdminRefreshToken{client: c}
}

func (r *AdminRefreshToken) Create(ctx context.Context, adminUserID, tokenHash string, expiresAt time.Time) error {
	_, err := r.client.AdminRefreshToken.
		Create().
		SetTokenHash(tokenHash).
		SetExpiresAt(expiresAt).
		SetAdminUserID(adminUserID).
		Save(ctx)
	return err
}

func (r *AdminRefreshToken) FindByTokenHash(ctx context.Context, tokenHash string) (*ent.AdminRefreshToken, error) {
	return r.client.AdminRefreshToken.
		Query().
		Where(adminrefreshtoken.TokenHashEQ(tokenHash)).
		WithAdminUser().
		Only(ctx)
}

func (r *AdminRefreshToken) MarkUsedAndRevoked(ctx context.Context, id string, at time.Time) error {
	_, err := r.client.AdminRefreshToken.
		UpdateOneID(id).
		SetUsedAt(at).
		SetRevokedAt(at).
		Save(ctx)
	return err
}
