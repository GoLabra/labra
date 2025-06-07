package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoLabra/labra/src/api/constants"
	"github.com/GoLabra/labra/src/api/entgql/ent"
	"github.com/mitchellh/mapstructure"
)

type File struct {
	client *ent.Client
}

func NewFile(c *ent.Client) *File {
	return &File{client: c}
}

func (r *File) Get(ctx context.Context, where *ent.FileWhereInput, orderBy *ent.FileOrder, skip *int, first *int, last *int) ([]*ent.File, error) {
	var (
		query = r.client.File.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.FileOrder{
			Direction: ent.DefaultFileOrder.Direction,
			Field:     ent.DefaultFileOrder.Field,
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

func (r *File) Connection(ctx context.Context, where *ent.FileWhereInput, orderBy *ent.FileOrder, skip *int, first *int, last *int) (*ent.FileConnection, error) {
	var (
		query = r.client.File.Query()
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

func (r *File) GetOne(ctx context.Context, where ent.FileWhereUniqueInput) (*ent.File, error) {
	query := r.client.File.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}

func (r *File) GetTx(ctx context.Context, tx *ent.Tx, where *ent.FileWhereInput, orderBy *ent.FileOrder, skip *int, first *int, last *int) ([]*ent.File, error) {
	var (
		query = tx.File.Query()
		err   error
	)
	if where != nil {
		query, err = where.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if orderBy == nil {
		orderBy = &ent.FileOrder{
			Direction: ent.DefaultFileOrder.Direction,
			Field:     ent.DefaultFileOrder.Field,
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

func (r *File) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereUniqueInput) (*ent.File, error) {
	query := tx.File.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	return query.First(ctx)
}
func (r *File) Create(ctx context.Context, data ent.CreateFileInput) (*ent.File, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)

	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	createdInput, err := repository.File.CreateTx(ctx, tx, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return createdInput, err
}

func (r *File) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateFileInput) (*ent.File, error) {
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

			data.CreatedByID = &toConnect.ID
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
	}

	createdInput, err := tx.File.Create().SetInput(data).Save(ctx)

	if err != nil {
		return nil, err
	}

	return createdInput.Unwrap(), err
}

func (r *File) CreateMany(ctx context.Context, data []ent.CreateFileInput) ([]*ent.File, error) {
	var createMap []*ent.FileCreate
	for _, createInput := range data {
		createMap = append(createMap, r.client.File.Create().SetInput(createInput))
	}
	createdItems, err := r.client.File.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}

func (r *File) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateFileInput) ([]*ent.File, error) {
	var createMap []*ent.FileCreate
	for _, createInput := range data {
		createMap = append(createMap, tx.File.Create().SetInput(createInput))
	}
	createdItems, err := tx.File.CreateBulk(createMap...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating items: %v", err)
	}
	return createdItems, nil
}
func (r *File) Update(ctx context.Context, where ent.FileWhereUniqueInput, data ent.UpdateFileInput) (*ent.File, error) {
	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	tx, err := repository.Tx.Create(ctx)

	if err != nil {
		return nil, err
	}

	updatedInput, err := repository.File.UpdateTx(ctx, tx, where, data)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return updatedInput.Unwrap(), nil
}

func (r *File) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereUniqueInput, data ent.UpdateFileInput) (*ent.File, error) {

	repository, ok := ctx.Value(constants.RepositoryContextValue).(*Repository)
	if !ok {
		return nil, errors.New(ErrRepositoryNotSetInContext)
	}

	query := tx.File.Query()
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

			data.CreatedByID = &toConnect.ID
		}
		if data.CreatedBy.Unset != nil && *data.CreatedBy.Unset {
			data.ClearCreatedBy = true
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
		if data.UpdatedBy.Unset != nil && *data.UpdatedBy.Unset {
			data.ClearUpdatedBy = true
		}
	}

	updatedInput, err := tx.File.UpdateOne(item).SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error updating item: %v", err)
	}

	return updatedInput, nil
}

func (r *File) UpdateMany(ctx context.Context, where ent.FileWhereInput, data ent.UpdateFileInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := r.client.File.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *File) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereInput, data ent.UpdateFileInput) (int, error) {
	ps, err := where.P()

	if err != nil {
		return 0, err
	}

	updatedRows, err := tx.File.Update().Where(ps).SetInput(data).Save(ctx)
	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

func (r *File) Upsert(ctx context.Context, data ent.CreateFileInput) (upsertedFile *ent.File, err error) {
	if CompareUniqueInput(data, ent.FileWhereUniqueInput{}) {
		var where ent.FileWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}
		err = r.client.File.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		upsertedFile, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		return upsertedFile, nil
	}
	upsertedFile, err = r.client.File.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedFile, nil
}

func (r *File) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateFileInput) (upsertedFile *ent.File, err error) {
	if CompareUniqueInput(data, ent.FileWhereUniqueInput{}) {
		var where ent.FileWhereUniqueInput
		err = mapstructure.Decode(data, &where)
		if err != nil {
			return nil, fmt.Errorf("error decoding where condition: %v", err)
		}
		err = tx.File.Create().SetInput(data).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("error upserting item: %v", err)
		}
		upsertedFile, err = r.GetOne(ctx, where)
		if err != nil {
			return nil, fmt.Errorf("error getting upserted item: %v", err)
		}
		return upsertedFile, nil
	}
	upsertedFile, err = tx.File.Create().SetInput(data).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error upserting item: %v", err)
	}
	return upsertedFile, nil
}

func (r *File) UpsertMany(ctx context.Context, data []ent.CreateFileInput) (int, error) {
	var upsertMap []*ent.FileCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, r.client.File.Create().SetInput(upsertInput))
	}
	err := r.client.File.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *File) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateFileInput) (int, error) {
	var upsertMap []*ent.FileCreate
	for _, upsertInput := range data {
		upsertMap = append(upsertMap, tx.File.Create().SetInput(upsertInput))
	}
	err := tx.File.CreateBulk(upsertMap...).OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error upserting items: %v", err)
	}
	return len(data), nil
}

func (r *File) Delete(ctx context.Context, where ent.FileWhereUniqueInput) (*ent.File, error) {
	query := r.client.File.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = r.client.File.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *File) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereUniqueInput) (*ent.File, error) {
	query := tx.File.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %v", err)
	}
	item, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting item to delete: %v", err)
	}

	err = tx.File.DeleteOne(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error deleting item: %v", err)
	}
	return item, err
}

func (r *File) DeleteMany(ctx context.Context, where ent.FileWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := r.client.File.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}

func (r *File) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereInput) (int, error) {
	ps, err := where.P()
	if err != nil {
		return 0, err
	}

	deletedRows, err := tx.File.Delete().Where(ps).Exec(ctx)

	if err != nil {
		return 0, err
	}
	return deletedRows, nil
}
