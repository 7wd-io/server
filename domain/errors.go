package domain

import "7wd.io/rr"

// room errors
var (
	ErrRoomNotFound             = rr.New("room not found")
	ErrInvalidRoomOptions       = rr.New("invalid room options")
	ErrOneRoomPerPlayer         = rr.New("one room per player")
	ErrRoomIsFull               = rr.New("room is full")
	ErrAlreadyJoined            = rr.New("already joined")
	ErrOnlyHostCanRemoveRoom    = rr.New("only host can remove room")
	ErrCantRemoveInProgressRoom = rr.New("cant remove in progress room")
	ErrCantLeaveInProgressRoom  = rr.New("cant leave in progress room")
	ErrJoinToTheRoomRestricted  = rr.New("join to the room restricted")
	ErrRoomPlayerNotFound       = rr.New("room player not found")
	ErrOnlyHostKick             = rr.New("only host can kick")
	ErrGameClockNotFound        = rr.New("game clock not found")
	ErrGameIsOver               = rr.New("game is over")
	ErrActionNotAllowed         = rr.New("action not allowed")
)

// account errors
var (
	ErrUserNotFound        = rr.New("user not found")
	ErrGameNotFound        = rr.New("game not found")
	ErrSessionNotFound     = rr.New("session not found")
	errCredentialsNotFound = rr.New("credentials not found")
)
