package consumer

import (
	"sync"

	"github.com/msantosfelipe/financial-chat/app/websocket"
	"github.com/msantosfelipe/financial-chat/infra/amqp"
)

var consumerInstance ConsumerService
var once sync.Once

type QueueMessage struct {
	Stock string `json:"stock"`
	Room  string `json:"room"`
}

func GetStockInstance() ConsumerService {
	once.Do(func() {
		consumerInstance = NewConsumer(
			amqp.GetInstance(),
			websocket.GetWSInstance(),
		)
	})

	return consumerInstance
}

type consumerService struct {
	amqpService      amqp.AmqpService
	websocketService websocket.WebsocketService
}

type ConsumerService interface {
	queueSubscriber
	serviceCleaner
}

type queueSubscriber interface {
	// Subscribe to queue and listen to new messages
	SubscribeToQueue(queue string)
}

type serviceCleaner interface {
	// Close connections and channels
	Clean()
}
