package websocket

import (
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
	amqpMock "github.com/msantosfelipe/financial-chat/infra/amqp/mocks"
	cacheMock "github.com/msantosfelipe/financial-chat/infra/cache/mocks"

	"github.com/stretchr/testify/assert"
)

func newMockInstance() WebsocketService {
	return NewInstance(
		make(map[string]map[*websocket.Conn]bool),
		make(chan ChatMessage),
		websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		&cacheMock.CacheService{},
		&amqpMock.AmqpService{},
	)
}

func TestNewInstance(t *testing.T) {
	wsInstance := newMockInstance()
	assert.Implements(t, (*WebsocketService)(nil), wsInstance)
}
