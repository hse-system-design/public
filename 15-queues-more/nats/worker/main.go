package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	uri := os.Getenv("NATS_URI")
	var err error
	var nc *nats.Conn

	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(uri)
		if err == nil {
			break
		}

		fmt.Println("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}
	fmt.Println("Connected to NATS at:", nc.ConnectedUrl())

	_, err = nc.Subscribe("tasks", func(m *nats.Msg) {
		fmt.Println("Got task request on:", m.Subject)

		data := m.Data
		result := base64.StdEncoding.EncodeToString(data)

		err := nc.Publish(m.Reply, []byte(result))
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listen for tasks from %v\n", nc.ConnectedUrl())

	runtime.Goexit()
}