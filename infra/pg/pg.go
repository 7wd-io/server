package pg

import (
	"7wd.io/config"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
