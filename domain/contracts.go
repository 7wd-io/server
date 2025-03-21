package domain

import (
	"context"
	swde "github.com/7wd-io/engine"
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
		Online(context.Context) ([]Nickname, error)
	}

	Dispatcher interface {
		Dispatch(ctx context.Context, event EventId, payload interface{})
	}

	Bot interface {
		GetMove(*Game) (swde.Mutator, error)
	}

	Mover interface {
		Move(ctx context.Context, u Nickname, id GameId, m swde.Mutator) (*Game, error)
	}

	Analyst interface {
		Top(context.Context) (Top, error)
		UpdateRatings(context.Context, *User) error
		Ratings(context.Context, ...Nickname) (UsersPreview, error)
		GamesReport(context.Context, Nickname) (*GamesReport, error)
		GamesReportVersus(context.Context, Nickname, Nickname) (*GamesReport, error)
		Rank(context.Context, Nickname) (int, error)
		Rating(context.Context, Nickname) (Rating, error)
	}

	PlayAgainStore interface {
		Create(context.Context, Game, RoomOptions) error
		Update(context.Context, GameId, Nickname, bool) (*PlayAgainAgreement, error)
	}

	GameCreator interface {
		Create(ctx context.Context, host User, guest User, o RoomOptions) (*Game, error)
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
		Delete(ctx context.Context, client uuid.UUID) (*Session, error)
		Find(ctx context.Context, client uuid.UUID) (*Session, error)
	}

	RoomRepo interface {
		Save(context.Context, *Room) error
		Delete(context.Context, RoomId) (*Room, error)
		Find(context.Context, RoomId) (*Room, error)
		FindByGame(context.Context, GameId) (*Room, error)
		FindAll(ctx context.Context) ([]*Room, error)
	}

	GameRepo interface {
		Save(context.Context, *Game, ...GameOption) error
		Update(context.Context, *Game, ...GameOption) error
		Find(context.Context, ...GameOption) (*Game, error)
	}

	GameClockRepo interface {
		Save(context.Context, *GameClock) error
		Delete(context.Context, GameId) error
		Find(context.Context, GameId) (*GameClock, error)
	}
)
