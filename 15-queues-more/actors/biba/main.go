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

type bibaActor struct {
	Count int
	kukaPID *actor.PID
}

func NewBibaActor(kukaPID *actor.PID) actor.Actor {
	return &bibaActor{
		kukaPID: kukaPID,
	}
}

func (a *bibaActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		log.Printf("Biba actor started " + context.Self().String())
		time.Sleep(1 * time.Second)

		context.Send(a.kukaPID, &messages.HelloResponse{Message: "Biba"})
	case *messages.HelloResponse:
		a.handle(context, context.Message().(*messages.HelloResponse))
	}
}

func (a *bibaActor) handle(context actor.Context, msg *messages.HelloResponse) {
	log.Printf("Biba actor: %v\n", msg.Message)

	a.Count++
	fmt.Printf("Actor handle count %v", a.Count)

	message := msg.Message + " biba"

	time.Sleep(1 * time.Second)

	resp := &messages.HelloResponse{Message: message}
	context.Send(a.kukaPID, resp)
}

func main() {
	system := actor.NewActorSystem()

	config := remote.Configure("0.0.0.0", 8080)
	remoting := remote.NewRemote(system, config)
	remoting.Start()

	kukaPid := actor.NewPID("kuka:8080", "kukaActor")
	bibaProps := actor.PropsFromProducer(func() actor.Actor {
		return NewBibaActor(kukaPid)
	})
	bibaPID, err := system.Root.SpawnNamed(bibaProps, "bibaActor")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Biba actor has been started with PID %v", bibaPID.String())

	runtime.Goexit()
}
