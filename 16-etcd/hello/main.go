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
	resp1, err := cli.Put(ctx, "biba", "kuka")
	if err != nil {
		panic(err)
	}

	fmt.Println(resp1.Header.String())

	resp2, err := cli.Get(ctx, "biba")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Number of keys: %d\n", resp2.Count)
	fmt.Printf("Key: %v\n", resp2.Kvs)

	resp3, err := cli.Delete(ctx, "biba")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Delete: %v\n", resp3)

	resp4, err := cli.Get(ctx, "biba")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Get after delete: %v\n", resp4)
}