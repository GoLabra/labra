package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/GoLabra/labra/jwtrefresh"
	"github.com/GoLabra/labra/tokenrevocation"

	jwt_hs "github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type ChangeSessionRoleRequest struct {
	Role string
}

func ChangeSessionRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(constants.UserContextValue).(*ent.AdminUser)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	role, ok := ctx.Value(constants.RoleContextValue).(*ent.Role)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	appConfig, ok := r.Context().Value("config").(*config.AppConfig)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	adminEntClient, _ := r.Context().Value(constants.AdminEntClientContextValue).(*ent.Client)
	if adminEntClient != nil {
		tokenString := extractJWTFromRequest(r)
		if tokenString != "" {
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

			expiresAt, err := claimUnixTime(claims, "exp")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				log.Printf("missing/invalid exp claim: %v", err)
				return
			}

			subjectType := claimStringValue(claims, jwtrefresh.ClaimSubjectType)
			if subjectType == "" {
				subjectType = jwtrefresh.SubjectTypeAdmin
			}

			if err := tokenrevocation.RevokeToken(
				r.Context(),
				adminEntClient,
				adminEntClient.DialectName(),
				tokenString,
				user.Email,
				subjectType,
				expiresAt,
				time.Now().UTC(),
			); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("failed revoking current session token: %v", err)
				return
			}
		}
	}

	signedToken, _, err := jwtrefresh.IssueAccessToken(
		appConfig.SecretKey,
		user.Email,
		role.Name,
		jwtrefresh.SubjectTypeAdmin,
		appConfig.AccessTokenTTL,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(map[string]string{
		"token": signedToken,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
