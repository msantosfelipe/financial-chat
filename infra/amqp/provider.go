package amqp

import (
	"log"
	"sync"

	amqpgo "github.com/rabbitmq/amqp091-go"

	"github.com/msantosfelipe/financial-chat/app"
)

var instance AmqpService
var once sync.Once

func GetInstance() AmqpService {
	once.Do(func() {
		host := app.ENV.AmqpServerURL
		log.Println("Creating amqp service, host:", host)
		instance = New(
			host,
			app.ENV.AmqpChatQueueName,
		)
	})
	return instance
}

type AmqpService interface {
	ampqSubscriber
	messageHandler
	serviceCleaner
}

type ampqSubscriber interface {
	SubscribeToQueue(queue string) <-chan amqpgo.Delivery
}

type messageHandler interface {
	PublishMessage(msg []byte, queue string) error
}

type serviceCleaner interface {
	Clean()
}
