package handler

import (
	"app/domain/svc"
	"app/ent"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	adminSvc "github.com/GoLabra/labra/entgql/domain/svc"
	adminEnt "github.com/GoLabra/labra/entgql/ent"
	"github.com/golang-jwt/jwt"
)

type LoginFormData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var (
		loginFormData LoginFormData
	)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error reading body: %v", err)
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &loginFormData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error unmarshaling body: %v", err)
		return
	}

	service, ok := r.Context().Value(constants.ServiceContextValue).(*svc.Service)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(adminSvc.ErrServiceNotSetInContext)
		return
	}

	iCtx := context.WithValue(r.Context(), constants.IsInternalOperationContextValue, true)

	user, err := service.User.GetOne(iCtx, ent.UserWhereUniqueInput{
		Email: &loginFormData.Email,
	})

	if err != nil && !adminEnt.IsNotFound(err) {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error getting user: %v", err)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("user not found")
		return
	}

	role, err := user.DefaultRole(iCtx)
	if err != nil && !ent.IsNotFound(err) {
		log.Printf("error getting role: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if role == nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("user default role not found")
		return
	}

	if user.Password != loginFormData.Password { // THIS IS A BASIC EXAMPLE. DO NOT STORE PASSWORDS IN PLAIN TEXT
		log.Printf("passwords do not match: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
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
		return
	}

	response, _ := json.Marshal(map[string]string{
		"token": signedToken,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
