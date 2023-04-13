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
	MaxRooms               int    `env:"MAX_ROOMS"`
	MaxClientsPerRoom      int    `env:"MAX_CLIENTS_PER_ROOM"`
	MaxMessagesPerRoom     int    `env:"MAX_MESSAGES_PER_ROOM"`
}

var ENV Environment

func init() {
	gotenv.Load() // load .env file (if exists)
	if _, err := env.UnmarshalFromEnviron(&ENV); err != nil {
		log.Fatal("Fatal error unmarshalling environment config: ", err)
	}
}
