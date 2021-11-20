package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"os"
	"strconv"
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

	totalKey := "total_workers"
	workerKey := fmt.Sprintf("worker_%s", name)
	current := 0

	barrier := recipe.NewDoubleBarrier(s, "/before-start/", 3)
	fmt.Println("First barrier - ", name)
	barrier.Enter()

	cli.Put(ctx, totalKey, "0")

	fmt.Println("Second barrier - ", name)
	barrier.Leave()

	for {
		fmt.Println("Try commit tx - ", name)
		tx, err := cli.Txn(ctx).If(
			clientv3.Compare(clientv3.Value(totalKey), "=", strconv.Itoa(current)),
		).Then(
			clientv3.OpPut(totalKey, strconv.Itoa(current+1)),
			clientv3.OpPut(workerKey,  strconv.Itoa(current)),
		).Else(
			clientv3.OpGet(totalKey),
		).Commit()
		if err != nil {
			fmt.Printf("Error during tx %v\n", err)
			continue
		}

		if !tx.Succeeded {
			fmt.Println("Tx is not succeeded - ", name)
			current, err = strconv.Atoi(string(tx.Responses[0].GetResponseRange().Kvs[0].Value))
			if err != nil {
				panic(err)
			}
			fmt.Println("Retry with latest total value - ", current)
			continue
		}

		fmt.Println("Successfully commit tx! My number is ", current)
		break
	}


	watchCh := cli.Watch(context.TODO(), totalKey)
	r, _ := cli.Get(ctx, totalKey)
	total := r.Kvs[0].Value
	fmt.Printf("Current total = %s, watching for updates - %s\n", total, name)
	for res := range watchCh {
		total = res.Events[0].Kv.Value
		fmt.Printf("Get update! Current total = %s - %s\n", total, name)
	}
}