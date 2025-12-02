package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
)

func TestChangeSessionRole(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func() context.Context
		expectedStatus int
		validateToken  bool
	}{
		{
			name: "successful role change",
			setupContext: func() context.Context {
				cfg := &config.Config{
					SecretKey: "test-secret-key-that-is-long-enough-for-validation-minimum-16-chars",
				}
				user := CreateTestAdminUser("user-1", "test@example.com", "hashed-password")
				role := CreateTestRole("role-1", "Admin")

				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.UserContextValue, user)
				ctx = context.WithValue(ctx, constants.RoleContextValue, role)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusOK,
			validateToken:  true,
		},
		{
			name: "user not in context",
			setupContext: func() context.Context {
				cfg := &config.Config{SecretKey: "test-secret-key-that-is-long-enough"}
				role := CreateTestRole("role-1", "Admin")

				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.RoleContextValue, role)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusInternalServerError,
			validateToken:  false,
		},
		{
			name: "role not in context",
			setupContext: func() context.Context {
				cfg := &config.Config{SecretKey: "test-secret-key-that-is-long-enough"}
				user := CreateTestAdminUser("user-1", "test@example.com", "hashed-password")

				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.UserContextValue, user)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusInternalServerError,
			validateToken:  false,
		},
		{
			name: "config not in context",
			setupContext: func() context.Context {
				user := CreateTestAdminUser("user-1", "test@example.com", "hashed-password")
				role := CreateTestRole("role-1", "Admin")

				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.UserContextValue, user)
				ctx = context.WithValue(ctx, constants.RoleContextValue, role)
				// Don't add config to context
				return ctx
			},
			expectedStatus: http.StatusInternalServerError,
			validateToken:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			req := httptest.NewRequest(http.MethodPost, "/change-session-role", nil).WithContext(ctx)
			w := httptest.NewRecorder()

			ChangeSessionRole(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.validateToken && w.Code == http.StatusOK {
				var response map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}
				if response["token"] == "" {
					t.Error("expected token in response")
				}
			}
		})
	}
}
