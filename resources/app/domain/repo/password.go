package repo

import (
	"app/ent"
	"context"
	"fmt"

	"github.com/GoLabra/labra/constants"
)

func (r *User) UpdatePassword(ctx context.Context, where ent.UserWhereUniqueInput, hashedPassword string) (*ent.User, error) {
	iCtx := context.WithValue(ctx, constants.IsInternalOperationContextValue, true)
	query := r.client.User.Query()
	query, err := where.Filter(query)
	if err != nil {
		return nil, fmt.Errorf("error applying unique where condition: %w", err)
	}
	item, err := query.First(iCtx)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return r.client.User.UpdateOne(item).SetPassword(hashedPassword).Save(ctx)
}
