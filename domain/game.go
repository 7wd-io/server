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
