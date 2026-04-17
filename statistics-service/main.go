package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskEvent struct {
	Title string `json:"title"`
}

func main() {
	// Подключаемся к RabbitMQ
	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	// Ждем немного, пока RabbitMQ точно прогрузится
	var conn *amqp.Connection
	var err error
	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		fmt.Println("Statistics: ждем RabbitMQ...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("Статистика не смогла подключиться: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Слушаем очередь
	msgs, err := ch.Consume(
		"tasks_events", // имя очереди
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Бесконечный цикл чтения
	slog.Info(" [***] Статистика запущена и слушает очередь...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var event TaskEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				slog.Error("Ошибка парсинга сообщения", "err", err)
				continue
			}
			fmt.Printf(" [x] Статистика: Получена новая задача: %s\n", event.Title)
		}
	}()

	<-forever
}