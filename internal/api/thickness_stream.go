package api

import (
	"github.com/wyw14/cry-161/internal/model"
	"github.com/wyw14/cry-161/internal/service"
	"net/http"
)

func thickness(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { _ = svc; writeJSON(w, []model.ThicknessSample{}) }
}
