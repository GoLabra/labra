package repo

import (
	"context"
	"errors"
	"fmt"
	"app/ent"

	"github.com/GoLabra/labrago/src/api/constants"
	"github.com/mitchellh/mapstructure"
)

type Role struct {
	client *ent.Client
}

func NewRole(c *ent.Client) *Role {
	return &Role{client: c}
}

func (r *Role) Get(ctx context.Context, where *ent.RoleWhereInput, orderBy *ent.RoleOrder, skip *int, first *int, last *int) ([]*ent.Role, error) {
	var (
        query = r.client.Role.Query()
        err   error
    )
    if where != nil {
        query, err = where.Filter(query)
        if err != nil {
            return nil, err
        }
    }

    if orderBy == nil {
		orderBy = &ent.RoleOrder{
			Direction: ent.DefaultRoleOrder.Direction,
			Field:     ent.DefaultRoleOrder.Field,
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

func (r *Role) Connection(ctx context.Context, where *ent.RoleWhereInput, orderBy *ent.RoleOrder, skip *int, first *int, last *int) (*ent.RoleConnection, error) {
	var (
		query = r.client.Role.Query()
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

func (r *Role) GetOne(ctx context.Context, where ent.RoleWhereUniqueInput) (*ent.Role, error) {
	query := r.client.Role.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}

func (r *Role) GetTx(ctx context.Context, tx *ent.Tx, where *ent.RoleWhereInput, orderBy *ent.RoleOrder, skip *int, first *int, last *int) ([]*ent.Role, error) {
	var (
		query = tx.Role.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.RoleOrder{
			Direction: ent.DefaultRoleOrder.Direction,
			Field:     ent.DefaultRoleOrder.Field,
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

func (r *Role) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereUniqueInput) (*ent.Role, error) {
	query := tx.Role.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}


func (r *Role) Create(ctx context.Context, data ent.CreateRoleInput) (*ent.Role, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)

	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	createdInput, err := repository.Role.CreateTx(ctx, tx, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return createdInput, err
}

func (r *Role) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateRoleInput) (*ent.Role, error) {
	var err error
    repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}
    if data.CreatedBy != nil {
        
    if data.CreatedBy.Connect != nil {    
        toConnect, err := repository.User.GetOneTx(ctx, tx, *data.CreatedBy.Connect)

        if err != nil {
            return nil, err
        }

        data.CreatedByID =  &toConnect.ID
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
    if data.UpdatedBy != nil {
        
    if data.UpdatedBy.Connect != nil {    
        toConnect, err := repository.User.GetOneTx(ctx, tx, *data.UpdatedBy.Connect)

        if err != nil {
            return nil, err
        }

        data.UpdatedByID =  &toConnect.ID
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

	createdInput, err := tx.Role.Create().SetInput(data).Save(ctx)

	if err != nil {
		return nil, err
	}

	return createdInput.Unwrap(), err
}

func (r *Role) CreateMany(ctx context.Context, data []ent.CreateRoleInput) ([]*ent.Role, error) {
	var createMap []*ent.RoleCreate
	for _, createInput := range data {
        createMap = append(createMap, r.client.Role.Create().SetInput(createInput))
    }
	createdItems, err := r.client.Role.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}

func (r *Role) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateRoleInput) ([]*ent.Role, error) {
	var createMap []*ent.RoleCreate
	for _, createInput := range data {
        createMap = append(createMap, tx.Role.Create().SetInput(createInput))
    }
	createdItems, err := tx.Role.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}

func (r *Role) Update(ctx context.Context, where ent.RoleWhereUniqueInput, data ent.UpdateRoleInput) (*ent.Role, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	updatedInput, err := repository.Role.UpdateTx(ctx, tx, where, data)

	if err != nil {
		return nil, err
	}
	
	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return updatedInput.Unwrap(), nil
}

func (r *Role) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereUniqueInput, data ent.UpdateRoleInput) (*ent.Role, error) {
	
    repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}
	
	query := tx.Role.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to update: %v", err)
	}
    if data.CreatedBy != nil {
        
    if data.CreatedBy.Connect != nil {    
        toConnect, err := repository.User.GetOneTx(ctx, tx, *data.CreatedBy.Connect)

        if err != nil {
            return nil, err
        }

        data.CreatedByID =  &toConnect.ID
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
    if data.UpdatedBy != nil {
        
    if data.UpdatedBy.Connect != nil {    
        toConnect, err := repository.User.GetOneTx(ctx, tx, *data.UpdatedBy.Connect)

        if err != nil {
            return nil, err
        }

        data.UpdatedByID =  &toConnect.ID
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

	updatedInput, err := tx.Role.UpdateOne(item).SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error updating item: %v", err)
	}

	return updatedInput, nil
}

func (r *Role) UpdateMany(ctx context.Context, where ent.RoleWhereInput, data ent.UpdateRoleInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := r.client.Role.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *Role) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereInput, data ent.UpdateRoleInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := tx.Role.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *Role) Upsert(ctx context.Context, data ent.CreateRoleInput) (upsertedRole *ent.Role, err error) {
	if CompareUniqueInput(data, ent.RoleWhereUniqueInput{}) {
		var where ent.RoleWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}
		err = r.client.Role.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		upsertedRole, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		return upsertedRole, nil
	}
	upsertedRole, err = r.client.Role.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedRole, nil
}

func (r *Role) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateRoleInput) (upsertedRole *ent.Role, err error) {
	if CompareUniqueInput(data, ent.RoleWhereUniqueInput{}) {
		var where ent.RoleWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}
		err = tx.Role.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		upsertedRole, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		return upsertedRole, nil
	}
	upsertedRole, err = tx.Role.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedRole, nil
}

func (r *Role) UpsertMany(ctx context.Context, data []ent.CreateRoleInput) (int, error) {
	var upsertMap []*ent.RoleCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, r.client.Role.Create().SetInput(upsertInput))
	}
	err := r.client.Role.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *Role) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateRoleInput) (int, error) {
	var upsertMap []*ent.RoleCreate
	for _, upsertInput := range data {
        upsertMap = append(upsertMap, tx.Role.Create().SetInput(upsertInput))
    }
	err := tx.Role.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *Role) Delete(ctx context.Context, where ent.RoleWhereUniqueInput) (*ent.Role, error) {
	query := r.client.Role.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = r.client.Role.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *Role) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereUniqueInput) (*ent.Role, error) {
	query := tx.Role.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = tx.Role.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *Role) DeleteMany(ctx context.Context, where ent.RoleWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := r.client.Role.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}

func (r *Role) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.RoleWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := tx.Role.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}