package bot

import (
	"7wd.io/domain"
	"bytes"
	"encoding/json"
	swde "github.com/7wd-io/engine"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func New(endpoint string) B {
	return B{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		endpoint: endpoint,
	}
}

type B struct {
	client   http.Client
	endpoint string
}

func (dst B) GetMove(g *domain.Game) (swde.Mutator, error) {
	req, err := json.Marshal(request{
		Id: g.Id,
		Host: player{
			Name:   g.HostNickname,
			Rating: g.HostRating,
			Points: g.HostPoints,
		},
		Guest: player{
			Name:   g.GuestNickname,
			Rating: g.GuestRating,
			Points: g.GuestPoints,
		},
		State:    g.State(),
		Finished: g.Winner != nil,
		Log:      g.Log[1:],
	})

	if err != nil {
		return nil, err
	}

	resp, err := dst.client.Post(
		dst.endpoint,
		"application/json",
		bytes.NewBuffer(req),
	)

	if err != nil {
		slog.Error("call bot server failed")
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		slog.Error("read answer from bot server failed")
		return nil, err
	}

	return domain.UnmarshalMove(body)
}

type player struct {
	Name   domain.Nickname `json:"name"`
	Rating domain.Rating   `json:"rating"`
	Points int             `json:"points"`
}

type request struct {
	Id       domain.GameId     `json:"id"`
	Host     player            `json:"host"`
	Guest    player            `json:"guest"`
	Clock    *domain.GameClock `json:"clock,omitempty"`
	State    *swde.State       `json:"state"`
	Finished bool              `json:"finished"`
	Log      domain.GameLog    `json:"log"`
}
