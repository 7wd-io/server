package domain

const (
	EventGameCreated EventId = iota + 1
	EventRoomCreated
	EventRoomUpdated
	EventRoomDeleted
)

type EventId int

type GameCreatedPayload struct {
	Game *Game
}

type RoomCreatedPayload struct {
	Room *Room
}

type RoomDeletedPayload struct {
	Room *Room
}

type RoomUpdatedPayload struct {
	Room *Room
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
