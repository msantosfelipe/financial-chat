package chatbot

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/msantosfelipe/financial-chat/app/client"
	"github.com/msantosfelipe/financial-chat/infra/amqp"
)

const (
	chatbotUser = "bot"
	prefixStock = "/stock"
	prefixHelp  = "/help"
)

func New() ChatbotService {
	return &chatbotService{
		amqpService: amqp.GetInstance(),
	}
}

func (s *chatbotService) SubscribeToQueue(queue string) {
	// todo maybe criar um channel

	messages := s.amqpService.SubscribeToQueue(queue)

	for message := range messages {
		var queueMessage QueueMessage
		log.Printf(" > Received message: %s\n", message.Body)
		if err := json.Unmarshal(message.Body, &queueMessage); err != nil {
			log.Println("error reading message: ", err)
		}

		switch {
		case strings.HasPrefix(queueMessage.Text, prefixStock):
			stock := strings.Split(queueMessage.Text, prefixStock)[1]
			s.StockHandler(stock)
		case strings.HasPrefix(queueMessage.Text, prefixHelp):
			text := "*** usage: \"/stock='stock_code'\""
			client.SendBotMessage(chatbotUser, queueMessage.Room, text)
		default:
			text := fmt.Sprintf("*** invalid bot command %s", queueMessage.Text)
			client.SendBotMessage(chatbotUser, queueMessage.Room, text)
		}
	}
}

func (s *chatbotService) StockHandler(stock string) {
	GetStockInstance().RequestStock(stock)
}

func (s *chatbotService) Clean() {
	s.amqpService.Clean()
}
