package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/msantosfelipe/financial-chat/infra/amqp"
	"github.com/msantosfelipe/financial-chat/infra/cache"
)

var instance WebsocketService
var once sync.Once

func GetWSInstance() WebsocketService {
	once.Do(func() {
		instance = NewInstance(
			make(map[string]map[*websocket.Conn]bool),
			make(chan ChatMessage),
			websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			cache.GetInstance(),
			amqp.GetInstance(),
		)
	})

	return instance
}

type websocketService struct {
	usersByRoom  map[string]map[*websocket.Conn]bool
	broadcaster  chan ChatMessage
	upgrader     websocket.Upgrader
	cacheService cache.CacheService
	amqpService  amqp.AmqpService
	mutex        sync.Mutex
}

type ChatMessage struct {
	Username  string `json:"username"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
	Room      string `json:"room"`
}

type QueueStockMessage struct {
	Stock string `json:"stock"`
	Room  string `json:"room"`
}

type WebsocketService interface {
	connRegister
	userRegister
	MessageSender
	MessageReceiver
	ServiceCleaner
	chatbotHandler
}

type connRegister interface {
	RegisterWSConnection(w http.ResponseWriter, r *http.Request) *websocket.Conn
}

type userRegister interface {
	AddUserToRoom(wsConn *websocket.Conn, room string) error
}

type MessageSender interface {
	SendPreviousCachedMessages(wsConn *websocket.Conn, room string)
	ListenAndSendMessage(wsConn *websocket.Conn, room string)
	SendMessage(user, room, text string)
	SendBotMessage(room, text string)
	PublishMessageToQueue(msg []byte, queue string) error
}

type MessageReceiver interface {
	HandleReceivedMessages()
}

type ServiceCleaner interface {
	Clean()
}

type chatbotHandler interface {
	HandleBotMessage(text, room string)
	StockHandler(stock, room string) error
}
