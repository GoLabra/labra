package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/domain/svc"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/GoLabra/labra/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
)

func toWSURL(httpURL string) string {
	if strings.HasPrefix(httpURL, "https://") {
		return "wss://" + strings.TrimPrefix(httpURL, "https://")
	}
	return "ws://" + strings.TrimPrefix(httpURL, "http://")
}

func signHS256(t *testing.T, key string, claims jwt.MapClaims) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte(key))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

func dialWS(t *testing.T, wsURL string, hdr http.Header) (*websocket.Conn, *http.Response, error) {
	t.Helper()
	return websocket.DefaultDialer.Dial(wsURL, hdr)
}

func TestAuthenticator_WebSocketHandshake_RequiresJWT(t *testing.T) {
	secretKey := "test-secret-key-that-is-long-enough-for-validation-minimum-16-chars"

	tests := []struct {
		name           string
		wsPath         string
		header         http.Header
		setupMocks     func(*mocks.MockAdminUser, *mocks.MockRole)
		expectedOK     bool // true => 101, false => 401
		expectNextCall bool
	}{
		{
			name:   "no token -> 401",
			wsPath: "/query",
			header: http.Header{},
			setupMocks: func(_ *mocks.MockAdminUser, _ *mocks.MockRole) {
				// should fail before DB calls
			},
			expectedOK:     false,
			expectNextCall: false,
		},
		{
			name:   "invalid token -> 401",
			wsPath: "/query",
			header: http.Header{"Cookie": []string{"jwt=definitely.invalid.token"}},
			setupMocks: func(_ *mocks.MockAdminUser, _ *mocks.MockRole) {
				// should fail before DB calls
			},
			expectedOK:     false,
			expectNextCall: false,
		},
		{
			name:   "valid JWT cookie -> 101",
			wsPath: "/query",
			header: func() http.Header {
				j := signHS256(t, secretKey, jwt.MapClaims{
					"exp":  time.Now().Add(5 * time.Minute).Unix(),
					"sub":  "test@example.com",
					"role": "Admin",
				})
				return http.Header{"Cookie": []string{"jwt=" + j}}
			}(),
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				user := CreateTestAdminUser("user-1", "test@example.com", "hashed-password")
				role := CreateTestRole("role-1", "Admin")

				mockAdminUser.EXPECT().
					GetOne(gomock.Any(), gomock.Any()).
					Return(user, nil)

				mockRole.EXPECT().
					GetOne(gomock.Any(), gomock.Any()).
					Return(role, nil)

				// IMPORTANT: use gomock.Nil() for typed nil args, not raw nil
				mockAdminUser.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Nil(), gomock.Nil(), gomock.Any(), gomock.Nil()).
					Return([]*ent.AdminUser{user}, nil).
					AnyTimes()
			},
			expectedOK:     true,
			expectNextCall: true,
		},
		{
			name: "valid JWT in query param -> 101",
			wsPath: func() string {
				j := signHS256(t, secretKey, jwt.MapClaims{
					"exp":  time.Now().Add(5 * time.Minute).Unix(),
					"sub":  "test@example.com",
					"role": "Admin",
				})
				return "/query?token=" + j
			}(),
			header: http.Header{},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				user := CreateTestAdminUser("user-1", "test@example.com", "hashed-password")
				role := CreateTestRole("role-1", "Admin")

				mockAdminUser.EXPECT().
					GetOne(gomock.Any(), gomock.Any()).
					Return(user, nil)

				mockRole.EXPECT().
					GetOne(gomock.Any(), gomock.Any()).
					Return(role, nil)

				// IMPORTANT: use gomock.Nil() for typed nil args, not raw nil
				mockAdminUser.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Nil(), gomock.Nil(), gomock.Any(), gomock.Nil()).
					Return([]*ent.AdminUser{user}, nil).
					AnyTimes()
			},
			expectedOK:     true,
			expectNextCall: true,
		},
		{
			name:   "ws upgrade with Connection list and no token -> 401",
			wsPath: "/query",
			// websocket.Dialer will set Connection/Upgrade itself, but we still validate our JWT extraction;
			// this case just ensures we don't accidentally treat WS upgrades differently.
			header: http.Header{},
			setupMocks: func(_ *mocks.MockAdminUser, _ *mocks.MockRole) {
				// should fail before DB calls
			},
			expectedOK:     false,
			expectNextCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAdminUser := mocks.NewMockAdminUser(ctrl)
			mockRole := mocks.NewMockRole(ctrl)
			tt.setupMocks(mockAdminUser, mockRole)

			cfg := &config.AppConfig{
				Secrets: config.Secrets{
					SecretKey: secretKey,
				},
			}
			service := &svc.Service{
				AdminUser: mockAdminUser,
				Role:      mockRole,
			}

			var nextCalled atomic.Bool
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled.Store(true)

				// prove auth context exists on upgrade request
				if r.Context().Value(constants.UserContextValue) == nil {
					t.Errorf("expected user in context on ws upgrade request")
				}
				if r.Context().Value(constants.RoleContextValue) == nil {
					t.Errorf("expected role in context on ws upgrade request")
				}

				up := websocket.Upgrader{}
				_, err := up.Upgrade(w, r, nil)
				if err != nil {
					t.Errorf("upgrade failed: %v", err)
					return
				}
			})

			// inject required context values then run Authenticator
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				Authenticator(next).ServeHTTP(w, r.WithContext(ctx))
			})

			srv := httptest.NewServer(h)
			defer srv.Close()

			wsURL := toWSURL(srv.URL) + tt.wsPath
			c, resp, err := dialWS(t, wsURL, tt.header)

			if tt.expectedOK {
				if err != nil {
					code := 0
					if resp != nil {
						code = resp.StatusCode
					}
					t.Fatalf("expected ws handshake OK, got err=%v status=%d", err, code)
				}
				_ = c.Close()
			} else {
				if err == nil {
					_ = c.Close()
					t.Fatalf("expected handshake to fail but it succeeded")
				}
				if resp == nil || resp.StatusCode != http.StatusUnauthorized {
					code := 0
					if resp != nil {
						code = resp.StatusCode
					}
					t.Fatalf("expected 401 on bad handshake; got status=%d err=%v", code, err)
				}
			}

			if nextCalled.Load() != tt.expectNextCall {
				t.Fatalf("expected nextCalled=%v got %v", tt.expectNextCall, nextCalled.Load())
			}
		})
	}
}

