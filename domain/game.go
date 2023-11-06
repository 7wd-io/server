package domain

import (
	"context"
	"encoding/json"
	"fmt"
	swde "github.com/7wd-io/engine"
	"log"
	"math"
	"strconv"
	"time"
)

const (
	timeBankDefault = TimeBank(10 * time.Minute)
	timeBankFast    = TimeBank(3 * time.Minute)
	timeBankBot     = TimeBank(30 * time.Minute)
	timeBankIncr    = TimeBank(5 * time.Second)
)

type GameId int

func NewGame(host *User, guest *User, now time.Time) *Game {
	return &Game{
		HostNickname:  host.Nickname,
		HostRating:    host.Rating,
		HostPoints:    Elo(host.Rating, guest.Rating),
		GuestNickname: guest.Nickname,
		GuestRating:   guest.Rating,
		GuestPoints:   Elo(guest.Rating, host.Rating),
		Log: GameLog{
			GameLogRecord{
				Move: swde.NewMovePrepare(
					swde.Nickname(host.Nickname),
					swde.Nickname(guest.Nickname),
				),
			},
		},
		StartedAt: now,
	}
}

type Game struct {
	Id            GameId
	HostNickname  Nickname
	HostRating    Rating
	HostPoints    int
	GuestNickname Nickname
	GuestRating   Rating
	GuestPoints   int
	Winner        *Nickname
	Victory       *swde.Victory
	Log           GameLog
	StartedAt     time.Time
	FinishedAt    *time.Time
}

func (dst *Game) State() *swde.State {
	s := new(swde.State)

	for _, item := range dst.Log {
		_ = item.Move.Mutate(s)
	}

	return s
}

type GameOptions struct {
	Tx Tx
	Id GameId
}

type GameOption func(o *GameOptions)

func WithGameId(v GameId) GameOption {
	return func(o *GameOptions) {
		o.Id = v
	}
}

type GameService struct {
	clock         Clock
	gameRepo      GameRepo
	gameClockRepo GameClockRepo
	pusher        Pusher
}

func (dst GameService) Create(
	ctx context.Context,
	host *User,
	guest *User,
	o RoomOptions,
) error {
	now := dst.clock.Now()

	g := NewGame(host, guest, now)

	if err := dst.gameRepo.Save(ctx, g); err != nil {
		return err
	}

	gc := &GameClock{
		Id:         g.Id,
		LastMoveAt: now,
		Turn:       g.State().Me.Name,
		Values: map[Nickname]TimeBank{
			host.Nickname:  o.Clock(),
			guest.Nickname: o.Clock(),
		},
	}

	if err := dst.gameClockRepo.Save(ctx, gc); err != nil {
		return err
	}

	// @TODO push

	return nil
}

type GameClock struct {
	Id         GameId                `json:"id"`
	LastMoveAt time.Time             `json:"lastMoveAt"`
	Turn       Nickname              `json:"turn"`
	Values     map[Nickname]TimeBank `json:"values"`
}

type GameResult struct {
	Winner  Nickname     `json:"winner"`
	Loser   Nickname     `json:"loser"`
	Victory swde.Victory `json:"victory"`
	Points  int          `json:"points"`
}

type GameLog []GameLogRecord

func (dst *GameLog) UnmarshalJSON(bytes []byte) error {
	var messages []*json.RawMessage

	if err := json.Unmarshal(bytes, &messages); err != nil {
		panic("moves unmarshal fail")
	}

	var record struct {
		Move map[string]interface{} `json:"move"`
		Meta MoveMeta               `json:"meta"`
	}

	out := make(GameLog, len(messages))

	for index, message := range messages {
		if err := json.Unmarshal(*message, &record); err != nil {
			log.Fatalln(err)
		}

		rawMove, err := json.Marshal(record.Move)

		if err != nil {
			log.Fatalln(err)
		}

		switch swde.MoveId(record.Move["id"].(float64)) {
		case swde.MovePrepare:
			var m1 swde.PrepareMove

			if err := json.Unmarshal(rawMove, &m1); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = GameLogRecord{
				Move: m1,
				Meta: record.Meta,
			}
		case swde.MovePickWonder:
			var m2 swde.PickWonderMove

			if err := json.Unmarshal(rawMove, &m2); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = GameLogRecord{
				Move: m2,
				Meta: record.Meta,
			}
		case swde.MovePickBoardToken:
			var m3 swde.PickBoardTokenMove

			if err := json.Unmarshal(rawMove, &m3); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = GameLogRecord{
				Move: m3,
				Meta: record.Meta,
			}
		case swde.MoveConstructCard:
			var m4 swde.ConstructCardMove

			if err := json.Unmarshal(rawMove, &m4); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = GameLogRecord{
				Move: m4,
				Meta: record.Meta,
			}
		case swde.MoveConstructWonder:
			var m5 swde.ConstructWonderMove

			if err := json.Unmarshal(rawMove, &m5); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = GameLogRecord{
				Move: m5,
				Meta: record.Meta,
			}
		case swde.MoveDiscardCard:
			var m6 swde.DiscardCardMove

			if err := json.Unmarshal(rawMove, &m6); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = GameLogRecord{
				Move: m6,
				Meta: record.Meta,
			}
		case swde.MoveSelectWhoBeginsTheNextAge:
			var m7 swde.SelectWhoBeginsTheNextAgeMove

			if err := json.Unmarshal(rawMove, &m7); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = GameLogRecord{
				Move: m7,
				Meta: record.Meta,
			}
		case swde.MoveBurnCard:
			var m8 swde.BurnCardMove

			if err := json.Unmarshal(rawMove, &m8); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = GameLogRecord{
				Move: m8,
				Meta: record.Meta,
			}
		case swde.MovePickRandomToken:
			var m9 swde.PickRandomTokenMove

			if err := json.Unmarshal(rawMove, &m9); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = GameLogRecord{
				Move: m9,
				Meta: record.Meta,
			}
		case swde.MovePickTopLineCard:
			var m10 swde.PickTopLineCardMove

			if err := json.Unmarshal(rawMove, &m10); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = GameLogRecord{
				Move: m10,
				Meta: record.Meta,
			}
		case swde.MovePickDiscardedCard:
			var m11 swde.PickDiscardedCardMove

			if err := json.Unmarshal(rawMove, &m11); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = GameLogRecord{
				Move: m11,
				Meta: record.Meta,
			}
		case swde.MovePickReturnedCards:
			var m12 swde.PickReturnedCardsMove

			if err := json.Unmarshal(rawMove, &m12); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = GameLogRecord{
				Move: m12,
				Meta: record.Meta,
			}
		case swde.MoveOver:
			var m13 swde.OverMove

			if err := json.Unmarshal(rawMove, &m13); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = GameLogRecord{
				Move: m13,
				Meta: record.Meta,
			}
		default:
			panic("unknown move")
		}
	}

	*dst = out

	return nil
}

type GameLogRecord struct {
	Move swde.Mutator `json:"move"`
	Meta MoveMeta     `json:"meta"`
}

type MoveMeta struct {
	Actor Nickname `json:"actor"`
}

type TimeBank time.Duration

func (dst TimeBank) MarshalJSON() ([]byte, error) {
	return json.Marshal(math.Round(time.Duration(dst).Seconds()))
}

func (dst *TimeBank) UnmarshalJSON(bytes []byte) error {
	seconds, err := strconv.Atoi(string(bytes))

	if err != nil {
		return err
	}

	*dst = TimeBank(time.Duration(seconds) * time.Second)

	return nil
}
