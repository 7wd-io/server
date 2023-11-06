package repo

import (
	"7wd.io/adapter/repo/internal/pg"
	"7wd.io/domain"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

func NewGame(c *pgxpool.Pool) GameRepo {
	return GameRepo{
		R: pg.R{
			PG: c,
			QB: pg.QB{
				TableName: "game",
				Columns: []string{
					"id",
					"host_nickname",
					"host_rating",
					"host_points",
					"guest_nickname",
					"guest_rating",
					"guest_points",
					"winner",
					"victory",
					"log",
					"started_at",
					"finished_at",
				},
			},
		},
	}
}

type GameRepo struct {
	pg.R
}

func (dst GameRepo) conn(t domain.Tx) pg.Conn {
	if t == nil {
		return dst.PG
	}

	tx, ok := t.Value().(pgx.Tx)

	if ok {
		return tx
	}

	slog.Warn("tx extract !ok")

	return dst.PG
}

func (dst GameRepo) findOneBy(ctx context.Context, o ...domain.GameOption) (*domain.Game, error) {
	var err error
	out := new(domain.Game)

	c, w := dst.c(o...)

	err = c.
		QueryRow(
			ctx,
			dst.QB.SelectWhere(w),
			w.Values()...,
		).
		Scan(
			&out.Id,
			&out.HostNickname,
			&out.HostRating,
			&out.HostPoints,
			&out.GuestNickname,
			&out.GuestRating,
			&out.GuestPoints,
			&out.Winner,
			&out.Victory,
			&out.Log,
			&out.StartedAt,
			&out.FinishedAt,
		)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrGameNotFound
		}

		return nil, err
	}

	return out, nil
}

func (dst GameRepo) Save(ctx context.Context, in *domain.Game, o ...domain.GameOption) error {
	c, _ := dst.c(o...)

	return c.
		QueryRow(
			ctx,
			dst.QB.Insert(),
			in.HostNickname,
			in.HostRating,
			in.HostPoints,
			in.GuestNickname,
			in.GuestRating,
			in.GuestPoints,
			in.Winner,
			in.Victory,
			in.Log,
			in.StartedAt,
			in.FinishedAt,
		).Scan(&in.Id)
}

func (dst GameRepo) Update(ctx context.Context, in *domain.Game, o ...domain.GameOption) error {
	c, _ := dst.c(o...)

	_, err := c.
		Exec(
			ctx,
			dst.QB.Update(),
			in.Id,
			in.HostNickname,
			in.HostRating,
			in.HostPoints,
			in.GuestNickname,
			in.GuestRating,
			in.GuestPoints,
			in.Winner,
			in.Victory,
			in.Log,
			in.StartedAt,
			in.FinishedAt,
		)

	return err
}

func (dst GameRepo) Find(ctx context.Context, o ...domain.GameOption) (*domain.Game, error) {
	return dst.findOneBy(ctx, o...)
}

func (dst GameRepo) c(o ...domain.GameOption) (c pg.Conn, w pg.Where) {
	opts := new(domain.GameOptions)

	for _, v := range o {
		v(opts)
	}

	if opts.Id != 0 {
		w = append(w, pg.F{
			Expr:  "id",
			Value: opts.Id,
		})
	}

	return dst.conn(opts.Tx), w
}
