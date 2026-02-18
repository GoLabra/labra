package handler

import (
	"context"
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

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// skip the authentication for subscriptions
		if strings.EqualFold(r.Header.Get("Connection"), "Upgrade") &&
			strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
			next.ServeHTTP(w, r)
			return
		}

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

		var tokenString string
		usedCookieAuth := false

		if token := r.Header.Get("Authorization"); token != "" {
			tokenString = strings.TrimPrefix(token, "Bearer ")
		} else if token, err := r.Cookie("jwt"); err == nil {
			usedCookieAuth = true
			tokenString = token.Value
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("token not found")
			return
		}

		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("token not found")
			return
		}

		token, err := jwt.ParseString(
			tokenString,
			jwt.WithValidate(true),
			jwt.WithKey(jwt_hs.HS256, []byte(appConfig.SecretKey)),
			jwt.WithVerify(true),
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

		userEmail := claims["sub"].(string)

		// Use internal context for authentication operations (bypasses permission checks)
		iCtx := context.WithValue(r.Context(), constants.IsInternalOperationContextValue, true)

		user, err := service.AdminUser.GetOne(iCtx, ent.AdminUserWhereUniqueInput{Email: &userEmail})
		if err != nil && !ent.IsNotFound(err) {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("admin user not found: %v", err)
			return
		}

		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("user not found")
			return
		}

		roleName := claims["role"].(string)
		role, err := service.Role.GetOne(iCtx, ent.RoleWhereUniqueInput{Name: &roleName})
		if err != nil {
			log.Printf("role not found: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// TODO validate if the user has these roles

		ctx := r.Context()

		ctx = context.WithValue(ctx, constants.UserContextValue, user)
		ctx = context.WithValue(ctx, constants.RoleContextValue, role)

		// If auth came from cookie-based JWT, ensure CSRF cookie exists (double-submit cookie pattern).
		if usedCookieAuth {
			if _, err := r.Cookie("csrf_token"); err != nil {
				csrf, err := NewCSRFToken()
				if err == nil {
					secure := r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
					SetCSRFCookie(w, csrf, secure)
				}
			}
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
