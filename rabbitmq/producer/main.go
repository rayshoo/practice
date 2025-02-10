package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@10.254.1.51:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer func() { _ = ch.Close() }()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello World!"
	for i := 0; i < 5; i++ {
		ctx, _ := context.WithCancel(context.Background())
		go func() {
			err = ch.PublishWithContext(
				ctx,
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			failOnError(err, "Failed to publish a message")
			log.Printf(" [x] Sent %s\n", body)
		}()
	}
	time.Sleep(5 * time.Second)
}
