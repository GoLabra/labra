package handler

import (
	"app/domain/svc"
	"app/ent"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	adminSvc "github.com/GoLabra/labra/entgql/domain/svc"
	adminEnt "github.com/GoLabra/labra/entgql/ent"
	"github.com/GoLabra/labra/jwtrefresh"
	"github.com/GoLabra/labra/refreshtoken"
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

	appConfig, ok := r.Context().Value("config").(*config.AppConfig)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	entClient, ok := r.Context().Value(constants.EntClientContextValue).(*ent.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ent client not found")
		return
	}

	pair, err := jwtrefresh.IssueTokenPair(
		appConfig.SecretKey,
		user.Email,
		role.Name,
		jwtrefresh.SubjectTypeUser,
		appConfig.AccessTokenTTL,
		appConfig.RefreshTokenTTL,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to issue token pair: %v", err)
		return
	}

	refreshTokenService := refreshtoken.NewService(refreshtoken.NewRepository(entClient, appConfig.DBDialect))
	err = refreshTokenService.Save(
		r.Context(),
		pair.RefreshToken,
		user.Email,
		jwtrefresh.SubjectTypeUser,
		role.Name,
		pair.RefreshExpiresAt,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to save refresh token: %v", err)
		return
	}

	response, _ := json.Marshal(map[string]string{
		"token":         pair.AccessToken,
		"refresh_token": pair.RefreshToken,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
