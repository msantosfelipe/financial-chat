package websocket

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/msantosfelipe/financial-chat/app"
	"github.com/msantosfelipe/financial-chat/infra/amqp"
	"github.com/msantosfelipe/financial-chat/infra/cache"
)

const botMsgPrefix = "/"

var chatbot ChatbotService

func init() {
	chatbot = GetChatbotInstance()
}

func New() WebsocketService {
	return &websocketService{
		usersByRoom: make(map[string]map[*websocket.Conn]bool),
		broadcaster: make(chan ChatMessage),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		cacheService: cache.GetInstance(),
		amqpService:  amqp.GetInstance(),
	}
}

func (w *websocketService) RegisterWSConnection(rw http.ResponseWriter, r *http.Request) *websocket.Conn {
	wsConn, err := w.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Fatal("error creating new websocket connection: ", err)
	}

	return wsConn
}

func (w *websocketService) AddUserToRoom(wsConn *websocket.Conn, room string) error {
	w.mutex.Lock()
	usersInRoom, err := w.getUsersInRoom(room)
	if err != nil {
		return err
	}

	usersInRoom[wsConn] = true
	w.usersByRoom[room] = usersInRoom
	w.mutex.Unlock()
	return nil
}

func (w *websocketService) SendPreviousCachedMessages(wsConn *websocket.Conn, room string) {
	if !w.cacheService.ExistsChatKey(room) {
		return
	}

	chatMessages := w.cacheService.GetPreviousChatMessages(room)
	for _, chatMessage := range chatMessages {
		var msg ChatMessage
		json.Unmarshal([]byte(chatMessage), &msg)
		w.messageClient(wsConn, msg)
	}
}

func (w *websocketService) ListenAndSendMessage(wsConn *websocket.Conn, room string) {
	for {
		var msg ChatMessage
		err := wsConn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error sending message on room %s: %v\n", room, err)
			delete(w.usersByRoom[room], wsConn)
			break
		}

		if strings.HasPrefix(msg.Text, botMsgPrefix) {
			chatbot.HandleBotMessage(msg.Text, room)
			continue
		}

		w.SendMessage(msg.Username, room, msg.Text)
	}
}

func (w *websocketService) SendMessage(user, room, text string) {
	msg := ChatMessage{
		Username:  user,
		Room:      room,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Text:      text,
	}

	w.broadcaster <- msg
}

func (w *websocketService) SendBotMessage(room, text string) {
	w.SendMessage(app.ENV.ChatbotUsername, room, text)
}

func (w *websocketService) PublishMessageToQueue(msg []byte, queue string) error {
	if err := w.amqpService.PublishMessage(msg, app.ENV.AmqpChatQueueName); err != nil {
		log.Println("error publishing message: ", err)
		return err
	}
	return nil
}

func (w *websocketService) HandleReceivedMessages() {
	for {
		msg, ok := <-w.broadcaster
		if !ok {
			return
		}

		w.cacheMessage(msg)
		w.messageClients(msg)
	}
}

func (w *websocketService) Clean() {
	close(w.broadcaster)
}

func (w *websocketService) cacheMessage(msg ChatMessage) {
	w.cacheService.HandleChatSize(msg.Room)

	json, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}

	w.cacheService.StoreChatMessage(json, msg.Room)
}

func (w *websocketService) messageClients(msg ChatMessage) {
	w.mutex.Lock()
	for wsConn := range w.usersByRoom[msg.Room] {
		w.messageClient(wsConn, msg)
	}
	w.mutex.Unlock()
}

func (w *websocketService) messageClient(wsConn *websocket.Conn, msg ChatMessage) {
	err := wsConn.WriteJSON(msg)
	if err != nil && w.unsafeError(err) {
		log.Printf("error sending message on room %s: %v\n", msg.Room, err)
		wsConn.Close()
		delete(w.usersByRoom[msg.Room], wsConn)
	}
}

func (w *websocketService) getUsersInRoom(room string) (map[*websocket.Conn]bool, error) {
	if len(w.usersByRoom) >= app.ENV.MaxRooms {
		return nil, errors.New("max number of rooms exceeded")
	}

	usersInRoom := w.usersByRoom[room]
	if usersInRoom == nil {
		usersInRoom = make(map[*websocket.Conn]bool)
	}

	if len(usersInRoom) >= app.ENV.MaxClientsPerRoom {
		return nil, errors.New("max number of clients per room exceeded")
	}

	return usersInRoom, nil
}

// If a message is sent while a client is closing, ignore the error
func (w *websocketService) unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}
