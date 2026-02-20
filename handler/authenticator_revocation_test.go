package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/domain/svc"
	"github.com/GoLabra/labra/jwtrefresh"
	"github.com/GoLabra/labra/mocks"
	"github.com/GoLabra/labra/tokenrevocation"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
)

func TestAuthenticator_RejectsRevokedToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdminUser := mocks.NewMockAdminUser(ctrl)
	mockRole := mocks.NewMockRole(ctrl)

	secretKey := "test-secret-key-that-is-long-enough-for-validation-minimum-16-chars"
	adminEntClient := newTestAdminEntClient(t)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":          time.Now().Add(30 * time.Minute).Unix(),
		"iat":          time.Now().Add(-time.Minute).Unix(),
		"sub":          "revoked@example.com",
		"role":         "Admin",
		"typ":          jwtrefresh.TokenTypeAccess,
		"subject_type": jwtrefresh.SubjectTypeAdmin,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("failed signing token: %v", err)
	}

	if err := tokenrevocation.RevokeToken(
		context.Background(),
		adminEntClient,
		adminEntClient.DialectName(),
		tokenString,
		"revoked@example.com",
		jwtrefresh.SubjectTypeAdmin,
		time.Now().Add(30*time.Minute),
		time.Now(),
	); err != nil {
		t.Fatalf("failed revoking token: %v", err)
	}

	cfg := &config.AppConfig{
		Secrets: config.Secrets{SecretKey: secretKey},
	}
	service := &svc.Service{
		AdminUser: mockAdminUser,
		Role:      mockRole,
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
	ctx = context.WithValue(ctx, constants.AdminEntClientContextValue, adminEntClient)
	ctx = context.WithValue(ctx, "config", cfg)

	req := httptest.NewRequest(http.MethodGet, "/query", nil).WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	w := httptest.NewRecorder()
	nextCalled := false
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		nextCalled = true
	})

	Authenticator(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
	if nextCalled {
		t.Fatalf("next handler must not be called for revoked tokens")
	}
}
