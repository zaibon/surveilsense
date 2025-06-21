package web

import (
	"context"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/tochemey/goakt/v3/actor"
	"github.com/zaibon/surveilsense/actors"
)

type Camera struct {
	CameraID string     `json:"camera_id"`
	DeviceID int        `json:"device_id"`
	PID      *actor.PID `json:"-"`
}

type Server struct {
	mux          *http.ServeMux
	actorSystem  actor.ActorSystem
	frameProcPID *actor.PID
	cameras      map[string]Camera // Track CameraFeedActor PIDs
}

func NewServer(actorSystem actor.ActorSystem, frameProcPID *actor.PID) *Server {
	mux := http.NewServeMux()
	server := &Server{mux: mux, actorSystem: actorSystem, frameProcPID: frameProcPID, cameras: make(map[string]Camera)}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})
	mux.HandleFunc("/clips", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/clips.html")
	})
	mux.Handle("/clips/", http.StripPrefix("/clips/", http.FileServer(http.Dir("clips"))))

	mux.HandleFunc("/api/cameras", server.camerasHandler)
	mux.HandleFunc("/api/cameras/", server.cameraHandler)
	mux.HandleFunc("/api/clips", clipsHandler)

	return server
}

func (s *Server) Start() {
	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", s.mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Server) camerasHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := make([]Camera, 0, len(s.cameras))
		for _, cam := range s.cameras {
			list = append(list, cam)
		}
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var cam Camera
		if err := json.NewDecoder(r.Body).Decode(&cam); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Spawn a new CameraFeedActor for this camera
		pid, err := s.actorSystem.Spawn(r.Context(), cam.CameraID, actors.NewCameraFeedActorWithConfig(cam.CameraID, cam.DeviceID, s.frameProcPID))
		if err != nil {
			log.Printf("Failed to spawn CameraFeedActor: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cam.PID = pid // Track actor PID associated with this camera
		s.cameras[cam.CameraID] = cam

		log.Printf("Spawned CameraFeedActor %s (device %d) with PID %v", cam.CameraID, cam.DeviceID, pid)
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Server) cameraHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	if r.Method == http.MethodDelete {
		// Stop the actor for this camera if it exists
		camera, ok := s.cameras[id]
		if ok && camera.PID != nil {
			err := camera.PID.Shutdown(context.Background())
			if err != nil {
				log.Printf("Failed to stop CameraFeedActor %s: %v", id, err)
			}
			delete(s.cameras, id)
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

type clip struct {
	Filename string `json:"filename"`
	CameraID string `json:"camera_id"`
}

func clipsHandler(w http.ResponseWriter, r *http.Request) {
	files := []clip{}
	_ = filepath.WalkDir("clips", func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && (filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".jpeg") {
			dir := filepath.Base(filepath.Dir(path))
			files = append(files, clip{CameraID: dir, Filename: filepath.Join(dir, filepath.Base(path))})
		}
		return nil
	})
	json.NewEncoder(w).Encode(files)
}
