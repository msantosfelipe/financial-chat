package cache

import (
	"log"
	"sync"
	"time"

	"github.com/msantosfelipe/financial-chat/app"
)

var instance CacheService
var once sync.Once

func GetInstance() CacheService {
	once.Do(func() {
		host := app.ENV.RedisURL
		log.Println("Creating redis service, host:", host)
		instance = New(
			host,
			time.Duration(app.ENV.RedisExpirationMinutes)*time.Minute)
	})
	return instance
}

type CacheService interface {
	cacheChecker
	cacheSetter
}

type cacheChecker interface {
	ExistsChatKey(room string) bool
	GetPreviousChatMessages(room string) []string
	HandleChatSize(room string)
}

type cacheSetter interface {
	StoreChatMessage(message []byte, room string)
}
