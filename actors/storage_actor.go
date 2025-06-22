package actors

import (
	"context"
	"log"
	"time"

	"github.com/tochemey/goakt/v3/actor"
	"github.com/zaibon/surveilsense/proto"
)

type StorageBackend interface {
	SaveMetadata(ctx context.Context, cameraID string, timestamp time.Time, detections []*proto.Detection) error
	SaveClip(ctx context.Context, cameraID string, timestamp time.Time, imageClip []byte) error
}

type StorageActor struct {
	backend StorageBackend
}

var _ actor.Actor = (*StorageActor)(nil)

func NewStorageActor(backend StorageBackend) *StorageActor {
	return &StorageActor{backend: backend}
}

func (a *StorageActor) PreStart(ctx *actor.Context) error {
	return nil
}

func (a *StorageActor) Receive(ctx *actor.ReceiveContext) {
	msg := ctx.Message()
	event, ok := msg.(*proto.DetectionEvent)
	if !ok {
		ctx.Unhandled()
		return
	}

	ts := time.UnixMilli(event.Timestamp)
	if err := a.backend.SaveMetadata(ctx.Context(), event.CameraId, ts, event.Detections); err != nil {
		log.Printf("StorageActor: failed to save metadata for camera %s: %v", event.CameraId, err)
	} else {
		if err := a.backend.SaveClip(ctx.Context(), event.CameraId, ts, event.ImageClip); err != nil {
			log.Printf("StorageActor: failed to save clip for camera %s: %v", event.CameraId, err)
		}
	}
}

func (a *StorageActor) PostStop(ctx *actor.Context) error {
	type closer interface {
		Close() error
	}
	if c, ok := a.backend.(closer); ok {
		return c.Close()
	}
	return nil
}
