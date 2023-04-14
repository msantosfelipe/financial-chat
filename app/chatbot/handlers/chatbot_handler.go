package handlers

import (
	"log"

	"github.com/msantosfelipe/financial-chat/app"
	"github.com/msantosfelipe/financial-chat/app/chatbot"
)

var chatbotService chatbot.ChatbotService

func init() {
	chatbotService = chatbot.GetInstance()
}

func HandleMessageConsumer() {
	go chatbotService.SubscribeToQueue(app.ENV.AmqpChatQueueName)
}

func Clean() {
	log.Println("cleanning chatbot tasks...")
	chatbotService.Clean()
	log.Println("chatbot stopped")
}
