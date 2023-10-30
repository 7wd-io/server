package rds

import (
	"7wd.io/config"
	"context"
	"github.com/redis/go-redis/v9"
)

func MustNew() *redis.Client {
	var err error

	opt, err := redis.ParseURL(config.C.RedisUrl())

	if err != nil {
		panic("unable to parse redis url")
	}

	rds := redis.NewClient(opt)

	if err = rds.Ping(context.Background()).Err(); err != nil {
		panic("unable to connect redis")
	}

	return rds
}
