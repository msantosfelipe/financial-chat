package consumer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/msantosfelipe/financial-chat/app"
	"github.com/msantosfelipe/financial-chat/app/websocket"
	"github.com/msantosfelipe/financial-chat/infra/amqp"
)

const stockValueTitle = "Close"
const invalidStock = "N/D"

func NewConsumer(
	amqpService amqp.AmqpService,
	websocketService websocket.WebsocketService,
) ConsumerService {
	return &consumerService{
		amqpService:      amqpService,
		websocketService: websocketService,
	}
}

func (s *consumerService) SubscribeToQueue(queue string) {
	messages := s.amqpService.SubscribeToQueue(queue)

	go func() {
		for message := range messages {
			var queueMessage QueueMessage
			log.Printf(" > Received message: %s\n", message.Body)
			if err := json.Unmarshal(message.Body, &queueMessage); err != nil {
				msg := fmt.Sprintf("error: %v", err)
				s.websocketService.SendBotMessage(queueMessage.Room, msg)
				continue
			}

			csvResponse, err := requestStockAPI(queueMessage.Stock)
			if err != nil {
				msg := fmt.Sprintf("error: %v", err)
				s.websocketService.SendBotMessage(queueMessage.Room, msg)
				continue
			}

			stockValue, err := parseCsvResponse(csvResponse, queueMessage.Stock)
			if err != nil {
				msg := fmt.Sprintf("error: %v", err)
				s.websocketService.SendBotMessage(queueMessage.Room, msg)
				continue
			}

			msg := fmt.Sprintf("%s quote is $%v per share.", queueMessage.Stock, stockValue)
			s.websocketService.SendBotMessage(queueMessage.Room, msg)
		}
	}()
}

func (s *consumerService) Clean() {
	s.amqpService.Clean()
}

func parseCsvResponse(csvResponse [][]string, stock string) (float64, error) {
	var stockValue float64
	var valueIndex int
	for indexRows, row := range csvResponse {
		if indexRows == 0 {
			for indexRow, column := range row {
				if column == stockValueTitle {
					valueIndex = indexRow
					break
				}
			}
		} else {
			for indexRow, column := range row {
				if indexRow == valueIndex {
					if column == invalidStock {
						return 0, fmt.Errorf("could not get values for stock %s", stock)
					}
					value, err := strconv.ParseFloat(column, 64)
					if err != nil {
						return 0, fmt.Errorf("error converting response of stock %s: %s", stock, err)
					}
					stockValue = value
					break
				}
			}
		}
	}

	return toFixed(stockValue, 2), nil
}

func requestStockAPI(stock string) ([][]string, error) {
	url := fmt.Sprintf("%s?s=%s&f=sd2t2ohlcv&h&e=csv", app.ENV.StockApiURL, stock)

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	client := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println("error calling stock api: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		log.Println("error parsing return from stock api: ", err)
		return nil, err
	}

	return records, nil
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
