package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/mitchellh/mapstructure"
)

type Permission struct {
	client *ent.Client
}

func NewPermission(c *ent.Client) *Permission {
	return &Permission{client: c}
}

func (r *Permission) Get(ctx context.Context, where *ent.PermissionWhereInput, orderBy *ent.PermissionOrder, skip *int, first *int, last *int) ([]*ent.Permission, error) {
	var (
		query = r.client.Permission.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.PermissionOrder{
			Direction: ent.DefaultPermissionOrder.Direction,
			Field:     ent.DefaultPermissionOrder.Field,
		}
	}

	if last != nil {
		orderBy.Direction = orderBy.Direction.Reverse()
	}

	query = query.Order(OrderFunc(orderBy.Direction, orderBy.Field.String()))

	if first != nil {
		query.Limit(*first)
	} else if last != nil {
		query.Limit(*last)
	}

	if skip != nil {
		query.Offset(*skip)
	}

	return query.All(ctx)
}

func (r *Permission) Connection(ctx context.Context, where *ent.PermissionWhereInput, orderBy *ent.PermissionOrder, skip *int, first *int, last *int) (*ent.PermissionConnection, error) {
	var (
		query = r.client.Permission.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy != nil {
		query = query.Order(OrderFunc(orderBy.Direction, orderBy.Field.String()))
	}

	if skip != nil {
		query.Offset(*skip)
	}

	return query.Paginate(ctx, nil, first, nil, last)
}

func (r *Permission) GetOne(ctx context.Context, where ent.PermissionWhereUniqueInput) (*ent.Permission, error) {
	query := r.client.Permission.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %w", err)
	}
	return query.First(ctx)
}

func (r *Permission) GetTx(ctx context.Context, tx *ent.Tx, where *ent.PermissionWhereInput, orderBy *ent.PermissionOrder, skip *int, first *int, last *int) ([]*ent.Permission, error) {
	var (
		query = tx.Permission.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.PermissionOrder{
			Direction: ent.DefaultPermissionOrder.Direction,
			Field:     ent.DefaultPermissionOrder.Field,
		}
	}

	if last != nil {
		orderBy.Direction = orderBy.Direction.Reverse()
	}

	query = query.Order(OrderFunc(orderBy.Direction, orderBy.Field.String()))

	if first != nil {
		query.Limit(*first)
	} else if last != nil {
		query.Limit(*last)
	}

	if skip != nil {
		query.Offset(*skip)
	}

	return query.All(ctx)
}

func (r *Permission) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereUniqueInput) (*ent.Permission, error) {
	query := tx.Permission.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %w", err)
	}
	return query.First(ctx)
}

func (r *Permission) Create(ctx context.Context, data ent.CreatePermissionInput) (*ent.Permission, error) {
	repository, ok := ctx.Value(constants.AdminRepositoryContextValue).(*Repository)

	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	createdInput, err := repository.Permission.CreateTx(ctx, tx, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return createdInput, err
}

func (r *Permission) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreatePermissionInput) (*ent.Permission, error) {
	var err error
	repository, ok := ctx.Value(constants.AdminRepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}
	if data.AdminCreatedBy != nil {

		if data.AdminCreatedBy.Connect != nil {
			toConnect, err := repository.AdminUser.GetOneTx(ctx, tx, *data.AdminCreatedBy.Connect)

			if err != nil {
				return nil, err
			}

			data.AdminCreatedByID = &toConnect.ID
		}
		if data.AdminCreatedBy.Create != nil {
			var createInput ent.CreateUserInput
			err = mapstructure.Decode(data.AdminCreatedBy.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.User.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.AdminCreatedByID = &toConnect.ID
		}
	}
	if data.AdminUpdatedBy != nil {

		if data.AdminUpdatedBy.Connect != nil {
			toConnect, err := repository.AdminUser.GetOneTx(ctx, tx, *data.AdminUpdatedBy.Connect)

			if err != nil {
				return nil, err
			}

			data.AdminUpdatedByID = &toConnect.ID
		}
		if data.AdminUpdatedBy.Create != nil {
			var createInput ent.CreateUserInput
			err = mapstructure.Decode(data.AdminUpdatedBy.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.User.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.AdminUpdatedByID = &toConnect.ID
		}
	}
	if data.Role != nil {

		if data.Role.Connect != nil {
			toConnect, err := repository.Role.GetOneTx(ctx, tx, *data.Role.Connect)

			if err != nil {
				return nil, err
			}

			data.RoleID = &toConnect.ID
		}
		if data.Role.Create != nil {
			var createInput ent.CreateRoleInput
			err = mapstructure.Decode(data.Role.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.Role.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.RoleID = &toConnect.ID
		}
	}

	createdInput, err := tx.Permission.Create().SetInput(data).Save(ctx)

	if err != nil {
		return nil, err
	}

	return createdInput.Unwrap(), err
}

func (r *Permission) CreateMany(ctx context.Context, data []ent.CreatePermissionInput) ([]*ent.Permission, error) {
	var createMap []*ent.PermissionCreate
	for _, createInput := range data {
		createMap = append(createMap, r.client.Permission.Create().SetInput(createInput))
	}
	createdItems, err := r.client.Permission.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %w", err)
	}
	return createdItems, nil
}

func (r *Permission) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreatePermissionInput) ([]*ent.Permission, error) {
	var createMap []*ent.PermissionCreate
	for _, createInput := range data {
		createMap = append(createMap, tx.Permission.Create().SetInput(createInput))
	}
	createdItems, err := tx.Permission.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %w", err)
	}
	return createdItems, nil
}

