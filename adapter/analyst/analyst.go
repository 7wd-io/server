package analyst

import (
	"7wd.io/domain"
	"context"
)

func New() A {
	return A{}
}

type A struct {
}

func (dst A) Top(ctx context.Context) (domain.Top, error) {
	//TODO implement me
	panic("implement me")
}

func (dst A) Update(ctx context.Context, result domain.GameResult) error {
	//TODO implement me
	panic("implement me")
}

func (dst A) Ratings(ctx context.Context, nickname ...domain.Nickname) (domain.UsersPreview, error) {
	//TODO implement me
	panic("implement me")
}
