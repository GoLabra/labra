package svc

import (
	"context"

	"github.com/GoLabra/labra/entgql/domain/repo"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/GoLabra/labra/jwtrefresh"
)

type Role struct {
	repository *repo.Repository
}

func NewRole(r *repo.Repository) *Role {
	return &Role{repository: r}
}

func (s *Role) Get(ctx context.Context, where *ent.RoleWhereInput, orderBy *ent.RoleOrder, skip *int, first *int, last *int) ([]*ent.Role, error) {
	return s.repository.Role.Get(ctx, where, orderBy, skip, first, last)
}

func (s *Role) Connection(ctx context.Context, where *ent.RoleWhereInput, orderBy *ent.RoleOrder, skip *int, first *int, last *int) (*ent.RoleConnection, error) {
	return s.repository.Role.Connection(ctx, where, orderBy, skip, first, last)
}

func (s *Role) GetOne(ctx context.Context, where ent.RoleWhereUniqueInput) (*ent.Role, error) {
	return s.repository.Role.GetOne(ctx, where)
}

func (s *Role) GetTx(ctx context.Context, tx *ent.Tx, where *ent.RoleWhereInput, orderBy *ent.RoleOrder, skip *int, first *int, last *int) ([]*ent.Role, error) {
	return s.repository.Role.GetTx(ctx, tx, where, orderBy, skip, first, last)
}

func (s *Role) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereUniqueInput) (*ent.Role, error) {
	return s.repository.Role.GetOneTx(ctx, tx, where)
}

func (s *Role) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateRoleInput) (*ent.Role, error) {
	createdInput, err := s.repository.Role.CreateTx(ctx, tx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *Role) Create(ctx context.Context, data ent.CreateRoleInput) (*ent.Role, error) {
	createdInput, err := s.repository.Role.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return createdInput, err
}

func (s *Role) CreateMany(ctx context.Context, data []ent.CreateRoleInput) ([]*ent.Role, error) {
	return s.repository.Role.CreateMany(ctx, data)
}

func (s *Role) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateRoleInput) ([]*ent.Role, error) {
	return s.repository.Role.CreateManyTx(ctx, tx, data)
}

func (s *Role) Update(ctx context.Context, where ent.RoleWhereUniqueInput, data ent.UpdateRoleInput) (*ent.Role, error) {
	role, err := s.repository.Role.GetOne(ctx, where)
	if err != nil {
		return nil, err
	}

	subjects, err := s.subjectsForRoleIDs(ctx, nil, []string{role.ID})
	if err != nil {
		return nil, err
	}

	updated, err := s.repository.Role.Update(ctx, where, data)
	if err != nil {
		return nil, err
	}

	if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Role) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereUniqueInput, data ent.UpdateRoleInput) (*ent.Role, error) {
	role, err := s.repository.Role.GetOneTx(ctx, tx, where)
	if err != nil {
		return nil, err
	}

	subjects, err := s.subjectsForRoleIDs(ctx, tx, []string{role.ID})
	if err != nil {
		return nil, err
	}

	updated, err := s.repository.Role.UpdateTx(ctx, tx, where, data)
	if err != nil {
		return nil, err
	}

	if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Role) UpdateMany(ctx context.Context, where ent.RoleWhereInput, data ent.UpdateRoleInput) (int, error) {
	roles, err := s.repository.Role.Get(ctx, &where, nil, nil, nil, nil)
	if err != nil {
		return 0, err
	}

	roleIDs := make([]string, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	subjects, err := s.subjectsForRoleIDs(ctx, nil, roleIDs)
	if err != nil {
		return 0, err
	}

	updatedRows, err := s.repository.Role.UpdateMany(ctx, where, data)
	if err != nil {
		return 0, err
	}

	if updatedRows > 0 {
		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return 0, err
		}
	}

	return updatedRows, nil
}

func (s *Role) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereInput, data ent.UpdateRoleInput) (int, error) {
	roles, err := s.repository.Role.GetTx(ctx, tx, &where, nil, nil, nil, nil)
	if err != nil {
		return 0, err
	}

	roleIDs := make([]string, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	subjects, err := s.subjectsForRoleIDs(ctx, tx, roleIDs)
	if err != nil {
		return 0, err
	}

	updatedRows, err := s.repository.Role.UpdateManyTx(ctx, tx, where, data)
	if err != nil {
		return 0, err
	}

	if updatedRows > 0 {
		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return 0, err
		}
	}

	return updatedRows, nil
}

func (s *Role) Upsert(ctx context.Context, data ent.CreateRoleInput) (*ent.Role, error) {
	return s.repository.Role.Upsert(ctx, data)
}

