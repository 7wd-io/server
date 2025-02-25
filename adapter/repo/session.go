package repo

import (
	"7wd.io/adapter/repo/internal/rds"
	"7wd.io/domain"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewSession(c *redis.Client) SessionRepo {
	return SessionRepo{
		R: rds.R{Rds: c},
	}
}

type SessionRepo struct {
	rds.R
}

func (dst SessionRepo) Save(ctx context.Context, s *domain.Session, ttl time.Duration) error {
	return dst.Set(ctx, dst.k(s.Client), s, ttl)
}

func (dst SessionRepo) Delete(ctx context.Context, client uuid.UUID) (*domain.Session, error) {
	s, err := dst.Find(ctx, client)

	if err != nil {
		return nil, err
	}

	if s == nil {
		return nil, domain.ErrSessionNotFound
	}

	err = dst.Rds.Del(ctx, dst.k(s.Client)).Err()

	if errors.Is(err, redis.Nil) {
		return s, nil
	}

	return s, err
}

func (dst SessionRepo) Find(ctx context.Context, client uuid.UUID) (*domain.Session, error) {
	s := new(domain.Session)
	err := dst.Get(ctx, dst.k(client), s)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	return s, nil
}

func (dst SessionRepo) k(fp uuid.UUID) string {
	return fmt.Sprintf("session:%s", fp)
}
