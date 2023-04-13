package websocket

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/msantosfelipe/financial-chat/app"
	"github.com/msantosfelipe/financial-chat/infra/cache"
)

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
	usersInRoom, err := w.getUsersInRoom(room)
	if err != nil {
		return err
	}

	usersInRoom[wsConn] = true
	w.usersByRoom[room] = usersInRoom
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

		msg.Room = room
		msg.Timestamp = time.Now().Format("2006-01-02 15:04:05")

		w.broadcaster <- msg
	}
}

func (w *websocketService) HandleReceivedMessages() {
	for {
		// grab any next message from channel
		msg := <-w.broadcaster

		w.cacheMessage(msg)
		w.messageClients(msg)
	}
}

func (w *websocketService) Clean() {
	log.Println("closing ws channel...")
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
	// send to every client currently connected
	for wsConn := range w.usersByRoom[msg.Room] {
		w.messageClient(wsConn, msg)
	}
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
