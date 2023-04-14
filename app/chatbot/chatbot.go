package chatbot

import (
	"sync"

	"github.com/msantosfelipe/financial-chat/infra/amqp"
)

var instance ChatbotService
var once sync.Once

type QueueMessage struct {
	Text string `json:"text"`
	Room string `json:"room"`
}

func GetInstance() ChatbotService {
	once.Do(func() {
		instance = New()
	})

	return instance
}

type chatbotService struct {
	amqpService amqp.AmqpService
}

type ChatbotService interface {
	queueSubscriber
	messageHandler
	serviceCleaner
}

type queueSubscriber interface {
	SubscribeToQueue(queue string)
}

type messageHandler interface {
	StockHandler(stockName string)
}

type serviceCleaner interface {
	Clean()
}
