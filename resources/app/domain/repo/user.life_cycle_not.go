package repo

import (
	"context"
	"errors"
	"fmt"
	"app/ent"

	"github.com/GoLabra/labra/entgql/entity"

	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/ext/mapstructure"
)

type LifeCycleNot struct {
	client *ent.Client
}

func NewLifeCycleNot(c *ent.Client) *LifeCycleNot {
	return &LifeCycleNot{client: c}
}

func (r *LifeCycleNot) Get(ctx context.Context, where *ent.LifeCycleNotWhereInput, orderBy *ent.LifeCycleNotOrder, skip *int, first *int, last *int) ([]*ent.LifeCycleNot, error) {
	var (
        query = r.client.LifeCycleNot.Query()
        err   error
    )
    if where != nil {
        query, err = where.Filter(query)
        if err != nil {
            return nil, err
        }
    }

    if orderBy == nil {
		orderBy = &ent.LifeCycleNotOrder{
			Direction: ent.DefaultLifeCycleNotOrder.Direction,
			Field:     ent.DefaultLifeCycleNotOrder.Field,
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

func (r *LifeCycleNot) Connection(ctx context.Context, where *ent.LifeCycleNotWhereInput, orderBy *ent.LifeCycleNotOrder, skip *int, first *int, last *int) (*ent.LifeCycleNotConnection, error) {
	var (
		query = r.client.LifeCycleNot.Query()
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

func (r *LifeCycleNot) GetOne(ctx context.Context, where ent.LifeCycleNotWhereUniqueInput) (*ent.LifeCycleNot, error) {
	query := r.client.LifeCycleNot.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}

func (r *LifeCycleNot) GetTx(ctx context.Context, tx *ent.Tx, where *ent.LifeCycleNotWhereInput, orderBy *ent.LifeCycleNotOrder, skip *int, first *int, last *int) ([]*ent.LifeCycleNot, error) {
	var (
		query = tx.LifeCycleNot.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.LifeCycleNotOrder{
			Direction: ent.DefaultLifeCycleNotOrder.Direction,
			Field:     ent.DefaultLifeCycleNotOrder.Field,
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

func (r *LifeCycleNot) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereUniqueInput) (*ent.LifeCycleNot, error) {
	query := tx.LifeCycleNot.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}
func (r *LifeCycleNot) Create(ctx context.Context, data ent.CreateLifeCycleNotInput) (*ent.LifeCycleNot, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)

	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	createdInput, err := repository.LifeCycleNot.CreateTx(ctx, tx, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return createdInput, err
}

func (r *LifeCycleNot) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateLifeCycleNotInput) (*ent.LifeCycleNot, error) {
	var err error
    repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}
    if data.CreatedBy != nil {
        
    if data.CreatedBy.Connect != nil {    
        toConnect, err := repository.User.GetOne(ctx, *data.CreatedBy.Connect)

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
        
        toConnect, err := repository.User.Create(ctx, createInput)

        if err != nil {
            return nil, err
        }

        data.CreatedByID = &toConnect.ID
    }
    }
    if data.UpdatedBy != nil {
        
    if data.UpdatedBy.Connect != nil {    
        toConnect, err := repository.User.GetOne(ctx, *data.UpdatedBy.Connect)

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
        
        toConnect, err := repository.User.Create(ctx, createInput)

        if err != nil {
            return nil, err
        }

        data.UpdatedByID = &toConnect.ID
    }
    }
    if data.AdminCreatedBy != nil {
        
    if data.AdminCreatedBy.Connect != nil {    
        toConnect, err := repository.AdminUser.GetOne(ctx, *data.AdminCreatedBy.Connect)

        if err != nil {
            return nil, err
        }

        data.AdminCreatedByID =  &toConnect.ID
    }
    }
    if data.AdminUpdatedBy != nil {
        
    if data.AdminUpdatedBy.Connect != nil {    
        toConnect, err := repository.AdminUser.GetOne(ctx, *data.AdminUpdatedBy.Connect)

        if err != nil {
            return nil, err
        }

        data.AdminUpdatedByID =  &toConnect.ID
    }
    }

	createdInput, err := tx.LifeCycleNot.Create().SetInput(data).Save(ctx)

	if err != nil {
		return nil, err
	}

	return createdInput.Unwrap(), err
}

func (r *LifeCycleNot) CreateMany(ctx context.Context, data []ent.CreateLifeCycleNotInput) ([]*ent.LifeCycleNot, error) {
	var createMap []*ent.LifeCycleNotCreate
	for _, createInput := range data {
        createMap = append(createMap, r.client.LifeCycleNot.Create().SetInput(createInput))
    }
	createdItems, err := r.client.LifeCycleNot.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}

func (r *LifeCycleNot) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateLifeCycleNotInput) ([]*ent.LifeCycleNot, error) {
	var createMap []*ent.LifeCycleNotCreate
	for _, createInput := range data {
        createMap = append(createMap, tx.LifeCycleNot.Create().SetInput(createInput))
    }
	createdItems, err := tx.LifeCycleNot.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}
func (r *LifeCycleNot) Update(ctx context.Context, where ent.LifeCycleNotWhereUniqueInput, data ent.UpdateLifeCycleNotInput) (*ent.LifeCycleNot, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	updatedInput, err := repository.LifeCycleNot.UpdateTx(ctx, tx, where, data)

	if err != nil {
		return nil, err
	}
	
	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return updatedInput.Unwrap(), nil
}
func (r *LifeCycleNot) UpdateState(ctx context.Context, where ent.LifeCycleNotWhereUniqueInput, state entity.EntityState) (*ent.LifeCycleNot, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}
	tx, err := repository.Tx.Create(ctx)
	if err != nil {
		return nil, err
	}
	updated, err := repository.LifeCycleNot.UpdateStateTx(ctx, tx, where, state)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return updated.Unwrap(), nil
}

func (r *LifeCycleNot) UpdateStateTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereUniqueInput, state entity.EntityState) (*ent.LifeCycleNot, error) {
	query := tx.LifeCycleNot.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to update state: %v", err)
	}
	return tx.LifeCycleNot.UpdateOne(item).SetEntityState(state).Save(ctx)
}

func (r *LifeCycleNot) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereUniqueInput, data ent.UpdateLifeCycleNotInput) (*ent.LifeCycleNot, error) {
	var iCtx = context.WithValue(ctx, constants.IsInternalOperationContextValue, true)

    repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}
	
	query := tx.LifeCycleNot.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(iCtx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to update: %v", err)
	}
	created_byToDelete := ent.UserWhereUniqueInput{}
    if data.CreatedBy != nil {
        
        if data.CreatedBy.Connect != nil {    
            toConnect, err := repository.User.GetOne(ctx, *data.CreatedBy.Connect)

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
            
            toConnect, err := repository.User.Create(ctx, createInput)

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
	updated_byToDelete := ent.UserWhereUniqueInput{}
    if data.UpdatedBy != nil {
        
        if data.UpdatedBy.Connect != nil {    
            toConnect, err := repository.User.GetOne(ctx, *data.UpdatedBy.Connect)

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
            
            toConnect, err := repository.User.Create(ctx, createInput)

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
    if data.AdminCreatedBy != nil {
        
        if data.AdminCreatedBy.Connect != nil {    
            toConnect, err := repository.AdminUser.GetOne(ctx, *data.AdminCreatedBy.Connect)

            if err != nil {
                return nil, err
            }

            data.AdminCreatedByID = &toConnect.ID
        }
		if data.AdminCreatedBy.Unset != nil && *data.AdminCreatedBy.Unset {
			data.ClearAdminCreatedBy = true
		}
    }
    if data.AdminUpdatedBy != nil {
        
        if data.AdminUpdatedBy.Connect != nil {    
            toConnect, err := repository.AdminUser.GetOne(ctx, *data.AdminUpdatedBy.Connect)

            if err != nil {
                return nil, err
            }

            data.AdminUpdatedByID = &toConnect.ID
        }
		if data.AdminUpdatedBy.Unset != nil && *data.AdminUpdatedBy.Unset {
			data.ClearAdminUpdatedBy = true
		}
    }

	updatedInput, err := tx.LifeCycleNot.UpdateOne(item).SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error updating item: %v", err)
	} 
	if data.CreatedBy != nil && data.CreatedBy.Delete != nil && *data.CreatedBy.Delete {
		_, err := repository.User.Delete(ctx, created_byToDelete)
		if err != nil {
			return nil, err
		}
	} 
	if data.UpdatedBy != nil && data.UpdatedBy.Delete != nil && *data.UpdatedBy.Delete {
		_, err := repository.User.Delete(ctx, updated_byToDelete)
		if err != nil {
			return nil, err
		}
	}

	return updatedInput, nil
}

func (r *LifeCycleNot) UpdateMany(ctx context.Context, where ent.LifeCycleNotWhereInput, data ent.UpdateLifeCycleNotInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := r.client.LifeCycleNot.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *LifeCycleNot) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereInput, data ent.UpdateLifeCycleNotInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := tx.LifeCycleNot.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *LifeCycleNot) Upsert(ctx context.Context, data ent.CreateLifeCycleNotInput) (upsertedLifeCycleNot *ent.LifeCycleNot, err error) {
	if CompareUniqueInput(data, ent.LifeCycleNotWhereUniqueInput{}) {
		var where ent.LifeCycleNotWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}

		upsertedLifeCycleNot, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		if upsertedLifeCycleNot != nil {
			updateLifeCycleNotInput := ent.UpdateLifeCycleNotInput{}
			err = mapstructure.Decode(data, &updateLifeCycleNotInput)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			upsertedLifeCycleNot, err = r.client.LifeCycleNot.UpdateOne(upsertedLifeCycleNot).SetInput(updateLifeCycleNotInput).Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			return upsertedLifeCycleNot, nil
		}

		err = r.client.LifeCycleNot.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		
		return upsertedLifeCycleNot, nil
	}
	upsertedLifeCycleNot, err = r.client.LifeCycleNot.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedLifeCycleNot, nil
}

