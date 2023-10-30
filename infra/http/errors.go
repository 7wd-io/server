package http

import "7wd.io/rr"

var (
	errInvalidToken   = rr.New("invalid token")
	errUnauthorized   = rr.New("unauthorized")
	errInvalidRequest = rr.New("invalid request")
)
