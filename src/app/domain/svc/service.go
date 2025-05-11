package svc

import (
	"app/domain/repo"
	"app/interfaces/svc"
	
	labraSvc "github.com/GoLabra/labrago/src/api/domain/svc"
	"github.com/GoLabra/labrago/src/api/entgql/generator"
	labraISvc "github.com/GoLabra/labrago/src/api/interfaces/svc"
)

var (
	ErrServiceNotSetInContext = "service is not set in context"
	ServiceContextValue       = "service"
)

type Service struct {
	// Node                  svc.Node
	Entity labraISvc.Entity
	Migration              svc.Migration
	Role              svc.Role
	User              svc.User
}

func New(repository *repo.Repository, schemaManager generator.SchemaManager) *Service {
	return &Service{
		// Node:                  NewNode(repository),
		Entity:                   labraSvc.NewEntity(schemaManager),
        Migration:              NewMigration(repository),
        Role:              NewRole(repository),
        User:              NewUser(repository),
	}
}
