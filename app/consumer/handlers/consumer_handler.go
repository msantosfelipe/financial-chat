package handlers

import (
	"log"

	"github.com/msantosfelipe/financial-chat/app"
	"github.com/msantosfelipe/financial-chat/app/consumer"
)

var stockService consumer.ConsumerService

func New() {
	stockService = consumer.GetStockInstance()
}

func HandleMessageConsumer() {
	stockService.SubscribeToQueue(app.ENV.AmqpChatQueueName)
}

func Clean() {
	log.Println("cleanning chatbot tasks...")
	if stockService != nil {
		stockService.Clean()
	}
	log.Println("chatbot stopped")
}
