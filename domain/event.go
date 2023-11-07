package domain

const (
	EventGameCreated EventId = iota + 1
)

type EventId int

type GameCreatedPayload struct {
	Game *Game
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
