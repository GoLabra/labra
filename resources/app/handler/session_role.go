package handler

import (
	"app/ent"
	"encoding/json"
	"net/http"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/jwtrefresh"
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

	signedToken, _, err := jwtrefresh.IssueAccessToken(
		appConfig.SecretKey,
		user.Email,
		role.Name,
		jwtrefresh.SubjectTypeUser,
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
