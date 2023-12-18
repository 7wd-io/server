package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

var C c

func init() {
	envconfig.MustProcess("", &C)
}

type c struct {
	Env    string `required:"true" envconfig:"SWD_ENV"`
	Port   int    `required:"true" envconfig:"SWD_PORT"`
	Secret string `required:"true" envconfig:"SWD_SECRET"`
	Domain string `required:"true" envconfig:"SWD_DOMAIN"`
	Path   string `required:"true" envconfig:"SWD_PATH"`
	Bot    struct {
		Endpoint string `required:"true" envconfig:"SWD_BOT_ENDPOINT"`
	}
	ClientOrigin string `required:"true" envconfig:"SWD_CLIENT_ORIGIN"`
	Pg           struct {
		Host     string `envconfig:"SWD_PG_HOST"`
		Port     int    `envconfig:"SWD_PG_PORT"`
		DbName   string `envconfig:"SWD_PG_DBNAME"`
		User     string `envconfig:"SWD_PG_USER"`
		Password string `envconfig:"SWD_PG_PASSWORD"`
	}
	Redis struct {
		Port int `envconfig:"SWD_REDIS_PORT"`
	}
	Centrifugo struct {
		Endpoint string `required:"true" envconfig:"SWD_CENTRIFUGO_ENDPOINT"`
		ApiKey   string `required:"true" envconfig:"SWD_CENTRIFUGO_API_KEY"`
	}
	Mailer struct {
		Server   string `required:"true" envconfig:"SWD_MAILER_SERVER"`
		Port     int    `required:"true" envconfig:"SWD_MAILER_PORT"`
		User     string `required:"true" envconfig:"SWD_MAILER_LOGIN"`
		Password string `required:"true" envconfig:"SWD_MAILER_PASSWORD"`
	}
}

func (dst c) IsTest() bool {
	return dst.Env == "test"
}

func (dst c) PgDsn() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		dst.Pg.User,
		dst.Pg.Password,
		dst.Pg.Host,
		dst.Pg.Port,
		dst.Pg.DbName,
	)
}

func (dst c) RedisUrl() string {
	return fmt.Sprintf("redis://localhost:%d", dst.Redis.Port)
}
