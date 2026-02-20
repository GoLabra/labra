package tokenrevocation

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

const (
	defaultCleanupInterval = time.Hour
)

type SQLExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func HashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}

func EnsureTokenRevocationTables(ctx context.Context, db SQLExecutor) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS token_blacklist (
	token_hash VARCHAR(64) PRIMARY KEY,
	subject_email VARCHAR(254) NOT NULL,
	subject_type VARCHAR(16) NOT NULL,
	expires_at_unix BIGINT NOT NULL,
	revoked_at_unix BIGINT NOT NULL,
	created_at_unix BIGINT NOT NULL
)`,
		`CREATE INDEX IF NOT EXISTS idx_token_blacklist_subject ON token_blacklist(subject_email, subject_type)`,
		`CREATE INDEX IF NOT EXISTS idx_token_blacklist_expires_at ON token_blacklist(expires_at_unix)`,
		`CREATE TABLE IF NOT EXISTS subject_token_revocations (
	subject_email VARCHAR(254) NOT NULL,
	subject_type VARCHAR(16) NOT NULL,
	revoked_after_unix BIGINT NOT NULL,
	updated_at_unix BIGINT NOT NULL,
	PRIMARY KEY (subject_email, subject_type)
)`,
	}

	for _, stmt := range statements {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}

	return nil
}

func RevokeToken(ctx context.Context, db SQLExecutor, dialect, rawToken, subjectEmail, subjectType string, expiresAt, revokedAt time.Time) error {
	if err := EnsureTokenRevocationTables(ctx, db); err != nil {
		return err
	}

	if strings.TrimSpace(rawToken) == "" {
		return fmt.Errorf("token is required")
	}

	revokedAtUnix := revokedAt.UTC().Unix()
	createdAtUnix := time.Now().UTC().Unix()
	expiresAtUnix := expiresAt.UTC().Unix()
	if expiresAtUnix <= 0 {
		expiresAtUnix = revokedAt.UTC().Add(24 * time.Hour).Unix()
	}

	query := tokenUpsertQuery(dialect)
	_, err := db.ExecContext(
		ctx,
		query,
		HashToken(rawToken),
		subjectEmail,
		subjectType,
		expiresAtUnix,
		revokedAtUnix,
		createdAtUnix,
	)
	return err
}

func RevokeSubjectTokens(ctx context.Context, db SQLExecutor, dialect, subjectEmail, subjectType string, revokedAt time.Time) error {
	if err := EnsureTokenRevocationTables(ctx, db); err != nil {
		return err
	}
	if strings.TrimSpace(subjectEmail) == "" {
		return fmt.Errorf("subject email is required")
	}
	if strings.TrimSpace(subjectType) == "" {
		return fmt.Errorf("subject type is required")
	}

	// JWT iat is second-precision in this codebase. Bumping the cutoff by one second
	// guarantees that tokens issued in the revocation second are invalidated too.
	revokedAfterUnix := revokedAt.UTC().Unix() + 1
	updatedAtUnix := time.Now().UTC().Unix()

	query := subjectRevocationUpsertQuery(dialect)
	_, err := db.ExecContext(
		ctx,
		query,
		subjectEmail,
		subjectType,
		revokedAfterUnix,
		updatedAtUnix,
	)
	return err
}

func IsTokenRevoked(ctx context.Context, db SQLExecutor, dialect, rawToken string) (bool, error) {
	if err := EnsureTokenRevocationTables(ctx, db); err != nil {
		return false, err
	}
	if strings.TrimSpace(rawToken) == "" {
		return false, nil
	}

	query := fmt.Sprintf(
		"SELECT 1 FROM token_blacklist WHERE token_hash = %s LIMIT 1",
		bindVar(dialect, 1),
	)

	rows, err := db.QueryContext(ctx, query, HashToken(rawToken))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), rows.Err()
}

func IsSubjectTokenRevoked(ctx context.Context, db SQLExecutor, dialect, subjectEmail, subjectType string, issuedAt time.Time) (bool, error) {
	if err := EnsureTokenRevocationTables(ctx, db); err != nil {
		return false, err
	}
	if strings.TrimSpace(subjectEmail) == "" || strings.TrimSpace(subjectType) == "" || issuedAt.IsZero() {
		return false, nil
	}

	query := fmt.Sprintf(
		"SELECT revoked_after_unix FROM subject_token_revocations WHERE subject_email = %s AND subject_type = %s LIMIT 1",
		bindVar(dialect, 1),
		bindVar(dialect, 2),
	)

	rows, err := db.QueryContext(ctx, query, subjectEmail, subjectType)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if !rows.Next() {
		return false, rows.Err()
	}

	var revokedAfterUnix int64
	if err := rows.Scan(&revokedAfterUnix); err != nil {
		return false, err
	}
	if err := rows.Err(); err != nil {
		return false, err
	}

	// Use strict comparison to avoid invalidating a token issued in the same second.
	return issuedAt.UTC().Unix() < revokedAfterUnix, nil
}

func CleanupExpiredRevocations(ctx context.Context, db SQLExecutor, dialect string, now time.Time) (int64, error) {
	if err := EnsureTokenRevocationTables(ctx, db); err != nil {
		return 0, err
	}

	query := fmt.Sprintf(
		"DELETE FROM token_blacklist WHERE expires_at_unix <= %s",
		bindVar(dialect, 1),
	)

	result, err := db.ExecContext(ctx, query, now.UTC().Unix())
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func StartCleanupJob(ctx context.Context, db SQLExecutor, dialect string, interval time.Duration, logf func(string, ...any)) {
	if interval <= 0 {
		interval = defaultCleanupInterval
	}

	go func() {
		// Initial cleanup to handle stale data after deploy/restart.
		if _, err := CleanupExpiredRevocations(ctx, db, dialect, time.Now().UTC()); err != nil && logf != nil {
			logf("token revocation cleanup failed: %v", err)
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := CleanupExpiredRevocations(ctx, db, dialect, time.Now().UTC()); err != nil && logf != nil {
					logf("token revocation cleanup failed: %v", err)
				}
			}
		}
	}()
}

func tokenUpsertQuery(dialect string) string {
	switch normalizeDialect(dialect) {
	case "postgres":
		return `
