package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:30001", "localhost:30002", "localhost:30003"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	ctx := context.Background()


	for i := 1; ; i++ {
		key, value := fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i)
		_, err := cli.Put(ctx, key, value)
		if err != nil {
			fmt.Printf("Error during write: %v\n", err)
		} else {
			fmt.Printf("Successfull write %s -> %s\n", key, value)
		}

		r, err := cli.Get(ctx, key)
		if err != nil {
			fmt.Printf("Error during read: %v\n", err)
		} else {
			fmt.Printf("Successfull read %s -> %s\n", r.Kvs[0].Key, r.Kvs[0].Value)
		}

		time.Sleep(1 * time.Second)
	}
}