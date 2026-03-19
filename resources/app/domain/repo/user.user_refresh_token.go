package repo

import (
	"app/ent"
	"app/ent/userrefreshtoken"
	"context"
	"time"
)

type UserRefreshToken struct {
	client *ent.Client
}

func NewUserRefreshToken(c *ent.Client) *UserRefreshToken {
	return &UserRefreshToken{client: c}
}

func (r *UserRefreshToken) Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.client.UserRefreshToken.
		Create().
		SetTokenHash(tokenHash).
		SetExpiresAt(expiresAt).
		SetUserID(userID).
		Save(ctx)
	return err
}

func (r *UserRefreshToken) FindByTokenHash(ctx context.Context, tokenHash string) (*ent.UserRefreshToken, error) {
	return r.client.UserRefreshToken.
		Query().
		Where(userrefreshtoken.TokenHashEQ(tokenHash)).
		WithUser().
		Only(ctx)
}

func (r *UserRefreshToken) MarkUsedAndRevoked(ctx context.Context, id int, at time.Time) error {
	_, err := r.client.UserRefreshToken.
		UpdateOneID(id).
		SetUsedAt(at).
		SetRevokedAt(at).
		Save(ctx)
	return err
}
