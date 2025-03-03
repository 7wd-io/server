package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	swde "github.com/7wd-io/engine"
	"log"
	"log/slog"
	"math"
	"strconv"
	"time"
)

const (
	TimeBankDefault = TimeBank(20 * time.Minute)
	TimeBankFast    = TimeBank(7 * time.Minute)
	TimeBankBot     = TimeBank(30 * time.Minute)
	TimeBankIncr    = TimeBank(5 * time.Second)
)

type GameId int

func newGame(host User, guest User, now time.Time, o RoomOptions) *Game {
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
					swde.Options{
						PromoWonders: o.PromoWonders,
					},
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

func (dst *Game) Move(u Nickname, move swde.Mutator) (*swde.State, error) {
	if dst.IsOver() {
		return nil, ErrGameIsOver
	}

	s := dst.State()

	switch move := move.(type) {
	// overMove no have priority
	case swde.OverMove:
		if move.Loser != s.Me.Name && move.Loser != s.Enemy.Name {
			return nil, ErrActionNotAllowed
		}
	default:
		if u != Nickname(s.Me.Name) {
			return nil, ErrActionNotAllowed
		}
	}

	if err := move.Mutate(s); err != nil {
		return nil, err
	}

	dst.Log = append(dst.Log, GameLogRecord{
		Move: move,
		Meta: MoveMeta{
			Actor: u,
		},
	})

	return s, nil
}

func (dst *Game) Over(s *swde.State, t time.Time) GameResult {
	dst.Winner = (*Nickname)(s.Winner)
	dst.Victory = s.Victory
	dst.FinishedAt = &t

	r := GameResult{
		Winner:  Nickname(*s.Winner),
		Victory: *s.Victory,
	}

	if *dst.Winner == Nickname(s.Me.Name) {
		r.Loser = Nickname(s.Enemy.Name)
	} else {
		r.Loser = Nickname(s.Me.Name)
	}

	if r.Winner == dst.HostNickname {
		r.Points = dst.HostPoints
	} else {
		r.Points = dst.GuestPoints
	}

	return r
}

func (dst *Game) IsOver() bool {
	return dst.Winner != nil
}

type GameOptions struct {
	Tx Tx

	Lock bool

	Id    GameId
	IdSet bool
}

type GameOption func(o *GameOptions)

func WithGameTx(v Tx) GameOption {
	return func(o *GameOptions) {
		o.Tx = v
	}
}

func WithGameLock() GameOption {
	return func(o *GameOptions) {
		o.Lock = true
	}
}

func WithGameId(v GameId) GameOption {
	return func(o *GameOptions) {
		o.Id = v
		o.IdSet = true
	}
}

func NewGameService(
	clock Clock,
	roomRepo RoomRepo,
	gameRepo GameRepo,
	gameClockRepo GameClockRepo,
	userRepo UserRepo,
	dispatcher Dispatcher,
	tx Txer,
) GameService {
	return GameService{
		clock:         clock,
		roomRepo:      roomRepo,
		gameRepo:      gameRepo,
		gameClockRepo: gameClockRepo,
		userRepo:      userRepo,
		dispatcher:    dispatcher,
		tx:            tx,
	}
}

type GameService struct {
	clock         Clock
	roomRepo      RoomRepo
	gameRepo      GameRepo
	gameClockRepo GameClockRepo
	userRepo      UserRepo
	pusher        Pusher
	dispatcher    Dispatcher
	tx            Txer
}

func (dst GameService) Get(ctx context.Context, id GameId) (*Game, error) {
	return dst.gameRepo.Find(ctx, WithGameId(id))
}

func (dst GameService) Clock(ctx context.Context, id GameId) (*GameClock, error) {
	return dst.gameClockRepo.Find(ctx, id)
}

func (dst GameService) State(ctx context.Context, id GameId, index int) (*swde.State, error) {
	game, err := dst.gameRepo.Find(ctx, WithGameId(id))

	if err != nil {
		return nil, err
	}

	s := new(swde.State)

	for _, item := range game.Log[:index+1] {
		_ = item.Move.Mutate(s)
	}

	return s, nil
}

func (dst GameService) Create(
	ctx context.Context,
	host User,
	guest User,
	o RoomOptions,
) (*Game, error) {
	now := dst.clock.Now()

	game := newGame(host, guest, now, o)

	if err := dst.gameRepo.Save(ctx, game); err != nil {
		return nil, err
	}

	gc := &GameClock{
		Id:         game.Id,
		LastMoveAt: now,
		Turn:       Nickname(game.State().Me.Name),
		Values: map[Nickname]TimeBank{
			host.Nickname:  o.Clock(),
			guest.Nickname: o.Clock(),
		},
	}

	if err := dst.gameClockRepo.Save(ctx, gc); err != nil {
		return nil, err
	}

	dst.dispatcher.Dispatch(
		ctx,
		EventGameCreated,
		GameCreatedPayload{Game: game},
	)

	return game, nil
}

func (dst GameService) CreateWithBot(ctx context.Context, pass Passport) error {
	if err := dst.alreadyJoin(ctx, pass); err != nil {
		return err
	}

	user, err := dst.userRepo.Find(ctx, WithUserNickname(pass.Nickname))

	if err != nil {
		return err
	}

	bot, err := dst.userRepo.Find(ctx, WithUserNickname(BotNickname))

	if err != nil {
		return err
	}

	options := RoomOptions{
		PromoWonders: false,
		TimeBank:     TimeBankBot,
	}

	game, err := dst.Create(ctx, *user, *bot, options)

	if err != nil {
		return err
	}

	room := &Room{
		Host:        game.HostNickname,
		HostRating:  game.HostRating,
		Guest:       game.GuestNickname,
		GuestRating: game.GuestRating,
		Options:     options,
		GameId:      game.Id,
	}

	if err = dst.roomRepo.Save(ctx, room); err != nil {
		return err
	}

	dst.dispatcher.Dispatch(
		ctx,
		EventRoomCreated,
		RoomCreatedPayload{Room: *room},
	)

	return nil
}

func (dst GameService) Move(ctx context.Context, u Nickname, id GameId, m swde.Mutator) (*Game, error) {
	var err error

	tx, err := dst.tx.Tx(ctx)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			slog.Error("GameService.Move tx.Rollback", slog.String("err", err.Error()))
		}
	}()

	g, err := dst.gameRepo.Find(
		ctx,
		WithGameId(id),
		WithGameLock(),
		WithGameTx(tx),
	)

	if err != nil {
		return nil, err
	}

	if g == nil {
		return nil, ErrGameNotFound
	}

	if g.IsOver() {
		return nil, ErrGameIsOver
	}

	gameClock, err := dst.gameClockRepo.Find(ctx, id)

	if err != nil {
		return nil, err
	}

	now := dst.clock.Now()

	timePassed := TimeBank(now.Sub(gameClock.LastMoveAt))
	gameClock.Values[u] -= timePassed

	var s *swde.State

	if gameClock.Values[u] > 0 {
		gameClock.Values[u] += TimeBankIncr
		s, err = g.Move(u, m)
		gameClock.LastMoveAt = now
	} else {
		m = swde.NewMoveOver(swde.Nickname(u), swde.Timeout)
		gameClock.Values[u] = 0
		s, err = g.Move(u, m)
	}

	if err != nil {
		return nil, err
	}

	gameClock.Turn = Nickname(s.Me.Name)

	if s.IsOver() {
		result := g.Over(s, now)

		room, err := dst.roomRepo.FindByGame(ctx, g.Id)

		if err != nil {
			return nil, err
		}

		dst.dispatcher.Dispatch(
			ctx,
			EventGameOver,
			GameOverPayload{
				Game:    *g,
				Result:  result,
				Options: room.Options,
			},
		)
	}

	if err = dst.gameRepo.Update(ctx, g, WithGameTx(tx)); err != nil {
		return nil, err
	}

	if err = dst.gameClockRepo.Save(ctx, gameClock); err != nil {
		return nil, err
	}

	dst.dispatcher.Dispatch(ctx, EventGameUpdated, GameUpdatedPayload{
		Id:       g.Id,
		State:    s,
		Clock:    gameClock,
		LastMove: g.Log[len(g.Log)-1],
	})

	dst.dispatcher.Dispatch(ctx, EventAfterGameMove, AfterGameMovePayload{
		Game: g,
	})

	return g, tx.Commit(ctx)
}

