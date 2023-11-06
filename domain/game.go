package domain

import (
	"encoding/json"
	"fmt"
	swde "github.com/7wd-io/engine"
	"log"
	"time"
)

const (
	timeBankDefault = 10 * time.Minute
	timeBankFast    = 3 * time.Minute
	timeBankBot     = 30 * time.Minute
	timeBankIncr    = 5 * time.Second
)

type GameId int

type Game struct {
	Id            GameId
	HostNickname  Nickname
	HostRating    int
	HostPoints    int
	GuestNickname Nickname
	GuestRating   int
	GuestPoints   int
	Winner        *Nickname
	Victory       *swde.Victory
	Log           GameLog
	StartedAt     time.Time
	FinishedAt    *time.Time
}

type GameClock struct {
	Id         GameId           `json:"id"`
	LastMoveAt time.Time        `json:"lastMoveAt"`
	Turn       Nickname         `json:"turn"`
	Values     map[Nickname]int `json:"values"`
}

type GameResult struct {
	Winner  Nickname     `json:"winner"`
	Loser   Nickname     `json:"loser"`
	Victory swde.Victory `json:"victory"`
	Points  int          `json:"points"`
}

type GameLog []GameLogRecord

type GameLogRecord struct {
	Move swde.Mutator `json:"move"`
	Meta MoveMeta     `json:"meta"`
}

type MoveMeta struct {
	Actor Nickname `json:"actor"`
}

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
			var m1 prepareMove

			if err := json.Unmarshal(rawMove, &m1); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = LogRecord{
				Move: m1,
				Meta: record.Meta,
			}
		case mPickWonder:
			var m2 pickWonderMove

			if err := json.Unmarshal(rawMove, &m2); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = LogRecord{
				Move: m2,
				Meta: record.Meta,
			}
		case mPickBoardToken:
			var m3 pickBoardTokenMove

			if err := json.Unmarshal(rawMove, &m3); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = LogRecord{
				Move: m3,
				Meta: record.Meta,
			}
		case mConstructCard:
			var m4 constructCardMove

			if err := json.Unmarshal(rawMove, &m4); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = LogRecord{
				Move: m4,
				Meta: record.Meta,
			}
		case mConstructWonder:
			var m5 constructWonderMove

			if err := json.Unmarshal(rawMove, &m5); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = LogRecord{
				Move: m5,
				Meta: record.Meta,
			}
		case mDiscardCard:
			var m6 discardCardMove

			if err := json.Unmarshal(rawMove, &m6); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = LogRecord{
				Move: m6,
				Meta: record.Meta,
			}
		case mSelectWhoBeginsTheNextAge:
			var m7 selectWhoBeginsTheNextAgeMove

			if err := json.Unmarshal(rawMove, &m7); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = LogRecord{
				Move: m7,
				Meta: record.Meta,
			}
		case mBurnCard:
			var m8 burnCardMove

			if err := json.Unmarshal(rawMove, &m8); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = LogRecord{
				Move: m8,
				Meta: record.Meta,
			}
		case mPickRandomToken:
			var m9 pickRandomTokenMove

			if err := json.Unmarshal(rawMove, &m9); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = LogRecord{
				Move: m9,
				Meta: record.Meta,
			}
		case mPickTopLineCard:
			var m10 pickTopLineCardMove

			if err := json.Unmarshal(rawMove, &m10); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = LogRecord{
				Move: m10,
				Meta: record.Meta,
			}
		case mPickDiscardedCard:
			var m11 pickDiscardedCardMove

			if err := json.Unmarshal(rawMove, &m11); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = LogRecord{
				Move: m11,
				Meta: record.Meta,
			}
		case mPickReturnedCards:
			var m12 pickReturnedCardsMove

			if err := json.Unmarshal(rawMove, &m12); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = LogRecord{
				Move: m12,
				Meta: record.Meta,
			}
		case mOver:
			var m13 overMove

			if err := json.Unmarshal(rawMove, &m13); err != nil {
				panic(fmt.Errorf("moves unmarshal fail: %w", err))
			}

			out[index] = LogRecord{
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
