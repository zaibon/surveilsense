# SurveilSense

A modular, actor-based smart surveillance system in Go. SurveilSense captures video from multiple cameras, processes frames for detections, stores clips, and provides a modern web UI for management and browsing—all powered by the actor model (goakt) for scalability and resilience.

---

## Features
- **Actor Model**: Built with [goakt](https://github.com/tochemey/goakt) for robust, concurrent camera and processing management.
- **Multi-Camera Support**: Dynamically add/remove camera feeds via the web UI or REST API.
- **Frame Processing**: Real-time frame analysis (face/human detection, pluggable).
- **Notifications**: Actor-based notification pipeline (extensible).
- **Clip Storage**: Per-camera clip storage, organized and browsable.
- **Web UI**: Modern, responsive UI with [TailwindCSS](https://tailwindcss.com/) and [htmx](https://htmx.org/) for live updates.
- **REST API**: Manage cameras, browse clips, and fetch live frames programmatically.
- **Extensible**: Add new actors for analytics, notifications, or storage backends.

---

## Architecture

See [`flowchart.mmd`](flowchart.mmd) for a Mermaid diagram of actor communication.

```
flowchart TD
    subgraph Camera
        CameraFeed[CameraFeedActor]
        FrameProcessor[FrameProcessorActor]
    end

    Notification[NotificationActor]
    Storage[StorageActor]

    %% Relationships
    CameraFeed -- FrameData --> FrameProcessor
    FrameProcessor -- DetectionEvent --> Notification
    FrameProcessor -- DetectionEvent --> Storage
```

---

## Getting Started

### Prerequisites
- Go 1.20+
- [gocv](https://gocv.io/getting-started/) (OpenCV bindings for Go)
- [protoc](https://grpc.io/docs/protoc-installation/) (for generating protobuf code)
- A webcam or RTSP camera (for live feeds)

### Clone & Build
```sh
git clone https://github.com/zaibon/surveilsense.git
cd surveilsense
# Install Go dependencies
go mod tidy
# Generate protobuf code
go generate ./...
# Build
go build -o surveilsense main.go
```

### Run
```sh
./surveilsense
```
- The web UI will be available at [http://localhost:8080](http://localhost:8080)

---

## Usage

### Web UI
- **Add Camera**: Enter a camera ID and device ID (e.g., 0 for default webcam) and click "Add Camera".
- **Remove Camera**: Click "Remove" next to a camera.
- **Live Feeds**: View the latest frame from each active camera.
- **Browse Clips**: Click "View Clips" to see recorded clips, organized by camera.

### REST API
- `GET /api/cameras` — List cameras (HTML for htmx)
- `POST /api/cameras` — Add a camera (form data: `camera_id`, `device_id`)
- `DELETE /api/cameras/{id}` — Remove a camera
- `GET /api/cameras/frames` — Get HTML for all live camera frames
- `GET /api/clips` — List all recorded clips (HTML for htmx)

---

## Development

- **Actors**: See the `actors/` directory for all actor implementations.
- **Protobuf**: Messages defined in `proto/messages.proto`.
- **Web**: UI and server logic in `web/`.
- **Flowchart**: Update `flowchart.mmd` for architecture diagrams.

### Testing
```sh
go test ./...
```

---

## Contributing

Contributions are welcome! Please open issues or pull requests for bug fixes, features, or improvements.

- Fork the repo
- Create a feature branch
- Commit your changes
- Open a pull request

---

## Acknowledgements
- [goakt](https://github.com/tochemey/goakt) — Actor model for Go
- [gocv](https://gocv.io/) — Go bindings for OpenCV
- [TailwindCSS](https://tailwindcss.com/)
- [htmx](https://htmx.org/)

---

## Contact

For questions or support, open an issue or contact the maintainer at [github.com/zaibon](https://github.com/zaibon).
