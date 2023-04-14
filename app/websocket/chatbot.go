package websocket

import (
	"sync"

	"github.com/msantosfelipe/financial-chat/infra/amqp"
)

var chatbotInstance ChatbotService
var chatbotOnce sync.Once

func GetChatbotInstance() ChatbotService {
	chatbotOnce.Do(func() {
		chatbotInstance = NewChatbot()
	})

	return chatbotInstance
}

type chatbotService struct {
	amqpService amqp.AmqpService
}

type ChatbotService interface {
	messageHandler
	// queueSubscriber
	// serviceCleaner
}

type messageHandler interface {
	HandleBotMessage(text, room string)
	StockHandler(stock, room string) error
}

type queueSubscriber interface {
	SubscribeToQueue(queue string)
}

type serviceCleaner interface {
	Clean()
}
