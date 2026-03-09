package svc

import (
	"github.com/GoLabra/labra/entgql/domain/repo"
	"github.com/GoLabra/labra/entgql/generator"
	"github.com/GoLabra/labra/entgql/interfaces/svc"
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
	AdminUser  svc.AdminUser
	File       svc.File
	Cron       svc.Cron
}

func New(repository *repo.Repository, schemaManager generator.SchemaManager) *Service {
	return &Service{
		// Node:                  NewNode(repository),
		Permission: NewPermission(repository),
		Role:       NewRole(repository),
		User:       NewUser(repository),
		AdminUser:  NewAdminUser(repository),
		Entity:     NewEntity(schemaManager),
		File:       NewFile(repository),
		Cron:       NewCron(repository),
	}
}
