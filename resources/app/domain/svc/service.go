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
	*CustomService
	AdminUser              adminInterfaces.AdminUser
	Cycle              svc.Cycle
	File              adminInterfaces.File
	ForPermission              svc.ForPermission
	LifeCycleNot              svc.LifeCycleNot
	Miau              svc.Miau
	Role              adminInterfaces.Role
	User              svc.User
}

func New(repository *repo.Repository, adminRepo *adminRepo.Repository) *Service {
	return &Service{
		CustomService: NewCustomService(repository),
        AdminUser:              adminSvc.NewAdminUser(adminRepo),
		Cycle:              NewCycle(repository),
        File:              adminSvc.NewFile(adminRepo),
		ForPermission:              NewForPermission(repository),
		LifeCycleNot:              NewLifeCycleNot(repository),
		Miau:              NewMiau(repository),
        Role:              adminSvc.NewRole(adminRepo),
		User:              NewUser(repository),
	}
}
