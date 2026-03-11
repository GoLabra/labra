package jwtrefresh

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	ClaimSubject     = "sub"
	ClaimRole        = "role"
	ClaimType        = "typ"
	ClaimSubjectType = "subject_type"

	TokenTypeAccess = "access"

	SubjectTypeAdmin = "admin"
	SubjectTypeUser  = "user"

	DefaultAccessTokenTTL = 15 * time.Minute
)

func EffectiveAccessTTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return DefaultAccessTokenTTL
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
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}

	return signed, expiresAt, nil
}
