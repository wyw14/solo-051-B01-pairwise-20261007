package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/wyw14/cry-161/internal/model"
	"github.com/wyw14/cry-161/internal/service"
	"net/http"
	"os"
)

func checkpointBook() *service.CheckpointBook {
	dir := os.Getenv("OPTICCOAT_DATA")
	if dir == "" {
		dir = "data"
	}
	return service.NewCheckpointBook(dir)
}

func captureCheckpoint(svc *service.Service, runtime *service.Runtime, book *service.CheckpointBook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runID, ok := model.ParseID(chi.URLParam(r, "runID"))
		if !ok {
			http.Error(w, "invalid run id", http.StatusBadRequest)
			return
		}
		checkpoint, err := book.Capture(svc, runtime, runID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if err := book.Validate(checkpoint); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, checkpoint)
	}
}

func listCheckpoints(svc *service.Service, book *service.CheckpointBook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = svc
		runID, ok := model.ParseID(chi.URLParam(r, "runID"))
		if !ok {
			http.Error(w, "invalid run id", http.StatusBadRequest)
			return
		}
		writeJSON(w, book.List(runID))
	}
}

func restoreSnapshot(svc *service.Service, book *service.CheckpointBook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = svc
		snapshot, err := book.Restore()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, snapshot)
	}
}
