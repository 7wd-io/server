package domain

import "math"

const (
	DefaultElo Rating = 1500
)

const (
	eloKFactor   = 20
	eloDeviation = 800
)

type Rating int

func Elo(winner, loser Rating) int {
	points := 1 / (1 + math.Pow(10, (float64(loser)-float64(winner))/eloDeviation))

	return int(eloKFactor * (1 - points))
}