func TestAuthenticator_WebSocketHandshake_CentrifugoTokenValidation(t *testing.T) {
	secretKey := "test-secret-key-that-is-long-enough-for-validation-minimum-16-chars"
	centrifugoKey := "test-centrifugo-key-that-is-long-enough-32-chars"

	makeAppJWT := func(email, role string) string {
		return signHS256(t, secretKey, jwt.MapClaims{
			"exp":  time.Now().Add(5 * time.Minute).Unix(),
			"sub":  email,
			"role": role,
		})
	}

	makeCentrifugoJWT := func(sub string) string {
		// Your Authenticator validates cf_token as HS256 using appConfig.Secrets.CentrifugoKey
		return signHS256(t, centrifugoKey, jwt.MapClaims{
			"exp": time.Now().Add(5 * time.Minute).Unix(),
			"sub": sub,
		})
	}

	tests := []struct {
		name           string
		wsPath         string
		header         http.Header
		setupMocks     func(*mocks.MockAdminUser, *mocks.MockRole)
		expectedOK     bool
		expectNextCall bool
	}{
		{
			name: "invalid cf_token -> 401",
			wsPath: func() string {
				appJWT := makeAppJWT("test@example.com", "Admin")
				return "/query?cf_token=definitely.invalid.token&token=" + appJWT
			}(),
			header: http.Header{},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				user := CreateTestAdminUser("user-1", "test@example.com", "hashed-password")
				role := CreateTestRole("role-1", "Admin")

				mockAdminUser.EXPECT().
					GetOne(gomock.Any(), gomock.Any()).
					Return(user, nil)

				mockRole.EXPECT().
					GetOne(gomock.Any(), gomock.Any()).
					Return(role, nil)

				// IMPORTANT: use gomock.Nil() for typed nil args, not raw nil
				mockAdminUser.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Nil(), gomock.Nil(), gomock.Any(), gomock.Nil()).
					Return([]*ent.AdminUser{user}, nil).
					AnyTimes()
			},
			expectedOK:     false,
			expectNextCall: false,
		},
		{
			name: "valid cf_token with matching sub(email) -> 101",
			wsPath: func() string {
				appJWT := makeAppJWT("test@example.com", "Admin")
				cfJWT := makeCentrifugoJWT("test@example.com")
				return "/query?cf_token=" + cfJWT + "&token=" + appJWT
			}(),
			header: http.Header{},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				user := CreateTestAdminUser("user-1", "test@example.com", "hashed-password")
				role := CreateTestRole("role-1", "Admin")

				mockAdminUser.EXPECT().
					GetOne(gomock.Any(), gomock.Any()).
					Return(user, nil)

				mockRole.EXPECT().
					GetOne(gomock.Any(), gomock.Any()).
					Return(role, nil)

				// IMPORTANT: use gomock.Nil() for typed nil args, not raw nil
				mockAdminUser.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Nil(), gomock.Nil(), gomock.Any(), gomock.Nil()).
					Return([]*ent.AdminUser{user}, nil).
					AnyTimes()
			},
			expectedOK:     true,
			expectNextCall: true,
		},
		{
			name: "valid cf_token but sub mismatch -> 401",
			wsPath: func() string {
				appJWT := makeAppJWT("test@example.com", "Admin")
				cfJWT := makeCentrifugoJWT("someoneelse@example.com")
				return "/query?cf_token=" + cfJWT + "&token=" + appJWT
			}(),
			header: http.Header{},
			setupMocks: func(mockAdminUser *mocks.MockAdminUser, mockRole *mocks.MockRole) {
				user := CreateTestAdminUser("user-1", "test@example.com", "hashed-password")
				role := CreateTestRole("role-1", "Admin")

				mockAdminUser.EXPECT().
					GetOne(gomock.Any(), gomock.Any()).
					Return(user, nil)

				mockRole.EXPECT().
					GetOne(gomock.Any(), gomock.Any()).
					Return(role, nil)

				// IMPORTANT: use gomock.Nil() for typed nil args, not raw nil
				mockAdminUser.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Nil(), gomock.Nil(), gomock.Any(), gomock.Nil()).
					Return([]*ent.AdminUser{user}, nil).
					AnyTimes()
			},
			expectedOK:     false,
			expectNextCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAdminUser := mocks.NewMockAdminUser(ctrl)
			mockRole := mocks.NewMockRole(ctrl)
			tt.setupMocks(mockAdminUser, mockRole)

			cfg := &config.AppConfig{
				Secrets: config.Secrets{
					SecretKey:     secretKey,
					CentrifugoKey: centrifugoKey,
				},
			}
			service := &svc.Service{
				AdminUser: mockAdminUser,
				Role:      mockRole,
			}

			var nextCalled atomic.Bool
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled.Store(true)

				if r.Context().Value(constants.UserContextValue) == nil {
					t.Errorf("expected user in context on ws upgrade request")
				}
				if r.Context().Value(constants.RoleContextValue) == nil {
					t.Errorf("expected role in context on ws upgrade request")
				}

				up := websocket.Upgrader{}
				_, err := up.Upgrade(w, r, nil)
				if err != nil {
					t.Errorf("upgrade failed: %v", err)
					return
				}
			})

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.Background()
				ctx = context.WithValue(ctx, constants.AdminServiceContextValue, service)
				ctx = context.WithValue(ctx, "config", cfg)
				Authenticator(next).ServeHTTP(w, r.WithContext(ctx))
			})

			srv := httptest.NewServer(h)
			defer srv.Close()

			wsURL := toWSURL(srv.URL) + tt.wsPath
			c, resp, err := dialWS(t, wsURL, tt.header)

			if tt.expectedOK {
				if err != nil {
					code := 0
					if resp != nil {
						code = resp.StatusCode
					}
					t.Fatalf("expected ws handshake OK, got err=%v status=%d", err, code)
				}
				_ = c.Close()
			} else {
				if err == nil {
					_ = c.Close()
					t.Fatalf("expected handshake to fail but it succeeded")
				}
				if resp == nil || resp.StatusCode != http.StatusUnauthorized {
					code := 0
					if resp != nil {
						code = resp.StatusCode
					}
					t.Fatalf("expected 401 on bad handshake; got status=%d err=%v", code, err)
				}
			}

			if nextCalled.Load() != tt.expectNextCall {
				t.Fatalf("expected nextCalled=%v got %v", tt.expectNextCall, nextCalled.Load())
			}
		})
	}
}
