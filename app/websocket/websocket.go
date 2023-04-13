package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/msantosfelipe/financial-chat/infra/cache"
)

var instance WebsocketService
var once sync.Once

func GetInstance() WebsocketService {
	once.Do(func() {
		log.Println("Creating websocket instance")
		instance = New()
	})

	return instance
}

type websocketService struct {
	UsersByRoom  map[string]map[*websocket.Conn]bool
	Broadcaster  chan ChatMessage
	Upgrader     websocket.Upgrader
	cacheService cache.CacheService
}

type ChatMessage struct {
	Username  string `json:"username"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
	Room      string `json:"room"`
}

type WebsocketService interface {
	connRegister
	userRegister
	MessageSender
	MessageReceiver
	ServiceCleaner
}

type connRegister interface {
	RegisterWSConnection(w http.ResponseWriter, r *http.Request) *websocket.Conn
}

type userRegister interface {
	AddUserToRoom(wsConn *websocket.Conn, room string)
}

type MessageSender interface {
	SendPreviousCachedMessages(wsConn *websocket.Conn, room string)
	ListenToAndSendMessage(wsConn *websocket.Conn, room string)
}

type MessageReceiver interface {
	HandleReceivedMessages()
}

type ServiceCleaner interface {
	Clean()
}
