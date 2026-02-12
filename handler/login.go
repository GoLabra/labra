package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/domain/svc"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/golang-jwt/jwt"
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

	service, ok := r.Context().Value(constants.AdminServiceContextValue).(*svc.Service)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(svc.ErrServiceNotSetInContext)
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

	secure := r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")

	// JWT cookie (HttpOnly)
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// CSRF cookie (readable by JS) - double-submit cookie pattern
	csrf, err := NewCSRFToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	SetCSRFCookie(w, csrf, secure)

	response, _ := json.Marshal(map[string]string{
		"token": signedToken,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
