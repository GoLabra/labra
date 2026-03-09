package svc

import (
	"context"
	"strings"
	"time"

	"github.com/GoLabra/labra/entgql/domain/repo"
)

type subjectRevocation struct {
	Email       string
	SubjectType string
}

func revokeSubjects(ctx context.Context, repository *repo.Repository, subjects []subjectRevocation) error {
	if len(subjects) == 0 {
		return nil
	}
	if repository == nil || repository.TokenRevocation == nil {
		return nil
	}

	seen := make(map[string]struct{}, len(subjects))
	revokedAt := time.Now().UTC()

	for _, subject := range subjects {
		email := normalizeRevocationEmail(subject.Email)
		subjectType := strings.TrimSpace(subject.SubjectType)
		if email == "" || subjectType == "" {
			continue
		}

		key := subjectType + "\x00" + email
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		if err := repository.TokenRevocation.RevokeSubjectTokens(ctx, email, subjectType, revokedAt); err != nil {
			return err
		}
	}

	return nil
}

func normalizeRevocationEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
