package cache

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/msantosfelipe/financial-chat/app"
)

var (
	chatKeyPrefix = "chat_messages"
)

type redisService struct {
	client *redis.Client
	expire time.Duration
}

func New(host string, expire time.Duration) CacheService {
	client := redis.NewClient(&redis.Options{
		Addr: host,
	})
	if err := client.Ping().Err(); err != nil {
		log.Println("error on redis connection: ", err)
	}

	return &redisService{
		client: client,
		expire: expire,
	}
}

func (s *redisService) ExistsChatKey(room string) bool {
	return s.client.Exists(getChatKey(room)).Val() != 0
}

func (s *redisService) GetPreviousChatMessages(room string) []string {
	chatMessages, err := s.client.LRange(getChatKey(room), 0, -1).Result()
	if err != nil {
		log.Println(fmt.Sprintf("cache: error retreiving previous chat messages from room %s :", room), err)
	}

	return revertOrder(chatMessages)
}

func (s *redisService) StoreChatMessage(message []byte, room string) {
	redisKey := getChatKey(room)
	if err := s.client.LPush(redisKey, message).Err(); err != nil {
		log.Println(fmt.Sprintf("cache: error caching message of room %s :", room), err)
	}

	if err := s.client.Expire(redisKey, s.expire); err.Err() != nil {
		log.Println(fmt.Sprintf("cache: error setting expiration time for message of room %s - all messages will be deleted :", room), err)
		if err := s.client.Del(redisKey); err != nil {
			log.Println(fmt.Sprintf("cache: error deleting message of room %s :", room), err)
		}
	}
}

func (s *redisService) HandleChatSize(room string) {
	chatMessages := s.GetPreviousChatMessages(room)
	numberOfMessages := len(chatMessages)
	if numberOfMessages >= app.ENV.MaxMessagesPerRoom {
		if err := s.client.LRem(getChatKey(room), 1, chatMessages[0]).Err(); err != nil {
			log.Println(fmt.Sprintf("cache: error deleting max message of room %s :", room), err)
		}
	}
}

func getChatKey(room string) string {
	return fmt.Sprintf("%s:%s", chatKeyPrefix, room)
}

func revertOrder(chatMessages []string) []string {
	last := len(chatMessages) - 1
	for i := 0; i < len(chatMessages)/2; i++ {
		chatMessages[i], chatMessages[last-i] = chatMessages[last-i], chatMessages[i]
	}
	return chatMessages
}
