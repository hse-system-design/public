package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"io"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
)

var (
	app *cli.App
)

const (
	ENCODE = "encode"
	DECODE = "decode"
)

func init() {
	// Initialise a CLI app
	app = cli.NewApp()
	app.Name = "machinery"
	app.Usage = "machinery worker and send example tasks with machinery send"
	app.Version = "0.0.0"
}

func main() {
	// Set the CLI app commands
	app.Commands = []cli.Command{
		{
			Name:  "worker",
			Usage: "launch machinery worker",
			Action: func(c *cli.Context) error {
				fmt.Println("Start worker...")
				if err := worker(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
		{
			Name:  "send",
			Usage: "send example tasks ",
			Action: func(c *cli.Context) error {
				fmt.Println("Send command...")
				taskType, data, password := c.Args().Get(0), c.Args().Get(1), c.Args().Get(2)
				if err := send(taskType, data, password); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
	}

	// Run the CLI app
	app.Run(os.Args)
}

func encodeTaskFunc(data, password string) (string, error) {
	complexity := 2048
	blockSize := 8
	parallelRate := 1
	keyLength := 32

	salt := []byte("1234")
	rawPass := []byte(password)
	key, err := scrypt.Key(rawPass, salt, complexity, blockSize, parallelRate, keyLength)
	if err != nil {
		return "", err
	}

	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	rawData := []byte(data)
	return base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, rawData, nil)), nil
}

func decodeTaskFunc(data, password string) (string, error) {
	complexity := 2048
	blockSize := 8
	parallelRate := 1
	keyLength := 32

	salt := []byte("1234")
	rawPass := []byte(password)
	key, err := scrypt.Key(rawPass, salt, complexity, blockSize, parallelRate, keyLength)
	if err != nil {
		return "", err
	}

	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", err
	}

	rawEncData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	nonce := rawEncData[:gcm.NonceSize()]
	encryptedData := rawEncData[gcm.NonceSize():]

	rawData, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func startServer() (*machinery.Server, error) {
	cnf := &config.Config{
		DefaultQueue:    "machinery_tasks",
		ResultsExpireIn: 3600,
		Broker:          "redis://localhost:6379",
		ResultBackend:   "redis://localhost:6379",
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, err
	}

	// Register tasks
	tasks := map[string]interface{}{
		"encode": encodeTaskFunc,
		"decode": decodeTaskFunc,
	}

	return server, server.RegisterTasks(tasks)
}

func worker() error {
	consumerTag := "machinery_worker"

	server, err := startServer()
	if err != nil {
		return err
	}

	worker := server.NewWorker(consumerTag, 0)

	errorhandler := func(err error) {
		log.ERROR.Println("Something went wrong:", err)
	}

	worker.SetErrorHandler(errorhandler)

	return worker.Launch()
}

func send(taskType, data, password string) error {
	server, err := startServer()
	if err != nil {
		return err
	}

	var task tasks.Signature

	switch taskType {
	case ENCODE:
		task = tasks.Signature{
			Name: "encode",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: data,
				},
				{
					Type: "string",
					Value: password,
				},
			},
		}
	case DECODE:
		task = tasks.Signature{
			Name: "decode",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: data,
				},
				{
					Type: "string",
					Value: password,
				},
			},
		}
	default:
		panic(fmt.Errorf("no such command %s", taskType))
	}

	asyncResult, err := server.SendTaskWithContext(context.Background(), &task)
	if err != nil {
		return fmt.Errorf("could not send task: %s", err.Error())
	}

	results, err := asyncResult.Get(time.Duration(1 * time.Second))
	if err != nil {
		return fmt.Errorf("getting task result failed with error: %s", err.Error())
	}
	log.INFO.Printf("%v\n", tasks.HumanReadableResults(results))

	return nil
}
