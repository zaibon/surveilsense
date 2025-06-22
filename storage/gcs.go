package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"time"

	"cloud.google.com/go/storage"

	"github.com/zaibon/surveilsense/proto"
)

type GCSStorage struct {
	bucketName string
	client     *storage.Client
}

func NewGCSStorage(ctx context.Context, bucketName string) (*GCSStorage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GCSStorage{bucketName: bucketName, client: client}, nil
}

func gcsObjectPath(base, cameraID string, timestamp int64, ext string) string {
	t := time.UnixMilli(timestamp)
	return path.Join(
		base,
		cameraID,
		t.Format("2006"),
		t.Format("01"),
		t.Format("02"),
		t.Format("15"),
		t.Format("04"),
		fmt.Sprintf("%d.%s", timestamp, ext),
	)
}

func (g *GCSStorage) SaveMetadata(ctx context.Context, cameraID string, timestamp time.Time, detections []*proto.Detection) error {
	meta := map[string]interface{}{
		"camera_id":  cameraID,
		"timestamp":  timestamp,
		"detections": detections,
	}
	b, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	objectPath := gcsObjectPath("metadata", cameraID, timestamp.UnixMilli(), "json")
	w := g.client.Bucket(g.bucketName).Object(objectPath).NewWriter(ctx)
	w.ContentType = "application/json"
	if _, err := w.Write(b); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

func (g *GCSStorage) SaveClip(ctx context.Context, cameraID string, timestamp time.Time, imageClip []byte) error {
	if len(imageClip) == 0 {
		return nil
	}
	objectPath := gcsObjectPath("clips", cameraID, timestamp.UnixMilli(), "jpg")
	w := g.client.Bucket(g.bucketName).Object(objectPath).NewWriter(ctx)
	w.ContentType = "image/jpeg"
	if _, err := w.Write(imageClip); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

func (g *GCSStorage) Close() error {
	return g.client.Close()
}
