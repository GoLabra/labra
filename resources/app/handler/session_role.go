package handler

import (
	"app/ent"
	"net/http"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	sharedHandler "github.com/GoLabra/labra/handler"
)

type ChangeSessionRoleRequest struct {
	Role string
}

func ChangeSessionRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(constants.UserContextValue).(*ent.User)
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

	signedToken, err := issueUserAccessToken(appConfig, user.Email, role.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sharedHandler.SetJWTCookieWithTTL(w, r, signedToken, "", accessTokenTTL(appConfig))
	writeTokenResponse(w, signedToken)
}
