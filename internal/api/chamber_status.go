package api

import (
	"github.com/wyw14/cry-161/internal/service"
	"net/http"
)

func chambers(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { writeJSON(w, svc.ChamberSnapshot()) }
}
