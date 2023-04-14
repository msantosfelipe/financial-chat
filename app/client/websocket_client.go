package client

import "github.com/msantosfelipe/financial-chat/app/websocket"

func SendBotMessage(botUser, room, text string) {
	ws := websocket.GetInstance()
	ws.SendBotMessage(botUser, room, text)
}
