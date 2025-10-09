package repo

import (
	"app/ent"
	"context"
)

type Tx struct {
	client *ent.Client
}

func NewTx(c *ent.Client) *Tx {
	return &Tx{client: c}
}

func (r *Tx) Create(ctx context.Context) (*ent.Tx, error) {
	return r.client.BeginTx(ctx, nil)
}

func (r *Tx) Commit(ctx context.Context, tx *ent.Tx) error {
	return tx.Commit()
}
