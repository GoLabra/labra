package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

func TestSignup(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    SignupFormData
		setupMocks     func(*mocks.MockAdminUser, *mocks.MockRole)
		setupContext   func(*mocks.MockAdminUser, *mocks.MockRole) context.Context
		expectedStatus int
	}{
		{
			name: "successful signup with existing super admin role",
			requestBody: SignupFormData{
				Email:     "newuser@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				superAdminRole := CreateTestRole("role-1", "SuperAdmin")
				mockRole.EXPECT().
					GetOne(gomock.Any(), ent.RoleWhereUniqueInput{Name: stringPtr("SuperAdmin")}).
					Return(superAdminRole, nil)

				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				newUser := CreateTestAdminUser("user-1", "newuser@example.com", string(hashedPassword))
				newUser.FirstName = "John"
				newUser.LastName = "Doe"

				mockAdminUser.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(newUser, nil)
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{Secrets: config.Secrets{SecretKey: "test-secret-key-that-is-long-enough"}}
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
		},
		{
			name: "successful signup creating super admin role",
			requestBody: SignupFormData{
				Email:     "newuser@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				// First call - role not found
				mockRole.EXPECT().
					GetOne(gomock.Any(), ent.RoleWhereUniqueInput{Name: stringPtr("SuperAdmin")}).
					Return(nil, &ent.NotFoundError{})

				// Create the role
				superAdminRole := CreateTestRole("role-1", "SuperAdmin")
				mockRole.EXPECT().
					Create(gomock.Any(), ent.CreateRoleInput{Name: "SuperAdmin"}).
					Return(superAdminRole, nil)

				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				newUser := CreateTestAdminUser("user-1", "newuser@example.com", string(hashedPassword))
				newUser.FirstName = "John"
				newUser.LastName = "Doe"

				mockAdminUser.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(newUser, nil)
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{Secrets: config.Secrets{SecretKey: "test-secret-key-that-is-long-enough"}}
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
		},
		{
			name:        "invalid json body",
			requestBody: SignupFormData{},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				// No mocks needed
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{Secrets: config.Secrets{SecretKey: "test-secret-key-that-is-long-enough"}}
				service := &svc.Service{
					AdminUser: mockAdminUser,
					Role:      mockRole,
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "service not in context",
			requestBody: SignupFormData{
				Email:     "newuser@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				// No mocks needed
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{Secrets: config.Secrets{SecretKey: "test-secret-key-that-is-long-enough"}}
				ctx := context.Background()
				ctx = context.WithValue(ctx, "config", cfg)
				// Don't add service to context
				return ctx
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "error getting super admin role",
			requestBody: SignupFormData{
				Email:     "newuser@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				mockRole.EXPECT().
					GetOne(gomock.Any(), ent.RoleWhereUniqueInput{Name: stringPtr("SuperAdmin")}).
					Return(nil, errors.New("database error"))
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{Secrets: config.Secrets{SecretKey: "test-secret-key-that-is-long-enough"}}
				service := &svc.Service{
					AdminUser: mockAdminUser,
					Role:      mockRole,
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "error creating user",
			requestBody: SignupFormData{
				Email:     "newuser@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				superAdminRole := CreateTestRole("role-1", "SuperAdmin")
				mockRole.EXPECT().
					GetOne(gomock.Any(), ent.RoleWhereUniqueInput{Name: stringPtr("SuperAdmin")}).
					Return(superAdminRole, nil)

				mockAdminUser.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			setupContext: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) context.Context {
				cfg := &config.AppConfig{Secrets: config.Secrets{SecretKey: "test-secret-key-that-is-long-enough"}}
				service := &svc.Service{
					AdminUser: mockAdminUser,
					Role:      mockRole,
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				return ctx
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAdminUser := mocks.NewMockAdminUser(ctrl)
			mockRole := mocks.NewMockRole(ctrl)
			tt.setupMocks(mockAdminUser, mockRole)

			// Create request body
			var body []byte
			if tt.name == "invalid json body" {
				body = []byte("invalid json")
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			ctx := tt.setupContext(mockAdminUser, mockRole)
			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body)).WithContext(ctx)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			Signup(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
