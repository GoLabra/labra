package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app/domain/svc"
	"app/ent"
	"github.com/GoLabra/labra/config"
)

func TestRefresh(t *testing.T) {
	originalRotate := rotateUserRefreshSession
	defer func() {
		rotateUserRefreshSession = originalRotate
	}()

	t.Run("successful refresh", func(t *testing.T) {
		rotateUserRefreshSession = func(ctx context.Context, rawRefreshToken string, appConfig *config.AppConfig) (*ent.User, *ent.Role, string, string, time.Duration, time.Duration, error) {
			return nil, nil, "new-access-token", "new-refresh-token", 15 * time.Minute, 7 * 24 * time.Hour, nil
		}

		req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "old-refresh-token"})
		req = req.WithContext(context.WithValue(req.Context(), "config", &config.AppConfig{
			Config: config.Config{
				AccessTokenTTLMin:    15,
				RefreshTokenTTLHours: 168,
			},
			Secrets: config.Secrets{
				SecretKey: "test-secret-key-that-is-long-enough",
			},
		}))

		w := httptest.NewRecorder()
		Refresh(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		cookies := w.Result().Cookies()
		if len(cookies) < 2 {
			t.Fatalf("expected refreshed cookies, got %d", len(cookies))
		}
	})

	t.Run("missing refresh cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
		req = req.WithContext(context.WithValue(req.Context(), "config", &config.AppConfig{
			Secrets: config.Secrets{
				SecretKey: "test-secret-key-that-is-long-enough",
			},
		}))

		w := httptest.NewRecorder()
		Refresh(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		rotateUserRefreshSession = func(ctx context.Context, rawRefreshToken string, appConfig *config.AppConfig) (*ent.User, *ent.Role, string, string, time.Duration, time.Duration, error) {
			return nil, nil, "", "", 0, 0, svc.ErrInvalidRefreshToken
		}

		req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "used-refresh-token"})
		req = req.WithContext(context.WithValue(req.Context(), "config", &config.AppConfig{
			Secrets: config.Secrets{
				SecretKey: "test-secret-key-that-is-long-enough",
			},
		}))

		w := httptest.NewRecorder()
		Refresh(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}