func (r *Permission) Update(ctx context.Context, where ent.PermissionWhereUniqueInput, data ent.UpdatePermissionInput) (*ent.Permission, error) {
	repository, ok := ctx.Value(constants.AdminRepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	updatedInput, err := repository.Permission.UpdateTx(ctx, tx, where, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return updatedInput.Unwrap(), nil
}

func (r *Permission) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereUniqueInput, data ent.UpdatePermissionInput) (*ent.Permission, error) {

	repository, ok := ctx.Value(constants.AdminRepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	query := tx.Permission.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %w", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to update: %w", err)
	}
	created_byToDelete := ent.UserWhereUniqueInput{}
	if data.AdminCreatedBy != nil {

		if data.AdminCreatedBy.Connect != nil {
			toConnect, err := repository.AdminUser.GetOneTx(ctx, tx, *data.AdminCreatedBy.Connect)

			if err != nil {
				return nil, err
			}

			data.AdminCreatedByID = &toConnect.ID
		}
		if data.AdminCreatedBy.Unset != nil && *data.AdminCreatedBy.Unset {
			data.ClearAdminCreatedBy = true
		}
		if data.AdminCreatedBy.Create != nil {
			var createInput ent.CreateUserInput
			err = mapstructure.Decode(data.AdminCreatedBy.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.User.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.AdminCreatedByID = &toConnect.ID
		}
		if data.AdminCreatedBy.Delete != nil && *data.AdminCreatedBy.Delete {
			itemToDelete, err := item.AdminCreatedBy(ctx)
			if err != nil {
				return nil, err
			}
			created_byToDelete.ID = &itemToDelete.ID
		}
	}
	updated_byToDelete := ent.UserWhereUniqueInput{}
	if data.AdminUpdatedBy != nil {

		if data.AdminUpdatedBy.Connect != nil {
			toConnect, err := repository.AdminUser.GetOneTx(ctx, tx, *data.AdminUpdatedBy.Connect)

			if err != nil {
				return nil, err
			}

			data.AdminUpdatedByID = &toConnect.ID
		}
		if data.AdminUpdatedBy.Unset != nil && *data.AdminUpdatedBy.Unset {
			data.ClearAdminUpdatedBy = true
		}
		if data.AdminUpdatedBy.Create != nil {
			var createInput ent.CreateUserInput
			err = mapstructure.Decode(data.AdminUpdatedBy.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.User.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.AdminUpdatedByID = &toConnect.ID
		}
		if data.AdminUpdatedBy.Delete != nil && *data.AdminUpdatedBy.Delete {
			itemToDelete, err := item.AdminUpdatedBy(ctx)
			if err != nil {
				return nil, err
			}
			updated_byToDelete.ID = &itemToDelete.ID
		}
	}
	roleToDelete := ent.RoleWhereUniqueInput{}
	if data.Role != nil {

		if data.Role.Connect != nil {
			toConnect, err := repository.Role.GetOneTx(ctx, tx, *data.Role.Connect)

			if err != nil {
				return nil, err
			}

			data.RoleID = &toConnect.ID
		}
		if data.Role.Unset != nil && *data.Role.Unset {
			data.ClearRole = true
		}
		if data.Role.Create != nil {
			var createInput ent.CreateRoleInput
			err = mapstructure.Decode(data.Role.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.Role.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.RoleID = &toConnect.ID
		}
		if data.Role.Delete != nil && *data.Role.Delete {
			itemToDelete, err := item.Role(ctx)
			if err != nil {
				return nil, err
			}
			roleToDelete.ID = &itemToDelete.ID
		}
	}

	updatedInput, err := tx.Permission.UpdateOne(item).SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error updating item: %w", err)
	}
	if data.AdminCreatedBy != nil && data.AdminCreatedBy.Delete != nil && *data.AdminCreatedBy.Delete {
		_, err := repository.User.DeleteTx(ctx, tx, created_byToDelete)
		if err != nil {
			return nil, err
		}
	}
	if data.AdminUpdatedBy != nil && data.AdminUpdatedBy.Delete != nil && *data.AdminUpdatedBy.Delete {
		_, err := repository.User.DeleteTx(ctx, tx, updated_byToDelete)
		if err != nil {
			return nil, err
		}
	}
	if data.Role != nil && data.Role.Delete != nil && *data.Role.Delete {
		_, err := repository.Role.DeleteTx(ctx, tx, roleToDelete)
		if err != nil {
			return nil, err
		}
	}

	return updatedInput, nil
}

func (r *Permission) UpdateMany(ctx context.Context, where ent.PermissionWhereInput, data ent.UpdatePermissionInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := r.client.Permission.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *Permission) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereInput, data ent.UpdatePermissionInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := tx.Permission.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *Permission) Upsert(ctx context.Context, data ent.CreatePermissionInput) (upsertedPermission *ent.Permission, err error) {
	if CompareUniqueInput(data, ent.PermissionWhereUniqueInput{}) {
		var where ent.PermissionWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %w", err)
		}
		err = r.client.Permission.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %w", err)
		}
		upsertedPermission, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %w", err)
		}
		return upsertedPermission, nil
	}
	upsertedPermission, err = r.client.Permission.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %w", err)
	}
	return upsertedPermission, nil
}

