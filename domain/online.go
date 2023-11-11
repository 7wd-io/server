package domain

import "context"

func NewOnlineService(
	client Onliner,
	analyst Analyst,
) OnlineService {
	return OnlineService{
		client:  client,
		analyst: analyst,
	}
}

type OnlineService struct {
	client  Onliner
	analyst Analyst
}

func (dst OnlineService) GetAll(ctx context.Context) (UsersPreview, error) {
	online, err := dst.client.Online(ctx)

	if err != nil {
		return nil, err
	}

	if len(online) == 0 {
		return nil, nil
	}

	return dst.analyst.Ratings(ctx, online...)
}
