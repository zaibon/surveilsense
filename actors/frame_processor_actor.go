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
	notificationPID *actor.PID
	storagePID      *actor.PID
	detector        detection.Detector
}

var _ actor.Actor = (*FrameProcessorActor)(nil)

// NewFrameProcessorActor creates a FrameProcessorActor with the required PIDs
func NewFrameProcessorActor(notificationPID, storagePID *actor.PID, detector detection.Detector) *FrameProcessorActor {
	return &FrameProcessorActor{
		notificationPID: notificationPID,
		storagePID:      storagePID,
		detector:        detector,
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

	// Send DetectionEvent to NotificationActor and StorageActor

	if a.notificationPID != nil {
		if err := actor.Tell(ctx.Context(), a.notificationPID, detectionEvent); err != nil {
			log.Printf("FrameProcessorActor: failed to send detection event to notification actor: %v", err)
		} else {
			log.Printf("FrameProcessorActor: sent detection event to notification actor for camera %s", frame.CameraId)
		}
	}
	if a.storagePID != nil {
		if err := actor.Tell(ctx.Context(), a.storagePID, detectionEvent); err != nil {
			log.Printf("FrameProcessorActor: failed to send detection event to storage actor: %v", err)
		} else {
			log.Printf("FrameProcessorActor: sent detection event to storage actor for camera %s", frame.CameraId)
		}
	}
}

func (a *FrameProcessorActor) PostStop(ctx *actor.Context) error {
	return nil
}
