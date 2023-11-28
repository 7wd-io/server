package analyst

import (
	"7wd.io/domain"
	"context"
	"encoding/json"
	"fmt"
	swde "github.com/7wd-io/engine"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"
)

func New(
	rds *redis.Client,
	pg *pgxpool.Pool,
) A {
	return A{
		rds: rds,
		pg:  pg,
		key: "ratings",
	}
}

type A struct {
	rds *redis.Client
	pg  *pgxpool.Pool
	key string
}

func (dst A) Top(ctx context.Context) (domain.Top, error) {
	members, err := dst.rds.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
		Key:   dst.key,
		Start: 0,
		Stop:  7, // show top 7, but +1 slot for bot
		Rev:   true,
	}).Result()

	if err != nil {
		return nil, err
	}

	var top domain.Top

	for _, m := range members {
		nickname := domain.Nickname(m.Member.(string))

		// bot is not in any ratings, skip
		if nickname == domain.BotNickname {
			continue
		}

		top = append(top, domain.TopMember{
			Name:   nickname,
			Rating: domain.Rating(m.Score),
		})
	}

	return top, nil
}

func (dst A) UpdateRatings(ctx context.Context, u *domain.User) error {
	if err := dst.rds.Del(ctx, dst.kGames(u.Nickname)).Err(); err != nil {
		return nil
	}

	return dst.rds.ZAdd(ctx, dst.key, redis.Z{
		Score:  float64(u.Rating),
		Member: string(u.Nickname),
	}).Err()
}

func (dst A) Ratings(ctx context.Context, u ...domain.Nickname) (domain.UsersPreview, error) {
	if err := dst.checkOrRefresh(ctx); err != nil {
		return nil, err
	}

	var members = make([]string, len(u))

	for k, v := range u {
		members[k] = string(v)
	}

	scores, err := dst.rds.ZMScore(ctx, dst.key, members...).Result()

	if err != nil {
		return nil, err
	}

	var up = make(domain.UsersPreview, len(u))

	for k, v := range u {
		up[v] = domain.Rating(scores[k])
	}

	return up, nil
}

func (dst A) GamesReport(ctx context.Context, u domain.Nickname) (*domain.GamesReport, error) {
	key := dst.kGames(u)

	found, err := dst.rds.Exists(ctx, key).Result()

	if err != nil {
		return nil, err
	}

	var gr *domain.GamesReport

	if found == 0 {
		gr, err = dst.gamesReport(ctx, u)

		if err != nil {
			return nil, err
		}

		if err = dst.setValue(ctx, key, gr, time.Hour*24*90); err != nil {
			return nil, err
		}
	}

	gr = new(domain.GamesReport)

	if err = dst.getValue(ctx, key, gr); err != nil {
		return nil, err
	}

	return gr, nil
}

func (dst A) GamesReportVersus(ctx context.Context, me domain.Nickname, enemy domain.Nickname) (*domain.GamesReport, error) {
	const sql = `
WITH games as (
    SELECT *
    FROM game
    WHERE
        (host_nickname = $1 AND guest_nickname = $2)
        OR (host_nickname = $2 AND guest_nickname = $1)
		AND winner IS NOT NULL
), won as (
    SELECT
       COUNT(id) as total,
       COUNT(CASE WHEN victory = $3 THEN id END) as points,
       COUNT(CASE WHEN victory = $4 THEN id END) as military,
       COUNT(CASE WHEN victory = $5 THEN id END) as science,
       COUNT(CASE WHEN victory = $6 THEN id END) as resign,
       COUNT(CASE WHEN victory = $7 THEN id END) as timeout
    FROM games
    WHERE winner = $1
), lose as (
    SELECT
       COUNT(id) as total,
       COUNT(CASE WHEN victory = $3 THEN id END) as points,
       COUNT(CASE WHEN victory = $4 THEN id END) as military,
       COUNT(CASE WHEN victory = $5 THEN id END) as science,
       COUNT(CASE WHEN victory = $6 THEN id END) as resign,
       COUNT(CASE WHEN victory = $7 THEN id END) as timeout
    FROM games
    WHERE winner != $1
)
    SELECT * FROM won, lose
`
	gr := new(domain.GamesReport)

	err := dst.pg.QueryRow(
		ctx,
		sql,
		me,
		enemy,
		swde.Civilian,
		swde.MilitarySupremacy,
		swde.ScienceSupremacy,
		swde.Resign,
		swde.Timeout,
	).
		Scan(
			&gr.Won.Total,
			&gr.Won.Points,
			&gr.Won.Military,
			&gr.Won.Science,
			&gr.Won.Resign,
			&gr.Won.Timeout,
			&gr.Lose.Total,
			&gr.Lose.Points,
			&gr.Lose.Military,
			&gr.Lose.Science,
			&gr.Lose.Resign,
			&gr.Lose.Timeout,
		)

	if err != nil {
		return nil, err
	}

	return gr, nil
}

