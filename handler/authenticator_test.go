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
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/GoLabra/labra/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
)

func TestAuthenticator(t *testing.T) {
	secretKey := "test-secret-key-that-is-long-enough-for-validation-minimum-16-chars"

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		setupMocks     func(*mocks.MockAdminUser, *mocks.MockRole)
		setupContext   func(*mocks.MockAdminUser, *mocks.MockRole) context.Context
		expectedStatus int
		nextCalled     bool
	}{
		{
			name: "successful authentication with bearer token",
			setupRequest: func() *http.Request {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"exp":  time.Now().Add(24 * time.Hour).Unix(),
					"sub":  "test@example.com",
					"role": "Admin",
				})
				tokenString, _ := token.SignedString([]byte(secretKey))

				req := httptest.NewRequest(http.MethodGet, "/query", nil)
				req.Header.Set("Authorization", "Bearer "+tokenString)
				return req
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				user := CreateTestAdminUser("user-1", "test@example.com", "hashed-password")
				role := CreateTestRole("role-1", "Admin")

				mockAdminUser.EXPECT().
					GetOne(gomock.Any(), ent.AdminUserWhereUniqueInput{Email: stringPtr("test@example.com")}).
					Return(user, nil)

				mockRole.EXPECT().
					GetOne(gomock.Any(), ent.RoleWhereUniqueInput{Name: stringPtr("Admin")}).
					Return(role, nil)
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{
					Secrets: config.Secrets{SecretKey: secretKey},
				}
				service := &svc.Service{
					AdminUser: mockAdminUser,
					Role:      mockRole,
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusOK,
			nextCalled:     true,
		},
		{
			name: "successful authentication with cookie",
			setupRequest: func() *http.Request {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"exp":  time.Now().Add(24 * time.Hour).Unix(),
					"sub":  "test@example.com",
					"role": "Admin",
				})
				tokenString, _ := token.SignedString([]byte(secretKey))

				req := httptest.NewRequest(http.MethodGet, "/query", nil)
				req.AddCookie(&http.Cookie{
					Name:  "jwt",
					Value: tokenString,
				})
				return req
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				user := CreateTestAdminUser("user-1", "test@example.com", "hashed-password")
				role := CreateTestRole("role-1", "Admin")

				mockAdminUser.EXPECT().
					GetOne(gomock.Any(), ent.AdminUserWhereUniqueInput{Email: stringPtr("test@example.com")}).
					Return(user, nil)

				mockRole.EXPECT().
					GetOne(gomock.Any(), ent.RoleWhereUniqueInput{Name: stringPtr("Admin")}).
					Return(role, nil)
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{
					Secrets: config.Secrets{SecretKey: secretKey},
				}
				service := &svc.Service{
					AdminUser: mockAdminUser,
					Role:      mockRole,
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusOK,
			nextCalled:     true,
		},
		{
			name: "skip authentication for websocket upgrade",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/query", nil)
				req.Header.Set("Connection", "Upgrade")
				req.Header.Set("Upgrade", "websocket")
				return req
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				// No mocks needed - should skip authentication
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{
					Secrets: config.Secrets{SecretKey: secretKey},
				}
				service := &svc.Service{
					AdminUser: mockAdminUser,
					Role:      mockRole,
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusOK,
			nextCalled:     true,
		},
		{
			name: "no token provided",
			setupRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/query", nil)
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				// No mocks needed
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{
					Secrets: config.Secrets{SecretKey: secretKey},
				}
				service := &svc.Service{
					AdminUser: mockAdminUser,
					Role:      mockRole,
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusUnauthorized,
			nextCalled:     false,
		},
		{
			name: "invalid token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/query", nil)
				req.Header.Set("Authorization", "Bearer invalid-token")
				return req
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				// No mocks needed - token parsing will fail
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{
					Secrets: config.Secrets{SecretKey: secretKey},
				}
				service := &svc.Service{
					AdminUser: mockAdminUser,
					Role:      mockRole,
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusUnauthorized,
			nextCalled:     false,
		},
		{
			name: "expired token",
			setupRequest: func() *http.Request {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"exp":  time.Now().Add(-1 * time.Hour).Unix(), // Expired
					"sub":  "test@example.com",
					"role": "Admin",
				})
				tokenString, _ := token.SignedString([]byte(secretKey))

				req := httptest.NewRequest(http.MethodGet, "/query", nil)
				req.Header.Set("Authorization", "Bearer "+tokenString)
				return req
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				// No mocks needed - token validation will fail
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{
					Secrets: config.Secrets{SecretKey: secretKey},
				}
				service := &svc.Service{
					AdminUser: mockAdminUser,
					Role:      mockRole,
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusUnauthorized,
			nextCalled:     false,
		},
		{
			name: "user not found",
			setupRequest: func() *http.Request {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"exp":  time.Now().Add(24 * time.Hour).Unix(),
					"sub":  "nonexistent@example.com",
					"role": "Admin",
				})
				tokenString, _ := token.SignedString([]byte(secretKey))

				req := httptest.NewRequest(http.MethodGet, "/query", nil)
				req.Header.Set("Authorization", "Bearer "+tokenString)
				return req
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				mockAdminUser.EXPECT().
					GetOne(gomock.Any(), ent.AdminUserWhereUniqueInput{Email: stringPtr("nonexistent@example.com")}).
					Return(nil, &ent.NotFoundError{})
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{
					Secrets: config.Secrets{SecretKey: secretKey},
				}
				service := &svc.Service{
					AdminUser: mockAdminUser,
					Role:      mockRole,
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusUnauthorized,
			nextCalled:     false,
		},
		{
			name: "role not found",
			setupRequest: func() *http.Request {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"exp":  time.Now().Add(24 * time.Hour).Unix(),
					"sub":  "test@example.com",
					"role": "NonExistentRole",
				})
				tokenString, _ := token.SignedString([]byte(secretKey))

				req := httptest.NewRequest(http.MethodGet, "/query", nil)
				req.Header.Set("Authorization", "Bearer "+tokenString)
				return req
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				user := CreateTestAdminUser("user-1", "test@example.com", "hashed-password")

				mockAdminUser.EXPECT().
					GetOne(gomock.Any(), ent.AdminUserWhereUniqueInput{Email: stringPtr("test@example.com")}).
					Return(user, nil)

				mockRole.EXPECT().
					GetOne(gomock.Any(), ent.RoleWhereUniqueInput{Name: stringPtr("NonExistentRole")}).
					Return(nil, &ent.NotFoundError{})
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{
					Secrets: config.Secrets{SecretKey: secretKey},
				}
				service := &svc.Service{
					AdminUser: mockAdminUser,
					Role:      mockRole,
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusUnauthorized,
			nextCalled:     false,
		},
		{
			name: "service not in context",
			setupRequest: func() *http.Request {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"exp":  time.Now().Add(24 * time.Hour).Unix(),
					"sub":  "test@example.com",
					"role": "Admin",
				})
				tokenString, _ := token.SignedString([]byte(secretKey))

				req := httptest.NewRequest(http.MethodGet, "/query", nil)
				req.Header.Set("Authorization", "Bearer "+tokenString)
				return req
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				// No mocks needed
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{
					Secrets: config.Secrets{SecretKey: secretKey},
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, "config", cfg)
				// Don't add service to context
				return ctx
			},
			expectedStatus: http.StatusInternalServerError,
			nextCalled:     false,
		},
		{
			name: "config not in context",
			setupRequest: func() *http.Request {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"exp":  time.Now().Add(24 * time.Hour).Unix(),
					"sub":  "test@example.com",
					"role": "Admin",
				})
				tokenString, _ := token.SignedString([]byte(secretKey))

				req := httptest.NewRequest(http.MethodGet, "/query", nil)
				req.Header.Set("Authorization", "Bearer "+tokenString)
				return req
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				// No mocks needed
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				service := &svc.Service{
					AdminUser: mockAdminUser,
					Role:      mockRole,
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				// Don't add config to context
				return ctx
			},
			expectedStatus: http.StatusInternalServerError,
			nextCalled:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAdminUser := mocks.NewMockAdminUser(ctrl)
			mockRole := mocks.NewMockRole(ctrl)
			tt.setupMocks(mockAdminUser, mockRole)

			req := tt.setupRequest()
			ctx := tt.setupContext(mockAdminUser, mockRole)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			nextCalled := false
			testName := tt.name // Capture test name for use in closure
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true

				// For websocket upgrade, user/role may not be in context (that's expected)
				if testName == "skip authentication for websocket upgrade" {
					w.WriteHeader(http.StatusOK)
					return
				}

				// Verify user and role are in context (only for authenticated requests)
				user, userOk := r.Context().Value(constants.UserContextValue).(*ent.AdminUser)
				role, roleOk := r.Context().Value(constants.RoleContextValue).(*ent.Role)

				// For other authenticated requests, user and role should be in context
				// Skip this check for websocket upgrade test
				if tt.nextCalled && testName != "skip authentication for websocket upgrade" {
					if !userOk {
						t.Error("expected user in context")
					}
					if !roleOk {
						t.Error("expected role in context")
					}
				}

				if user != nil && role != nil {
					w.WriteHeader(http.StatusOK)
				}
			})

			authenticator := Authenticator(nextHandler)
			authenticator.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if nextCalled != tt.nextCalled {
				t.Errorf("expected nextCalled=%v, got %v", tt.nextCalled, nextCalled)
			}
		})
	}
}
