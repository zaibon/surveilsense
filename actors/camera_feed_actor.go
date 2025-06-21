package actors

import (
	"context"
	"image"
	"log"
	"time"

	"gocv.io/x/gocv"

	"github.com/tochemey/goakt/v3/actor"
	"github.com/zaibon/surveilsense/proto"
)

const frameRate = time.Second // ~1 FPS

// CameraFeedActor captures frames and sends FrameData messages
type CameraFeedActor struct {
	cameraID  string
	deviceID  int
	capture   *gocv.VideoCapture
	imgMat    gocv.Mat
	quit      chan struct{}
	processor *actor.PID
}

var _ actor.Actor = (*CameraFeedActor)(nil)

// NewCameraFeedActor creates a CameraFeedActor with required dependencies
func NewCameraFeedActor(processor *actor.PID) *CameraFeedActor {
	return &CameraFeedActor{
		cameraID:  "cam-1",
		deviceID:  0,
		quit:      make(chan struct{}),
		processor: processor,
	}
}

// NewCameraFeedActorWithConfig creates a CameraFeedActor for a specific camera ID and device ID
func NewCameraFeedActorWithConfig(cameraID string, deviceID int, processor *actor.PID) *CameraFeedActor {
	return &CameraFeedActor{
		cameraID:  cameraID,
		deviceID:  deviceID,
		quit:      make(chan struct{}),
		processor: processor,
	}
}

func (a *CameraFeedActor) PreStart(ctx *actor.Context) error {
	var err error
	a.capture, err = gocv.OpenVideoCapture(a.deviceID)
	if err != nil {
		return err
	}
	if !a.capture.IsOpened() {
		return err
	}
	a.imgMat = gocv.NewMat()
	go a.captureLoop()

	return nil
}

func (a *CameraFeedActor) captureLoop() {
	for {
		select {
		case <-a.quit:
			return
		default:
			if ok := a.capture.Read(&a.imgMat); !ok || a.imgMat.Empty() {
				time.Sleep(frameRate) // Wait before retrying
				continue
			}
			// Resize to standard resolution
			gocv.Resize(a.imgMat, &a.imgMat, image.Pt(640, 480), 0, 0, gocv.InterpolationDefault)
			// Encode as JPEG
			buf, err := gocv.IMEncode(gocv.JPEGFileExt, a.imgMat)
			if err != nil {
				log.Printf("failed to encode frame: %v", err)
				continue
			}
			frame := &proto.FrameData{
				CameraId:  a.cameraID,
				Timestamp: time.Now().UnixMilli(),
				ImageData: buf.GetBytes(),
			}
			// Send to FrameProcessorActor
			if a.processor != nil {
				if err := actor.Tell(context.Background(), a.processor, frame); err != nil {
					log.Printf("failed to send frame to processor: %v", err)
				} else {
					log.Printf("sent frame to processor: %s", a.cameraID)
				}
			}
			buf.Close()
			time.Sleep(frameRate)
		}
	}
}

func (a *CameraFeedActor) Receive(ctx *actor.ReceiveContext) {
	// No message handling for now
	ctx.Unhandled()
}

func (a *CameraFeedActor) PostStop(ctx *actor.Context) error {
	close(a.quit)
	if a.capture != nil {
		a.capture.Close()
	}
	a.imgMat.Close()

	return nil
}
