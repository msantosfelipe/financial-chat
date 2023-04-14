package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/msantosfelipe/financial-chat/app"
	"github.com/msantosfelipe/financial-chat/infra/amqp"
)

var ws WebsocketService

const (
	prefixStock = "/stock="
	prefixHelp  = "/help"
)

func init() {
	ws = GetWSInstance()
}

func NewChatbot() ChatbotService {
	return &chatbotService{
		amqpService: amqp.GetInstance(),
	}
}

func (s *chatbotService) HandleBotMessage(text, room string) {
	switch {
	case strings.HasPrefix(text, prefixStock):
		err := s.StockHandler(text, room)
		if err != nil {
			msg := fmt.Sprintf("*** error: %v", err)
			ws.SendBotMessage(room, msg)
		}
	case strings.HasPrefix(text, prefixHelp):
		msg := "*** usage: \"/stock='stock_code'\""
		ws.SendBotMessage(room, msg)
	default:
		msg := fmt.Sprintf("*** invalid command %s", text)
		ws.SendBotMessage(room, msg)
	}
}

func (s *chatbotService) StockHandler(text, room string) error {
	if len(strings.Split(text, prefixStock)) == 1 {
		return fmt.Errorf("invalid stock name %s", text)
	}

	stock := strings.ToUpper(strings.Split(text, prefixStock)[1])

	bytes, err := json.Marshal(QueueStockMessage{
		Stock: stock,
		Room:  room,
	})
	if err != nil {
		log.Println("error marshaling message: ", err)
	}

	if err := ws.PublishMessageToQueue(bytes, app.ENV.AmqpChatQueueName); err != nil {
		return err
	}

	return nil
}
