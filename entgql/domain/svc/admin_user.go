package svc

import (
	"context"
	"fmt"

	"github.com/GoLabra/labra/entgql/domain/repo"
	"github.com/GoLabra/labra/entgql/ent"
	"golang.org/x/crypto/bcrypt"
)

type AdminUser struct {
	repository *repo.Repository
}

func NewAdminUser(r *repo.Repository) *AdminUser {
	return &AdminUser{repository: r}
}

func (s *AdminUser) Get(ctx context.Context, where *ent.AdminUserWhereInput, orderBy *ent.AdminUserOrder, skip *int, first *int, last *int) ([]*ent.AdminUser, error) {
	return s.repository.AdminUser.Get(ctx, where, orderBy, skip, first, last)
}

func (s *AdminUser) Connection(ctx context.Context, where *ent.AdminUserWhereInput, orderBy *ent.AdminUserOrder, skip *int, first *int, last *int) (*ent.AdminUserConnection, error) {
	return s.repository.AdminUser.Connection(ctx, where, orderBy, skip, first, last)
}

func (s *AdminUser) GetOne(ctx context.Context, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error) {
	return s.repository.AdminUser.GetOne(ctx, where)
}

func (s *AdminUser) GetTx(ctx context.Context, tx *ent.Tx, where *ent.AdminUserWhereInput, orderBy *ent.AdminUserOrder, skip *int, first *int, last *int) ([]*ent.AdminUser, error) {
	return s.repository.AdminUser.GetTx(ctx, tx, where, orderBy, skip, first, last)
}

func (s *AdminUser) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error) {
	return s.repository.AdminUser.GetOneTx(ctx, tx, where)
}

func (s *AdminUser) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateAdminUserInput) (*ent.AdminUser, error) {
	createdInput, err := s.repository.AdminUser.CreateTx(ctx, tx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *AdminUser) Create(ctx context.Context, data ent.CreateAdminUserInput) (*ent.AdminUser, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)

	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	data.Password = string(hashedPassword)

	createdInput, err := s.repository.AdminUser.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *AdminUser) CreateMany(ctx context.Context, data []ent.CreateAdminUserInput) ([]*ent.AdminUser, error) {
	return s.repository.AdminUser.CreateMany(ctx, data)
}

func (s *AdminUser) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateAdminUserInput) ([]*ent.AdminUser, error) {
	return s.repository.AdminUser.CreateManyTx(ctx, tx, data)
}

func (s *AdminUser) Update(ctx context.Context, where ent.AdminUserWhereUniqueInput, data ent.UpdateAdminUserInput) (*ent.AdminUser, error) {
	return s.repository.AdminUser.Update(ctx, where, data)
}

func (s *AdminUser) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereUniqueInput, data ent.UpdateAdminUserInput) (*ent.AdminUser, error) {
	return s.repository.AdminUser.UpdateTx(ctx, tx, where, data)
}

func (s *AdminUser) UpdateMany(ctx context.Context, where ent.AdminUserWhereInput, data ent.UpdateAdminUserInput) (int, error) {
	return s.repository.AdminUser.UpdateMany(ctx, where, data)
}

func (s *AdminUser) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereInput, data ent.UpdateAdminUserInput) (int, error) {
	return s.repository.AdminUser.UpdateManyTx(ctx, tx, where, data)
}

func (s *AdminUser) Upsert(ctx context.Context, data ent.CreateAdminUserInput) (*ent.AdminUser, error) {
	return s.repository.AdminUser.Upsert(ctx, data)
}

func (s *AdminUser) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateAdminUserInput) (*ent.AdminUser, error) {
	return s.repository.AdminUser.UpsertTx(ctx, tx, data)
}

func (s *AdminUser) UpsertMany(ctx context.Context, data []ent.CreateAdminUserInput) (int, error) {
	return s.repository.AdminUser.UpsertMany(ctx, data)
}

func (s *AdminUser) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateAdminUserInput) (int, error) {
	return s.repository.AdminUser.UpsertManyTx(ctx, tx, data)
}

func (s *AdminUser) Delete(ctx context.Context, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error) {
	return s.repository.AdminUser.Delete(ctx, where)
}

func (s *AdminUser) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error) {
	return s.repository.AdminUser.DeleteTx(ctx, tx, where)
}

func (s *AdminUser) DeleteMany(ctx context.Context, where ent.AdminUserWhereInput) (int, error) {
	return s.repository.AdminUser.DeleteMany(ctx, where)
}

func (s *AdminUser) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereInput) (int, error) {
	return s.repository.AdminUser.DeleteManyTx(ctx, tx, where)
}

// HasActiveAdminUserWithRole checks if there is at least one active admin user with the specified role in the system.
// Returns true if an admin user with the role exists, false otherwise.
func (s *AdminUser) HasActiveAdminUserWithRole(ctx context.Context, roleName string) (bool, error) {
	if roleName == "" {
		return false, fmt.Errorf("role name cannot be empty")
	}

	where := &ent.AdminUserWhereInput{
		HasRolesWith: []*ent.RoleWhereInput{
			{
				Name: &roleName,
			},
		},
	}

	limit := 1
	adminUsers, err := s.Get(ctx, where, nil, nil, &limit, nil)
	if err != nil {
		return false, fmt.Errorf("error checking for admin users with role %s: %w", roleName, err)
	}

	return len(adminUsers) > 0, nil
}
