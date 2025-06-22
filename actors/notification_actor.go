package actors

import (
	"log"

	"github.com/tochemey/goakt/v3/actor"
	"github.com/zaibon/surveilsense/proto"
)

type Notifier interface {
	Notify(event *proto.DetectionEvent) error
}

// NotificationActor receives DetectionEvent and notifies via all configured Notifiers
type NotificationActor struct {
	notifiers []Notifier
}

// NewNotificationActor creates a NotificationActor with the given notifiers
func NewNotificationActor(notifiers ...Notifier) *NotificationActor {
	return &NotificationActor{notifiers: notifiers}
}

var _ actor.Actor = (*NotificationActor)(nil)

func (a *NotificationActor) PreStart(ctx *actor.Context) error {
	return nil
}

func (a *NotificationActor) Receive(ctx *actor.ReceiveContext) {
	msg := ctx.Message()
	event, ok := msg.(*proto.DetectionEvent)
	if !ok {
		ctx.Unhandled()
		return
	}
	for _, notifier := range a.notifiers {
		if err := notifier.Notify(event); err != nil {
			log.Printf("NotificationActor: failed to notify: %v", err)
		}
	}
}

func (a *NotificationActor) PostStop(ctx *actor.Context) error {
	return nil
}
