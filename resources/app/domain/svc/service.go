package svc

import (
	"app/domain/repo"
	"app/interfaces/svc"

	adminRepo "github.com/GoLabra/labra/entgql/domain/repo"
	adminSvc "github.com/GoLabra/labra/entgql/domain/svc"
	adminInterfaces "github.com/GoLabra/labra/entgql/interfaces/svc"
)

var (
	ErrServiceNotSetInContext = "service is not set in context"
)

type Service struct {
	// Node                  svc.Node
	AdminUser              adminInterfaces.AdminUser
	File              adminInterfaces.File
	ForPermission              svc.ForPermission
	Role              adminInterfaces.Role
	User              svc.User
}

func New(repository *repo.Repository, adminRepo *adminRepo.Repository) *Service {
	return &Service{
		// Node:                  NewNode(repository),
		//  <no value>
		//  <no value>
        AdminUser:              adminSvc.NewAdminUser(adminRepo),
		//  <no value>
		//  <no value>
        File:              adminSvc.NewFile(adminRepo),
		//  User
		//  map[Caption:For Permission DisplayField:name Owner:User]
		ForPermission:              NewForPermission(repository),
		//  <no value>
		//  <no value>
        Role:              adminSvc.NewRole(adminRepo),
		//  User
		//  map[Caption:User DisplayField:email Owner:User]
		User:              NewUser(repository),
	}
}
