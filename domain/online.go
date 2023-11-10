package domain

import "context"

func NewOnlineService(client Onliner) OnlineService {
	return OnlineService{
		client: client,
	}
}

type OnlineService struct {
	client Onliner
}

func (dst OnlineService) GetAll(ctx context.Context) (UsersPreview, error) {
	users, err := dst.client.Online(ctx)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
