package amqp

import (
	"errors"
	"testing"

	amqpMock "github.com/msantosfelipe/financial-chat/infra/amqp/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAmqpService_PublishMessage(t *testing.T) {
	amqpServiceMock := &amqpMock.AmqpService{}
	msg := []byte("testing")

	amqpServiceMock.On("PublishMessage", []byte("testing"), "test_queue").Return(nil)

	err := amqpServiceMock.PublishMessage(msg, "test_queue")

	assert.NoError(t, err)
}

func TestAmqpService_PublishMessage_Error(t *testing.T) {
	amqpServiceMock := &amqpMock.AmqpService{}
	msg := []byte("testing")
	expectedErr := errors.New("error publishing message")

	amqpServiceMock.On("PublishMessage", []byte("testing"), "test_queue").Return(expectedErr)

	err := amqpServiceMock.PublishMessage(msg, "test_queue")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
}
