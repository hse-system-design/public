package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"os"
	"time"
)

func main() {
	var name = os.Getenv("NAME")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:30001", "localhost:30002", "localhost:30003"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	ctx := context.Background()

	s, _:= concurrency.NewSession(cli)
	defer s.Close()

	barrier := recipe.NewDoubleBarrier(s, "/before-start/", 3)
	fmt.Println("First barrier - ", name)
	barrier.Enter()

	election := concurrency.NewElection(s, "/my-election-3/")

	fmt.Println("Second barrier - ", name)
	barrier.Leave()

	fmt.Println("Run campaign - ", name)
	ctxTTL, _ := context.WithTimeout(ctx, 5 * time.Second)
	err = election.Campaign(ctxTTL, name)
	if err != nil {
		panic(err)
	}

	fmt.Println("Look for results - ", name)
	r, _ := election.Leader(ctx)
	fmt.Println("result - ", string(r.Kvs[0].Value))
	for obs := range election.Observe(ctx) {
		fmt.Printf("New leader = %s\n", obs.Kvs[0].Value)
	}
}