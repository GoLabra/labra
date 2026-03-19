package handler

import (
	"app/domain/svc"
	"app/ent"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	adminSvc "github.com/GoLabra/labra/entgql/domain/svc"
	adminEnt "github.com/GoLabra/labra/entgql/ent"
	sharedHandler "github.com/GoLabra/labra/handler"
	"golang.org/x/crypto/bcrypt"
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginFormData.Password)); err != nil {
		log.Printf("passwords do not match: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	appConfig, ok := r.Context().Value("config").(*config.AppConfig)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, accessTTL, refreshTTL, err := issueUserSession(iCtx, appConfig, user, role)
	if err != nil {
		log.Printf("error issuing session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	secure := r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")

	csrf, err := sharedHandler.NewCSRFToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sharedHandler.SetCSRFCookie(w, csrf, secure)
	sharedHandler.SetJWTCookieWithTTL(w, r, accessToken, "", accessTTL)
	sharedHandler.SetRefreshTokenCookie(w, r, refreshToken, "", refreshTTL)
	writeTokenResponse(w, accessToken)
}
