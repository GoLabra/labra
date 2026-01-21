package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/domain/svc"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/GoLabra/labra/utils"
	"github.com/golang-jwt/jwt"
)

type ActivateFormData struct {
	Email         string `json:"email"`
	ActivationKey string `json:"activationKey"`
}

// ActivateAdminUser activates an admin user using their email and activation key.
// If the activation key is valid, returns a JWT token.
// This is a public endpoint that does not require authentication.
func ActivateAdminUser(w http.ResponseWriter, r *http.Request) {
	var activateFormData ActivateFormData

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error reading body: %v", err)
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &activateFormData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error unmarshaling body: %v", err)
		return
	}

	if activateFormData.Email == "" || activateFormData.ActivationKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("email and activation key are required")
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
		Email: &activateFormData.Email,
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

	// Check if activation key is set
	if user.ActivationKey == nil || *user.ActivationKey == "" {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("activation key not set for user")
		return
	}

	// Verify the activation key
	if !utils.VerifyActivationKey(activateFormData.ActivationKey, *user.ActivationKey) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("invalid activation key")
		return
	}

	// Get the user's default role
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

	// Generate JWT token
	config, ok := r.Context().Value("config").(*config.Config)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"sub":  user.Email,
		"role": role.Name,
	})

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



