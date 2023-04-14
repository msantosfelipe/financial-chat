package amqp

import (
	"context"
	"log"
	"time"

	amqpgo "github.com/rabbitmq/amqp091-go"
)

type amqpService struct {
	client  *amqpgo.Connection
	channel *amqpgo.Channel
}

func New(host, queue string) AmqpService {
	connectRabbitMQ, err := amqpgo.Dial(host)
	if err != nil {
		panic(err)
	}

	channel, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}

	_, err = channel.QueueDeclare(
		queue, // queue name
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}

	return &amqpService{
		client:  connectRabbitMQ,
		channel: channel,
	}
}

func (s *amqpService) Clean() {
	s.channel.Close()
	s.client.Close()
}

func (s *amqpService) SubscribeToQueue(queue string) <-chan amqpgo.Delivery {
	messages, err := s.channel.Consume(
		queue, // queue name
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	return messages
}

func (s *amqpService) PublishMessage(msg []byte, queue string) error {
	message := amqpgo.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(1*time.Second))
	if err := s.channel.PublishWithContext(
		ctx,
		"",      // exchange
		queue,   // queue name
		false,   // mandatory
		false,   // immediate
		message, // message to publish
	); err != nil {
		return err
	}

	return nil
}
