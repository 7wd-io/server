package repo

import (
	"7wd.io/adapter/repo/internal/pg"
	"7wd.io/domain"
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewUser(c *pgxpool.Pool) UserRepo {
	return UserRepo{
		R: pg.R{
			PG:        c,
			TableName: `"user"`,
			Columns: []string{
				"id",
				"email",
				"nickname",
				"password",
				"settings",
				"rating",
				"created_at",
			},
		},
	}
}

type UserRepo struct {
	pg.R
}

func (dst UserRepo) Save(ctx context.Context, in *domain.User, o ...domain.UserOption) error {
	opts := dst.opts(o...)

	q, args, err := sq.
		Insert(dst.TableName).
		// skip id column
		Columns(dst.Columns[1:]...).
		Values(
			in.Email,
			in.Nickname,
			in.Password,
			in.Settings,
			in.Rating,
			in.CreatedAt,
		).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		return err
	}

	return dst.Conn(opts.Tx).QueryRow(ctx, q, args...).Scan(&in.Id)
}

func (dst UserRepo) Update(ctx context.Context, in *domain.User, o ...domain.UserOption) error {
	opts := dst.opts(o...)

	q, args, err := sq.
		Update(dst.TableName).
		Where(sq.Eq{"id": in.Id}).
		Set("email", in.Email).
		Set("nickname", in.Nickname).
		Set("password", in.Password).
		Set("settings", in.Settings).
		Set("rating", in.Rating).
		Set("created_at", in.CreatedAt).
		ToSql()

	if err != nil {
		return err
	}

	_, err = dst.Conn(opts.Tx).Exec(ctx, q, args...)

	return err
}

func (dst UserRepo) Find(ctx context.Context, o ...domain.UserOption) (*domain.User, error) {
	users, err := dst.FindMany(ctx, o...)

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, domain.ErrUserNotFound
	}

	return users[0], nil
}

func (dst UserRepo) FindMany(ctx context.Context, o ...domain.UserOption) ([]*domain.User, error) {
	opts := dst.opts(o...)

	sql, args, err := dst.selectb(opts).ToSql()

	if err != nil {
		return nil, err
	}

	var out []*domain.User

	if err = pgxscan.Select(ctx, dst.Conn(opts.Tx), &out, sql, args...); err != nil {
		return nil, err
	}

	return out, nil
}

func (dst UserRepo) opts(opts ...domain.UserOption) domain.UserOptions {
	o := new(domain.UserOptions)

	for _, v := range opts {
		v(o)
	}

	return *o
}

func (dst UserRepo) selectb(o domain.UserOptions) sq.SelectBuilder {
	sb := sq.Select(dst.Columns...).
		From(dst.TableName).
		Where(dst.where(o)).
		OrderBy("id ASC")

	if o.Lock {
		sb.Suffix("FOR UPDATE")
	}

	return sb
}

func (UserRepo) where(o domain.UserOptions) sq.Eq {
	w := sq.Eq{}

	if o.IdSet {
		w["id"] = o.Id
	}

	if o.EmailSet {
		w["email"] = o.Email
	}

	if o.NicknameSet {
		w["nickname"] = o.Nickname
	}

	return w
}
