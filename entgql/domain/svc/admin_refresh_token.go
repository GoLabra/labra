package svc

import (
	"context"
	"errors"
	"time"

	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/domain/repo"
	"github.com/GoLabra/labra/entgql/ent"
)

var ErrInvalidRefreshToken = errors.New("invalid refresh token")

type AdminRefreshToken struct {
	repository *repo.Repository
}

func NewAdminRefreshToken(r *repo.Repository) *AdminRefreshToken {
	return &AdminRefreshToken{repository: r}
}

func (s *AdminRefreshToken) Persist(ctx context.Context, adminUserID, tokenHash string, expiresAt time.Time) error {
	return s.repository.AdminRefreshToken.Create(constants.WithInternalOperation(ctx), adminUserID, tokenHash, expiresAt)
}

func (s *AdminRefreshToken) Use(ctx context.Context, tokenHash string) (*ent.AdminUser, *ent.Role, error) {
	iCtx := constants.WithInternalOperation(ctx)

	storedToken, err := s.repository.AdminRefreshToken.FindByTokenHash(iCtx, tokenHash)
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

	adminUser, err := storedToken.Edges.AdminUserOrErr()
	if err != nil {
		return nil, nil, err
	}

	role, err := adminUser.DefaultRole(iCtx)
	if err != nil {
		return nil, nil, err
	}
	if role == nil {
		return nil, nil, ErrInvalidRefreshToken
	}

	if err := s.repository.AdminRefreshToken.MarkUsedAndRevoked(iCtx, storedToken.ID, now); err != nil {
		return nil, nil, err
	}

	return adminUser, role, nil
}
