package consumer

import (
	"testing"

	wsMock "github.com/msantosfelipe/financial-chat/app/websocket/mocks"
	amqpMock "github.com/msantosfelipe/financial-chat/infra/amqp/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewConsumer(t *testing.T) {
	amqpServiceMock := &amqpMock.AmqpService{}
	websocketMock := &wsMock.WebsocketService{}

	consumer := NewConsumer(amqpServiceMock, websocketMock)

	assert.Implements(t, (*ConsumerService)(nil), consumer)
}

func TestStockService_Clean(t *testing.T) {
	amqpServiceMock := &amqpMock.AmqpService{}
	websocketMock := &wsMock.WebsocketService{}
	consumer := NewConsumer(amqpServiceMock, websocketMock)

	// Arrange
	amqpServiceMock.On("Clean")

	// Act
	consumer.Clean()

	// Assert
	amqpServiceMock.AssertExpectations(t)
}
