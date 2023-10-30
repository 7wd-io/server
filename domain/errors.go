package domain

import "7wd.io/rr"

var (
	ErrUserNotFound        = rr.New("user not found")
	ErrRoomNotFound        = rr.New("room not found")
	ErrSessionNotFound     = rr.New("session not found")
	ErrAlreadyJoined       = rr.New("already joined")
	ErrRoomIsFull          = rr.New("room is full")
	errCredentialsNotFound = rr.New("credentials not found")
)
