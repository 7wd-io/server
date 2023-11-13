package repo

import (
	"7wd.io/adapter/repo/internal/rds"
	"7wd.io/domain"
	"context"
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
	return dst.Set(ctx, dst.k(s.Fingerprint), s, ttl)
}

func (dst SessionRepo) Delete(ctx context.Context, fingerprint uuid.UUID) (*domain.Session, error) {
	s, err := dst.Find(ctx, fingerprint)

	if err != nil {
		return nil, err
	}

	if s == nil {
		return nil, domain.ErrSessionNotFound
	}

	return s, dst.Rds.Del(ctx, dst.k(s.Fingerprint)).Err()
}

func (dst SessionRepo) Find(ctx context.Context, fingerprint uuid.UUID) (*domain.Session, error) {
	s := new(domain.Session)
	err := dst.Get(ctx, dst.k(fingerprint), s)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func (dst SessionRepo) k(fp uuid.UUID) string {
	return fmt.Sprintf("session:%s", fp)
}
