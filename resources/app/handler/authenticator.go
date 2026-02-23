package handler

import (
	"app/domain/svc"
	"app/ent"
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	adminSvc "github.com/GoLabra/labra/entgql/domain/svc"
	adminEnt "github.com/GoLabra/labra/entgql/ent"
	jwt_hs "github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func isWebSocketUpgrade(r *http.Request) bool {
	conn := r.Header.Get("Connection")
	upg := r.Header.Get("Upgrade")
	return strings.Contains(strings.ToLower(conn), "upgrade") &&
		strings.EqualFold(upg, "websocket")
}

// Only allow query token for WS handshakes (so regular HTTP requests still require header/cookie).
func extractJWTFromRequest(r *http.Request) string {
	// 1) Authorization header
	if h := strings.TrimSpace(r.Header.Get("Authorization")); h != "" {
		// support "Bearer <token>" with/without extra spaces
		if strings.HasPrefix(strings.ToLower(h), "bearer") {
			return strings.TrimSpace(strings.TrimPrefix(h, "Bearer"))
		}
		return strings.TrimSpace(strings.TrimPrefix(h, "Bearer"))
	}

	// 2) Cookie
	if c, err := r.Cookie("jwt"); err == nil && strings.TrimSpace(c.Value) != "" {
		return strings.TrimSpace(c.Value)
	}

	// 3) Query params (WS-friendly) — ONLY for WS upgrades
	if isWebSocketUpgrade(r) {
		q := r.URL.Query()
		for _, key := range []string{"token", "jwt", "access_token"} {
			if v := strings.TrimSpace(q.Get(key)); v != "" {
				return v
			}
		}
	}

	return ""
}

// Try to read centrifugo secret key from config in a compile-safe way.
// Supports either:
//
//	appCfg.CentrifugoKey
//	appCfg.Secrets.CentrifugoKey
func getCentrifugoSecret(appCfg *config.AppConfig) string {
	if appCfg == nil {
		return ""
	}

	v := reflect.ValueOf(appCfg)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ""
	}

	// appCfg.Secrets.CentrifugoKey
	secrets := v.FieldByName("Secrets")
	if secrets.IsValid() {
		if secrets.Kind() == reflect.Pointer {
			if !secrets.IsNil() {
				secrets = secrets.Elem()
			}
		}
		if secrets.IsValid() && secrets.Kind() == reflect.Struct {
			ck := secrets.FieldByName("CentrifugoKey")
			if ck.IsValid() && ck.Kind() == reflect.String {
				return strings.TrimSpace(ck.String())
			}
		}
	}

	// appCfg.CentrifugoKey
	ck := v.FieldByName("CentrifugoKey")
	if ck.IsValid() && ck.Kind() == reflect.String {
		return strings.TrimSpace(ck.String())
	}

	return ""
}

// Optional: validate centrifugo connection token (if your client sends it via query).
func validateCentrifugoTokenIfPresent(r *http.Request, appCfg *config.AppConfig, userEmail string, userID any) error {
	q := r.URL.Query()

	cfTok := strings.TrimSpace(q.Get("cf_token"))
	if cfTok == "" {
		cfTok = strings.TrimSpace(q.Get("centrifugo_token"))
	}
	if cfTok == "" {
		return nil // no centrifugo token provided
	}

	secret := getCentrifugoSecret(appCfg)
	if secret == "" {
		return fmt.Errorf("centrifugo token provided but centrifugo secret not configured")
	}

	parsed, err := jwt.ParseString(
		cfTok,
		jwt.WithValidate(true),
		jwt.WithVerify(true),
		jwt.WithKey(jwt_hs.HS256, []byte(secret)),
	)
	if err != nil {
		return fmt.Errorf("invalid centrifugo token: %w", err)
	}

	claims, err := parsed.AsMap(context.Background())
	if err != nil {
		return fmt.Errorf("invalid centrifugo claims: %w", err)
	}

	sub, _ := claims["sub"].(string)
	if strings.TrimSpace(sub) == "" {
		return fmt.Errorf("centrifugo token missing sub claim")
	}

	uid := fmt.Sprint(userID)
	if sub != userEmail && sub != uid {
		return fmt.Errorf("centrifugo token subject does not match authenticated user")
	}

	return nil
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ✅ NO websocket bypass here. WS upgrades must authenticate.

		adminService, ok := r.Context().Value(constants.AdminServiceContextValue).(*adminSvc.Service)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(adminSvc.ErrServiceNotSetInContext)
			return
		}

		service, ok := r.Context().Value(constants.ServiceContextValue).(*svc.Service)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(adminSvc.ErrServiceNotSetInContext)
			return
		}

		appConfig, ok := r.Context().Value("config").(*config.AppConfig)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("config not found")
			return
		}

		tokenString := extractJWTFromRequest(r)
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("token not found")
			return
		}

		token, err := jwt.ParseString(
			tokenString,
			jwt.WithValidate(true),
			jwt.WithVerify(true),
			jwt.WithKey(jwt_hs.HS256, []byte(appConfig.SecretKey)),
		)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("error parsing token: %v", err)
			return
		}

		claims, err := token.AsMap(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("error parsing claims: %v", err)
			return
		}

		userEmail, _ := claims["sub"].(string)
		userEmail = strings.TrimSpace(userEmail)
		if userEmail == "" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("missing sub claim")
			return
		}

		roleName, _ := claims["role"].(string)
		roleName = strings.TrimSpace(roleName)
		if roleName == "" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("missing role claim")
			return
		}

		// Use internal context for authentication operations (bypasses permission checks)
		iCtx := context.WithValue(r.Context(), constants.IsInternalOperationContextValue, true)

		// 1) Normal user
		user, err := service.User.GetOne(iCtx, ent.UserWhereUniqueInput{Email: &userEmail})
		if err != nil && !ent.IsNotFound(err) {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("error getting user by email: %v", err)
			return
		}

		if user != nil {
			role, err := adminService.Role.GetOne(iCtx, adminEnt.RoleWhereUniqueInput{Name: &roleName})
			if err != nil {
				log.Printf("role not found: %v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// ✅ If WS: validate centrifugo token too (if provided)
			if isWebSocketUpgrade(r) {
				if err := validateCentrifugoTokenIfPresent(r, appConfig, user.Email, user.ID); err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					log.Printf("centrifugo token validation failed: %v", err)
					return
				}
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.UserContextValue, user)
			ctx = context.WithValue(ctx, constants.RoleContextValue, role)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// 2) Admin user
		adminUser, err := adminService.AdminUser.GetOne(iCtx, adminEnt.AdminUserWhereUniqueInput{Email: &userEmail})
		if err != nil && !adminEnt.IsNotFound(err) {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("error getting admin user by email: %v", err)
			return
		}

		if adminUser != nil {
			role, err := adminService.Role.GetOne(iCtx, adminEnt.RoleWhereUniqueInput{Name: &roleName})
			if err != nil {
				log.Printf("role not found: %v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// ✅ If WS: validate centrifugo token too (if provided)
			if isWebSocketUpgrade(r) {
				if err := validateCentrifugoTokenIfPresent(r, appConfig, adminUser.Email, adminUser.ID); err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					log.Printf("centrifugo token validation failed: %v", err)
					return
				}
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.UserContextValue, adminUser)
			ctx = context.WithValue(ctx, constants.RoleContextValue, role)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		log.Println("user not found")
	})
}
