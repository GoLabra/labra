package svc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/GoLabra/labra/entgql/domain/repo"
	"github.com/GoLabra/labra/tokenrevocation"
)

type subjectRevocation struct {
	Email       string
	SubjectType string
}

func revokeSubjects(ctx context.Context, repository *repo.Repository, subjects []subjectRevocation) error {
	if len(subjects) == 0 {
		return nil
	}
	if repository == nil || repository.Client == nil {
		return fmt.Errorf("repository client is not initialized")
	}

	dialect := repository.Client.DialectName()
	if strings.TrimSpace(dialect) == "" {
		return fmt.Errorf("database dialect is not available")
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

		if err := tokenrevocation.RevokeSubjectTokens(ctx, repository.Client, dialect, email, subjectType, revokedAt); err != nil {
			return err
		}
	}

	return nil
}

func normalizeRevocationEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
