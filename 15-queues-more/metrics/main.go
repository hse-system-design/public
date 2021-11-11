package main

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"time"
)


/*
import "experimental/aggregate"


from(bucket: "mybucket")
  |> range(start: -5m)
  |> aggregate.rate(every: 30s, unit: 1s)

 */

func main() {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient("http://localhost:8186", "mytoken")
	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking("myorg", "mybucket")

	fmt.Println("Start...")

	cnt := 0.0
	for {
		p := influxdb2.NewPointWithMeasurement("stat").
			AddTag("unit", "point").
			AddField("value", cnt).
			SetTime(time.Now())

		err := writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			fmt.Println(err)
		}
		cnt++
		if int(cnt) % 1000 == 0 {
			fmt.Printf("Wrote %d points\n", int(cnt))
		}
	}
}
