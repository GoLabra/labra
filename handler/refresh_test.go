package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/entgql/domain/svc"
	"github.com/GoLabra/labra/entgql/ent"
)

func TestRefresh(t *testing.T) {
	originalRotate := rotateAdminRefreshSession
	defer func() {
		rotateAdminRefreshSession = originalRotate
	}()

	t.Run("successful refresh", func(t *testing.T) {
		rotateAdminRefreshSession = func(ctx context.Context, rawRefreshToken string, appConfig *config.AppConfig) (*ent.AdminUser, *ent.Role, string, string, time.Duration, time.Duration, error) {
			return nil, nil, "new-access-token", "new-refresh-token", 15 * time.Minute, 7 * 24 * time.Hour, nil
		}

		req := httptest.NewRequest(http.MethodPost, "/admin/refresh", nil)
		req.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: "old-refresh-token"})
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
		req := httptest.NewRequest(http.MethodPost, "/admin/refresh", nil)
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
		rotateAdminRefreshSession = func(ctx context.Context, rawRefreshToken string, appConfig *config.AppConfig) (*ent.AdminUser, *ent.Role, string, string, time.Duration, time.Duration, error) {
			return nil, nil, "", "", 0, 0, svc.ErrInvalidRefreshToken
		}

		req := httptest.NewRequest(http.MethodPost, "/admin/refresh", nil)
		req.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: "used-refresh-token"})
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
