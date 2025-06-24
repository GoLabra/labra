package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/GoLabra/labra/src/api/config"
	"github.com/GoLabra/labra/src/api/constants"
	"github.com/GoLabra/labra/src/api/entgql/domain/svc"
	"github.com/GoLabra/labra/src/api/entgql/ent"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type LoginFormData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TODO: 1. sanitize error messages; 2. move to api; 3. add logs;
func Login(w http.ResponseWriter, r *http.Request) {
	var (
		loginFormData LoginFormData
	)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &loginFormData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	service, ok := r.Context().Value(constants.AdminServiceContextValue).(*svc.Service) // TODO @David use admin service and move to api
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := service.User.GetOne(r.Context(), ent.UserWhereUniqueInput{
		Email: &loginFormData.Email,
	})

	if err != nil && !ent.IsNotFound(err) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginFormData.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	role, err := user.DefaultRole(r.Context())
	if err != nil || role == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"sub":  user.Email,
		"role": role.Name,
	})

	config, ok := r.Context().Value("config").(*config.Config)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	signedToken, err := token.SignedString([]byte(config.SecretKey))
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
