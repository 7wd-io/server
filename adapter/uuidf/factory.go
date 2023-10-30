package uuidf

import "github.com/google/uuid"

func New() F {
	return F{}
}

type F struct{}

func (dst F) Uuid() uuid.UUID {
	return uuid.New()
}
