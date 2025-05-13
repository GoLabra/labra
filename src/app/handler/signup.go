package handler

import (
	"app/domain/svc"
	"app/ent"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/GoLabra/labra/src/api/constants"
	"golang.org/x/crypto/bcrypt"
)

type SignupFormData struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var superAdmin = "SuperAdmin"

// TODO: 1. sanitize error messages; 2. move to api; 3. add logs;
func Signup(w http.ResponseWriter, r *http.Request) {
	var (
		signupFormData SignupFormData
		superAdminRole *ent.Role
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

	service, ok := r.Context().Value(constants.ServiceContextValue).(*svc.Service)
	if !ok {
		fmt.Printf("error getting service: %v", err)
		writeErrorResponse(w, "unable to get service", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupFormData.Password), 14)

	if err != nil {
		fmt.Printf("error hashing password: %v", err)
		writeErrorResponse(w, "unable to hash password", http.StatusInternalServerError)
		return
	}

	superAdminRole, err = service.Role.GetOne(r.Context(), ent.RoleWhereUniqueInput{Name: &superAdmin})
	if err != nil && !ent.IsNotFound(err) {
		fmt.Printf("error getting super admin role: %v", err)
		writeErrorResponse(w, "unable to get super admin role", http.StatusInternalServerError)
		return
	} else if err != nil && ent.IsNotFound(err) {
		superAdminRole, err = service.Role.Create(r.Context(), ent.CreateRoleInput{
			Name: superAdmin,
		})
		if err != nil {
			fmt.Printf("error creating super admin role: %v", err)
			writeErrorResponse(w, "unable to create super admin role", http.StatusInternalServerError)
			return
		}
	}

	user, err := service.User.Create(r.Context(), ent.CreateUserInput{
		Email:     signupFormData.Email,
		Password:  string(hashedPassword),
		FirstName: signupFormData.FirstName,
		LastName:  signupFormData.LastName,
		Roles: &ent.CreateManyRoleInput{
			Connect: []*ent.RoleWhereUniqueInput{
				{
					ID: &superAdminRole.ID,
				},
			},
		},
		DefaultRole: &ent.CreateOneRoleInput{
			Connect: &ent.RoleWhereUniqueInput{
				ID: &superAdminRole.ID,
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
