package app

import (
	"log"

	"github.com/Netflix/go-env"
	"github.com/subosito/gotenv"
)

type Environment struct {
	Port                   string `env:"PORT"`
	RedisURL               string `env:"REDIS_URL"`
	RedisExpirationMinutes int    `env:"REDIS_EXPIRE_MINUTES"`
}

var ENV Environment

func init() {
	gotenv.Load() // load .env file (if exists)
	if _, err := env.UnmarshalFromEnviron(&ENV); err != nil {
		log.Fatal("Fatal error unmarshalling environment config: ", err)
	}
}
