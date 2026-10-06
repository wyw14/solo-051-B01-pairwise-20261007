package api

import (
	"encoding/json"
	"github.com/wyw14/cry-161/internal/service"
	"net/http"
)

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
func incidents(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { _ = svc; writeJSON(w, []map[string]string{}) }
}
