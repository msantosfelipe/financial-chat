package handlers

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/msantosfelipe/financial-chat/app/websocket"
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Create connection
	ws := websocket.GetInstance()

	go ws.HandleReceivedMessages()
	go listenForShutdown(ws)

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
	ws.AddUserToRoom(wsConn, room)

	// if it's zero, no messages were ever sent/saved
	ws.SendPreviousCachedMessages(wsConn, room)

	ws.ListenToAndSendMessage(wsConn, room)
}

func listenForShutdown(ws websocket.WebsocketService) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("would run cleanup tasks...")
	ws.Clean()
	log.Println("shutting down application...")
	os.Exit(0)
}
