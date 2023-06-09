package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/msantosfelipe/financial-chat/app"
	consumerHanler "github.com/msantosfelipe/financial-chat/app/consumer/handlers"
	wsHandler "github.com/msantosfelipe/financial-chat/app/websocket/handlers"
)

func main() {
	consumerHanler.New()
	wsHandler.New()

	// inianlt amqp consumer (bot)
	consumerHanler.HandleMessageConsumer()

	// init websocket (chat)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/websocket", wsHandler.HandleConnections)

	port := app.ENV.Port
	log.Println("Server starting at localhost:", port)

	go listenForShutdown()

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}

func listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("cleanning before exit...")
	wsHandler.Clean()
	consumerHanler.Clean()
	log.Println("stopping application")
	os.Exit(0)
}
