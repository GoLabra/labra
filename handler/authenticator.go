package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/domain/svc"
	"github.com/GoLabra/labra/entgql/ent"

	jwt_hs "github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func isWebSocketUpgrade(r *http.Request) bool {
	// RFC allows comma-separated values, so use Contains checks.
	conn := r.Header.Get("Connection")
	upg := r.Header.Get("Upgrade")
	return strings.Contains(strings.ToLower(conn), "upgrade") &&
		strings.EqualFold(upg, "websocket")
}

func extractJWTFromRequest(r *http.Request) string {
	// 1) Authorization header
	if h := r.Header.Get("Authorization"); h != "" {
		return strings.TrimSpace(strings.TrimPrefix(h, "Bearer"))
	}

	// 2) Cookie
	if c, err := r.Cookie("jwt"); err == nil && c.Value != "" {
		return c.Value
	}

	// 3) Query params (WS-friendly)
	q := r.URL.Query()
	for _, key := range []string{"token", "jwt", "access_token"} {
		if v := strings.TrimSpace(q.Get(key)); v != "" {
			return v
		}
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
		return nil // not provided → nothing to validate here
	}

	// Centrifugo connection tokens are JWTs signed with your Centrifugo secret. :contentReference[oaicite:5]{index=5}
	parsed, err := jwt.ParseString(
		cfTok,
		jwt.WithValidate(true),
		jwt.WithVerify(true),
		jwt.WithKey(jwt_hs.HS256, []byte(appCfg.Secrets.CentrifugoKey)),
	)
	if err != nil {
		return fmt.Errorf("invalid centrifugo token: %w", err)
	}

	claims, err := parsed.AsMap(context.Background())
	if err != nil {
		return fmt.Errorf("invalid centrifugo claims: %w", err)
	}

	sub, _ := claims["sub"].(string)
	if sub == "" {
		return fmt.Errorf("centrifugo token missing sub claim")
	}

	// Match against authenticated user (support either email or id-style subjects).
	uid := fmt.Sprint(userID)
	if sub != userEmail && sub != uid {
		return fmt.Errorf("centrifugo token subject does not match authenticated user")
	}

	return nil
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service, ok := r.Context().Value(constants.AdminServiceContextValue).(*svc.Service)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(svc.ErrServiceNotSetInContext)
			return
		}

		appConfig, ok := r.Context().Value("config").(*config.AppConfig)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("config not found")
			return
		}

		// ✅ DO NOT bypass auth for WebSocket upgrades.
		tokenString := extractJWTFromRequest(r)
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("token not found")
			return
		}

		// Validate app JWT
		appJWT, err := jwt.ParseString(
			tokenString,
			jwt.WithValidate(true),
			jwt.WithVerify(true),
			jwt.WithKey(jwt_hs.HS256, []byte(appConfig.Secrets.SecretKey)),
		)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("error parsing token: %v", err)
			return
		}

		claims, err := appJWT.AsMap(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("error parsing claims: %v", err)
			return
		}

		userEmail, _ := claims["sub"].(string)
		if userEmail == "" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("missing sub claim")
			return
		}

		// Use internal context for authentication operations (bypasses permission checks)
		iCtx := context.WithValue(r.Context(), constants.IsInternalOperationContextValue, true)

		user, err := service.AdminUser.GetOne(iCtx, ent.AdminUserWhereUniqueInput{Email: &userEmail})
		if err != nil && !ent.IsNotFound(err) {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("admin user lookup error: %v", err)
			return
		}
		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("user not found")
			return
		}

		roleName, _ := claims["role"].(string)
		if roleName == "" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("missing role claim")
			return
		}

		role, err := service.Role.GetOne(iCtx, ent.RoleWhereUniqueInput{Name: &roleName})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("role not found: %v", err)
			return
		}

		// ✅ If this is a WS upgrade, optionally validate Centrifugo token too (if present).
		if isWebSocketUpgrade(r) {
			if err := validateCentrifugoTokenIfPresent(r, appConfig, user.Email, user.ID); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				log.Printf("centrifugo token validation failed: %v", err)
				return
			}
		}

		// Attach auth context
		ctx := r.Context()
		ctx = context.WithValue(ctx, constants.UserContextValue, user)
		ctx = context.WithValue(ctx, constants.RoleContextValue, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
