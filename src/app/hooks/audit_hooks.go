package hooks

import (
	"app/ent"
	"context"

	"github.com/GoLabra/labrago/src/api/constants"
)

func CreatedByUpdatedByHook(next ent.Mutator) ent.Mutator {
	var currentUserId string
	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		currentUser, ok := ctx.Value(constants.UserContextValue).(*ent.User)
		if ok {
			currentUserId = currentUser.ID
		} else {
			return next.Mutate(ctx, m) // TODO clear up operation without user
		}
		if m.Op().Is(ent.OpCreate) {
			if ml, ok := m.(interface{ SetCreatedByID(string) }); ok {
				ml.SetCreatedByID(currentUserId)
			}
		} else if m.Op().Is(ent.OpUpdateOne) || m.Op().Is(ent.OpUpdate) {
			if ml, ok := m.(interface{ SetUpdatedByID(string) }); ok {
				ml.SetUpdatedByID(currentUserId)
			}
		}
		return next.Mutate(ctx, m)
	})
}
