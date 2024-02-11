package pg

import (
	"7wd.io/config"
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lann/builder"
)

func init() {
	// global setup squirrel to work with postgres
	sq.StatementBuilder = sq.StatementBuilderType(builder.EmptyBuilder).PlaceholderFormat(sq.Dollar)
}

func MustNew(ctx context.Context) *pgxpool.Pool {
	poolcfg, err := pgxpool.ParseConfig(config.C.PgDsn())

	if err != nil {
		panic("pg unable to parse config")
	}

	pg, err := pgxpool.NewWithConfig(ctx, poolcfg)

	if err != nil {
		panic("pg unable to create connection pool")
	}

	return pg
}
