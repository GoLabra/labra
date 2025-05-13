package svc

import (
	"app/domain/repo"
	"app/interfaces/svc"
	
	"github.com/GoLabra/labra/src/api/entgql/generator"
)

var (
	ErrServiceNotSetInContext = "service is not set in context"
	ServiceContextValue       = "service"
)

type Service struct {
	// Node                  svc.Node
	Permission              svc.Permission
	Role              svc.Role
	User              svc.User
}

func New(repository *repo.Repository, schemaManager generator.SchemaManager) *Service {
	return &Service{
		// Node:                  NewNode(repository),
        Permission:              NewPermission(repository),
        Role:              NewRole(repository),
        User:              NewUser(repository),
	}
}
