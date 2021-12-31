package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

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

	// Create a channel to interact with RabbitMQ.
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare the queue we're sending to - this is idempotent.
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Create some messages.
	for i := 0; i < 10; i++ {
		body := strconv.Itoa(rand.Intn(1000))
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf("Published Message: %s", body)
	}
}
