package stock

import (
	"sync"

	"github.com/msantosfelipe/financial-chat/app/websocket"
	"github.com/msantosfelipe/financial-chat/infra/amqp"
)

var stockInstance StockService
var stockOnce sync.Once

type QueueStockMessage struct {
	Stock string `json:"stock"`
	Room  string `json:"room"`
}

func GetStockInstance() StockService {
	stockOnce.Do(func() {
		stockInstance = NewStock()
	})

	return stockInstance
}

type stockService struct {
	amqpService      amqp.AmqpService
	websocketService websocket.WebsocketService
}

type StockService interface {
	queueSubscriber
	serviceCleaner
}

type queueSubscriber interface {
	SubscribeToQueue(queue string)
}

type serviceCleaner interface {
	Clean()
}
