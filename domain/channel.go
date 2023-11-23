package domain

import "fmt"

const ChOnline = "service:online"

const ChRoomUpdate = "room:update"
const ChRoomCreate = "room:create"
const ChRoomDelete = "room:delete"

var ChGameUpdate = func(id GameId) string {
	return fmt.Sprintf("game:update_%d", id)
}

var ChPlayAgainUpdate = func(id GameId) string {
	return fmt.Sprintf("play-again:update_%d", id)
}

var ChPlayAgainApprove = func(id GameId) string {
	return fmt.Sprintf("play-again:approve_%d", id)
}
