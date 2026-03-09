package jwtrefresh

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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

func newTokenID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
