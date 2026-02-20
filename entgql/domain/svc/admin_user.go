package svc

import (
	"context"
	"fmt"

	"github.com/GoLabra/labra/entgql/domain/repo"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/GoLabra/labra/jwtrefresh"
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

func shouldRevokeAdminUserTokensOnUpdate(data ent.UpdateAdminUserInput) bool {
	if data.Password != nil {
		return true
	}

	// Any role assignment/default-role mutation must invalidate active access tokens.
	if data.Roles != nil || data.DefaultRole != nil {
		return true
	}
	if data.DefaultRoleID != nil || data.ClearDefaultRole {
		return true
	}
	if data.ClearRoles || len(data.AddRoleIDs) > 0 || len(data.RemoveRoleIDs) > 0 {
		return true
	}

	return false
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
	shouldRevoke := shouldRevokeAdminUserTokensOnUpdate(data)

	var before *ent.AdminUser
	var err error
	if shouldRevoke {
		before, err = s.repository.AdminUser.GetOne(ctx, where)
		if err != nil {
			return nil, err
		}
	}

	if err := hashPasswordAdminUserUpdate(&data); err != nil {
		return nil, err
	}

	updated, err := s.repository.AdminUser.Update(ctx, where, data)
	if err != nil {
		return nil, err
	}

	if shouldRevoke {
		subjects := []subjectRevocation{{
			Email:       before.Email,
			SubjectType: jwtrefresh.SubjectTypeAdmin,
		}}

		if data.Email != nil {
			subjects = append(subjects, subjectRevocation{
				Email:       *data.Email,
				SubjectType: jwtrefresh.SubjectTypeAdmin,
			})
		}

		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return nil, err
		}
	}

	return updated, nil
}

func (s *AdminUser) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereUniqueInput, data ent.UpdateAdminUserInput) (*ent.AdminUser, error) {
	shouldRevoke := shouldRevokeAdminUserTokensOnUpdate(data)

	var before *ent.AdminUser
	var err error
	if shouldRevoke {
		before, err = s.repository.AdminUser.GetOneTx(ctx, tx, where)
		if err != nil {
			return nil, err
		}
	}

	if err := hashPasswordAdminUserUpdate(&data); err != nil {
		return nil, err
	}

	updated, err := s.repository.AdminUser.UpdateTx(ctx, tx, where, data)
	if err != nil {
		return nil, err
	}

	if shouldRevoke {
		subjects := []subjectRevocation{{
			Email:       before.Email,
			SubjectType: jwtrefresh.SubjectTypeAdmin,
		}}

		if data.Email != nil {
			subjects = append(subjects, subjectRevocation{
				Email:       *data.Email,
				SubjectType: jwtrefresh.SubjectTypeAdmin,
			})
		}

		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return nil, err
		}
	}

	return updated, nil
}

func (s *AdminUser) UpdateMany(ctx context.Context, where ent.AdminUserWhereInput, data ent.UpdateAdminUserInput) (int, error) {
	shouldRevoke := shouldRevokeAdminUserTokensOnUpdate(data)

	var subjects []subjectRevocation
	var err error
	if shouldRevoke {
		users, getErr := s.repository.AdminUser.Get(ctx, &where, nil, nil, nil, nil)
		if getErr != nil {
			return 0, getErr
		}
		for _, user := range users {
			subjects = append(subjects, subjectRevocation{
				Email:       user.Email,
				SubjectType: jwtrefresh.SubjectTypeAdmin,
			})
		}
	}

	if err := hashPasswordAdminUserUpdate(&data); err != nil {
		return 0, err
	}

	updatedRows, err := s.repository.AdminUser.UpdateMany(ctx, where, data)
	if err != nil {
		return 0, err
	}

	if shouldRevoke && updatedRows > 0 {
		if data.Email != nil {
			subjects = append(subjects, subjectRevocation{
				Email:       *data.Email,
				SubjectType: jwtrefresh.SubjectTypeAdmin,
			})
		}
		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return 0, err
		}
	}

	return updatedRows, nil
}

