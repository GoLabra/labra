package svc

import (
	"context"
	"fmt"

	"github.com/GoLabra/labra/entgql/domain/repo"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/GoLabra/labra/jwtrefresh"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	repository *repo.Repository
}

func NewUser(r *repo.Repository) *User {
	return &User{repository: r}
}

// hashPasswordUserUpdate hashes password only if present in Update input (Password != nil).
func hashPasswordUserUpdate(data *ent.UpdateUserInput) error {
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

// hashPasswordUserCreate hashes password in Create input (Password is required string).
func hashPasswordUserCreate(data *ent.CreateUserInput) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcryptCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	data.Password = string(hashedPassword)
	return nil
}

func shouldRevokeUserTokensOnUpdate(data ent.UpdateUserInput) bool {
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

func (s *User) Get(ctx context.Context, where *ent.UserWhereInput, orderBy *ent.UserOrder, skip *int, first *int, last *int) ([]*ent.User, error) {
	return s.repository.User.Get(ctx, where, orderBy, skip, first, last)
}

func (s *User) Connection(ctx context.Context, where *ent.UserWhereInput, orderBy *ent.UserOrder, skip *int, first *int, last *int) (*ent.UserConnection, error) {
	return s.repository.User.Connection(ctx, where, orderBy, skip, first, last)
}

func (s *User) GetOne(ctx context.Context, where ent.UserWhereUniqueInput) (*ent.User, error) {
	return s.repository.User.GetOne(ctx, where)
}

func (s *User) GetTx(ctx context.Context, tx *ent.Tx, where *ent.UserWhereInput, orderBy *ent.UserOrder, skip *int, first *int, last *int) ([]*ent.User, error) {
	return s.repository.User.GetTx(ctx, tx, where, orderBy, skip, first, last)
}

func (s *User) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereUniqueInput) (*ent.User, error) {
	return s.repository.User.GetOneTx(ctx, tx, where)
}

func (s *User) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateUserInput) (*ent.User, error) {
	if err := hashPasswordUserCreate(&data); err != nil {
		return nil, err
	}

	createdInput, err := s.repository.User.CreateTx(ctx, tx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *User) Create(ctx context.Context, data ent.CreateUserInput) (*ent.User, error) {
	if err := hashPasswordUserCreate(&data); err != nil {
		return nil, err
	}

	createdInput, err := s.repository.User.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *User) CreateMany(ctx context.Context, data []ent.CreateUserInput) ([]*ent.User, error) {
	return s.repository.User.CreateMany(ctx, data)
}

func (s *User) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateUserInput) ([]*ent.User, error) {
	return s.repository.User.CreateManyTx(ctx, tx, data)
}

func (s *User) Update(ctx context.Context, where ent.UserWhereUniqueInput, data ent.UpdateUserInput) (*ent.User, error) {
	shouldRevoke := shouldRevokeUserTokensOnUpdate(data)

	var before *ent.User
	var err error
	if shouldRevoke {
		before, err = s.repository.User.GetOne(ctx, where)
		if err != nil {
			return nil, err
		}
	}

	if err := hashPasswordUserUpdate(&data); err != nil {
		return nil, err
	}

	updated, err := s.repository.User.Update(ctx, where, data)
	if err != nil {
		return nil, err
	}

	if shouldRevoke {
		subjects := []subjectRevocation{{
			Email:       before.Email,
			SubjectType: jwtrefresh.SubjectTypeUser,
		}}

		if data.Email != nil {
			subjects = append(subjects, subjectRevocation{
				Email:       *data.Email,
				SubjectType: jwtrefresh.SubjectTypeUser,
			})
		}

		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return nil, err
		}
	}

	return updated, nil
}

func (s *User) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereUniqueInput, data ent.UpdateUserInput) (*ent.User, error) {
	shouldRevoke := shouldRevokeUserTokensOnUpdate(data)

	var before *ent.User
	var err error
	if shouldRevoke {
		before, err = s.repository.User.GetOneTx(ctx, tx, where)
		if err != nil {
			return nil, err
		}
	}

	if err := hashPasswordUserUpdate(&data); err != nil {
		return nil, err
	}

	updated, err := s.repository.User.UpdateTx(ctx, tx, where, data)
	if err != nil {
		return nil, err
	}

	if shouldRevoke {
		subjects := []subjectRevocation{{
			Email:       before.Email,
			SubjectType: jwtrefresh.SubjectTypeUser,
		}}

		if data.Email != nil {
			subjects = append(subjects, subjectRevocation{
				Email:       *data.Email,
				SubjectType: jwtrefresh.SubjectTypeUser,
			})
		}

		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return nil, err
		}
	}

	return updated, nil
}

