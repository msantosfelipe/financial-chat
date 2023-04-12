package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type ChatMessage struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	Room     string `json:"room"`
}

var (
	rdb                *redis.Client
	chatExpirationTime = 60 * time.Minute
)

var clientsByRoom = make(map[string]map[*websocket.Conn]bool)
var broadcaster = make(chan ChatMessage)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	room := r.URL.Query().Get("room")
	if room == "" {
		log.Printf("Room not specified")
		return
	}

	fmt.Println("Starting new connection on room:", room)

	clientsInTheRoom := clientsByRoom[room]
	if clientsInTheRoom == nil {
		clientsInTheRoom = make(map[*websocket.Conn]bool)
	}
	clientsInTheRoom[ws] = true

	// ensure connection close when function returns
	defer ws.Close()
	clientsByRoom[room] = clientsInTheRoom

	// if it's zero, no messages were ever sent/saved
	if rdb.Exists("chat_messages:"+room).Val() != 0 {
		sendPreviousMessages(ws, room)
	}

	for {
		var msg ChatMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("deu erro")
			delete(clientsByRoom[room], ws)
			break
		}
		msg.Room = room
		// send new message to the channel
		broadcaster <- msg
	}
}

func sendPreviousMessages(ws *websocket.Conn, room string) {
	chatMessages, err := rdb.LRange("chat_messages:"+room, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	// send previous messages
	for _, chatMessage := range chatMessages {
		var msg ChatMessage
		json.Unmarshal([]byte(chatMessage), &msg)
		messageClient(ws, msg)
	}
}

// If a message is sent while a client is closing, ignore the error
func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}

func handleMessages() {
	for {
		// grab any next message from channel
		msg := <-broadcaster

		storeInRedis(msg)
		messageClients(msg)
	}
}

func storeInRedis(msg ChatMessage) {
	json, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	redisKey := fmt.Sprintf("chat_messages:%s", msg.Room)
	if err := rdb.LPush(redisKey, json).Err(); err != nil {
		panic(err)
	}
	if err := rdb.Expire(redisKey, chatExpirationTime); err.Err() != nil {
		panic(err)
	}
}

func messageClients(msg ChatMessage) {
	// send to every client currently connected
	for client := range clientsByRoom[msg.Room] {
		messageClient(client, msg)
	}
}

func messageClient(client *websocket.Conn, msg ChatMessage) {
	err := client.WriteJSON(msg)
	if err != nil && unsafeError(err) {
		log.Printf("error: %v", err)
		client.Close()
		delete(clientsByRoom[msg.Room], client)

	}
}

func main() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	port := os.Getenv("PORT")
	initRedis()

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/websocket", handleConnections)

	// create a map to store the clients for each room
	clientsByRoom = make(map[string]map[*websocket.Conn]bool)

	go handleMessages()

	log.Print("Server starting at localhost:" + port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func initRedis() {
	redisURL := os.Getenv("REDIS_URL")

	rdb = redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
}
