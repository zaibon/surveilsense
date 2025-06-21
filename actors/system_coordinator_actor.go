package actors

import (
	"github.com/tochemey/goakt/v3/actor"
)

// SystemCoordinatorActor manages system state and orchestration
// (Coordination logic will be added later)
type SystemCoordinatorActor struct{}

var _ actor.Actor = (*SystemCoordinatorActor)(nil)

func NewSystemCoordinatorActor() *SystemCoordinatorActor {
	return &SystemCoordinatorActor{}
}

func (a *SystemCoordinatorActor) PreStart(ctx *actor.Context) error {
	return nil
}

func (a *SystemCoordinatorActor) Receive(ctx *actor.ReceiveContext) {
	// TODO: Implement coordination logic
	ctx.Unhandled()
}

func (a *SystemCoordinatorActor) PostStop(ctx *actor.Context) error {
	return nil
}
