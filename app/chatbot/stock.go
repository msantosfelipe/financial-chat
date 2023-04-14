package chatbot

import "sync"

var stockInstance StockService
var stockOnce sync.Once

func GetStockInstance() StockService {
	stockOnce.Do(func() {
		stockInstance = NewStock()
	})

	return stockInstance
}

type stockService struct {
}

type StockService interface {
	stockRequester
}

type stockRequester interface {
	RequestStock(stock string)
}
