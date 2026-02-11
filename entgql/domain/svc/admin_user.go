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

// hashPasswordAdminUserUpdate hashes password only if present in Update input (Password != nil).
func hashPasswordAdminUserUpdate(data *ent.UpdateAdminUserInput) error {
	if data.Password == nil {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcryptCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	*data.Password = string(hashedPassword)
	return nil
}

// hashPasswordAdminUserCreate hashes password in Create input (Password is required string).
func hashPasswordAdminUserCreate(data *ent.CreateAdminUserInput) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcryptCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	data.Password = string(hashedPassword)
	return nil
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
	if err := hashPasswordAdminUserCreate(&data); err != nil {
		return nil, err
	}

	createdInput, err := s.repository.AdminUser.CreateTx(ctx, tx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *AdminUser) Create(ctx context.Context, data ent.CreateAdminUserInput) (*ent.AdminUser, error) {
	if err := hashPasswordAdminUserCreate(&data); err != nil {
		return nil, err
	}

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
	if err := hashPasswordAdminUserUpdate(&data); err != nil {
		return nil, err
	}
	return s.repository.AdminUser.Update(ctx, where, data)
}

func (s *AdminUser) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereUniqueInput, data ent.UpdateAdminUserInput) (*ent.AdminUser, error) {
	if err := hashPasswordAdminUserUpdate(&data); err != nil {
		return nil, err
	}
	return s.repository.AdminUser.UpdateTx(ctx, tx, where, data)
}

func (s *AdminUser) UpdateMany(ctx context.Context, where ent.AdminUserWhereInput, data ent.UpdateAdminUserInput) (int, error) {
	if err := hashPasswordAdminUserUpdate(&data); err != nil {
		return 0, err
	}
	return s.repository.AdminUser.UpdateMany(ctx, where, data)
}

func (s *AdminUser) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereInput, data ent.UpdateAdminUserInput) (int, error) {
	if err := hashPasswordAdminUserUpdate(&data); err != nil {
		return 0, err
	}
	return s.repository.AdminUser.UpdateManyTx(ctx, tx, where, data)
}

func (s *AdminUser) Upsert(ctx context.Context, data ent.CreateAdminUserInput) (*ent.AdminUser, error) {
	if err := hashPasswordAdminUserCreate(&data); err != nil {
		return nil, err
	}
	return s.repository.AdminUser.Upsert(ctx, data)
}

func (s *AdminUser) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateAdminUserInput) (*ent.AdminUser, error) {
	if err := hashPasswordAdminUserCreate(&data); err != nil {
		return nil, err
	}
	return s.repository.AdminUser.UpsertTx(ctx, tx, data)
}

func (s *AdminUser) UpsertMany(ctx context.Context, data []ent.CreateAdminUserInput) (int, error) {
	for i := range data {
		if err := hashPasswordAdminUserCreate(&data[i]); err != nil {
			return 0, err
		}
	}
	return s.repository.AdminUser.UpsertMany(ctx, data)
}

func (s *AdminUser) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateAdminUserInput) (int, error) {
	for i := range data {
		if err := hashPasswordAdminUserCreate(&data[i]); err != nil {
			return 0, err
		}
	}
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