func (s *AdminUser) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereInput, data ent.UpdateAdminUserInput) (int, error) {
	shouldRevoke := shouldRevokeAdminUserTokensOnUpdate(data)

	var subjects []subjectRevocation
	if shouldRevoke {
		users, err := s.repository.AdminUser.GetTx(ctx, tx, &where, nil, nil, nil, nil)
		if err != nil {
			return 0, err
		}
		for _, user := range users {
			subjects = append(subjects, subjectRevocation{
				Email:       user.Email,
				SubjectType: jwtrefresh.SubjectTypeAdmin,
			})
		}
	}

	if err := hashPasswordAdminUserUpdate(&data); err != nil {
		return 0, err
	}

	updatedRows, err := s.repository.AdminUser.UpdateManyTx(ctx, tx, where, data)
	if err != nil {
		return 0, err
	}

	if shouldRevoke && updatedRows > 0 {
		if data.Email != nil {
			subjects = append(subjects, subjectRevocation{
				Email:       *data.Email,
				SubjectType: jwtrefresh.SubjectTypeAdmin,
			})
		}
		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return 0, err
		}
	}

	return updatedRows, nil
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

func (s *AdminUser) UpdatePassword(ctx context.Context, where ent.AdminUserWhereUniqueInput, password string) (*ent.AdminUser, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}
	return s.repository.AdminUser.UpdatePassword(ctx, where, string(hashedPassword))
}

func (s *AdminUser) Delete(ctx context.Context, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error) {
	user, err := s.repository.AdminUser.GetOne(ctx, where)
	if err != nil {
		return nil, err
	}

	deleted, err := s.repository.AdminUser.Delete(ctx, where)
	if err != nil {
		return nil, err
	}

	if err := revokeSubjects(ctx, s.repository, []subjectRevocation{{
		Email:       user.Email,
		SubjectType: jwtrefresh.SubjectTypeAdmin,
	}}); err != nil {
		return nil, err
	}

	return deleted, nil
}

func (s *AdminUser) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error) {
	user, err := s.repository.AdminUser.GetOneTx(ctx, tx, where)
	if err != nil {
		return nil, err
	}

	deleted, err := s.repository.AdminUser.DeleteTx(ctx, tx, where)
	if err != nil {
		return nil, err
	}

	if err := revokeSubjects(ctx, s.repository, []subjectRevocation{{
		Email:       user.Email,
		SubjectType: jwtrefresh.SubjectTypeAdmin,
	}}); err != nil {
		return nil, err
	}

	return deleted, nil
}

func (s *AdminUser) DeleteMany(ctx context.Context, where ent.AdminUserWhereInput) (int, error) {
	users, err := s.repository.AdminUser.Get(ctx, &where, nil, nil, nil, nil)
	if err != nil {
		return 0, err
	}

	deletedRows, err := s.repository.AdminUser.DeleteMany(ctx, where)
	if err != nil {
		return 0, err
	}

	if deletedRows > 0 {
		subjects := make([]subjectRevocation, 0, len(users))
		for _, user := range users {
			subjects = append(subjects, subjectRevocation{
				Email:       user.Email,
				SubjectType: jwtrefresh.SubjectTypeAdmin,
			})
		}
		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return 0, err
		}
	}

	return deletedRows, nil
}

func (s *AdminUser) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereInput) (int, error) {
	users, err := s.repository.AdminUser.GetTx(ctx, tx, &where, nil, nil, nil, nil)
	if err != nil {
		return 0, err
	}

	deletedRows, err := s.repository.AdminUser.DeleteManyTx(ctx, tx, where)
	if err != nil {
		return 0, err
	}

	if deletedRows > 0 {
		subjects := make([]subjectRevocation, 0, len(users))
		for _, user := range users {
			subjects = append(subjects, subjectRevocation{
				Email:       user.Email,
				SubjectType: jwtrefresh.SubjectTypeAdmin,
			})
		}
		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return 0, err
		}
	}

	return deletedRows, nil
}