INSERT INTO token_blacklist (token_hash, subject_email, subject_type, expires_at_unix, revoked_at_unix, created_at_unix)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (token_hash) DO UPDATE
SET subject_email = EXCLUDED.subject_email,
	subject_type = EXCLUDED.subject_type,
	expires_at_unix = EXCLUDED.expires_at_unix,
	revoked_at_unix = EXCLUDED.revoked_at_unix`
	case "mysql":
		return `
INSERT INTO token_blacklist (token_hash, subject_email, subject_type, expires_at_unix, revoked_at_unix, created_at_unix)
VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	subject_email = VALUES(subject_email),
	subject_type = VALUES(subject_type),
	expires_at_unix = VALUES(expires_at_unix),
	revoked_at_unix = VALUES(revoked_at_unix)`
	default:
		return `
INSERT INTO token_blacklist (token_hash, subject_email, subject_type, expires_at_unix, revoked_at_unix, created_at_unix)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(token_hash) DO UPDATE
SET subject_email = excluded.subject_email,
	subject_type = excluded.subject_type,
	expires_at_unix = excluded.expires_at_unix,
	revoked_at_unix = excluded.revoked_at_unix`
	}
}

func subjectRevocationUpsertQuery(dialect string) string {
	switch normalizeDialect(dialect) {
	case "postgres":
		return `
INSERT INTO subject_token_revocations (subject_email, subject_type, revoked_after_unix, updated_at_unix)
VALUES ($1, $2, $3, $4)
ON CONFLICT (subject_email, subject_type) DO UPDATE
SET revoked_after_unix = CASE
	WHEN subject_token_revocations.revoked_after_unix > EXCLUDED.revoked_after_unix
	THEN subject_token_revocations.revoked_after_unix
	ELSE EXCLUDED.revoked_after_unix
END,
updated_at_unix = EXCLUDED.updated_at_unix`
	case "mysql":
		return `
INSERT INTO subject_token_revocations (subject_email, subject_type, revoked_after_unix, updated_at_unix)
VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	revoked_after_unix = GREATEST(revoked_after_unix, VALUES(revoked_after_unix)),
	updated_at_unix = VALUES(updated_at_unix)`
	default:
		return `
INSERT INTO subject_token_revocations (subject_email, subject_type, revoked_after_unix, updated_at_unix)
VALUES (?, ?, ?, ?)
ON CONFLICT(subject_email, subject_type) DO UPDATE
SET revoked_after_unix = CASE
	WHEN subject_token_revocations.revoked_after_unix > excluded.revoked_after_unix
	THEN subject_token_revocations.revoked_after_unix
	ELSE excluded.revoked_after_unix
END,
updated_at_unix = excluded.updated_at_unix`
	}
}

func normalizeDialect(dialect string) string {
	return strings.ToLower(strings.TrimSpace(dialect))
}

func bindVar(dialect string, position int) string {
	if normalizeDialect(dialect) == "postgres" {
		return fmt.Sprintf("$%d", position)
	}
	return "?"
}
