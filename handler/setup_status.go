package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/domain/svc"
)

type SetupStatusResponse struct {
	NeedsSetup bool `json:"needsSetup"`
}

// SetupStatus returns whether the system needs initial setup (no active user with DEFAULT_USER_ROLE exists).
// This is a public endpoint that does not require authentication.
func SetupStatus(w http.ResponseWriter, r *http.Request) {
	service, ok := r.Context().Value(constants.AdminServiceContextValue).(*svc.Service)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(svc.ErrServiceNotSetInContext)
		return
	}

	cfg, ok := r.Context().Value("config").(*config.Config)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("config not found in context")
		json.NewEncoder(w).Encode(map[string]string{"error": "Configuration error"})
		return
	}

	// Get the default user role from config, default to "MANAGER" if not set
	roleName := cfg.DefaultUserRole
	if roleName == "" {
		roleName = "MANAGER"
	}

	// Use internal operation context to bypass permission checks
	iCtx := context.WithValue(r.Context(), constants.IsInternalOperationContextValue, true)

	hasUser, err := service.User.HasActiveUserWithRole(iCtx, roleName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error checking setup status: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to check setup status"})
		return
	}

	// Return true if no user with the role exists (setup needed), false otherwise
	response := SetupStatusResponse{
		NeedsSetup: !hasUser,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AdminSetupStatus returns whether the admin system needs initial setup (no active admin user with DEFAULT_ADMIN_USER_ROLE exists).
// This is a public endpoint that does not require authentication.
func AdminSetupStatus(w http.ResponseWriter, r *http.Request) {
	service, ok := r.Context().Value(constants.AdminServiceContextValue).(*svc.Service)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(svc.ErrServiceNotSetInContext)
		return
	}

	cfg, ok := r.Context().Value("config").(*config.Config)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("config not found in context")
		json.NewEncoder(w).Encode(map[string]string{"error": "Configuration error"})
		return
	}

	// Get the default admin user role from config, default to "SuperAdmin" if not set
	roleName := cfg.DefaultAdminUserRole
	if roleName == "" {
		roleName = "SuperAdmin"
	}

	// Use internal operation context to bypass permission checks
	iCtx := context.WithValue(r.Context(), constants.IsInternalOperationContextValue, true)

	hasAdminUser, err := service.AdminUser.HasActiveAdminUserWithRole(iCtx, roleName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error checking admin setup status: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to check admin setup status"})
		return
	}

	// Return true if no admin user with the role exists (setup needed), false otherwise
	response := SetupStatusResponse{
		NeedsSetup: !hasAdminUser,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


