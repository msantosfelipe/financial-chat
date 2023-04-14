package websocket

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/msantosfelipe/financial-chat/app"
	"github.com/msantosfelipe/financial-chat/infra/amqp"
)

var ws WebsocketService

const (
	chatbotUser = "bot"
	prefixStock = "/stock="
	prefixHelp  = "/help"

	stockValueTitle = "Close"
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
			ws.SendMessage(chatbotUser, room, msg)
		}
	case strings.HasPrefix(text, prefixHelp):
		msg := "*** usage: \"/stock='stock_code'\""
		ws.SendMessage(chatbotUser, room, msg)
	default:
		msg := fmt.Sprintf("*** invalid bot command %s", text)
		ws.SendMessage(chatbotUser, room, msg)
	}
}

func (s *chatbotService) StockHandler(text, room string) error {
	if len(strings.Split(text, prefixStock)) == 1 {
		return fmt.Errorf("invalid stock name %s", text)
	}

	stock := strings.ToUpper(strings.Split(text, prefixStock)[1])

	// TODO substituir para publish message

	csvResponse, err := requestStockAPI(stock)
	if err != nil {
		return err
	}

	stockValue, err := parseCsvResponse(csvResponse)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("%s quote is $%v per share.", stock, stockValue)
	fmt.Println(msg)
	ws.SendMessage(chatbotUser, room, msg)

	return nil
}

func parseCsvResponse(csvResponse [][]string) (float64, error) {
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
					value, err := strconv.ParseFloat(column, 64)
					if err != nil {
						log.Fatal(err)
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
