package domain

import "context"

type OnlineService struct {
	client Onliner
}

func (dst OnlineService) GetAll(ctx context.Context) {
	//users, err := dst.client.Online(ctx)
}
