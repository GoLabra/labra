package handler

import (
	"github.com/GoLabra/labra/entgql/ent"
)

// CreateTestAdminUser creates a test AdminUser entity
func CreateTestAdminUser(id, email, password string) *ent.AdminUser {
	return &ent.AdminUser{
		ID:       id,
		Email:    email,
		Password: password,
	}
}

// CreateTestRole creates a test Role entity
func CreateTestRole(id, name string) *ent.Role {
	return &ent.Role{
		ID:   id,
		Name: name,
	}
}

// CreateTestAdminUserWithDefaultRole creates a test AdminUser with DefaultRole in Edges
func CreateTestAdminUserWithDefaultRole(user *ent.AdminUser, role *ent.Role) *ent.AdminUser {
	// Set the DefaultRole in the Edges so DefaultRole() method works
	user.Edges.DefaultRole = role
	return user
}
