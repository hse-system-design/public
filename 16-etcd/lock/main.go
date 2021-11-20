package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	var name = os.Getenv("NAME")
	var gcWait, err = strconv.Atoi(os.Getenv("GC_WAIT"))
	if err != nil {
		panic(err)
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:30001", "localhost:30002", "localhost:30003"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	s, _:= concurrency.NewSession(cli, concurrency.WithTTL(10))
	defer s.Close()

	l := concurrency.NewMutex(s, "/distributed-lock/")
	ctx := context.Background()

	for {
		if err := l.Lock(ctx); err != nil {
			fmt.Printf("Error during locking - %v", err)
			continue
		}
		break
	}

	time.Sleep(time.Duration(rand.Float32() * 3) * time.Second)
	fmt.Println("Acquired lock for ", name)

	if gcWait > 0 {
		fmt.Println("Oh no, its stop the world in ", name)
		s.Orphan()  // do not monitor
		time.Sleep(time.Duration(gcWait) * time.Second) // wait for too long (e.g. gc)
	}

	fmt.Println("I DO SOME WORK IN ")

	if err := l.Unlock(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Released lock for ", name)
}