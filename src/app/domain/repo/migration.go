package repo

import (
	"context"
	"errors"
	"fmt"
	"app/ent"

	"github.com/GoLabra/labrago/src/api/constants"
	"github.com/mitchellh/mapstructure"
)

type Migration struct {
	client *ent.Client
}

func NewMigration(c *ent.Client) *Migration {
	return &Migration{client: c}
}

func (r *Migration) Get(ctx context.Context, where *ent.MigrationWhereInput, orderBy *ent.MigrationOrder, skip *int, first *int, last *int) ([]*ent.Migration, error) {
	var (
        query = r.client.Migration.Query()
        err   error
    )
    if where != nil {
        query, err = where.Filter(query)
        if err != nil {
            return nil, err
        }
    }

    if orderBy == nil {
		orderBy = &ent.MigrationOrder{
			Direction: ent.DefaultMigrationOrder.Direction,
			Field:     ent.DefaultMigrationOrder.Field,
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

func (r *Migration) Connection(ctx context.Context, where *ent.MigrationWhereInput, orderBy *ent.MigrationOrder, skip *int, first *int, last *int) (*ent.MigrationConnection, error) {
	var (
		query = r.client.Migration.Query()
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

func (r *Migration) GetOne(ctx context.Context, where ent.MigrationWhereUniqueInput) (*ent.Migration, error) {
	query := r.client.Migration.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}

func (r *Migration) GetTx(ctx context.Context, tx *ent.Tx, where *ent.MigrationWhereInput, orderBy *ent.MigrationOrder, skip *int, first *int, last *int) ([]*ent.Migration, error) {
	var (
		query = tx.Migration.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.MigrationOrder{
			Direction: ent.DefaultMigrationOrder.Direction,
			Field:     ent.DefaultMigrationOrder.Field,
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

func (r *Migration) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereUniqueInput) (*ent.Migration, error) {
	query := tx.Migration.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}


func (r *Migration) Create(ctx context.Context, data ent.CreateMigrationInput) (*ent.Migration, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)

	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	createdInput, err := repository.Migration.CreateTx(ctx, tx, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return createdInput, err
}

func (r *Migration) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateMigrationInput) (*ent.Migration, error) {
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

	createdInput, err := tx.Migration.Create().SetInput(data).Save(ctx)

	if err != nil {
		return nil, err
	}

	return createdInput.Unwrap(), err
}

func (r *Migration) CreateMany(ctx context.Context, data []ent.CreateMigrationInput) ([]*ent.Migration, error) {
	var createMap []*ent.MigrationCreate
	for _, createInput := range data {
        createMap = append(createMap, r.client.Migration.Create().SetInput(createInput))
    }
	createdItems, err := r.client.Migration.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}

func (r *Migration) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateMigrationInput) ([]*ent.Migration, error) {
	var createMap []*ent.MigrationCreate
	for _, createInput := range data {
        createMap = append(createMap, tx.Migration.Create().SetInput(createInput))
    }
	createdItems, err := tx.Migration.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}

func (r *Migration) Update(ctx context.Context, where ent.MigrationWhereUniqueInput, data ent.UpdateMigrationInput) (*ent.Migration, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	updatedInput, err := repository.Migration.UpdateTx(ctx, tx, where, data)

	if err != nil {
		return nil, err
	}
	
	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return updatedInput.Unwrap(), nil
}

func (r *Migration) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereUniqueInput, data ent.UpdateMigrationInput) (*ent.Migration, error) {
	
    repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}
	
	query := tx.Migration.Query()
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

	updatedInput, err := tx.Migration.UpdateOne(item).SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error updating item: %v", err)
	}

	return updatedInput, nil
}

func (r *Migration) UpdateMany(ctx context.Context, where ent.MigrationWhereInput, data ent.UpdateMigrationInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := r.client.Migration.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *Migration) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereInput, data ent.UpdateMigrationInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := tx.Migration.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *Migration) Upsert(ctx context.Context, data ent.CreateMigrationInput) (upsertedMigration *ent.Migration, err error) {
	if CompareUniqueInput(data, ent.MigrationWhereUniqueInput{}) {
		var where ent.MigrationWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}
		err = r.client.Migration.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		upsertedMigration, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		return upsertedMigration, nil
	}
	upsertedMigration, err = r.client.Migration.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedMigration, nil
}

func (r *Migration) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateMigrationInput) (upsertedMigration *ent.Migration, err error) {
	if CompareUniqueInput(data, ent.MigrationWhereUniqueInput{}) {
		var where ent.MigrationWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}
		err = tx.Migration.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		upsertedMigration, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		return upsertedMigration, nil
	}
	upsertedMigration, err = tx.Migration.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedMigration, nil
}

func (r *Migration) UpsertMany(ctx context.Context, data []ent.CreateMigrationInput) (int, error) {
	var upsertMap []*ent.MigrationCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, r.client.Migration.Create().SetInput(upsertInput))
	}
	err := r.client.Migration.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *Migration) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateMigrationInput) (int, error) {
	var upsertMap []*ent.MigrationCreate
	for _, upsertInput := range data {
        upsertMap = append(upsertMap, tx.Migration.Create().SetInput(upsertInput))
    }
	err := tx.Migration.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *Migration) Delete(ctx context.Context, where ent.MigrationWhereUniqueInput) (*ent.Migration, error) {
	query := r.client.Migration.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = r.client.Migration.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *Migration) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereUniqueInput) (*ent.Migration, error) {
	query := tx.Migration.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = tx.Migration.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *Migration) DeleteMany(ctx context.Context, where ent.MigrationWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := r.client.Migration.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}

func (r *Migration) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.MigrationWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := tx.Migration.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}