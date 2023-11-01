package domain

import "time"

const (
	timeBankDefault = 10 * time.Minute
	timeBankFast    = 3 * time.Minute
	timeBankBot     = 30 * time.Minute
	timeBankIncr    = 5 * time.Second
)

type GameId int

type Game struct {
	Id GameId
}
