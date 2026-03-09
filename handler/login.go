package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/domain/svc"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/GoLabra/labra/jwtrefresh"
	"github.com/GoLabra/labra/refreshtoken"
	"golang.org/x/crypto/bcrypt"
)

type LoginFormData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User interface {
	DefaultRole(ctx context.Context) (*ent.Role, error)
}

// TODO: 1. sanitize error messages; 2. move to api; 3. add logs;
func Login(w http.ResponseWriter, r *http.Request) {
	service, ok := r.Context().Value(constants.AdminServiceContextValue).(*svc.Service)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(svc.ErrServiceNotSetInContext)
		return
	}

	var loginFormData LoginFormData

	if err := decodeJSONLimited(w, r, &loginFormData, MaxBodyLoginBytes); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("invalid request body: %v", err)
		return
	}

	loginFormData.Sanitize()
	if err := loginFormData.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("validation failed: %v", err)
		return
	}

	// Use internal context for user lookup (bypasses permission checks)
	iCtx := context.WithValue(r.Context(), constants.IsInternalOperationContextValue, true)

	user, err := service.AdminUser.GetOne(iCtx, ent.AdminUserWhereUniqueInput{
		Email: &loginFormData.Email,
	})

	if err != nil && !ent.IsNotFound(err) {
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
		log.Println("error getting role")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if role == nil {
		log.Println("user default role not found")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginFormData.Password))
	if err != nil {
		log.Printf("error comparing password: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	appConfig, ok := r.Context().Value("config").(*config.AppConfig)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	adminEntClient, ok := r.Context().Value(constants.AdminEntClientContextValue).(*ent.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("admin ent client not found")
		return
	}

	pair, err := jwtrefresh.IssueTokenPair(
		appConfig.SecretKey,
		user.Email,
		role.Name,
		jwtrefresh.SubjectTypeAdmin,
		appConfig.AccessTokenTTL,
		appConfig.RefreshTokenTTL,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to issue token pair: %v", err)
		return
	}

	refreshTokenService := refreshtoken.NewService(refreshtoken.NewRepository(adminEntClient, appConfig.DBDialect))
	err = refreshTokenService.Save(
		r.Context(),
		pair.RefreshToken,
		user.Email,
		jwtrefresh.SubjectTypeAdmin,
		role.Name,
		pair.RefreshExpiresAt,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to save refresh token: %v", err)
		return
	}

	secure := r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")

	// CSRF cookie (readable by JS) - double-submit cookie pattern
	csrf, err := NewCSRFToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	SetCSRFCookie(w, csrf, secure)

	cookieDomain := ""
	SetJWTCookieWithTTL(w, r, pair.AccessToken, cookieDomain, jwtrefresh.EffectiveAccessTTL(appConfig.AccessTokenTTL))

	response, _ := json.Marshal(map[string]string{
		"token":         pair.AccessToken,
		"refresh_token": pair.RefreshToken,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
