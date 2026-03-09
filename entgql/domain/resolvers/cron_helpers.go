package resolvers

import (
	"context"
	"database/sql"
	"errors"

	"github.com/GoLabra/labra/constants"
)

type sqlExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	DialectName() string
}

func cronDB(ctx context.Context) (sqlExecutor, context.Context, error) {
	client, ok := ctx.Value(constants.AdminEntClientContextValue).(sqlExecutor)
	if !ok || client == nil {
		return nil, nil, errors.New("admin ent client not found in context")
	}
	return client, constants.WithInternalOperation(ctx), nil
}