func (s *Role) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateRoleInput) (*ent.Role, error) {
	return s.repository.Role.UpsertTx(ctx, tx, data)
}

func (s *Role) UpsertMany(ctx context.Context, data []ent.CreateRoleInput) (int, error) {
	return s.repository.Role.UpsertMany(ctx, data)
}

func (s *Role) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateRoleInput) (int, error) {
	return s.repository.Role.UpsertManyTx(ctx, tx, data)
}

func (s *Role) Delete(ctx context.Context, where ent.RoleWhereUniqueInput) (*ent.Role, error) {
	role, err := s.repository.Role.GetOne(ctx, where)
	if err != nil {
		return nil, err
	}

	subjects, err := s.subjectsForRoleIDs(ctx, nil, []string{role.ID})
	if err != nil {
		return nil, err
	}

	deleted, err := s.repository.Role.Delete(ctx, where)
	if err != nil {
		return nil, err
	}

	if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
		return nil, err
	}

	return deleted, nil
}

func (s *Role) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereUniqueInput) (*ent.Role, error) {
	role, err := s.repository.Role.GetOneTx(ctx, tx, where)
	if err != nil {
		return nil, err
	}

	subjects, err := s.subjectsForRoleIDs(ctx, tx, []string{role.ID})
	if err != nil {
		return nil, err
	}

	deleted, err := s.repository.Role.DeleteTx(ctx, tx, where)
	if err != nil {
		return nil, err
	}

	if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
		return nil, err
	}

	return deleted, nil
}

func (s *Role) DeleteMany(ctx context.Context, where ent.RoleWhereInput) (int, error) {
	roles, err := s.repository.Role.Get(ctx, &where, nil, nil, nil, nil)
	if err != nil {
		return 0, err
	}

	roleIDs := make([]string, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	subjects, err := s.subjectsForRoleIDs(ctx, nil, roleIDs)
	if err != nil {
		return 0, err
	}

	deletedRows, err := s.repository.Role.DeleteMany(ctx, where)
	if err != nil {
		return 0, err
	}

	if deletedRows > 0 {
		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return 0, err
		}
	}

	return deletedRows, nil
}

func (s *Role) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereInput) (int, error) {
	roles, err := s.repository.Role.GetTx(ctx, tx, &where, nil, nil, nil, nil)
	if err != nil {
		return 0, err
	}

	roleIDs := make([]string, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	subjects, err := s.subjectsForRoleIDs(ctx, tx, roleIDs)
	if err != nil {
		return 0, err
	}

	deletedRows, err := s.repository.Role.DeleteManyTx(ctx, tx, where)
	if err != nil {
		return 0, err
	}

	if deletedRows > 0 {
		if err := revokeSubjects(ctx, s.repository, subjects); err != nil {
			return 0, err
		}
	}

	return deletedRows, nil
}

func (s *Role) subjectsForRoleIDs(ctx context.Context, tx *ent.Tx, roleIDs []string) ([]subjectRevocation, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}

	roleFilters := make([]*ent.RoleWhereInput, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		roleID := roleID
		roleFilters = append(roleFilters, &ent.RoleWhereInput{ID: &roleID})
	}

	adminUserWhere := &ent.AdminUserWhereInput{
		Or: []*ent.AdminUserWhereInput{
			{HasDefaultRoleWith: roleFilters},
			{HasRolesWith: roleFilters},
		},
	}
	userWhere := &ent.UserWhereInput{
		Or: []*ent.UserWhereInput{
			{HasDefaultRoleWith: roleFilters},
			{HasRolesWith: roleFilters},
		},
	}

	var (
		adminUsers []*ent.AdminUser
		users      []*ent.User
		err        error
	)

	if tx == nil {
		adminUsers, err = s.repository.AdminUser.Get(ctx, adminUserWhere, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		users, err = s.repository.User.Get(ctx, userWhere, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}
	} else {
		adminUsers, err = s.repository.AdminUser.GetTx(ctx, tx, adminUserWhere, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		users, err = s.repository.User.GetTx(ctx, tx, userWhere, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}
	}

	subjects := make([]subjectRevocation, 0, len(adminUsers)+len(users))
	for _, adminUser := range adminUsers {
		subjects = append(subjects, subjectRevocation{
			Email:       adminUser.Email,
			SubjectType: jwtrefresh.SubjectTypeAdmin,
		})
	}
	for _, user := range users {
		subjects = append(subjects, subjectRevocation{
			Email:       user.Email,
			SubjectType: jwtrefresh.SubjectTypeUser,
		})
	}

	return subjects, nil
}
