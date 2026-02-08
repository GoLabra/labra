package repo

import (
	"context"
	"errors"
	"fmt"
	"app/ent"

	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/ext/mapstructure"
)

type Miau struct {
	client *ent.Client
}

func NewMiau(c *ent.Client) *Miau {
	return &Miau{client: c}
}

func (r *Miau) Get(ctx context.Context, where *ent.MiauWhereInput, orderBy *ent.MiauOrder, skip *int, first *int, last *int) ([]*ent.Miau, error) {
	var (
        query = r.client.Miau.Query()
        err   error
    )
    if where != nil {
        query, err = where.Filter(query)
        if err != nil {
            return nil, err
        }
    }

    if orderBy == nil {
		orderBy = &ent.MiauOrder{
			Direction: ent.DefaultMiauOrder.Direction,
			Field:     ent.DefaultMiauOrder.Field,
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

func (r *Miau) Connection(ctx context.Context, where *ent.MiauWhereInput, orderBy *ent.MiauOrder, skip *int, first *int, last *int) (*ent.MiauConnection, error) {
	var (
		query = r.client.Miau.Query()
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

func (r *Miau) GetOne(ctx context.Context, where ent.MiauWhereUniqueInput) (*ent.Miau, error) {
	query := r.client.Miau.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}

func (r *Miau) GetTx(ctx context.Context, tx *ent.Tx, where *ent.MiauWhereInput, orderBy *ent.MiauOrder, skip *int, first *int, last *int) ([]*ent.Miau, error) {
	var (
		query = tx.Miau.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.MiauOrder{
			Direction: ent.DefaultMiauOrder.Direction,
			Field:     ent.DefaultMiauOrder.Field,
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

func (r *Miau) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereUniqueInput) (*ent.Miau, error) {
	query := tx.Miau.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}
func (r *Miau) Create(ctx context.Context, data ent.CreateMiauInput) (*ent.Miau, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)

	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	createdInput, err := repository.Miau.CreateTx(ctx, tx, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return createdInput, err
}

func (r *Miau) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateMiauInput) (*ent.Miau, error) {
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

	createdInput, err := tx.Miau.Create().SetInput(data).Save(ctx)

	if err != nil {
		return nil, err
	}

	return createdInput.Unwrap(), err
}

func (r *Miau) CreateMany(ctx context.Context, data []ent.CreateMiauInput) ([]*ent.Miau, error) {
	var createMap []*ent.MiauCreate
	for _, createInput := range data {
        createMap = append(createMap, r.client.Miau.Create().SetInput(createInput))
    }
	createdItems, err := r.client.Miau.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}

func (r *Miau) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateMiauInput) ([]*ent.Miau, error) {
	var createMap []*ent.MiauCreate
	for _, createInput := range data {
        createMap = append(createMap, tx.Miau.Create().SetInput(createInput))
    }
	createdItems, err := tx.Miau.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}
func (r *Miau) Update(ctx context.Context, where ent.MiauWhereUniqueInput, data ent.UpdateMiauInput) (*ent.Miau, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	updatedInput, err := repository.Miau.UpdateTx(ctx, tx, where, data)

	if err != nil {
		return nil, err
	}
	
	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return updatedInput.Unwrap(), nil
}

func (r *Miau) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereUniqueInput, data ent.UpdateMiauInput) (*ent.Miau, error) {
	var iCtx = context.WithValue(ctx, constants.IsInternalOperationContextValue, true)

    repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}
	
	query := tx.Miau.Query()
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

	updatedInput, err := tx.Miau.UpdateOne(item).SetInput(data).Save(ctx)
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

func (r *Miau) UpdateMany(ctx context.Context, where ent.MiauWhereInput, data ent.UpdateMiauInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := r.client.Miau.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *Miau) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereInput, data ent.UpdateMiauInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := tx.Miau.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *Miau) Upsert(ctx context.Context, data ent.CreateMiauInput) (upsertedMiau *ent.Miau, err error) {
	if CompareUniqueInput(data, ent.MiauWhereUniqueInput{}) {
		var where ent.MiauWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}

		upsertedMiau, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		if upsertedMiau != nil {
			updateMiauInput := ent.UpdateMiauInput{}
			err = mapstructure.Decode(data, &updateMiauInput)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			upsertedMiau, err = r.client.Miau.UpdateOne(upsertedMiau).SetInput(updateMiauInput).Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			return upsertedMiau, nil
		}

		err = r.client.Miau.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		
		return upsertedMiau, nil
	}
	upsertedMiau, err = r.client.Miau.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedMiau, nil
}

func (r *Miau) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateMiauInput) (upsertedMiau *ent.Miau, err error) {
	var where ent.MiauWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}

		upsertedMiau, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		if upsertedMiau != nil {
			updateMiauInput := ent.UpdateMiauInput{}
			err = mapstructure.Decode(data, &updateMiauInput)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			upsertedMiau, err = tx.Miau.UpdateOne(upsertedMiau).SetInput(updateMiauInput).Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			return upsertedMiau, nil
		}
	upsertedMiau, err = tx.Miau.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedMiau, nil
}

func (r *Miau) UpsertMany(ctx context.Context, data []ent.CreateMiauInput) (int, error) {
	var upsertMap []*ent.MiauCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, r.client.Miau.Create().SetInput(upsertInput))
	}
	err := r.client.Miau.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *Miau) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateMiauInput) (int, error) {
	var upsertMap []*ent.MiauCreate
	for _, upsertInput := range data {
        upsertMap = append(upsertMap, tx.Miau.Create().SetInput(upsertInput))
    }
	err := tx.Miau.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *Miau) Delete(ctx context.Context, where ent.MiauWhereUniqueInput) (*ent.Miau, error) {
	query := r.client.Miau.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = r.client.Miau.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *Miau) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereUniqueInput) (*ent.Miau, error) {
	query := tx.Miau.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = tx.Miau.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *Miau) DeleteMany(ctx context.Context, where ent.MiauWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := r.client.Miau.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}

func (r *Miau) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.MiauWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := tx.Miau.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}