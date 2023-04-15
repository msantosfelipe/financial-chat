package handlers

import (
	"log"

	"github.com/msantosfelipe/financial-chat/app"
	"github.com/msantosfelipe/financial-chat/app/consumer"
)

var consumerService consumer.ConsumerService

func New() {
	consumerService = consumer.GetStockInstance()
}

func HandleMessageConsumer() {
	consumerService.SubscribeToQueue(app.ENV.AmqpChatQueueName)
}

func Clean() {
	log.Println("cleanning chatbot tasks...")
	if consumerService != nil {
		consumerService.Clean()
	}
	log.Println("chatbot stopped")
}
