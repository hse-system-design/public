package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

func main() {
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

	_, err = channelRabbitMQ.QueueDeclare(
		"QueueService1", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		panic(err)
	}

	wait, err := strconv.Atoi(os.Getenv("WAIT"))
	if err != nil {
		wait = 1000
	}

	for i := 0; true; i++ {
		body := []byte(fmt.Sprintf("Message with number %v", i))
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		}

		if err := channelRabbitMQ.Publish(
			"",              // exchange
			"QueueService1", // queue name
			false,           // mandatory
			false,           // immediate
			message,         // message to publish
		); err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Duration(wait) * time.Millisecond)
	}
}