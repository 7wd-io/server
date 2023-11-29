package domain

import swde "github.com/7wd-io/engine"

const (
	EventGameCreated EventId = iota + 1
	EventGameUpdated
	EventGameOver
	EventAfterGameMove
	EventBotIsReadyToMove
)

const (
	EventRoomCreated = iota + 100
	EventRoomUpdated
	EventRoomDeleted
	EventRoomStarted
)

const (
	EventPlayAgainUpdated = iota + 1000
	EventPlayAgainApproved
)

type EventId int

type GameCreatedPayload struct {
	Game *Game
}

// @TODO Id -> Game?
type GameUpdatedPayload struct {
	Id       GameId        `json:"id"`
	State    *swde.State   `json:"state"`
	Clock    *GameClock    `json:"clock"`
	LastMove GameLogRecord `json:"lastMove"`
}

type GameOverPayload struct {
	Game    Game
	Result  GameResult
	Options RoomOptions
}

type AfterGameMovePayload struct {
	Game *Game
}

type BotIsReadyToMovePayload struct {
	Game GameId
	Move swde.Mutator
}

type RoomCreatedPayload struct {
	Room Room
}

type RoomDeletedPayload struct {
	Room Room
}

type RoomUpdatedPayload struct {
	Room Room
}

type RoomStartedPayload struct {
	Host  User
	Guest User
	Room  Room
}

type OnlineUpdatedPayload struct {
	Users []Nickname
}

type PlayAgainUpdatedPayload struct {
	Game   GameId   `json:"game"`
	User   Nickname `json:"user"`
	Answer bool     `json:"answer"`
}

type PlayAgainApprovedPayload struct {
	Id   GameId `json:"id"`
	Next GameId `json:"next"`
}

type RoomCreated struct {
	Room *Room
}

type RoomDeleted struct {
	Host Nickname
}

type RoomUpdated struct {
	Room *Room
}
