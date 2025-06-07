package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoLabra/labra/src/api/constants"
	"github.com/GoLabra/labra/src/api/entgql/ent"
	"github.com/mitchellh/mapstructure"
)

type User struct {
	client *ent.Client
}

func NewUser(c *ent.Client) *User {
	return &User{client: c}
}

func (r *User) Get(ctx context.Context, where *ent.UserWhereInput, orderBy *ent.UserOrder, skip *int, first *int, last *int) ([]*ent.User, error) {
	var (
		query = r.client.User.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.UserOrder{
			Direction: ent.DefaultUserOrder.Direction,
			Field:     ent.DefaultUserOrder.Field,
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

func (r *User) Connection(ctx context.Context, where *ent.UserWhereInput, orderBy *ent.UserOrder, skip *int, first *int, last *int) (*ent.UserConnection, error) {
	var (
		query = r.client.User.Query()
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

func (r *User) GetOne(ctx context.Context, where ent.UserWhereUniqueInput) (*ent.User, error) {
	query := r.client.User.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}

func (r *User) GetTx(ctx context.Context, tx *ent.Tx, where *ent.UserWhereInput, orderBy *ent.UserOrder, skip *int, first *int, last *int) ([]*ent.User, error) {
	var (
		query = tx.User.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.UserOrder{
			Direction: ent.DefaultUserOrder.Direction,
			Field:     ent.DefaultUserOrder.Field,
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

func (r *User) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereUniqueInput) (*ent.User, error) {
	query := tx.User.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}

func (r *User) Create(ctx context.Context, data ent.CreateUserInput) (*ent.User, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)

	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	createdInput, err := repository.User.CreateTx(ctx, tx, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return createdInput, err
}

func (r *User) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateUserInput) (*ent.User, error) {
	var err error
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}
	if data.RefCreatedBy != nil {

		if data.RefCreatedBy.Connect != nil {
			for _, connect := range data.RefCreatedBy.Connect {
				toConnect, err := repository.User.GetOneTx(ctx, tx, *connect)

				if err != nil {
					return nil, err
				}

				data.RefCreatedByIDs = append(data.RefCreatedByIDs, toConnect.ID)
			}
		}
		if data.RefCreatedBy.Create != nil {
			for _, create := range data.RefCreatedBy.Create {
				var createInput ent.CreateUserInput
				err = mapstructure.Decode(create, &createInput)
				if err != nil {
					return nil, err
				}

				toConnect, err := repository.User.CreateTx(ctx, tx, createInput)

				if err != nil {
					return nil, err
				}

				data.RefCreatedByIDs = append(data.RefCreatedByIDs, toConnect.ID)
			}
		}
	}
	if data.CreatedBy != nil {

		if data.CreatedBy.Connect != nil {
			toConnect, err := repository.User.GetOneTx(ctx, tx, *data.CreatedBy.Connect)

			if err != nil {
				return nil, err
			}

			data.CreatedByID = &toConnect.ID
		}
		if data.CreatedBy.Create != nil {
			var createInput ent.CreateUserInput
			err = mapstructure.Decode(data.CreatedBy.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.User.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.CreatedByID = &toConnect.ID
		}
	}
	if data.RefUpdatedBy != nil {

		if data.RefUpdatedBy.Connect != nil {
			for _, connect := range data.RefUpdatedBy.Connect {
				toConnect, err := repository.User.GetOneTx(ctx, tx, *connect)

				if err != nil {
					return nil, err
				}

				data.RefUpdatedByIDs = append(data.RefUpdatedByIDs, toConnect.ID)
			}
		}
		if data.RefUpdatedBy.Create != nil {
			for _, create := range data.RefUpdatedBy.Create {
				var createInput ent.CreateUserInput
				err = mapstructure.Decode(create, &createInput)
				if err != nil {
					return nil, err
				}

				toConnect, err := repository.User.CreateTx(ctx, tx, createInput)

				if err != nil {
					return nil, err
				}

				data.RefUpdatedByIDs = append(data.RefUpdatedByIDs, toConnect.ID)
			}
		}
	}
	if data.UpdatedBy != nil {

		if data.UpdatedBy.Connect != nil {
			toConnect, err := repository.User.GetOneTx(ctx, tx, *data.UpdatedBy.Connect)

			if err != nil {
				return nil, err
			}

			data.UpdatedByID = &toConnect.ID
		}
		if data.UpdatedBy.Create != nil {
			var createInput ent.CreateUserInput
			err = mapstructure.Decode(data.UpdatedBy.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.User.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.UpdatedByID = &toConnect.ID
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

	createdInput, err := tx.User.Create().SetInput(data).Save(ctx)

	if err != nil {
		return nil, err
	}

	return createdInput.Unwrap(), err
}

func (r *User) CreateMany(ctx context.Context, data []ent.CreateUserInput) ([]*ent.User, error) {
	var createMap []*ent.UserCreate
	for _, createInput := range data {
		createMap = append(createMap, r.client.User.Create().SetInput(createInput))
	}
	createdItems, err := r.client.User.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}

func (r *User) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateUserInput) ([]*ent.User, error) {
	var createMap []*ent.UserCreate
	for _, createInput := range data {
		createMap = append(createMap, tx.User.Create().SetInput(createInput))
	}
	createdItems, err := tx.User.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}

func (r *User) Update(ctx context.Context, where ent.UserWhereUniqueInput, data ent.UpdateUserInput) (*ent.User, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	updatedInput, err := repository.User.UpdateTx(ctx, tx, where, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return updatedInput.Unwrap(), nil
}

func (r *User) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereUniqueInput, data ent.UpdateUserInput) (*ent.User, error) {

	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	query := tx.User.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to update: %v", err)
	}
	if data.RefCreatedBy != nil {

		if data.RefCreatedBy.Connect != nil {
			for _, connect := range data.RefCreatedBy.Connect {
				toConnect, err := repository.User.GetOneTx(ctx, tx, *connect)

				if err != nil {
					return nil, err
				}

				data.AddRefCreatedByIDs = append(data.AddRefCreatedByIDs, toConnect.ID)
			}
		}
		if data.RefCreatedBy.Disconnect != nil {
			for _, disconnect := range data.RefCreatedBy.Disconnect {
				toDisconnect, err := repository.User.GetOneTx(ctx, tx, *disconnect)

				if err != nil {
					return nil, err
				}

				data.RemoveRefCreatedByIDs = append(data.RemoveRefCreatedByIDs, toDisconnect.ID)
			}
		}
		if data.RefCreatedBy.Create != nil {
			for _, create := range data.RefCreatedBy.Create {
				var createInput ent.CreateUserInput
				err = mapstructure.Decode(create, &createInput)
				if err != nil {
					return nil, err
				}

				toConnect, err := repository.User.CreateTx(ctx, tx, createInput)

				if err != nil {
					return nil, err
				}

				data.AddRefCreatedByIDs = append(data.AddRefCreatedByIDs, toConnect.ID)
			}
		}
		if data.RefCreatedBy.Delete != nil {
			for _, delete := range data.RefCreatedBy.Delete {
				_, err := repository.User.DeleteTx(ctx, tx, *delete)

				if err != nil {
					return nil, err
				}
			}
		}
	}
	created_byToDelete := ent.UserWhereUniqueInput{}
	if data.CreatedBy != nil {

		if data.CreatedBy.Connect != nil {
			toConnect, err := repository.User.GetOneTx(ctx, tx, *data.CreatedBy.Connect)

			if err != nil {
				return nil, err
			}

			data.CreatedByID = &toConnect.ID
		}
		if data.CreatedBy.Unset != nil && *data.CreatedBy.Unset {
			data.ClearCreatedBy = true
		}
		if data.CreatedBy.Create != nil {
			var createInput ent.CreateUserInput
			err = mapstructure.Decode(data.CreatedBy.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.User.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.CreatedByID = &toConnect.ID
		}
		if data.CreatedBy.Delete != nil && *data.CreatedBy.Delete {
			itemToDelete, err := item.CreatedBy(ctx)
			if err != nil {
				return nil, err
			}
			created_byToDelete.ID = &itemToDelete.ID
		}
	}
	if data.RefUpdatedBy != nil {

		if data.RefUpdatedBy.Connect != nil {
			for _, connect := range data.RefUpdatedBy.Connect {
				toConnect, err := repository.User.GetOneTx(ctx, tx, *connect)

				if err != nil {
					return nil, err
				}

				data.AddRefUpdatedByIDs = append(data.AddRefUpdatedByIDs, toConnect.ID)
			}
		}
		if data.RefUpdatedBy.Disconnect != nil {
			for _, disconnect := range data.RefUpdatedBy.Disconnect {
				toDisconnect, err := repository.User.GetOneTx(ctx, tx, *disconnect)

				if err != nil {
					return nil, err
				}

				data.RemoveRefUpdatedByIDs = append(data.RemoveRefUpdatedByIDs, toDisconnect.ID)
			}
		}
		if data.RefUpdatedBy.Create != nil {
			for _, create := range data.RefUpdatedBy.Create {
				var createInput ent.CreateUserInput
				err = mapstructure.Decode(create, &createInput)
				if err != nil {
					return nil, err
				}

				toConnect, err := repository.User.CreateTx(ctx, tx, createInput)

				if err != nil {
					return nil, err
				}

				data.AddRefUpdatedByIDs = append(data.AddRefUpdatedByIDs, toConnect.ID)
			}
		}
		if data.RefUpdatedBy.Delete != nil {
			for _, delete := range data.RefUpdatedBy.Delete {
				_, err := repository.User.DeleteTx(ctx, tx, *delete)

				if err != nil {
					return nil, err
				}
			}
		}
	}
	updated_byToDelete := ent.UserWhereUniqueInput{}
	if data.UpdatedBy != nil {

		if data.UpdatedBy.Connect != nil {
			toConnect, err := repository.User.GetOneTx(ctx, tx, *data.UpdatedBy.Connect)

			if err != nil {
				return nil, err
			}

			data.UpdatedByID = &toConnect.ID
		}
		if data.UpdatedBy.Unset != nil && *data.UpdatedBy.Unset {
			data.ClearUpdatedBy = true
		}
		if data.UpdatedBy.Create != nil {
			var createInput ent.CreateUserInput
			err = mapstructure.Decode(data.UpdatedBy.Create, &createInput)
			if err != nil {
				return nil, err
			}

			toConnect, err := repository.User.CreateTx(ctx, tx, createInput)

			if err != nil {
				return nil, err
			}

			data.UpdatedByID = &toConnect.ID
		}
		if data.UpdatedBy.Delete != nil && *data.UpdatedBy.Delete {
			itemToDelete, err := item.UpdatedBy(ctx)
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

	updatedInput, err := tx.User.UpdateOne(item).SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error updating item: %v", err)
	}
	if data.CreatedBy != nil && data.CreatedBy.Delete != nil && *data.CreatedBy.Delete {
		_, err := repository.User.DeleteTx(ctx, tx, created_byToDelete)
		if err != nil {
			return nil, err
		}
	}
	if data.UpdatedBy != nil && data.UpdatedBy.Delete != nil && *data.UpdatedBy.Delete {
		_, err := repository.User.DeleteTx(ctx, tx, updated_byToDelete)
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

func (r *User) UpdateMany(ctx context.Context, where ent.UserWhereInput, data ent.UpdateUserInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := r.client.User.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *User) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereInput, data ent.UpdateUserInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := tx.User.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *User) Upsert(ctx context.Context, data ent.CreateUserInput) (upsertedUser *ent.User, err error) {
	if CompareUniqueInput(data, ent.UserWhereUniqueInput{}) {
		var where ent.UserWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}
		err = r.client.User.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		upsertedUser, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		return upsertedUser, nil
	}
	upsertedUser, err = r.client.User.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedUser, nil
}

func (r *User) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateUserInput) (upsertedUser *ent.User, err error) {
	if CompareUniqueInput(data, ent.UserWhereUniqueInput{}) {
		var where ent.UserWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}
		err = tx.User.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		upsertedUser, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		return upsertedUser, nil
	}
	upsertedUser, err = tx.User.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedUser, nil
}

func (r *User) UpsertMany(ctx context.Context, data []ent.CreateUserInput) (int, error) {
	var upsertMap []*ent.UserCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, r.client.User.Create().SetInput(upsertInput))
	}
	err := r.client.User.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *User) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateUserInput) (int, error) {
	var upsertMap []*ent.UserCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, tx.User.Create().SetInput(upsertInput))
	}
	err := tx.User.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *User) Delete(ctx context.Context, where ent.UserWhereUniqueInput) (*ent.User, error) {
	query := r.client.User.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = r.client.User.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *User) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereUniqueInput) (*ent.User, error) {
	query := tx.User.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = tx.User.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *User) DeleteMany(ctx context.Context, where ent.UserWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := r.client.User.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}

func (r *User) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.UserWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := tx.User.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}
