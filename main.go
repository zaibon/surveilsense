package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tochemey/goakt/v3/actor"
	aktlog "github.com/tochemey/goakt/v3/log"
	"github.com/zaibon/surveilsense/actors"
	"github.com/zaibon/surveilsense/detection"
	"github.com/zaibon/surveilsense/storage"
	"github.com/zaibon/surveilsense/web"
)

func main() {
	ctx := context.Background()
	logger := aktlog.DefaultLogger

	faceDetector := detection.NewFaceDetector("detection/haarcascade_frontalface_default.xml")
	defer faceDetector.Close()

	// Create the actor system
	actorSystem, err := actor.NewActorSystem(
		"SurveilSenseSystem",
		actor.WithLogger(logger),
		actor.WithActorInitMaxRetries(3),
	)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	if err := actorSystem.Start(ctx); err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	fs, err := storage.NewFilesystemStorage()
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	// Spawn actors
	// Spawn NotificationActor and StorageActor first to get their PIDs
	notificationPID, _ := actorSystem.Spawn(ctx, "NotificationActor", actors.NewNotificationActor())
	storagePID, _ := actorSystem.Spawn(ctx, "StorageActor", actors.NewStorageActor(fs))
	// Spawn FrameProcessorActor with actorSystem, notificationPID, and storagePID
	frameProcessorPID, _ := actorSystem.Spawn(ctx, "FrameProcessorActor", actors.NewFrameProcessorActor(notificationPID, storagePID, faceDetector))
	// Pass actorSystem and frameProcessorPID to CameraFeedActor
	// _, _ = actorSystem.Spawn(ctx, "CameraFeedActor", actors.NewCameraFeedActor(frameProcessorPID))

	server := web.NewServer(actorSystem, frameProcessorPID)
	go server.Start()

	// Wait for interrupt signal to gracefully shutdown
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptSignal

	_ = actorSystem.Stop(ctx)
	os.Exit(0)
}
