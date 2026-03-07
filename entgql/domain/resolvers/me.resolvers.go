package resolvers

import (
	"context"
	"fmt"

	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/ent"
)

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*ent.AdminUser, error) {
	user, ok := ctx.Value(constants.UserContextValue).(*ent.AdminUser)
	if !ok || user == nil {
		return nil, fmt.Errorf("unauthorized: no authenticated admin user found")
	}
	return user, nil
}
