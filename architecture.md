# Smart Surveillance System: High-Level Actor Model Architecture

The core idea is to decompose the surveillance system into independent, communicating actors, each responsible for a specific part of the process. This promotes concurrency, resilience, and scalability.

## Key Actors and Their Responsibilities:

1.CameraFeedActor (or CameraSourceActor)
- Purpose: Represents a single camera feed (real or simulated).
- Responsibilities:
  - Continuously capture frames from its associated webcam.
  - Apply basic preprocessing (e.g., resizing to a standard resolution).f
  - Timestamp each frame.
  - Send processed frames as messages to FrameProcessorActor instances.
  - Messages Sent: FrameData(timestamp, image_data, camera_id)
2. FrameProcessorActor (or HumanDetectionActor)
- Purpose: Receives frames and performs human detection.
- Responsibilities:
  - Receive FrameData messages from CameraFeedActors.
  - Apply the human detection algorithm (e.g., using a pre-trained model like a simplified YOLO or HOG + SVM for - demonstration).
  -   If humans are detected, generate a DetectionEvent containing relevant details (bounding box coordinates, confidence - score, timestamp, camera ID).
  - Send DetectionEvent messages to NotificationActor and StorageActor.
  - Messages Received: FrameData
  - Messages Sent: DetectionEvent(timestamp, camera_id, detections, optional_image_clip)
3. NotificationActor
- Purpose: Handles the generation and dispatch of alerts based on detected events.
- Responsibilities:
  - Receive DetectionEvent messages from FrameProcessorActors.
  - Format alert messages (e.g., "Human detected at Camera X on YYYY-MM-DD HH:MM:SS").
      Dispatch notifications (e.g., print to console for a basic demo, or integrate with email/push notification services for       - a more advanced system).
  - Could also implement debounce logic to avoid too many notifications from a continuous detection.
  - Messages Received: DetectionEvent
4. StorageActor
- Purpose: Persists relevant data from detection events.
- Responsibilities:
  - Receive DetectionEvent messages from FrameProcessorActors.
  - Save detection metadata (timestamp, camera ID, detection details) to a database or log file.
  - Optionally, save small image clips or short video segments associated with the detection for review.
  - Messages Received: DetectionEvent
5. SupervisorActor (Implicit/Common in Actor Frameworks)
  - Purpose: Manages the lifecycle of child actors and handles failures.
   - Responsibilities:
     - Monitors child actors (CameraFeedActor, FrameProcessorActor, etc.).
     - Implements a supervision strategy (e.g., restart a FrameProcessorActor if it crashes, or stop a CameraFeedActor if its - - connection drops permanently).
6. SystemCoordinatorActor (Optional, for higher-level orchestration)
- Purpose: Manages the overall system state, registers cameras, and initializes other actors.
- Responsibilities:
  - Handles requests to add/remove cameras.
  - Spawns CameraFeedActors and their associated FrameProcessorActors.
  - Can provide an interface for a dashboard or user control.

How the Actor Model Helps:
- Concurrency: Multiple CameraFeedActors and FrameProcessorActors can run in parallel, allowing the system to handle many cameras simultaneously without blocking.
- Isolation & Fault Tolerance: Each actor encapsulates its own state. If a FrameProcessorActor fails (e.g., due to a bad frame or algorithm error), it can be restarted by its supervisor without affecting other cameras or the entire system.
- Asynchronous Communication: All interactions are message-driven. CameraFeedActors don't wait for FrameProcessorActors to finish; they just send frames and move on. This keeps the system responsive.
- Scalability: You can easily add more FrameProcessorActors to distribute the detection workload, or add more CameraFeedActors as you connect more cameras.
- Modularity: Each actor has a clear, single responsibility, making the system easier to understand, develop, and maintain.
- 
This architecture provides a robust foundation for building your smart surveillance system!

## Technology

Programming language: Go
Actor library: github.com/Tochemey/goakt