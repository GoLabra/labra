package refreshtoken

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/GoLabra/labra/jwtrefresh"
)

type SQLExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

type Repository struct {
	db      SQLExecutor
	dialect string
}

func NewRepository(db SQLExecutor, dialect string) *Repository {
	return &Repository{
		db:      db,
		dialect: dialect,
	}
}

func (r *Repository) EnsureTable(ctx context.Context) error {
	const createRefreshTokensTableQuery = `
CREATE TABLE IF NOT EXISTS refresh_tokens (
	token_hash VARCHAR(64) PRIMARY KEY,
	subject_email VARCHAR(254) NOT NULL,
	subject_type VARCHAR(16) NOT NULL,
	role_name VARCHAR(128) NOT NULL,
	expires_at_unix BIGINT NOT NULL,
	created_at_unix BIGINT NOT NULL,
	revoked_at_unix BIGINT NULL
)`

	_, err := r.db.ExecContext(ctx, createRefreshTokensTableQuery)
	return err
}

func (r *Repository) Save(ctx context.Context, rawToken, subjectEmail, subjectType, roleName string, expiresAt time.Time) error {
	if err := r.EnsureTable(ctx); err != nil {
		return err
	}

	insertQuery := fmt.Sprintf(
		"INSERT INTO refresh_tokens (token_hash, subject_email, subject_type, role_name, expires_at_unix, created_at_unix, revoked_at_unix) VALUES (%s, %s, %s, %s, %s, %s, NULL)",
		bindVar(r.dialect, 1),
		bindVar(r.dialect, 2),
		bindVar(r.dialect, 3),
		bindVar(r.dialect, 4),
		bindVar(r.dialect, 5),
		bindVar(r.dialect, 6),
	)

	_, err := r.db.ExecContext(
		ctx,
		insertQuery,
		jwtrefresh.HashToken(rawToken),
		subjectEmail,
		subjectType,
		roleName,
		expiresAt.Unix(),
		time.Now().UTC().Unix(),
	)
	return err
}

func (r *Repository) Load(ctx context.Context, rawToken string) (StoredRefreshToken, error) {
	var out StoredRefreshToken

	if err := r.EnsureTable(ctx); err != nil {
		return out, err
	}

	query := fmt.Sprintf(
		"SELECT subject_email, subject_type, role_name, expires_at_unix, revoked_at_unix FROM refresh_tokens WHERE token_hash = %s LIMIT 1",
		bindVar(r.dialect, 1),
	)

	rows, err := r.db.QueryContext(ctx, query, jwtrefresh.HashToken(rawToken))
	if err != nil {
		return out, err
	}
	defer rows.Close()

	if !rows.Next() {
		return out, ErrRefreshTokenNotFound
	}

	if err := rows.Scan(
		&out.SubjectEmail,
		&out.SubjectType,
		&out.RoleName,
		&out.ExpiresAtUnix,
		&out.RevokedAtUnix,
	); err != nil {
		return out, err
	}

	if err := rows.Err(); err != nil {
		return out, err
	}

	return out, nil
}

func (r *Repository) Revoke(ctx context.Context, rawToken string, revokedAt time.Time) (bool, error) {
	if err := r.EnsureTable(ctx); err != nil {
		return false, err
	}

	query := fmt.Sprintf(
		"UPDATE refresh_tokens SET revoked_at_unix = %s WHERE token_hash = %s AND revoked_at_unix IS NULL",
		bindVar(r.dialect, 1),
		bindVar(r.dialect, 2),
	)

	result, err := r.db.ExecContext(ctx, query, revokedAt.UTC().Unix(), jwtrefresh.HashToken(rawToken))
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func bindVar(dialect string, position int) string {
	if strings.EqualFold(strings.TrimSpace(dialect), "postgres") {
		return fmt.Sprintf("$%d", position)
	}
	return "?"
}