func (dst GameService) OnEventBotIsReadyToMove(ctx context.Context, payload interface{}) error {
	p, ok := payload.(BotIsReadyToMovePayload)

	if !ok {
		return errors.New("func (dst GameService) OnEventBotIsReadyToMove !ok := payload.(BotIsReadyToMovePayload)")
	}

	_, err := dst.Move(ctx, BotNickname, p.Game, p.Move)

	return err
}

func (dst GameService) OnRoomStarted(ctx context.Context, payload interface{}) error {
	p, ok := payload.(RoomStartedPayload)

	if !ok {
		return errors.New("func (dst GameService) OnRoomStarted !ok := payload.(RoomStartedPayload)")
	}

	game, err := dst.Create(ctx, p.Host, p.Guest, p.Room.Options)

	if err != nil {
		return err
	}

	p.Room.GameId = game.Id

	if err = dst.roomRepo.Save(ctx, &p.Room); err != nil {
		return err
	}

	dst.dispatcher.Dispatch(
		ctx,
		EventRoomUpdated,
		RoomUpdatedPayload{Room: p.Room},
	)

	return nil
}

func (dst GameService) alreadyJoin(ctx context.Context, pass Passport) error {
	rooms, err := dst.roomRepo.FindAll(ctx)

	if err != nil {
		return err
	}

	for _, v := range rooms {
		if pass.Nickname == v.Host || pass.Nickname == v.Guest {
			return ErrOneRoomPerPlayer
		}
	}

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

// UnmarshalMove @TODO полность переделать убрать пересечение с маршалингом списка ходов
func UnmarshalMove(move []byte) (swde.Mutator, error) {
	var err error

	var m map[string]interface{}

	if err = json.Unmarshal(move, &m); err != nil {
		slog.Error("unmarshal move fail", slog.String("raw", string(move)))
		return nil, err
	}

	switch swde.MoveId(m["id"].(float64)) {
	case swde.MovePrepare:
		var m1 swde.PrepareMove
		err = json.Unmarshal(move, &m1)

		return m1, err
	case swde.MovePickWonder:
		var m2 swde.PickWonderMove
		err = json.Unmarshal(move, &m2)

		return m2, err
	case swde.MovePickBoardToken:
		var m3 swde.PickBoardTokenMove
		err = json.Unmarshal(move, &m3)

		return m3, err
	case swde.MoveConstructCard:
		var m4 swde.ConstructCardMove
		err = json.Unmarshal(move, &m4)

		return m4, err
	case swde.MoveConstructWonder:
		var m5 swde.ConstructWonderMove
		err = json.Unmarshal(move, &m5)

		return m5, err
	case swde.MoveDiscardCard:
		var m6 swde.DiscardCardMove
		err = json.Unmarshal(move, &m6)

		return m6, err
	case swde.MoveSelectWhoBeginsTheNextAge:
		var m7 swde.SelectWhoBeginsTheNextAgeMove
		err = json.Unmarshal(move, &m7)

		return m7, err
	case swde.MoveBurnCard:
		var m8 swde.BurnCardMove
		err = json.Unmarshal(move, &m8)

		return m8, err
	case swde.MovePickRandomToken:
		var m9 swde.PickRandomTokenMove
		err = json.Unmarshal(move, &m9)

		return m9, err
	case swde.MovePickTopLineCard:
		var m10 swde.PickTopLineCardMove
		err = json.Unmarshal(move, &m10)

		return m10, err
	case swde.MovePickDiscardedCard:
		var m11 swde.PickDiscardedCardMove
		err = json.Unmarshal(move, &m11)

		return m11, err
	case swde.MovePickReturnedCards:
		var m12 swde.PickReturnedCardsMove
		err = json.Unmarshal(move, &m12)

		return m12, err
	case swde.MoveOver:
		var m13 swde.OverMove
		err = json.Unmarshal(move, &m13)

		return m13, err
	default:
		return nil, errors.New("unknown move")
	}
}
