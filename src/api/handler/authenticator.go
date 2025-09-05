package handler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/GoLabra/labra/src/api/config"
	"github.com/GoLabra/labra/src/api/constants"
	"github.com/GoLabra/labra/src/api/entgql/domain/svc"
	"github.com/GoLabra/labra/src/api/entgql/ent"
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
			log.Printf("user not found: %v", err)
			return
		}

		adminUser, err := service.AdminUser.GetOne(iCtx, ent.AdminUserWhereUniqueInput{Email: &userEmail})
		if err != nil && !ent.IsNotFound(err) {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("admin user not found: %v", err)
			return
		}

		if user == nil && adminUser == nil {
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
		if user != nil {
			ctx = context.WithValue(ctx, constants.UserContextValue, user)
		}

		if adminUser != nil {
			ctx = context.WithValue(ctx, constants.UserContextValue, adminUser)
		}

		ctx = context.WithValue(ctx, constants.RoleContextValue, role)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
