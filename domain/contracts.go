package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type (
	Clock interface {
		Now() time.Time
	}

	Txer interface {
		Tx(context.Context) (Tx, error)
	}

	Tx interface {
		Rollback(context.Context) error
		Commit(context.Context) error
		Value() any
	}

	Password interface {
		Hash(password string, cost int) (string, error)
		Check(hash string, password string) bool
	}

	Uuidf interface {
		Uuid() uuid.UUID
	}

	Tokenf interface {
		Token(*Passport) (string, error)
	}

	Pusher interface {
		Push(msg interface{})
	}

	Onliner interface {
		Online() []Nickname
	}
)

type (
	UserRepo interface {
		Save(context.Context, *User, ...UserOption) error
		Update(context.Context, *User, ...UserOption) error
		Find(context.Context, ...UserOption) (*User, error)
	}

	SessionRepo interface {
		Save(ctx context.Context, s *Session, ttl time.Duration) error
		Delete(ctx context.Context, fingerprint uuid.UUID) (*Session, error)
		Find(ctx context.Context, fingerprint uuid.UUID) (*Session, error)
	}

	RoomRepo interface {
		Save(context.Context, *Room) error
		Delete(context.Context, RoomId) (*Room, error)
		Find(context.Context, RoomId) (*Room, error)
		FindAll(ctx context.Context) ([]*Room, error)
	}
)