func (s *User) UpdateMany(ctx context.Context, where ent.UserWhereInput, data ent.UpdateUserInput) (int, error) {
	shouldRevoke := shouldRevokeUserTokensOnUpdate(data)

	var subjects []subjectRevocation
	if shouldRevoke {
		users, err := s.repository.User.Get(ctx, &where, nil, nil, nil, nil)
		if err != nil {
			return 0, err
		}
		for _, user := range users {
			subjects = append(subjects, subjectRevocation{
				Email:       user.Email,
				SubjectType: jwtrefresh.SubjectTypeUser,
			})
		}
	}

	if err := hashPasswordUserUpdate(&data); err != nil {
		return 0, err
	}

	updatedRows, err := s.repository.User.UpdateMany(ctx, where, data)
	if err != nil {
		return 0, err
	}

	if shouldRevoke && updatedRows > 0 {
		if data.Email != nil {
			subjects = append(subjects, subjectRevocation{
				Email:       *data.Email,
				SubjectType: jwtrefresh.SubjectTypeUser,
			})
		}
		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return 0, err
		}
	}

	return updatedRows, nil
}

func (s *User) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereInput, data ent.UpdateUserInput) (int, error) {
	shouldRevoke := shouldRevokeUserTokensOnUpdate(data)

	var subjects []subjectRevocation
	if shouldRevoke {
		users, err := s.repository.User.GetTx(ctx, tx, &where, nil, nil, nil, nil)
		if err != nil {
			return 0, err
		}
		for _, user := range users {
			subjects = append(subjects, subjectRevocation{
				Email:       user.Email,
				SubjectType: jwtrefresh.SubjectTypeUser,
			})
		}
	}

	if err := hashPasswordUserUpdate(&data); err != nil {
		return 0, err
	}

	updatedRows, err := s.repository.User.UpdateManyTx(ctx, tx, where, data)
	if err != nil {
		return 0, err
	}

	if shouldRevoke && updatedRows > 0 {
		if data.Email != nil {
			subjects = append(subjects, subjectRevocation{
				Email:       *data.Email,
				SubjectType: jwtrefresh.SubjectTypeUser,
			})
		}
		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return 0, err
		}
	}

	return updatedRows, nil
}

func (s *User) Upsert(ctx context.Context, data ent.CreateUserInput) (*ent.User, error) {
	if err := hashPasswordUserCreate(&data); err != nil {
		return nil, err
	}
	return s.repository.User.Upsert(ctx, data)
}

func (s *User) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateUserInput) (*ent.User, error) {
	if err := hashPasswordUserCreate(&data); err != nil {
		return nil, err
	}
	return s.repository.User.UpsertTx(ctx, tx, data)
}

func (s *User) UpsertMany(ctx context.Context, data []ent.CreateUserInput) (int, error) {
	for i := range data {
		if err := hashPasswordUserCreate(&data[i]); err != nil {
			return 0, err
		}
	}
	return s.repository.User.UpsertMany(ctx, data)
}

func (s *User) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateUserInput) (int, error) {
	for i := range data {
		if err := hashPasswordUserCreate(&data[i]); err != nil {
			return 0, err
		}
	}
	return s.repository.User.UpsertManyTx(ctx, tx, data)
}

func (s *User) Delete(ctx context.Context, where ent.UserWhereUniqueInput) (*ent.User, error) {
	user, err := s.repository.User.GetOne(ctx, where)
	if err != nil {
		return nil, err
	}

	deleted, err := s.repository.User.Delete(ctx, where)
	if err != nil {
		return nil, err
	}

	if err := revokeSubjects(ctx, s.repository, []subjectRevocation{{
		Email:       user.Email,
		SubjectType: jwtrefresh.SubjectTypeUser,
	}}); err != nil {
		return nil, err
	}

	return deleted, nil
}

func (s *User) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereUniqueInput) (*ent.User, error) {
	user, err := s.repository.User.GetOneTx(ctx, tx, where)
	if err != nil {
		return nil, err
	}

	deleted, err := s.repository.User.DeleteTx(ctx, tx, where)
	if err != nil {
		return nil, err
	}

	if err := revokeSubjects(ctx, s.repository, []subjectRevocation{{
		Email:       user.Email,
		SubjectType: jwtrefresh.SubjectTypeUser,
	}}); err != nil {
		return nil, err
	}

	return deleted, nil
}

func (s *User) DeleteMany(ctx context.Context, where ent.UserWhereInput) (int, error) {
	users, err := s.repository.User.Get(ctx, &where, nil, nil, nil, nil)
	if err != nil {
		return 0, err
	}

	deletedRows, err := s.repository.User.DeleteMany(ctx, where)
	if err != nil {
		return 0, err
	}

	if deletedRows > 0 {
		subjects := make([]subjectRevocation, 0, len(users))
		for _, user := range users {
			subjects = append(subjects, subjectRevocation{
				Email:       user.Email,
				SubjectType: jwtrefresh.SubjectTypeUser,
			})
		}
		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return 0, err
		}
	}

	return deletedRows, nil
}

func (s *User) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereInput) (int, error) {
	users, err := s.repository.User.GetTx(ctx, tx, &where, nil, nil, nil, nil)
	if err != nil {
		return 0, err
	}

	deletedRows, err := s.repository.User.DeleteManyTx(ctx, tx, where)
	if err != nil {
		return 0, err
	}

	if deletedRows > 0 {
		subjects := make([]subjectRevocation, 0, len(users))
		for _, user := range users {
			subjects = append(subjects, subjectRevocation{
				Email:       user.Email,
				SubjectType: jwtrefresh.SubjectTypeUser,
			})
		}
		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return 0, err
		}
	}

	return deletedRows, nil
}
