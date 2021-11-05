package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	rand.Seed(time.Now().Unix())

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	fmt.Println("Waiting 15 seconds for warming up...")
	time.Sleep(15 * time.Second)
	fmt.Println("Let's rock!")

	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	messages, err := channelRabbitMQ.Consume(
		"QueueService1", // queue name
		"",              // consumer
		false,            // auto-ack
		false,           // exclusive
		false,           // no local
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	id := os.Getenv("ID")
	if id == "" {
		id = "default"
	}

	for message := range messages {
		log.Printf("Consumer(%v) > Received message: %s\n", id, message.Body)
		time.Sleep(100 * time.Millisecond)

		if rand.Int() % 5 == 0 {
			fmt.Println("Something went wrong!")
			if err := message.Nack(false, true); err != nil {
				fmt.Println(err)
			}
			continue
		}

		if err := message.Ack(false); err != nil {
			fmt.Println(err)
		}
	}
}
