package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/jwtrefresh"
)

func TestRefresh_RotatesRefreshToken(t *testing.T) {
	client := newTestAdminEntClient(t)

	cfg := &config.AppConfig{
		Config: config.Config{
			DBDialect:       "sqlite",
			AccessTokenTTL:  20 * time.Minute,
			RefreshTokenTTL: 7 * 24 * time.Hour,
		},
		Secrets: config.Secrets{
			SecretKey: "test-secret-key-that-is-long-enough-for-validation-minimum-16-chars",
		},
	}

	pair, err := jwtrefresh.IssueTokenPair(
		cfg.SecretKey,
		"refresh-test@example.com",
		"Admin",
		jwtrefresh.SubjectTypeAdmin,
		cfg.AccessTokenTTL,
		cfg.RefreshTokenTTL,
	)
	if err != nil {
		t.Fatalf("failed issuing initial pair: %v", err)
	}

	err = jwtrefresh.SaveRefreshToken(
		context.Background(),
		client,
		cfg.DBDialect,
		pair.RefreshToken,
		"refresh-test@example.com",
		jwtrefresh.SubjectTypeAdmin,
		"Admin",
		pair.RefreshExpiresAt,
	)
	if err != nil {
		t.Fatalf("failed storing initial refresh token: %v", err)
	}

	requestBody, _ := json.Marshal(map[string]string{
		"refresh_token": pair.RefreshToken,
	})

	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", cfg)
	ctx = context.WithValue(ctx, constants.AdminEntClientContextValue, client)

	req := httptest.NewRequest(http.MethodPost, "/admin/refresh", bytes.NewBuffer(requestBody)).WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Refresh(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal refresh response: %v", err)
	}
	if response["token"] == "" {
		t.Fatal("expected token in refresh response")
	}
	if response["refresh_token"] == "" {
		t.Fatal("expected refresh_token in refresh response")
	}
	if response["refresh_token"] == pair.RefreshToken {
		t.Fatal("expected rotated refresh token to differ from original")
	}

	oldTokenBody, _ := json.Marshal(map[string]string{
		"refresh_token": pair.RefreshToken,
	})
	oldTokenReq := httptest.NewRequest(http.MethodPost, "/admin/refresh", bytes.NewBuffer(oldTokenBody)).WithContext(ctx)
	oldTokenReq.Header.Set("Content-Type", "application/json")
	oldTokenW := httptest.NewRecorder()

	Refresh(oldTokenW, oldTokenReq)

	if oldTokenW.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d for reused refresh token, got %d", http.StatusUnauthorized, oldTokenW.Code)
	}
}
