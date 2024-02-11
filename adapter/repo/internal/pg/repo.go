package pg

import (
	"7wd.io/domain"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Conn type
type Conn interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type R struct {
	PG        *pgxpool.Pool
	TableName string
	Columns   []string
}

func (dst R) Conn(t domain.Tx) Conn {
	if t == nil {
		return dst.PG
	}

	tx, ok := t.Value().(pgx.Tx)

	if ok {
		return tx
	}

	// log !ok?

	return dst.PG
}
