package chatbot

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/msantosfelipe/financial-chat/app"
)

func NewStock() StockService {
	return &stockService{}
}

func (s *stockService) RequestStock(stock string) {
	url := fmt.Sprintf("%s?s=%s&f=sd2t2ohlcv&h&e=csv", app.ENV.StockApiURL, stock)

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	client := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println("error calling stock api: ", err)
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, row := range records {
		fmt.Println(strings.Join(row, ","))
	}

}
