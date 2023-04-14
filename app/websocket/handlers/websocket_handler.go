package handlers

import (
	"log"
	"net/http"

	"github.com/msantosfelipe/financial-chat/app/websocket"
)

var ws websocket.WebsocketService

func init() {
	ws = websocket.GetInstance()
	go ws.HandleReceivedMessages()
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Create new connection
	wsConn := ws.RegisterWSConnection(w, r)
	defer wsConn.Close()

	// Get room from url
	room := r.URL.Query().Get("room")
	if room == "" {
		log.Println("room not specified")
		return
	}

	log.Println("New connection on room:", room)

	// Register user to the room
	if err := ws.AddUserToRoom(wsConn, room); err != nil {
		log.Println("error adding user to room: ", err)
		return
	}

	// if it's zero, no messages were ever sent/saved
	ws.SendPreviousCachedMessages(wsConn, room)

	ws.ListenAndSendMessage(wsConn, room)
}

func Clean() {
	log.Println("cleanning websocket tasks...")
	ws.Clean()
	log.Println("websocket stopped")
}
