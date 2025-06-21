package actors

import (
	"github.com/tochemey/goakt/v3/actor"
)

// NotificationActor handles alerts based on DetectionEvent
// (Notification logic will be added later)
type NotificationActor struct{}

var _ actor.Actor = (*NotificationActor)(nil)

func NewNotificationActor() *NotificationActor {
	return &NotificationActor{}
}

func (a *NotificationActor) PreStart(ctx *actor.Context) error {
	return nil
}

func (a *NotificationActor) Receive(ctx *actor.ReceiveContext) {
	// TODO: Implement notification logic
	ctx.Unhandled()
}

func (a *NotificationActor) PostStop(ctx *actor.Context) error {
	return nil
}
