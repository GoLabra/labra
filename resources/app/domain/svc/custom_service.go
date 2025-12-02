package svc

import "app/domain/repo"

type CustomService struct {
}

func NewCustomService(repository *repo.Repository) *CustomService {
	return &CustomService{}
}
