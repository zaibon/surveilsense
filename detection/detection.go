package detection

import (
	_ "embed"
	"image"

	"gocv.io/x/gocv"
)

type Detector interface {
	Detect(img gocv.Mat) []image.Rectangle
	Close()
}

func NewFaceDetector(data string) Detector {
	cascade := gocv.NewCascadeClassifier()
	if !cascade.Load(data) {
		panic("failed to load face detection model")
	}
	return &faceDetector{classifier: cascade}
}

type faceDetector struct {
	classifier gocv.CascadeClassifier
}

func (d *faceDetector) Detect(img gocv.Mat) []image.Rectangle {
	return d.classifier.DetectMultiScale(img)
}

func (d *faceDetector) Close() {
	d.classifier.Close()
}
