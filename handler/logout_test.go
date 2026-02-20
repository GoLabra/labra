package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/jwtrefresh"
	"github.com/GoLabra/labra/tokenrevocation"
	"github.com/golang-jwt/jwt"
)

func TestLogout_RevokesCurrentToken(t *testing.T) {
	client := newTestAdminEntClient(t)
	secretKey := "test-secret-key-that-is-long-enough-for-validation-minimum-16-chars"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":          time.Now().Add(30 * time.Minute).Unix(),
		"iat":          time.Now().Add(-time.Minute).Unix(),
		"sub":          "logout@example.com",
		"role":         "Admin",
		"typ":          jwtrefresh.TokenTypeAccess,
		"subject_type": jwtrefresh.SubjectTypeAdmin,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("failed signing token: %v", err)
	}

	cfg := &config.AppConfig{
		Secrets: config.Secrets{SecretKey: secretKey},
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", cfg)
	ctx = context.WithValue(ctx, constants.AdminEntClientContextValue, client)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil).WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	Logout(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	revoked, err := tokenrevocation.IsTokenRevoked(context.Background(), client, client.DialectName(), tokenString)
	if err != nil {
		t.Fatalf("failed checking revoked token: %v", err)
	}
	if !revoked {
		t.Fatalf("expected token to be revoked")
	}

	setCookie := strings.Join(w.Header().Values("Set-Cookie"), "\n")
	if !strings.Contains(setCookie, "jwt=") {
		t.Fatalf("expected jwt cookie to be cleared")
	}
	if !strings.Contains(setCookie, "csrf_token=") {
		t.Fatalf("expected csrf cookie to be cleared")
	}
}

func TestLogout_MissingToken(t *testing.T) {
	client := newTestAdminEntClient(t)
	cfg := &config.AppConfig{
		Secrets: config.Secrets{
			SecretKey: "test-secret-key-that-is-long-enough-for-validation-minimum-16-chars",
		},
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", cfg)
	ctx = context.WithValue(ctx, constants.AdminEntClientContextValue, client)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	Logout(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
