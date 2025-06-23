package main

import (
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DSN                 string        `envconfig:"DSN" default:"postgres://user_articles_feed:pass_articles_feed@localhost:5432/articles_feed?sslmode=disable"`
	DBMaxConns          int32         `envconfig:"DB_MAX_CONNS" default:"10"`
	DBMinConns          int32         `envconfig:"DB_MIN_CONNS" default:"2"`
	DBMaxConnLifetime   time.Duration `envconfig:"DB_MAX_CONN_LIFETIME" default:"1h"`
	DBMaxConnIdleTime   time.Duration `envconfig:"DB_MAX_CONN_IDLE_TIME" default:"30m"`
	DBHealthcheckPeriod time.Duration `envconfig:"DB_HEALTHCHECK_PERIOD" default:"1m"`

	HTTPReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"30s"`
	HTTPWriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"30s"`
	HTTPIdleTimeout  time.Duration `envconfig:"HTTP_IDLE_TIMEOUT" default:"120s"`
}

func getConfig() Config {
	cfg := Config{}
	once := sync.Once{}

	once.Do(func() {
		envconfig.MustProcess("", &cfg)
	})

	return cfg
}
