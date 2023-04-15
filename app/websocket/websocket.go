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
	// Register new websocket connection
	RegisterWSConnection(w http.ResponseWriter, r *http.Request) *websocket.Conn
}

type userRegister interface {
	// Add websocket connection to room
	AddUserToRoom(wsConn *websocket.Conn, room string) error
}

type MessageSender interface {
	// Send cached messages to websocket connection
	SendPreviousCachedMessages(wsConn *websocket.Conn, room string)
	// Listen for new messages in chat room and send to connected users in room
	ListenAndSendMessage(wsConn *websocket.Conn, room string)
	// Send new websocket message to connected users in room
	SendMessage(user, room, text string)
	// Send new websocket message from bot to connected users in room
	SendBotMessage(room, text string)
	// Publish message to AMQP queue
	PublishMessageToQueue(msg []byte, queue string) error
}

type MessageReceiver interface {
	// Handle received websocket messages
	HandleReceivedMessages()
}

type ServiceCleaner interface {
	// Close connections and channels
	Clean()
}

type chatbotHandler interface {
	// Handle bot command
	HandleBotMessage(text, room string)
	// Handle stock bot command
	StockHandler(stock, room string) error
}
