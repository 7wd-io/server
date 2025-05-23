package repo

import (
	"7wd.io/adapter/repo/internal/pg"
	"7wd.io/domain"
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewGame(c *pgxpool.Pool) GameRepo {
	return GameRepo{
		R: pg.R{
			PG:        c,
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
	}
}

type GameRepo struct {
	pg.R
}

func (dst GameRepo) Save(ctx context.Context, in *domain.Game, o ...domain.GameOption) error {
	opts := dst.opts(o...)

	q, args, err := sq.
		Insert(dst.TableName).
		// skip id column
		Columns(dst.Columns[1:]...).
		Values(
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
		).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		return err
	}

	return dst.Conn(opts.Tx).QueryRow(ctx, q, args...).Scan(&in.Id)
}

func (dst GameRepo) Update(ctx context.Context, in *domain.Game, o ...domain.GameOption) error {
	opts := dst.opts(o...)

	q, args, err := sq.
		Update(dst.TableName).
		Where(sq.Eq{"id": in.Id}).
		Set("host_nickname", in.HostNickname).
		Set("host_rating", in.HostRating).
		Set("host_points", in.HostPoints).
		Set("guest_nickname", in.GuestNickname).
		Set("guest_rating", in.GuestRating).
		Set("guest_points", in.GuestPoints).
		Set("winner", in.Winner).
		Set("victory", in.Victory).
		Set("log", in.Log).
		Set("started_at", in.StartedAt).
		Set("finished_at", in.FinishedAt).
		ToSql()

	if err != nil {
		return err
	}

	_, err = dst.Conn(opts.Tx).Exec(ctx, q, args...)

	return err
}

func (dst GameRepo) Find(ctx context.Context, o ...domain.GameOption) (*domain.Game, error) {
	games, err := dst.FindMany(ctx, o...)

	if err != nil {
		return nil, err
	}

	if len(games) == 0 {
		return nil, domain.ErrGameNotFound
	}

	return games[0], nil
}

func (dst GameRepo) FindMany(ctx context.Context, o ...domain.GameOption) ([]*domain.Game, error) {
	opts := dst.opts(o...)

	sql, args, err := dst.selectb(opts).ToSql()

	if err != nil {
		return nil, err
	}

	var out []*domain.Game

	if err = pgxscan.Select(ctx, dst.Conn(opts.Tx), &out, sql, args...); err != nil {
		return nil, err
	}

	return out, nil
}

func (dst GameRepo) opts(opts ...domain.GameOption) domain.GameOptions {
	o := new(domain.GameOptions)

	for _, v := range opts {
		v(o)
	}

	return *o
}

func (dst GameRepo) selectb(o domain.GameOptions) sq.SelectBuilder {
	sb := sq.Select(dst.Columns...).
		From(dst.TableName).
		Where(dst.where(o)).
		OrderBy("id ASC")

	if o.Lock {
		sb.Suffix("FOR UPDATE")
	}

	return sb
}

func (GameRepo) where(o domain.GameOptions) sq.Eq {
	w := sq.Eq{}

	if o.IdSet {
		w["id"] = o.Id
	}

	return w
}
