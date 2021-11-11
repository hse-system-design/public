package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"remoteactivate/messages"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

type kukaActor struct {
	Count   int
	bibaPID *actor.PID
}

func NewKukaActor(bibaPID *actor.PID) actor.Actor {
	return &kukaActor{
		bibaPID: bibaPID,
	}
}

func (a *kukaActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		log.Printf("Kuka actor started " + context.Self().String())
	case *messages.HelloResponse:
		a.handle(context, context.Message().(*messages.HelloResponse))
	}
}

func (a *kukaActor) handle(context actor.Context, msg *messages.HelloResponse) {
	log.Printf("Kuka actor: %v\n", msg.Message)

	a.Count++
	fmt.Printf("Actor handle count %v", a.Count)

	message := msg.Message + " kuka"

	time.Sleep(1 * time.Second)

	resp := &messages.HelloResponse{Message: message}
	context.Send(a.bibaPID, resp)
}

func main() {
	system := actor.NewActorSystem()

	config := remote.Configure("0.0.0.0", 8080)
	remoting := remote.NewRemote(system, config)
	remoting.Start()

	bibaPid := actor.NewPID("biba:8080", "bibaActor")
	kukaProps := actor.PropsFromProducer(func() actor.Actor {
		return NewKukaActor(bibaPid)
	})
	bibaPID, err := system.Root.SpawnNamed(kukaProps, "kukaActor")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Kuka actor has been started with PID %v", bibaPID.String())

	runtime.Goexit()
}
