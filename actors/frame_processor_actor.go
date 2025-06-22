package actors

import (
	"image/color"
	"log"
	"time"

	"gocv.io/x/gocv"

	"github.com/tochemey/goakt/v3/actor"
	"github.com/zaibon/surveilsense/detection"
	"github.com/zaibon/surveilsense/proto"
)

var blue = color.RGBA{0, 0, 255, 0}

// FrameProcessorActor receives FrameData and sends DetectionEvent
type FrameProcessorActor struct {
	detector detection.Detector
}

var _ actor.Actor = (*FrameProcessorActor)(nil)

// NewFrameProcessorActor creates a FrameProcessorActor with the required PIDs
func NewFrameProcessorActor(detector detection.Detector) *FrameProcessorActor {
	return &FrameProcessorActor{
		detector: detector,
	}
}

func (a *FrameProcessorActor) PreStart(ctx *actor.Context) error {
	return nil
}

func (a *FrameProcessorActor) Receive(ctx *actor.ReceiveContext) {
	msg := ctx.Message()
	frame, ok := msg.(*proto.FrameData)
	if !ok {
		ctx.Unhandled()
		return
	}

	// Decode JPEG image
	imgMat, err := gocv.IMDecode(frame.ImageData, gocv.IMReadColor)
	if err != nil || imgMat.Empty() {
		log.Printf("FrameProcessorActor: failed to decode image: %v", err)
		return
	}
	defer imgMat.Close()

	recs := a.detector.Detect(imgMat)
	if len(recs) == 0 {
		log.Printf("FrameProcessorActor: no objects detected in frame from camera %s", frame.CameraId)
		return
	}

	log.Printf("FrameProcessorActor: detected %d objects in frame from camera %s", len(recs), frame.CameraId)

	detections := make([]*proto.Detection, 0, len(recs))
	for _, rec := range recs {
		// draw a rectangle around each face on the original image
		gocv.Rectangle(&imgMat, rec, blue, 3)

		detections = append(detections, &proto.Detection{
			Confidence: 0.9, //TODO: set actual confidence if available
			X:          int32(rec.Min.X),
			Y:          int32(rec.Min.Y),
			Width:      int32(rec.Dx()),
			Height:     int32(rec.Dy()),
		})
	}

	// Optionally crop the detection region (for demo, send full frame)
	buf, err := gocv.IMEncode(gocv.JPEGFileExt, imgMat)
	var imageClip []byte
	if err == nil {
		imageClip = buf.GetBytes()
		buf.Close()
	}

	detectionEvent := &proto.DetectionEvent{
		CameraId:   frame.CameraId,
		Timestamp:  time.Now().UnixMilli(),
		Detections: detections,
		ImageClip:  imageClip,
	}

	// Send DetectionEvent to all NotificationActor and StorageActor instances
	a.sendDetectionEvent(ctx, detectionEvent)
}

func (a *FrameProcessorActor) PostStop(ctx *actor.Context) error {
	return nil
}

func (a *FrameProcessorActor) sendDetectionEvent(ctx *actor.ReceiveContext, event *proto.DetectionEvent) {
	pids := ctx.ActorSystem().Actors()
	for _, pid := range pids {
		switch pid.Actor().(type) {
		case *NotificationActor, *StorageActor:
			if err := actor.Tell(ctx.Context(), pid, event); err != nil {
				log.Printf("FrameProcessorActor: failed to send detection event to %s: %v", pid.Address(), err)
			} else {
				log.Printf("FrameProcessorActor: sent detection event to %s for camera %s", pid.Address(), event.CameraId)
			}
		default:
			continue
		}
	}
}
