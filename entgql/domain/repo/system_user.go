package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/mitchellh/mapstructure"
)

type AdminUser struct {
	client *ent.Client
}

func NewAdminUser(c *ent.Client) *AdminUser {
	return &AdminUser{client: c}
}

func (r *AdminUser) Get(ctx context.Context, where *ent.AdminUserWhereInput, orderBy *ent.AdminUserOrder, skip *int, first *int, last *int) ([]*ent.AdminUser, error) {
	var (
		query = r.client.AdminUser.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.AdminUserOrder{
			Direction: ent.DefaultAdminUserOrder.Direction,
			Field:     ent.DefaultAdminUserOrder.Field,
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

func (r *AdminUser) Connection(ctx context.Context, where *ent.AdminUserWhereInput, orderBy *ent.AdminUserOrder, skip *int, first *int, last *int) (*ent.AdminUserConnection, error) {
	var (
		query = r.client.AdminUser.Query()
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

func (r *AdminUser) GetOne(ctx context.Context, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error) {
	query := r.client.AdminUser.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}

func (r *AdminUser) GetTx(ctx context.Context, tx *ent.Tx, where *ent.AdminUserWhereInput, orderBy *ent.AdminUserOrder, skip *int, first *int, last *int) ([]*ent.AdminUser, error) {
	var (
		query = tx.AdminUser.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.AdminUserOrder{
			Direction: ent.DefaultAdminUserOrder.Direction,
			Field:     ent.DefaultAdminUserOrder.Field,
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

func (r *AdminUser) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error) {
	query := tx.AdminUser.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}

func (r *AdminUser) Create(ctx context.Context, data ent.CreateAdminUserInput) (*ent.AdminUser, error) {
	repository, ok := ctx.Value(constants.AdminRepositoryContextValue).(*Repository)

	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	createdInput, err := repository.AdminUser.CreateTx(ctx, tx, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return createdInput, err
}

func (r *AdminUser) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateAdminUserInput) (*ent.AdminUser, error) {
	var err error
	repository, ok := ctx.Value(constants.AdminRepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}
	if data.RefAdminCreatedBy != nil {

		if data.RefAdminCreatedBy.Connect != nil {
			for _, connect := range data.RefAdminCreatedBy.Connect {
				toConnect, err := repository.AdminUser.GetOneTx(ctx, tx, *connect)

				if err != nil {
					return nil, err
				}

				data.RefAdminCreatedByIDs = append(data.RefAdminCreatedByIDs, toConnect.ID)
			}
		}
		if data.RefAdminCreatedBy.Create != nil {
			for _, create := range data.RefAdminCreatedBy.Create {
				var createInput ent.CreateAdminUserInput
				err = mapstructure.Decode(create, &createInput)
				if err != nil {
					return nil, err
				}

				toConnect, err := repository.AdminUser.CreateTx(ctx, tx, createInput)

				if err != nil {
					return nil, err
				}

				data.RefAdminCreatedByIDs = append(data.RefAdminCreatedByIDs, toConnect.ID)
			}
		}
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
			var createInput ent.CreateAdminUserInput
			err = mapstructure.Decode(data.AdminCreatedBy.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.AdminUser.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.AdminCreatedByID = &toConnect.ID
		}
	}
	if data.RefAdminUpdatedBy != nil {

		if data.RefAdminUpdatedBy.Connect != nil {
			for _, connect := range data.RefAdminUpdatedBy.Connect {
				toConnect, err := repository.AdminUser.GetOneTx(ctx, tx, *connect)

				if err != nil {
					return nil, err
				}

				data.RefAdminUpdatedByIDs = append(data.RefAdminUpdatedByIDs, toConnect.ID)
			}
		}
		if data.RefAdminUpdatedBy.Create != nil {
			for _, create := range data.RefAdminUpdatedBy.Create {
				var createInput ent.CreateAdminUserInput
				err = mapstructure.Decode(create, &createInput)
				if err != nil {
					return nil, err
				}

				toConnect, err := repository.AdminUser.CreateTx(ctx, tx, createInput)

				if err != nil {
					return nil, err
				}

				data.RefAdminUpdatedByIDs = append(data.RefAdminUpdatedByIDs, toConnect.ID)
			}
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
			var createInput ent.CreateAdminUserInput
			err = mapstructure.Decode(data.AdminUpdatedBy.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.AdminUser.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.AdminUpdatedByID = &toConnect.ID
		}
	}
	if data.Roles != nil {

		if data.Roles.Connect != nil {
			for _, connect := range data.Roles.Connect {
				toConnect, err := repository.Role.GetOneTx(ctx, tx, *connect)

				if err != nil {
					return nil, err
				}

				data.RoleIDs = append(data.RoleIDs, toConnect.ID)
			}
		}
		if data.Roles.Create != nil {
			for _, create := range data.Roles.Create {
				var createInput ent.CreateRoleInput
				err = mapstructure.Decode(create, &createInput)
				if err != nil {
					return nil, err
				}

				toConnect, err := repository.Role.CreateTx(ctx, tx, createInput)

				if err != nil {
					return nil, err
				}

				data.RoleIDs = append(data.RoleIDs, toConnect.ID)
			}
		}
	}
	if data.DefaultRole != nil {

		if data.DefaultRole.Connect != nil {
			toConnect, err := repository.Role.GetOneTx(ctx, tx, *data.DefaultRole.Connect)

			if err != nil {
				return nil, err
			}

			data.DefaultRoleID = &toConnect.ID
		}
		if data.DefaultRole.Create != nil {
			var createInput ent.CreateRoleInput
			err = mapstructure.Decode(data.DefaultRole.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.Role.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.DefaultRoleID = &toConnect.ID
		}
	}

	createdInput, err := tx.AdminUser.Create().SetInput(data).Save(ctx)

	if err != nil {
		return nil, err
	}

	return createdInput.Unwrap(), err
}

func (r *AdminUser) CreateMany(ctx context.Context, data []ent.CreateAdminUserInput) ([]*ent.AdminUser, error) {
	var createMap []*ent.AdminUserCreate
	for _, createInput := range data {
		createMap = append(createMap, r.client.AdminUser.Create().SetInput(createInput))
	}
	createdItems, err := r.client.AdminUser.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}

func (r *AdminUser) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateAdminUserInput) ([]*ent.AdminUser, error) {
	var createMap []*ent.AdminUserCreate
	for _, createInput := range data {
		createMap = append(createMap, tx.AdminUser.Create().SetInput(createInput))
	}
	createdItems, err := tx.AdminUser.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}

func (r *AdminUser) Update(ctx context.Context, where ent.AdminUserWhereUniqueInput, data ent.UpdateAdminUserInput) (*ent.AdminUser, error) {
	repository, ok := ctx.Value(constants.AdminRepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	updatedInput, err := repository.AdminUser.UpdateTx(ctx, tx, where, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return updatedInput.Unwrap(), nil
}

func (r *AdminUser) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereUniqueInput, data ent.UpdateAdminUserInput) (*ent.AdminUser, error) {

	repository, ok := ctx.Value(constants.AdminRepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	query := tx.AdminUser.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to update: %v", err)
	}
	if data.RefAdminCreatedBy != nil {

		if data.RefAdminCreatedBy.Connect != nil {
			for _, connect := range data.RefAdminCreatedBy.Connect {
				toConnect, err := repository.AdminUser.GetOneTx(ctx, tx, *connect)

				if err != nil {
					return nil, err
				}

				data.AddRefAdminCreatedByIDs = append(data.AddRefAdminCreatedByIDs, toConnect.ID)
			}
		}
		if data.RefAdminCreatedBy.Disconnect != nil {
			for _, disconnect := range data.RefAdminCreatedBy.Disconnect {
				toDisconnect, err := repository.AdminUser.GetOneTx(ctx, tx, *disconnect)

				if err != nil {
					return nil, err
				}

				data.RemoveRefAdminCreatedByIDs = append(data.RemoveRefAdminCreatedByIDs, toDisconnect.ID)
			}
		}
		if data.RefAdminCreatedBy.Create != nil {
			for _, create := range data.RefAdminCreatedBy.Create {
				var createInput ent.CreateAdminUserInput
				err = mapstructure.Decode(create, &createInput)
				if err != nil {
					return nil, err
				}

				toConnect, err := repository.AdminUser.CreateTx(ctx, tx, createInput)

				if err != nil {
					return nil, err
				}

				data.AddRefAdminCreatedByIDs = append(data.AddRefAdminCreatedByIDs, toConnect.ID)
			}
		}
		if data.RefAdminCreatedBy.Delete != nil {
			for _, delete := range data.RefAdminCreatedBy.Delete {
				_, err := repository.AdminUser.DeleteTx(ctx, tx, *delete)

				if err != nil {
					return nil, err
				}
			}
		}
	}
	created_byToDelete := ent.AdminUserWhereUniqueInput{}
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
			var createInput ent.CreateAdminUserInput
			err = mapstructure.Decode(data.AdminCreatedBy.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.AdminUser.CreateTx(ctx, tx, createInput)

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
	if data.RefAdminUpdatedBy != nil {

		if data.RefAdminUpdatedBy.Connect != nil {
			for _, connect := range data.RefAdminUpdatedBy.Connect {
				toConnect, err := repository.AdminUser.GetOneTx(ctx, tx, *connect)

				if err != nil {
					return nil, err
				}

				data.AddRefAdminUpdatedByIDs = append(data.AddRefAdminUpdatedByIDs, toConnect.ID)
			}
		}
		if data.RefAdminUpdatedBy.Disconnect != nil {
			for _, disconnect := range data.RefAdminUpdatedBy.Disconnect {
				toDisconnect, err := repository.AdminUser.GetOneTx(ctx, tx, *disconnect)

				if err != nil {
					return nil, err
				}

				data.RemoveRefAdminUpdatedByIDs = append(data.RemoveRefAdminUpdatedByIDs, toDisconnect.ID)
			}
		}
		if data.RefAdminUpdatedBy.Create != nil {
			for _, create := range data.RefAdminUpdatedBy.Create {
				var createInput ent.CreateAdminUserInput
				err = mapstructure.Decode(create, &createInput)
				if err != nil {
					return nil, err
				}

				toConnect, err := repository.AdminUser.CreateTx(ctx, tx, createInput)

				if err != nil {
					return nil, err
				}

				data.AddRefAdminUpdatedByIDs = append(data.AddRefAdminUpdatedByIDs, toConnect.ID)
			}
		}
		if data.RefAdminUpdatedBy.Delete != nil {
			for _, delete := range data.RefAdminUpdatedBy.Delete {
				_, err := repository.AdminUser.DeleteTx(ctx, tx, *delete)

				if err != nil {
					return nil, err
				}
			}
		}
	}
	updated_byToDelete := ent.AdminUserWhereUniqueInput{}
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
			var createInput ent.CreateAdminUserInput
			err = mapstructure.Decode(data.AdminUpdatedBy.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.AdminUser.CreateTx(ctx, tx, createInput)

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
	if data.Roles != nil {

		if data.Roles.Connect != nil {
			for _, connect := range data.Roles.Connect {
				toConnect, err := repository.Role.GetOneTx(ctx, tx, *connect)

				if err != nil {
					return nil, err
				}

				data.AddRoleIDs = append(data.AddRoleIDs, toConnect.ID)
			}
		}
		if data.Roles.Disconnect != nil {
			for _, disconnect := range data.Roles.Disconnect {
				toDisconnect, err := repository.Role.GetOneTx(ctx, tx, *disconnect)

				if err != nil {
					return nil, err
				}

				data.RemoveRoleIDs = append(data.RemoveRoleIDs, toDisconnect.ID)
			}
		}
		if data.Roles.Create != nil {
			for _, create := range data.Roles.Create {
				var createInput ent.CreateRoleInput
				err = mapstructure.Decode(create, &createInput)
				if err != nil {
					return nil, err
				}

				toConnect, err := repository.Role.CreateTx(ctx, tx, createInput)

				if err != nil {
					return nil, err
				}

				data.AddRoleIDs = append(data.AddRoleIDs, toConnect.ID)
			}
		}
		if data.Roles.Delete != nil {
			for _, delete := range data.Roles.Delete {
				_, err := repository.Role.DeleteTx(ctx, tx, *delete)

				if err != nil {
					return nil, err
				}
			}
		}
	}
	default_roleToDelete := ent.RoleWhereUniqueInput{}
	if data.DefaultRole != nil {

		if data.DefaultRole.Connect != nil {
			toConnect, err := repository.Role.GetOneTx(ctx, tx, *data.DefaultRole.Connect)

			if err != nil {
				return nil, err
			}

			data.DefaultRoleID = &toConnect.ID
		}
		if data.DefaultRole.Unset != nil && *data.DefaultRole.Unset {
			data.ClearDefaultRole = true
		}
		if data.DefaultRole.Create != nil {
			var createInput ent.CreateRoleInput
			err = mapstructure.Decode(data.DefaultRole.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.Role.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.DefaultRoleID = &toConnect.ID
		}
		if data.DefaultRole.Delete != nil && *data.DefaultRole.Delete {
			itemToDelete, err := item.DefaultRole(ctx)
			if err != nil {
				return nil, err
			}
			default_roleToDelete.ID = &itemToDelete.ID
		}
	}

	updatedInput, err := tx.AdminUser.UpdateOne(item).SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error updating item: %v", err)
	}
	if data.AdminCreatedBy != nil && data.AdminCreatedBy.Delete != nil && *data.AdminCreatedBy.Delete {
		_, err := repository.AdminUser.DeleteTx(ctx, tx, created_byToDelete)
		if err != nil {
			return nil, err
		}
	}
	if data.AdminUpdatedBy != nil && data.AdminUpdatedBy.Delete != nil && *data.AdminUpdatedBy.Delete {
		_, err := repository.AdminUser.DeleteTx(ctx, tx, updated_byToDelete)
		if err != nil {
			return nil, err
		}
	}
	if data.DefaultRole != nil && data.DefaultRole.Delete != nil && *data.DefaultRole.Delete {
		_, err := repository.Role.DeleteTx(ctx, tx, default_roleToDelete)
		if err != nil {
			return nil, err
		}
	}

	return updatedInput, nil
}

func (r *AdminUser) UpdateMany(ctx context.Context, where ent.AdminUserWhereInput, data ent.UpdateAdminUserInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := r.client.AdminUser.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *AdminUser) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereInput, data ent.UpdateAdminUserInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := tx.AdminUser.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *AdminUser) Upsert(ctx context.Context, data ent.CreateAdminUserInput) (upsertedAdminUser *ent.AdminUser, err error) {
	if CompareUniqueInput(data, ent.AdminUserWhereUniqueInput{}) {
		var where ent.AdminUserWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}
		err = r.client.AdminUser.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		upsertedAdminUser, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		return upsertedAdminUser, nil
	}
	upsertedAdminUser, err = r.client.AdminUser.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedAdminUser, nil
}

func (r *AdminUser) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateAdminUserInput) (upsertedAdminUser *ent.AdminUser, err error) {
	if CompareUniqueInput(data, ent.AdminUserWhereUniqueInput{}) {
		var where ent.AdminUserWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}
		err = tx.AdminUser.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		upsertedAdminUser, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		return upsertedAdminUser, nil
	}
	upsertedAdminUser, err = tx.AdminUser.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedAdminUser, nil
}

func (r *AdminUser) UpsertMany(ctx context.Context, data []ent.CreateAdminUserInput) (int, error) {
	var upsertMap []*ent.AdminUserCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, r.client.AdminUser.Create().SetInput(upsertInput))
	}
	err := r.client.AdminUser.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *AdminUser) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateAdminUserInput) (int, error) {
	var upsertMap []*ent.AdminUserCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, tx.AdminUser.Create().SetInput(upsertInput))
	}
	err := tx.AdminUser.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *AdminUser) Delete(ctx context.Context, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error) {
	query := r.client.AdminUser.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = r.client.AdminUser.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *AdminUser) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereUniqueInput) (*ent.AdminUser, error) {
	query := tx.AdminUser.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = tx.AdminUser.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *AdminUser) DeleteMany(ctx context.Context, where ent.AdminUserWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := r.client.AdminUser.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}

func (r *AdminUser) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.AdminUserWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := tx.AdminUser.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}

func (r *AdminUser) UpdatePassword(ctx context.Context, where ent.AdminUserWhereUniqueInput, hashedPassword string) (*ent.AdminUser, error) {
	iCtx := context.WithValue(ctx, constants.IsInternalOperationContextValue, true)
	query := r.client.AdminUser.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(iCtx)
	if err != nil {
		return nil, fmt.Errorf("error getting admin user: %v", err)
	}
	return r.client.AdminUser.UpdateOne(item).SetPassword(hashedPassword).Save(ctx)
}
