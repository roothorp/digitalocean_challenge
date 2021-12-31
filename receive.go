package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// failOnError is a helper function, so I can type if err != nil less.
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Grab parameters
	username := os.Getenv("RABBITMQ_USERNAME")
	if username == "" {
		log.Fatalf("No Username env variable set")
	}
	password := os.Getenv("RABBITMQ_PASSWORD")
	if password == "" {
		log.Fatalf("No password env variable set")
	}
	service := os.Getenv("ROO_TEST_SERVICE_HOST")
	if service == "" {
		log.Fatalf("No service env variable set")
	}

	connection := fmt.Sprintf("amqp://%s:%s@%s", username, password, service)
	// Connect to RabbitMQ.
	conn, err := amqp.Dial(connection)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
