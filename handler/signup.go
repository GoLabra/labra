package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/domain/svc"
	"github.com/GoLabra/labra/entgql/ent"
)

type SignupFormData struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// TODO: 1. sanitize error messages; 2. move to api; 3. add logs;
func Signup(w http.ResponseWriter, r *http.Request) {
	var (
		signupFormData SignupFormData
		defaultRole    *ent.Role
	)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error reading body: %v", err)
		writeErrorResponse(w, "unable to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &signupFormData)
	if err != nil {
		fmt.Printf("error unmarshaling body: %v", err)
		writeErrorResponse(w, "unable to unmarshal body", http.StatusInternalServerError)
		return
	}

	service, ok := r.Context().Value(constants.AdminServiceContextValue).(*svc.Service)
	if !ok {
		fmt.Printf("error getting service: %v", err)
		writeErrorResponse(w, "unable to get service", http.StatusInternalServerError)
		return
	}

	cfg, ok := r.Context().Value("config").(*config.Config)
	if !ok {
		fmt.Printf("error getting config")
		writeErrorResponse(w, "unable to get config", http.StatusInternalServerError)
		return
	}

	// Get the default admin user role from config, default to "SuperAdmin" if not set
	roleName := cfg.DefaultAdminUserRole
	if roleName == "" {
		roleName = "SuperAdmin"
	}

	// Use internal context for role lookup (bypasses permission checks)
	iCtx := context.WithValue(r.Context(), constants.IsInternalOperationContextValue, true)
	defaultRole, err = service.Role.GetOne(iCtx, ent.RoleWhereUniqueInput{Name: &roleName})
	if err != nil && !ent.IsNotFound(err) {
		fmt.Printf("error getting default admin role: %v", err)
		writeErrorResponse(w, fmt.Sprintf("unable to get default admin role: %s", roleName), http.StatusInternalServerError)
		return
	} else if err != nil && ent.IsNotFound(err) {
		defaultRole, err = service.Role.Create(iCtx, ent.CreateRoleInput{
			Name: roleName,
		})
		if err != nil {
			fmt.Printf("error creating default admin role: %v", err)
			writeErrorResponse(w, fmt.Sprintf("unable to create default admin role: %s", roleName), http.StatusInternalServerError)
			return
		}
	}

	user, err := service.AdminUser.Create(iCtx, ent.CreateAdminUserInput{
		Email:     signupFormData.Email,
		Password:  signupFormData.Password,
		FirstName: signupFormData.FirstName,
		LastName:  signupFormData.LastName,
		Roles: &ent.CreateManyRoleInput{
			Connect: []*ent.RoleWhereUniqueInput{
				{
					ID: &defaultRole.ID,
				},
			},
		},
		DefaultRole: &ent.CreateOneRoleInput{
			Connect: &ent.RoleWhereUniqueInput{
				ID: &defaultRole.ID,
			},
		},
	})

	if err != nil || user == nil {
		fmt.Printf("error creating user: %v", err)
		writeErrorResponse(w, "unable to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
