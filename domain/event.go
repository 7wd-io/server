package domain

const (
	EventGameCreated EventId = iota + 1
	EventGameUpdated
)

const (
	EventRoomCreated = iota + 100
	EventRoomUpdated
	EventRoomDeleted
)

const (
	EventOnlineUpdated = iota + 1000
	EventPlayAgainUpdated
	EventPlayAgainApproved
)

type EventId int

type GameCreatedPayload struct {
	Game *Game
}

type GameUpdatedPayload struct {
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

type OnlineUpdatedPayload struct {
	Users []Nickname
}

type PlayAgainUpdatedPayload struct {
}

type PlayAgainApprovedPayload struct {
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
