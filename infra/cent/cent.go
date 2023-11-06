package cent

import (
	"7wd.io/config"
	"github.com/centrifugal/gocent/v3"
)

func New() *gocent.Client {
	return gocent.New(gocent.Config{
		Addr: config.C.Centrifugo.Endpoint,
		Key:  config.C.Centrifugo.ApiKey,
	})
}
