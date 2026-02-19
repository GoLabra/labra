package jwtrefresh

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	ClaimSubject     = "sub"
	ClaimRole        = "role"
	ClaimType        = "typ"
	ClaimSubjectType = "subject_type"

	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"

	SubjectTypeAdmin = "admin"
	SubjectTypeUser  = "user"

	DefaultAccessTokenTTL  = 15 * time.Minute
	DefaultRefreshTokenTTL = 7 * 24 * time.Hour
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type TokenPair struct {
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
}

type RefreshClaims struct {
	Subject     string
	Role        string
	SubjectType string
}

type StoredRefreshToken struct {
	SubjectEmail  string
	SubjectType   string
	RoleName      string
	ExpiresAtUnix int64
	RevokedAtUnix sql.NullInt64
}

type SQLExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func EffectiveAccessTTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return DefaultAccessTokenTTL
	}
	return ttl
}

func EffectiveRefreshTTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return DefaultRefreshTokenTTL
	}
	return ttl
}

func IssueAccessToken(secret, subject, role, subjectType string, ttl time.Duration) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(EffectiveAccessTTL(ttl))

	claims := jwt.MapClaims{
		"exp":            expiresAt.Unix(),
		"iat":            now.Unix(),
		ClaimSubject:     subject,
		ClaimRole:        role,
		ClaimType:        TokenTypeAccess,
		ClaimSubjectType: subjectType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}

func IssueTokenPair(secret, subject, role, subjectType string, accessTTL, refreshTTL time.Duration) (TokenPair, error) {
	var pair TokenPair

	accessToken, accessExpiresAt, err := IssueAccessToken(secret, subject, role, subjectType, accessTTL)
	if err != nil {
		return pair, err
	}

	now := time.Now().UTC()
	refreshExpiresAt := now.Add(EffectiveRefreshTTL(refreshTTL))
	tokenID, err := newTokenID()
	if err != nil {
		return pair, err
	}

	refreshClaims := jwt.MapClaims{
		"exp":            refreshExpiresAt.Unix(),
		"iat":            now.Unix(),
		"jti":            tokenID,
		ClaimSubject:     subject,
		ClaimRole:        role,
		ClaimType:        TokenTypeRefresh,
		ClaimSubjectType: subjectType,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefresh, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return pair, err
	}

	pair.AccessToken = accessToken
	pair.AccessExpiresAt = accessExpiresAt
	pair.RefreshToken = signedRefresh
	pair.RefreshExpiresAt = refreshExpiresAt
	return pair, nil
}

func ParseRefreshToken(tokenString, secret, expectedSubjectType string) (RefreshClaims, error) {
	var out RefreshClaims

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token == nil || token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return out, err
	}
	if !token.Valid {
		return out, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return out, fmt.Errorf("invalid token claims")
	}

	tokenType, err := claimString(claims, ClaimType)
	if err != nil {
		return out, err
	}
	if tokenType != TokenTypeRefresh {
		return out, fmt.Errorf("token is not a refresh token")
	}

	subject, err := claimString(claims, ClaimSubject)
	if err != nil {
		return out, err
	}
	role, err := claimString(claims, ClaimRole)
	if err != nil {
		return out, err
	}
	subjectType, err := claimString(claims, ClaimSubjectType)
	if err != nil {
		return out, err
	}

	if expectedSubjectType != "" && subjectType != expectedSubjectType {
		return out, fmt.Errorf("invalid subject type")
	}

	out.Subject = subject
	out.Role = role
	out.SubjectType = subjectType
	return out, nil
}

func HashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}

func EnsureRefreshTokenTable(ctx context.Context, db SQLExecutor) error {
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

	_, err := db.ExecContext(ctx, createRefreshTokensTableQuery)
	return err
}

func SaveRefreshToken(ctx context.Context, db SQLExecutor, dialect, rawToken, subjectEmail, subjectType, roleName string, expiresAt time.Time) error {
	if err := EnsureRefreshTokenTable(ctx, db); err != nil {
		return err
	}

	insertQuery := fmt.Sprintf(
		"INSERT INTO refresh_tokens (token_hash, subject_email, subject_type, role_name, expires_at_unix, created_at_unix, revoked_at_unix) VALUES (%s, %s, %s, %s, %s, %s, NULL)",
		bindVar(dialect, 1),
		bindVar(dialect, 2),
		bindVar(dialect, 3),
		bindVar(dialect, 4),
		bindVar(dialect, 5),
		bindVar(dialect, 6),
	)

	_, err := db.ExecContext(
		ctx,
		insertQuery,
		HashToken(rawToken),
		subjectEmail,
		subjectType,
		roleName,
		expiresAt.Unix(),
		time.Now().UTC().Unix(),
	)
	return err
}

func LoadRefreshToken(ctx context.Context, db SQLExecutor, dialect, rawToken string) (StoredRefreshToken, error) {
	var out StoredRefreshToken

	if err := EnsureRefreshTokenTable(ctx, db); err != nil {
		return out, err
	}

	query := fmt.Sprintf(
		"SELECT subject_email, subject_type, role_name, expires_at_unix, revoked_at_unix FROM refresh_tokens WHERE token_hash = %s LIMIT 1",
		bindVar(dialect, 1),
	)

	rows, err := db.QueryContext(ctx, query, HashToken(rawToken))
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

func RevokeRefreshToken(ctx context.Context, db SQLExecutor, dialect, rawToken string, revokedAt time.Time) (bool, error) {
	if err := EnsureRefreshTokenTable(ctx, db); err != nil {
		return false, err
	}

	query := fmt.Sprintf(
		"UPDATE refresh_tokens SET revoked_at_unix = %s WHERE token_hash = %s AND revoked_at_unix IS NULL",
		bindVar(dialect, 1),
		bindVar(dialect, 2),
	)

	result, err := db.ExecContext(ctx, query, revokedAt.UTC().Unix(), HashToken(rawToken))
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func claimString(claims jwt.MapClaims, key string) (string, error) {
	raw, exists := claims[key]
	if !exists {
		return "", fmt.Errorf("missing %s claim", key)
	}

	s, ok := raw.(string)
	if !ok || strings.TrimSpace(s) == "" {
		return "", fmt.Errorf("invalid %s claim", key)
	}

	return strings.TrimSpace(s), nil
}

func bindVar(dialect string, position int) string {
	if strings.EqualFold(strings.TrimSpace(dialect), "postgres") {
		return fmt.Sprintf("$%d", position)
	}
	return "?"
}

func newTokenID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
