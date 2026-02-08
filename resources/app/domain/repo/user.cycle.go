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

type Cycle struct {
	client *ent.Client
}

func NewCycle(c *ent.Client) *Cycle {
	return &Cycle{client: c}
}

func (r *Cycle) Get(ctx context.Context, where *ent.CycleWhereInput, orderBy *ent.CycleOrder, skip *int, first *int, last *int) ([]*ent.Cycle, error) {
	var (
        query = r.client.Cycle.Query()
        err   error
    )
    if where != nil {
        query, err = where.Filter(query)
        if err != nil {
            return nil, err
        }
    }

    if orderBy == nil {
		orderBy = &ent.CycleOrder{
			Direction: ent.DefaultCycleOrder.Direction,
			Field:     ent.DefaultCycleOrder.Field,
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

func (r *Cycle) Connection(ctx context.Context, where *ent.CycleWhereInput, orderBy *ent.CycleOrder, skip *int, first *int, last *int) (*ent.CycleConnection, error) {
	var (
		query = r.client.Cycle.Query()
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

func (r *Cycle) GetOne(ctx context.Context, where ent.CycleWhereUniqueInput) (*ent.Cycle, error) {
	query := r.client.Cycle.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}

func (r *Cycle) GetTx(ctx context.Context, tx *ent.Tx, where *ent.CycleWhereInput, orderBy *ent.CycleOrder, skip *int, first *int, last *int) ([]*ent.Cycle, error) {
	var (
		query = tx.Cycle.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.CycleOrder{
			Direction: ent.DefaultCycleOrder.Direction,
			Field:     ent.DefaultCycleOrder.Field,
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

func (r *Cycle) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereUniqueInput) (*ent.Cycle, error) {
	query := tx.Cycle.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}
func (r *Cycle) Create(ctx context.Context, data ent.CreateCycleInput) (*ent.Cycle, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)

	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	createdInput, err := repository.Cycle.CreateTx(ctx, tx, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return createdInput, err
}

func (r *Cycle) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateCycleInput) (*ent.Cycle, error) {
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

	createdInput, err := tx.Cycle.Create().SetInput(data).Save(ctx)

	if err != nil {
		return nil, err
	}

	return createdInput.Unwrap(), err
}

func (r *Cycle) CreateMany(ctx context.Context, data []ent.CreateCycleInput) ([]*ent.Cycle, error) {
	var createMap []*ent.CycleCreate
	for _, createInput := range data {
        createMap = append(createMap, r.client.Cycle.Create().SetInput(createInput))
    }
	createdItems, err := r.client.Cycle.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}

func (r *Cycle) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateCycleInput) ([]*ent.Cycle, error) {
	var createMap []*ent.CycleCreate
	for _, createInput := range data {
        createMap = append(createMap, tx.Cycle.Create().SetInput(createInput))
    }
	createdItems, err := tx.Cycle.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}
func (r *Cycle) Update(ctx context.Context, where ent.CycleWhereUniqueInput, data ent.UpdateCycleInput) (*ent.Cycle, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	updatedInput, err := repository.Cycle.UpdateTx(ctx, tx, where, data)

	if err != nil {
		return nil, err
	}
	
	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return updatedInput.Unwrap(), nil
}
func (r *Cycle) UpdateState(ctx context.Context, where ent.CycleWhereUniqueInput, state entity.EntityState) (*ent.Cycle, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}
	tx, err := repository.Tx.Create(ctx)
	if err != nil {
		return nil, err
	}
	updated, err := repository.Cycle.UpdateStateTx(ctx, tx, where, state)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return updated.Unwrap(), nil
}

func (r *Cycle) UpdateStateTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereUniqueInput, state entity.EntityState) (*ent.Cycle, error) {
	query := tx.Cycle.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to update state: %v", err)
	}
	return tx.Cycle.UpdateOne(item).SetEntityState(state).Save(ctx)
}

func (r *Cycle) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereUniqueInput, data ent.UpdateCycleInput) (*ent.Cycle, error) {
	var iCtx = context.WithValue(ctx, constants.IsInternalOperationContextValue, true)

    repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}
	
	query := tx.Cycle.Query()
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

	updatedInput, err := tx.Cycle.UpdateOne(item).SetInput(data).Save(ctx)
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

func (r *Cycle) UpdateMany(ctx context.Context, where ent.CycleWhereInput, data ent.UpdateCycleInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := r.client.Cycle.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *Cycle) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereInput, data ent.UpdateCycleInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := tx.Cycle.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *Cycle) Upsert(ctx context.Context, data ent.CreateCycleInput) (upsertedCycle *ent.Cycle, err error) {
	if CompareUniqueInput(data, ent.CycleWhereUniqueInput{}) {
		var where ent.CycleWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}

		upsertedCycle, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		if upsertedCycle != nil {
			updateCycleInput := ent.UpdateCycleInput{}
			err = mapstructure.Decode(data, &updateCycleInput)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			upsertedCycle, err = r.client.Cycle.UpdateOne(upsertedCycle).SetInput(updateCycleInput).Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			return upsertedCycle, nil
		}

		err = r.client.Cycle.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		
		return upsertedCycle, nil
	}
	upsertedCycle, err = r.client.Cycle.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedCycle, nil
}

func (r *Cycle) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateCycleInput) (upsertedCycle *ent.Cycle, err error) {
	var where ent.CycleWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}

		upsertedCycle, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		if upsertedCycle != nil {
			updateCycleInput := ent.UpdateCycleInput{}
			err = mapstructure.Decode(data, &updateCycleInput)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			upsertedCycle, err = tx.Cycle.UpdateOne(upsertedCycle).SetInput(updateCycleInput).Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			return upsertedCycle, nil
		}
	upsertedCycle, err = tx.Cycle.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedCycle, nil
}

func (r *Cycle) UpsertMany(ctx context.Context, data []ent.CreateCycleInput) (int, error) {
	var upsertMap []*ent.CycleCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, r.client.Cycle.Create().SetInput(upsertInput))
	}
	err := r.client.Cycle.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *Cycle) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateCycleInput) (int, error) {
	var upsertMap []*ent.CycleCreate
	for _, upsertInput := range data {
        upsertMap = append(upsertMap, tx.Cycle.Create().SetInput(upsertInput))
    }
	err := tx.Cycle.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *Cycle) Delete(ctx context.Context, where ent.CycleWhereUniqueInput) (*ent.Cycle, error) {
	query := r.client.Cycle.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = r.client.Cycle.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *Cycle) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereUniqueInput) (*ent.Cycle, error) {
	query := tx.Cycle.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = tx.Cycle.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *Cycle) DeleteMany(ctx context.Context, where ent.CycleWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := r.client.Cycle.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}

func (r *Cycle) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.CycleWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := tx.Cycle.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}