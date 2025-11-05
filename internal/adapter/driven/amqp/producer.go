package amqp

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type ProducersAMQP struct {
	AuthConsumer *amqp.Queue
	Channel      *amqp.Channel
}

func NewProducersAMQP(authConsumer *amqp.Queue,
	Channel *amqp.Channel) *ProducersAMQP {
	return &ProducersAMQP{
		AuthConsumer: authConsumer,
		Channel:      Channel,
	}
}

func (p *ProducersAMQP) Publish(ctx context.Context, queueName string, msg []byte) error {
	//body, err := json.Marshal(&Message{
	//	"nekruzrakhimov@icloud.com",
	//	"Появились новые товары",
	//	"Чудо-фен с блютузом",
	//})
	err := p.Channel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key (имя очереди)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         msg,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}

func InitAMQPProducer(host string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(host)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}

	return conn, ch
}

func InitQueue(ch *amqp.Channel, queueName string) (*amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare a queue:", err)
	}

	return &queue, nil
}
