package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/msantosfelipe/financial-chat/app"
	"github.com/msantosfelipe/financial-chat/app/websocket/handlers"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/websocket", handlers.HandleConnections)

	port := app.ENV.Port
	log.Println("Server starting at localhost:", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}
