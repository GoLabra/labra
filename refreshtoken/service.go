package refreshtoken

import (
	"context"
	"time"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Save(ctx context.Context, rawToken, subjectEmail, subjectType, roleName string, expiresAt time.Time) error {
	return s.repository.Save(ctx, rawToken, subjectEmail, subjectType, roleName, expiresAt)
}

func (s *Service) Load(ctx context.Context, rawToken string) (StoredRefreshToken, error) {
	return s.repository.Load(ctx, rawToken)
}

func (s *Service) Revoke(ctx context.Context, rawToken string, revokedAt time.Time) (bool, error) {
	return s.repository.Revoke(ctx, rawToken, revokedAt)
}