func (dst A) Rank(ctx context.Context, u domain.Nickname) (int, error) {
	// 0-based
	rank, err := dst.rds.ZRevRank(ctx, dst.key, string(u)).Result()

	if err != nil {
		return 0, err
	}

	return int(rank) + 1, nil
}

func (dst A) Rating(ctx context.Context, u domain.Nickname) (domain.Rating, error) {
	rating, err := dst.rds.ZScore(ctx, dst.key, string(u)).Result()

	if err != nil {
		return 0, err
	}

	return domain.Rating(rating), nil
}

func (dst A) gamesReport(ctx context.Context, u domain.Nickname) (*domain.GamesReport, error) {
	const sql = `
WITH games as (
    SELECT *
    FROM game
    WHERE
		($1 = $2 AND guest_nickname = $2)
		OR ($1 != $2 AND ((host_nickname = $1 AND guest_nickname != $2) OR guest_nickname = $1))
		AND winner IS NOT NULL
), won as (
    SELECT
       COUNT(id) as total,
       COUNT(CASE WHEN victory = $3 THEN id END) as points,
       COUNT(CASE WHEN victory = $4 THEN id END) as military,
       COUNT(CASE WHEN victory = $5 THEN id END) as science,
       COUNT(CASE WHEN victory = $6 THEN id END) as resign,
       COUNT(CASE WHEN victory = $7 THEN id END) as timeout
    FROM games
    WHERE winner = $1
), lose as (
    SELECT
       COUNT(id) as total,
       COUNT(CASE WHEN victory = $3 THEN id END) as points,
       COUNT(CASE WHEN victory = $4 THEN id END) as military,
       COUNT(CASE WHEN victory = $5 THEN id END) as science,
       COUNT(CASE WHEN victory = $6 THEN id END) as resign,
       COUNT(CASE WHEN victory = $7 THEN id END) as timeout
    FROM games
    WHERE winner != $1
)
    SELECT * FROM won, lose
`
	gr := new(domain.GamesReport)

	err := dst.pg.QueryRow(
		ctx,
		sql,
		u,
		domain.BotNickname,
		swde.Civilian,
		swde.MilitarySupremacy,
		swde.ScienceSupremacy,
		swde.Resign,
		swde.Timeout,
	).
		Scan(
			&gr.Won.Total,
			&gr.Won.Points,
			&gr.Won.Military,
			&gr.Won.Science,
			&gr.Won.Resign,
			&gr.Won.Timeout,
			&gr.Lose.Total,
			&gr.Lose.Points,
			&gr.Lose.Military,
			&gr.Lose.Science,
			&gr.Lose.Resign,
			&gr.Lose.Timeout,
		)

	if err != nil {
		return nil, err
	}

	return gr, nil
}

func (dst A) checkOrRefresh(ctx context.Context) error {
	found, err := dst.rds.Exists(ctx, dst.key).Result()

	if err != nil {
		return err
	}

	if found != 0 {
		return nil
	}

	return dst.refreshRatings(ctx)
}

func (dst A) refreshRatings(ctx context.Context) error {
	rows, err := dst.pg.Query(ctx, `SELECT "nickname","rating" from "user"`)
	defer rows.Close()

	if err != nil {
		return err
	}

	var nickname domain.Nickname
	var rating domain.Rating

	var members []redis.Z

	for rows.Next() {
		if err = rows.Scan(&nickname, &rating); err != nil {
			return err
		}

		members = append(members, redis.Z{
			Score:  float64(rating),
			Member: string(nickname),
		})
	}

	if err = dst.rds.ZAdd(ctx, dst.key, members...).Err(); err != nil {
		return err
	}

	return nil
}

func (dst A) setValue(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	v, err := json.Marshal(value)

	if err != nil {
		return nil
	}

	return dst.rds.Set(ctx, key, v, ttl).Err()
}

func (dst A) getValue(ctx context.Context, key string, dest interface{}) error {
	v, err := dst.rds.Get(ctx, key).Bytes()

	if err != nil {
		return err
	}

	return json.Unmarshal(v, dest)
}

func (dst A) kGames(u domain.Nickname) string {
	return fmt.Sprintf("user:%s:games", u)
}
