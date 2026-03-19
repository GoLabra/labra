package handler

import (
	"net/http"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/ent"
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

	signedToken, err := issueAccessToken(appConfig, user.Email, role.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	SetJWTCookieWithTTL(w, r, signedToken, "", accessTokenTTL(appConfig))
	writeTokenResponse(w, signedToken)
}
