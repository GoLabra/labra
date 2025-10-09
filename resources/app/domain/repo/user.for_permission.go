package repo

import (
	"context"
	"errors"
	"fmt"
	"app/ent"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/ext/mapstructure"
)

type ForPermission struct {
	client *ent.Client
}

func NewForPermission(c *ent.Client) *ForPermission {
	return &ForPermission{client: c}
}

func (r *ForPermission) Get(ctx context.Context, where *ent.ForPermissionWhereInput, orderBy *ent.ForPermissionOrder, skip *int, first *int, last *int) ([]*ent.ForPermission, error) {
	var (
        query = r.client.ForPermission.Query()
        err   error
    )
    if where != nil {
        query, err = where.Filter(query)
        if err != nil {
            return nil, err
        }
    }

    if orderBy == nil {
		orderBy = &ent.ForPermissionOrder{
			Direction: ent.DefaultForPermissionOrder.Direction,
			Field:     ent.DefaultForPermissionOrder.Field,
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

func (r *ForPermission) Connection(ctx context.Context, where *ent.ForPermissionWhereInput, orderBy *ent.ForPermissionOrder, skip *int, first *int, last *int) (*ent.ForPermissionConnection, error) {
	var (
		query = r.client.ForPermission.Query()
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

func (r *ForPermission) GetOne(ctx context.Context, where ent.ForPermissionWhereUniqueInput) (*ent.ForPermission, error) {
	query := r.client.ForPermission.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}

func (r *ForPermission) GetTx(ctx context.Context, tx *ent.Tx, where *ent.ForPermissionWhereInput, orderBy *ent.ForPermissionOrder, skip *int, first *int, last *int) ([]*ent.ForPermission, error) {
	var (
		query = tx.ForPermission.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.ForPermissionOrder{
			Direction: ent.DefaultForPermissionOrder.Direction,
			Field:     ent.DefaultForPermissionOrder.Field,
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

func (r *ForPermission) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereUniqueInput) (*ent.ForPermission, error) {
	query := tx.ForPermission.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}
func (r *ForPermission) Create(ctx context.Context, data ent.CreateForPermissionInput) (*ent.ForPermission, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)

	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	createdInput, err := repository.ForPermission.CreateTx(ctx, tx, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return createdInput, err
}

func (r *ForPermission) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateForPermissionInput) (*ent.ForPermission, error) {
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

	createdInput, err := tx.ForPermission.Create().SetInput(data).Save(ctx)

	if err != nil {
		return nil, err
	}

	return createdInput.Unwrap(), err
}

func (r *ForPermission) CreateMany(ctx context.Context, data []ent.CreateForPermissionInput) ([]*ent.ForPermission, error) {
	var createMap []*ent.ForPermissionCreate
	for _, createInput := range data {
        createMap = append(createMap, r.client.ForPermission.Create().SetInput(createInput))
    }
	createdItems, err := r.client.ForPermission.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}

func (r *ForPermission) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateForPermissionInput) ([]*ent.ForPermission, error) {
	var createMap []*ent.ForPermissionCreate
	for _, createInput := range data {
        createMap = append(createMap, tx.ForPermission.Create().SetInput(createInput))
    }
	createdItems, err := tx.ForPermission.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}
func (r *ForPermission) Update(ctx context.Context, where ent.ForPermissionWhereUniqueInput, data ent.UpdateForPermissionInput) (*ent.ForPermission, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	updatedInput, err := repository.ForPermission.UpdateTx(ctx, tx, where, data)

	if err != nil {
		return nil, err
	}
	
	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return updatedInput.Unwrap(), nil
}

func (r *ForPermission) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereUniqueInput, data ent.UpdateForPermissionInput) (*ent.ForPermission, error) {
	
    repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}
	
	query := tx.ForPermission.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
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

	updatedInput, err := tx.ForPermission.UpdateOne(item).SetInput(data).Save(ctx)
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

func (r *ForPermission) UpdateMany(ctx context.Context, where ent.ForPermissionWhereInput, data ent.UpdateForPermissionInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := r.client.ForPermission.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *ForPermission) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereInput, data ent.UpdateForPermissionInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := tx.ForPermission.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *ForPermission) Upsert(ctx context.Context, data ent.CreateForPermissionInput) (upsertedForPermission *ent.ForPermission, err error) {
	if CompareUniqueInput(data, ent.ForPermissionWhereUniqueInput{}) {
		var where ent.ForPermissionWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}

		upsertedForPermission, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		if upsertedForPermission != nil {
			updateForPermissionInput := ent.UpdateForPermissionInput{}
			err = mapstructure.Decode(data, &updateForPermissionInput)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			upsertedForPermission, err = r.client.ForPermission.UpdateOne(upsertedForPermission).SetInput(updateForPermissionInput).Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			return upsertedForPermission, nil
		}

		err = r.client.ForPermission.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		
		return upsertedForPermission, nil
	}
	upsertedForPermission, err = r.client.ForPermission.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedForPermission, nil
}

func (r *ForPermission) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateForPermissionInput) (upsertedForPermission *ent.ForPermission, err error) {
	var where ent.ForPermissionWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}

		upsertedForPermission, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		if upsertedForPermission != nil {
			updateForPermissionInput := ent.UpdateForPermissionInput{}
			err = mapstructure.Decode(data, &updateForPermissionInput)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			upsertedForPermission, err = tx.ForPermission.UpdateOne(upsertedForPermission).SetInput(updateForPermissionInput).Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("error upserting item: %v", err)
			}
			return upsertedForPermission, nil
		}
	upsertedForPermission, err = tx.ForPermission.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedForPermission, nil
}

func (r *ForPermission) UpsertMany(ctx context.Context, data []ent.CreateForPermissionInput) (int, error) {
	var upsertMap []*ent.ForPermissionCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, r.client.ForPermission.Create().SetInput(upsertInput))
	}
	err := r.client.ForPermission.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *ForPermission) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateForPermissionInput) (int, error) {
	var upsertMap []*ent.ForPermissionCreate
	for _, upsertInput := range data {
        upsertMap = append(upsertMap, tx.ForPermission.Create().SetInput(upsertInput))
    }
	err := tx.ForPermission.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *ForPermission) Delete(ctx context.Context, where ent.ForPermissionWhereUniqueInput) (*ent.ForPermission, error) {
	query := r.client.ForPermission.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = r.client.ForPermission.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *ForPermission) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereUniqueInput) (*ent.ForPermission, error) {
	query := tx.ForPermission.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = tx.ForPermission.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *ForPermission) DeleteMany(ctx context.Context, where ent.ForPermissionWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := r.client.ForPermission.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}

func (r *ForPermission) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.ForPermissionWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := tx.ForPermission.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}