package producer

import (
	"fmt"
	"github.com/patyukin/mbs-auth/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func New(cfg *config.Config) (*Producer, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		if err = conn.Close(); err != nil {
			return nil, fmt.Errorf("failed to close connection: %w", err)
		}

		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		cfg.RabbitMQ.QueueName,
		cfg.RabbitMQ.Durable,
		cfg.RabbitMQ.AutoDelete,
		cfg.RabbitMQ.Exclusive,
		cfg.RabbitMQ.NoWait,
		nil,
	)
	if err != nil {

		if err = ch.Close(); err != nil {
			return nil, fmt.Errorf("failed to close channel: %w", err)
		}

		if err = conn.Close(); err != nil {
			return nil, fmt.Errorf("failed to close connection: %w", err)
		}

		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &Producer{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (p *Producer) SendMessage(body []byte) error {
	err := p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Nessage sent: %s", body)
	return nil
}

func (p *Producer) Close() {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			log.Printf("failed to close channel: %v", err)
		}
	}

	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}
}
