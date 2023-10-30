package tx

import (
	"7wd.io/domain"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(c *pgxpool.Pool) *Tx {
	return &Tx{c: c}
}

type Tx struct {
	c *pgxpool.Pool
}

func (dst *Tx) Tx(ctx context.Context) (domain.Tx, error) {
	v, err := dst.c.Begin(ctx)

	if err != nil {
		return nil, err
	}

	return &tx{v: v}, nil
}

type tx struct {
	v pgx.Tx
}

func (dst *tx) Rollback(ctx context.Context) error {
	return dst.v.Rollback(ctx)
}

func (dst *tx) Commit(ctx context.Context) error {
	return dst.v.Commit(ctx)
}

func (dst *tx) Value() any {
	return dst.v
}
