package repo

import (
	"7wd.io/adapter/repo/internal/pg"
	"7wd.io/domain"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

func NewUser(c *pgxpool.Pool) UserRepo {
	return UserRepo{
		R: pg.R{
			PG: c,
			QB: pg.QB{
				TableName: `"user"`,
				Columns: []string{
					"id",
					"email",
					"nickname",
					"password",
					"created_at",
				},
			},
		},
	}
}

type UserRepo struct {
	pg.R
}

func (dst UserRepo) conn(t domain.Tx) pg.Conn {
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

func (dst UserRepo) findOneBy(ctx context.Context, o ...domain.UserOption) (*domain.User, error) {
	var err error
	out := new(domain.User)

	c, w := dst.c(o...)

	err = c.
		QueryRow(
			ctx,
			dst.QB.SelectWhere(w),
			w.Values()...,
		).
		Scan(
			&out.Id,
			&out.Email,
			&out.Nickname,
			&out.Password,
			&out.CreatedAt,
		)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return out, nil
}

func (dst UserRepo) Save(ctx context.Context, in *domain.User, o ...domain.UserOption) error {
	c, _ := dst.c(o...)

	return c.
		QueryRow(
			ctx,
			dst.QB.Insert(),
			in.Email,
			in.Nickname,
			in.Password,
			in.CreatedAt,
		).Scan(&in.Id)
}

func (dst UserRepo) Update(ctx context.Context, in *domain.User, o ...domain.UserOption) error {
	c, _ := dst.c(o...)

	_, err := c.
		Exec(
			ctx,
			dst.QB.Update(),
			in.Id,
			in.Email,
			in.Nickname,
			in.Password,
			in.CreatedAt,
		)

	return err
}

func (dst UserRepo) Find(ctx context.Context, o ...domain.UserOption) (*domain.User, error) {
	return dst.findOneBy(ctx, o...)
}

func (dst UserRepo) c(o ...domain.UserOption) (c pg.Conn, w pg.Where) {
	opts := new(domain.UserOptions)

	for _, v := range o {
		v(opts)
	}

	if opts.Id != 0 {
		w = append(w, pg.F{
			Expr:  "id",
			Value: opts.Id,
		})
	}

	if opts.Email != "" {
		w = append(w, pg.F{
			Expr:  "email",
			Value: opts.Email,
		})
	}

	if opts.Nickname != "" {
		w = append(w, pg.F{
			Expr:  "nickname",
			Value: opts.Nickname,
		})
	}

	return dst.conn(opts.Tx), w
}
