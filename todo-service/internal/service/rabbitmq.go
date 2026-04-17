package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"todo-proj/internal/models"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

// 1 и 2: Dial + Channel
func NewRabbitMQ(url string, queueName string) (*RabbitMQ, error) {
	var conn *amqp.Connection
	var err error

	// Пробуем подключиться 5 раз с паузой в 3 секунды
	for i := 1; i <= 5; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		slog.Warn("RabbitMQ еще не готов, ждем...", "attempt", i, "err", err)
		time.Sleep(3 * time.Second)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Декларируем очередь сразу при создании
	_, err = ch.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		queue:   queueName,
	}, nil
}

// 3: Publish
func (r *RabbitMQ) PublishTaskCreated(ctx context.Context, task models.Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	err = r.channel.PublishWithContext(ctx,
		"",      // exchange
		r.queue, // routing key (наша очередь)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	slog.Info("событие отправлено в RabbitMQ", "task_id", task.ID)
	return nil
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}