package main

import (
	"fmt"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
	"io/ioutil"
	"net"
)

// params.count / (params._interval / 1000)

func es() {
	log := logrus.New()
	configES, err := config.Parse("http://localhost:9200/index?sniff=false")
	if err != nil {
		log.Panic(err)
	}

	client, err := elastic.NewClientFromConfig(configES)
	if err != nil {
		log.Panic(err)
	}
	hook, err := elogrus.NewElasticHook(client, "localhost", logrus.DebugLevel, "mylog")
	if err != nil {
		log.Panic(err)
	}
	log.Hooks.Add(hook)
	log.Out = ioutil.Discard

	cnt := 0
	for {
		log.WithFields(logrus.Fields{
			"name": "joe",
			"age":  cnt,
		}).Error("Hello world!")
		cnt++
		if int(cnt) % 1000 == 0 {
			fmt.Printf("Wrote %d points\n", int(cnt))
		}
	}
}

func logstash() {
	log := logrus.New()
	conn, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Fatal(err)
	}
	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": "myappName"}))

	log.Hooks.Add(hook)
	log.Out = ioutil.Discard

	cnt := 0
	for {
		ctx := log.WithFields(logrus.Fields{
			"method": "main",
			"cnt": cnt,
		})
		cnt++
		ctx.Info("Hello World!")
		if int(cnt) % 1000 == 0 {
			fmt.Printf("Wrote %d points\n", int(cnt))
		}
	}
}

func main() {
	logstash()
}