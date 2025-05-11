package handler

import (
	"app/config"
	"app/domain/svc"
	"app/ent"
	"context"
	"net/http"
	"strings"

	"github.com/GoLabra/labrago/src/api/constants"
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
		
		service, ok := r.Context().Value(constants.ServiceContextValue).(*svc.Service)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		config, ok := r.Context().Value("config").(*config.Config)

		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var tokenString string

		if token := r.Header.Get("Authorization"); token != "" {
			tokenString = strings.TrimPrefix(token, "Bearer ")
		} else if token, err := r.Cookie("jwt"); err == nil {
			tokenString = token.Value
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
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
			return
		}

		claims, err := token.AsMap(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userEmail := claims["sub"].(string)
		user, err := service.User.GetOne(r.Context(), ent.UserWhereUniqueInput{Email: &userEmail})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		roleName := claims["role"].(string)
		role, err := service.Role.GetOne(r.Context(), ent.RoleWhereUniqueInput{Name: &roleName})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, constants.UserContextValue, user)
		ctx = context.WithValue(ctx, constants.RoleContextValue, role)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
