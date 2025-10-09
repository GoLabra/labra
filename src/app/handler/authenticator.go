package handler

import (
	"app/domain/svc"
	"app/ent"
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/GoLabra/labra/src/api/config"
	"github.com/GoLabra/labra/src/api/constants"
	adminSvc "github.com/GoLabra/labra/src/api/entgql/domain/svc"
	adminEnt "github.com/GoLabra/labra/src/api/entgql/ent"
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

		config, ok := r.Context().Value("config").(*config.Config)

		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("config not found")
			return
		}

		var tokenString string

		if token := r.Header.Get("Authorization"); token != "" {
			tokenString = strings.TrimPrefix(token, "Bearer ")
		} else if token, err := r.Cookie("jwt"); err == nil {
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
			jwt.WithKey(jwt_hs.HS256, []byte(config.SecretKey)),
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
		user, err := service.User.GetOne(iCtx, ent.UserWhereUniqueInput{Email: &userEmail})
		if err != nil && !ent.IsNotFound(err) {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("error getting user by email: %v", err)
			return
		}

		if user != nil {
			roleName := claims["role"].(string)
			role, err := adminService.Role.GetOne(iCtx, adminEnt.RoleWhereUniqueInput{Name: &roleName})
			if err != nil {
				log.Printf("role not found: %v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.UserContextValue, user)
			ctx = context.WithValue(ctx, constants.RoleContextValue, role)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}

		adminUser, err := adminService.AdminUser.GetOne(iCtx, adminEnt.AdminUserWhereUniqueInput{Email: &userEmail})
		if err != nil && !adminEnt.IsNotFound(err) {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("error getting user by email: %v", err)
			return
		}

		if adminUser != nil {
			roleName := claims["role"].(string)
			role, err := adminService.Role.GetOne(iCtx, adminEnt.RoleWhereUniqueInput{Name: &roleName})
			if err != nil {
				log.Printf("role not found: %v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.UserContextValue, adminUser)
			ctx = context.WithValue(ctx, constants.RoleContextValue, role)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}

		w.WriteHeader(http.StatusUnauthorized)
		log.Println("user not found")
	})
}
