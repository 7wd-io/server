package domain

type Top []TopMember

type TopMember struct {
	Name   Nickname `json:"name"`
	Rating Rating   `json:"rating"`
}

type UserScore struct {
	Rank   int    `json:"rank"`
	Rating Rating `json:"rating"`
}

type GamesReport struct {
	Won  GamesStat `json:"won"`
	Lose GamesStat `json:"lose"`
}

type GamesStat struct {
	Total    int `json:"total"`
	Points   int `json:"points"`
	Military int `json:"military"`
	Science  int `json:"science"`
	Resign   int `json:"resign"`
	Timeout  int `json:"timeout"`
}
