package actors

import (
	"context"
	"testing"
	"time"

	"github.com/tochemey/goakt/v3/actor"
	"github.com/zaibon/surveilsense/proto"
)

func TestExternalAsk(t *testing.T) {
	as, err := actor.NewActorSystem("TestSystem")
	if err != nil {
		t.Fatalf("Failed to create actor system: %v", err)
	}

	as.Start(context.Background())
	defer as.Stop(context.Background())

	pid, err := as.Spawn(context.Background(), "testActor", &testActor{})
	if err != nil {
		t.Fatalf("Failed to spawn actor: %v", err)
	}

	resp, err := actor.Ask(context.Background(), pid, &proto.GetLatestFrame{}, time.Second)
	if err != nil {
		t.Fatalf("Ask failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	<-time.After(time.Second * 10)
}

type testActor struct{}

func (a *testActor) PreStart(ctx *actor.Context) error {
	return nil
}

func (a *testActor) Receive(ctx *actor.ReceiveContext) {
	msg := ctx.Message()
	if msg == nil {
		ctx.Unhandled()
		return
	}

	ctx.Tell(ctx.Sender(), &proto.FrameData{})
}

func (a *testActor) PostStop(ctx *actor.Context) error {
	return nil
}
