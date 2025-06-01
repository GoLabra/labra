package svc

import (
	"github.com/GoLabra/labra/src/api/entgql/domain/repo"
	"github.com/GoLabra/labra/src/api/entgql/generator"
	"github.com/GoLabra/labra/src/api/entgql/interfaces/svc"
)

var (
	ErrServiceNotSetInContext = "service is not set in context"
)

type Service struct {
	// Node                  svc.Node
	Entity     *Entity
	Permission svc.Permission
	Role       svc.Role
	User       svc.User
}

func New(repository *repo.Repository, schemaManager generator.SchemaManager) *Service {
	return &Service{
		// Node:                  NewNode(repository),
		Permission: NewPermission(repository),
		Role:       NewRole(repository),
		User:       NewUser(repository),
		Entity:     NewEntity(schemaManager),
	}
}