func (r *Permission) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreatePermissionInput) (upsertedPermission *ent.Permission, err error) {
	if CompareUniqueInput(data, ent.PermissionWhereUniqueInput{}) {
		var where ent.PermissionWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %w", err)
		}
		err = tx.Permission.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %w", err)
		}
		upsertedPermission, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %w", err)
		}
		return upsertedPermission, nil
	}
	upsertedPermission, err = tx.Permission.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %w", err)
	}
	return upsertedPermission, nil
}

func (r *Permission) UpsertMany(ctx context.Context, data []ent.CreatePermissionInput) (int, error) {
	var upsertMap []*ent.PermissionCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, r.client.Permission.Create().SetInput(upsertInput))
	}
	err := r.client.Permission.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %w", err)
	}
	return len(data), nil
}

func (r *Permission) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreatePermissionInput) (int, error) {
	var upsertMap []*ent.PermissionCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, tx.Permission.Create().SetInput(upsertInput))
	}
	err := tx.Permission.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %w", err)
	}
	return len(data), nil
}

func (r *Permission) Delete(ctx context.Context, where ent.PermissionWhereUniqueInput) (*ent.Permission, error) {
	query := r.client.Permission.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %w", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %w", err)
	}

	err = r.client.Permission.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %w", err)
	}
	return item, err
}

func (r *Permission) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereUniqueInput) (*ent.Permission, error) {
	query := tx.Permission.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %w", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %w", err)
	}

	err = tx.Permission.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %w", err)
	}
	return item, err
}

func (r *Permission) DeleteMany(ctx context.Context, where ent.PermissionWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := r.client.Permission.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}

func (r *Permission) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.PermissionWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := tx.Permission.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}
