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
	File              adminInterfaces.File
	ForPermission              svc.ForPermission
	Role              adminInterfaces.Role
	User              svc.User
}

func New(repository *repo.Repository, adminRepo *adminRepo.Repository) *Service {
	return &Service{
		CustomService: NewCustomService(repository),
        AdminUser:              adminSvc.NewAdminUser(adminRepo),
        File:              adminSvc.NewFile(adminRepo),
		ForPermission:              NewForPermission(repository),
        Role:              adminSvc.NewRole(adminRepo),
		User:              NewUser(repository),
	}
}
