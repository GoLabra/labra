package handler

import (
	"context"
	"testing"
	"time"

	"app/ent"
	"github.com/GoLabra/labra/config"
)

func TestIssueUserSession(t *testing.T) {
	originalPersist := persistUserRefreshSession
	defer func() {
		persistUserRefreshSession = originalPersist
	}()

	var (
		gotUserID    string
		gotTokenHash string
		gotExpiresAt time.Time
	)

	persistUserRefreshSession = func(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
		gotUserID = userID
		gotTokenHash = tokenHash
		gotExpiresAt = expiresAt
		return nil
	}

	cfg := &config.AppConfig{
		Config: config.Config{
			AccessTokenTTLMin:    20,
			RefreshTokenTTLHours: 168,
		},
		Secrets: config.Secrets{
			SecretKey: "test-secret-key-that-is-long-enough",
		},
	}

	user := &ent.User{ID: "user-1", Email: "user@example.com"}
	role := &ent.Role{Name: "Member"}

	accessToken, refreshToken, accessTTL, refreshTTL, err := issueUserSession(context.Background(), cfg, user, role)
	if err != nil {
		t.Fatalf("issueUserSession returned error: %v", err)
	}

	if accessToken == "" {
		t.Fatal("expected non-empty access token")
	}
	if refreshToken == "" {
		t.Fatal("expected non-empty refresh token")
	}
	if accessTTL != 20*time.Minute {
		t.Fatalf("expected access ttl %v, got %v", 20*time.Minute, accessTTL)
	}
	if refreshTTL != 168*time.Hour {
		t.Fatalf("expected refresh ttl %v, got %v", 168*time.Hour, refreshTTL)
	}
	if gotUserID != "user-1" {
		t.Fatalf("expected persisted user id %q, got %q", "user-1", gotUserID)
	}
	if gotTokenHash == "" {
		t.Fatal("expected persisted token hash")
	}
	if gotTokenHash == refreshToken {
		t.Fatal("expected stored refresh token to be hashed, not raw")
	}
	if time.Until(gotExpiresAt) <= 0 {
		t.Fatal("expected refresh expiry in the future")
	}
}
