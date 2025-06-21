package actors

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/tochemey/goakt/v3/actor"
	"github.com/zaibon/surveilsense/proto"
)

// StorageActor persists DetectionEvent data to a log file
// For demo: writes JSON lines to detections.log and saves image clips if present

type StorageActor struct {
	logFile *os.File
}

var _ actor.Actor = (*StorageActor)(nil)

func NewStorageActor() *StorageActor {
	return &StorageActor{}
}

func (a *StorageActor) PreStart(ctx *actor.Context) error {
	f, err := os.OpenFile("detections.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	a.logFile = f
	return nil
}

func (a *StorageActor) Receive(ctx *actor.ReceiveContext) {
	msg := ctx.Message()
	event, ok := msg.(*proto.DetectionEvent)
	if !ok {
		ctx.Unhandled()
		return
	}

	// Write detection metadata as JSON line
	meta := map[string]interface{}{
		"camera_id":  event.CameraId,
		"timestamp":  event.Timestamp,
		"detections": event.Detections,
	}
	b, err := json.Marshal(meta)
	if err == nil && a.logFile != nil {
		a.logFile.Write(b)
		a.logFile.Write([]byte("\n"))
	}

	// Save image clip in a per-camera folder
	if len(event.ImageClip) > 0 {
		clipDir := filepath.Join("clips", event.CameraId)
		os.MkdirAll(clipDir, 0755)
		imgName := filepath.Join(clipDir, time.Now().Format("20060102_150405.000")+".jpg")
		err := os.WriteFile(imgName, event.ImageClip, 0644)
		if err != nil {
			log.Printf("StorageActor: failed to save image clip: %v", err)
		}
	}
}

func (a *StorageActor) PostStop(ctx *actor.Context) error {
	if a.logFile != nil {
		a.logFile.Close()
	}
	return nil
}
