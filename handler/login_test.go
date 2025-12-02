package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/domain/svc"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/GoLabra/labra/mocks"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    LoginFormData
		setupMocks     func(*mocks.MockAdminUser)
		setupContext   func(*mocks.MockAdminUser) context.Context
		expectedStatus int
		validateToken  bool
	}{
		{
			name: "successful login",
			requestBody: LoginFormData{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := CreateTestAdminUser("user-1", "test@example.com", string(hashedPassword))
				role := CreateTestRole("role-1", "Admin")
				testUser := CreateTestAdminUserWithDefaultRole(user, role)

				mockAdminUser.EXPECT().
					GetOne(gomock.Any(), ent.AdminUserWhereUniqueInput{Email: stringPtr("test@example.com")}).
					Return(testUser, nil)
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser) context.Context {
				cfg := &config.Config{
					SecretKey: "test-secret-key-that-is-long-enough-for-validation-minimum-16-chars",
				}
				service := &svc.Service{
					AdminUser: mockAdminUser,
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusOK,
			validateToken:  true,
		},
		{
			name: "user not found",
			requestBody: LoginFormData{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser) {
				mockAdminUser.EXPECT().
					GetOne(gomock.Any(), ent.AdminUserWhereUniqueInput{Email: stringPtr("nonexistent@example.com")}).
					Return(nil, &ent.NotFoundError{})
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser) context.Context {
				cfg := &config.Config{SecretKey: "test-secret-key-that-is-long-enough"}
				service := &svc.Service{AdminUser: mockAdminUser}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusUnauthorized,
			validateToken:  false,
		},
		{
			name: "invalid password",
			requestBody: LoginFormData{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := CreateTestAdminUser("user-1", "test@example.com", string(hashedPassword))
				role := CreateTestRole("role-1", "Admin")
				testUser := CreateTestAdminUserWithDefaultRole(user, role)

				mockAdminUser.EXPECT().
					GetOne(gomock.Any(), ent.AdminUserWhereUniqueInput{Email: stringPtr("test@example.com")}).
					Return(testUser, nil)
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser) context.Context {
				cfg := &config.Config{SecretKey: "test-secret-key-that-is-long-enough"}
				service := &svc.Service{AdminUser: mockAdminUser}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusUnauthorized,
			validateToken:  false,
		},
		{
			name:        "invalid json body",
			requestBody: LoginFormData{},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser) {
				// No mocks needed
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser) context.Context {
				cfg := &config.Config{SecretKey: "test-secret-key-that-is-long-enough"}
				service := &svc.Service{AdminUser: mockAdminUser}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusBadRequest,
			validateToken:  false,
		},
		{
			name: "service not in context",
			requestBody: LoginFormData{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser) {
				// No mocks needed
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser) context.Context {
				// Don't add service to context
				cfg := &config.Config{SecretKey: "test-secret-key-that-is-long-enough"}
				ctx := context.Background()
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusInternalServerError,
			validateToken:  false,
		},
		{
			name: "config not in context",
			requestBody: LoginFormData{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := CreateTestAdminUser("user-1", "test@example.com", string(hashedPassword))
				role := CreateTestRole("role-1", "Admin")
				testUser := CreateTestAdminUserWithDefaultRole(user, role)

				mockAdminUser.EXPECT().
					GetOne(gomock.Any(), ent.AdminUserWhereUniqueInput{Email: stringPtr("test@example.com")}).
					Return(testUser, nil)
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser) context.Context {
				service := &svc.Service{AdminUser: mockAdminUser}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				// Don't add config to context
				return ctx
			},
			expectedStatus: http.StatusInternalServerError,
			validateToken:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAdminUser := mocks.NewMockAdminUser(ctrl)
			tt.setupMocks(mockAdminUser)

			// Create request body
			var body []byte
			if tt.name == "invalid json body" {
				body = []byte("invalid json")
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			ctx := tt.setupContext(mockAdminUser)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body)).WithContext(ctx)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			Login(w, req)

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

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
