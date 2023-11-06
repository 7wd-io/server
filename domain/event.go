package domain

type RoomCreated struct {
	Room *Room
}

type RoomDeleted struct {
	Host Nickname
}

type RoomUpdated struct {
	Room *Room
}