func (r *LifeCycleNot) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateLifeCycleNotInput) (upsertedLifeCycleNot *ent.LifeCycleNot, err error) {
	var where ent.LifeCycleNotWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}

		upsertedLifeCycleNot, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		if upsertedLifeCycleNot != nil {
			updateLifeCycleNotInput := ent.UpdateLifeCycleNotInput{}
			err = mapstructure.Decode(data, &updateLifeCycleNotInput)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			upsertedLifeCycleNot, err = tx.LifeCycleNot.UpdateOne(upsertedLifeCycleNot).SetInput(updateLifeCycleNotInput).Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			return upsertedLifeCycleNot, nil
		}
	upsertedLifeCycleNot, err = tx.LifeCycleNot.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedLifeCycleNot, nil
}

func (r *LifeCycleNot) UpsertMany(ctx context.Context, data []ent.CreateLifeCycleNotInput) (int, error) {
	var upsertMap []*ent.LifeCycleNotCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, r.client.LifeCycleNot.Create().SetInput(upsertInput))
	}
	err := r.client.LifeCycleNot.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *LifeCycleNot) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateLifeCycleNotInput) (int, error) {
	var upsertMap []*ent.LifeCycleNotCreate
	for _, upsertInput := range data {
        upsertMap = append(upsertMap, tx.LifeCycleNot.Create().SetInput(upsertInput))
    }
	err := tx.LifeCycleNot.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *LifeCycleNot) Delete(ctx context.Context, where ent.LifeCycleNotWhereUniqueInput) (*ent.LifeCycleNot, error) {
	query := r.client.LifeCycleNot.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = r.client.LifeCycleNot.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *LifeCycleNot) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereUniqueInput) (*ent.LifeCycleNot, error) {
	query := tx.LifeCycleNot.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = tx.LifeCycleNot.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *LifeCycleNot) DeleteMany(ctx context.Context, where ent.LifeCycleNotWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := r.client.LifeCycleNot.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}

func (r *LifeCycleNot) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.LifeCycleNotWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := tx.LifeCycleNot.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}