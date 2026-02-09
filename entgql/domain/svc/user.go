package svc

import (
	"context"
	"fmt"

	"github.com/GoLabra/labra/entgql/domain/repo"
	"github.com/GoLabra/labra/entgql/ent"
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
	if err := hashPasswordUserUpdate(&data); err != nil {
		return nil, err
	}
	return s.repository.User.Update(ctx, where, data)
}

func (s *User) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereUniqueInput, data ent.UpdateUserInput) (*ent.User, error) {
	if err := hashPasswordUserUpdate(&data); err != nil {
		return nil, err
	}
	return s.repository.User.UpdateTx(ctx, tx, where, data)
}

func (s *User) UpdateMany(ctx context.Context, where ent.UserWhereInput, data ent.UpdateUserInput) (int, error) {
	if err := hashPasswordUserUpdate(&data); err != nil {
		return 0, err
	}
	return s.repository.User.UpdateMany(ctx, where, data)
}

func (s *User) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereInput, data ent.UpdateUserInput) (int, error) {
	if err := hashPasswordUserUpdate(&data); err != nil {
		return 0, err
	}
	return s.repository.User.UpdateManyTx(ctx, tx, where, data)
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
	return s.repository.User.Delete(ctx, where)
}

func (s *User) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereUniqueInput) (*ent.User, error) {
	return s.repository.User.DeleteTx(ctx, tx, where)
}

func (s *User) DeleteMany(ctx context.Context, where ent.UserWhereInput) (int, error) {
	return s.repository.User.DeleteMany(ctx, where)
}

func (s *User) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereInput) (int, error) {
	return s.repository.User.DeleteManyTx(ctx, tx, where)
}
