package domain

type Top []TopMember

type TopMember struct {
	Name   Nickname `json:"name"`
	Rating Rating   `json:"rating"`
}
