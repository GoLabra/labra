package handler

import (
	"app/ent"
	"encoding/json"
	"net/http"
	"time"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/golang-jwt/jwt"
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

	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"sub":  user.Email,
		"role": role.Name,
	})

	appConfig, ok := r.Context().Value("config").(*config.AppConfig)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	signedToken, err := token.SignedString([]byte(appConfig.SecretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	response, _ := json.Marshal(map[string]string{
		"token": signedToken,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
