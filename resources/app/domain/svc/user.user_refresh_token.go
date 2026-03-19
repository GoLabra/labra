package svc

import (
	"app/domain/repo"
	"app/ent"
	"context"
	"errors"
	"time"

	"github.com/GoLabra/labra/constants"
)

var ErrInvalidRefreshToken = errors.New("invalid refresh token")

type UserRefreshToken struct {
	repository *repo.Repository
}

func NewUserRefreshToken(r *repo.Repository) *UserRefreshToken {
	return &UserRefreshToken{repository: r}
}

func (s *UserRefreshToken) Persist(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	return s.repository.UserRefreshToken.Create(constants.WithInternalOperation(ctx), userID, tokenHash, expiresAt)
}

func (s *UserRefreshToken) Use(ctx context.Context, tokenHash string) (*ent.User, *ent.Role, error) {
	iCtx := constants.WithInternalOperation(ctx)

	storedToken, err := s.repository.UserRefreshToken.FindByTokenHash(iCtx, tokenHash)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil, ErrInvalidRefreshToken
		}
		return nil, nil, err
	}

	now := time.Now()
	if storedToken.RevokedAt != nil || storedToken.UsedAt != nil || now.After(storedToken.ExpiresAt) {
		return nil, nil, ErrInvalidRefreshToken
	}

	user, err := storedToken.Edges.UserOrErr()
	if err != nil {
		return nil, nil, err
	}

	role, err := user.DefaultRole(iCtx)
	if err != nil {
		return nil, nil, err
	}
	if role == nil {
		return nil, nil, ErrInvalidRefreshToken
	}

	if err := s.repository.UserRefreshToken.MarkUsedAndRevoked(iCtx, storedToken.ID, now); err != nil {
		return nil, nil, err
	}

	return user, role, nil
}
