package storage

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/zaibon/surveilsense/proto"
)

type FilesystemStorage struct {
	logFile *os.File
}

func NewFilesystemStorage() (*FilesystemStorage, error) {
	f, err := os.OpenFile("detections.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &FilesystemStorage{logFile: f}, nil
}

func (fs *FilesystemStorage) SaveMetadata(ctx context.Context, cameraID string, timestamp time.Time, detections []*proto.Detection) error {
	meta := map[string]interface{}{
		"camera_id":  cameraID,
		"timestamp":  timestamp.UnixMilli(),
		"detections": detections,
	}
	b, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	_, err = fs.logFile.Write(b)
	if err == nil {
		_, err = fs.logFile.Write([]byte("\n"))
	}
	return err
}

func (fs *FilesystemStorage) SaveClip(ctx context.Context, cameraID string, timestamp time.Time, imageClip []byte) error {
	if len(imageClip) == 0 {
		return nil
	}
	clipDir := filepath.Join("clips", cameraID)
	os.MkdirAll(clipDir, 0755)
	imgName := filepath.Join(clipDir, timestamp.Format("20060102_150405.000")+".jpg")
	return os.WriteFile(imgName, imageClip, 0644)
}

func (fs *FilesystemStorage) Close() error {
	if fs.logFile != nil {
		return fs.logFile.Close()
	}
	return nil
}
