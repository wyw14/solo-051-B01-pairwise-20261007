package api

import (
	"github.com/wyw14/cry-161/internal/service"
	"net/http"
)

func runtimeStatus(svc *service.Service, runtime *service.Runtime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = svc
		writeJSON(w, runtime.Summary())
	}
}

func auditStatus(svc *service.Service, runtime *service.Runtime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		report := svc.Audit(r.Context(), runtime)
		if !report.Passed {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		writeJSON(w, report)
	}
}

func metricsStatus(svc *service.Service, runtime *service.Runtime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := service.BuildOperationsMetrics(svc, runtime)
		if !metrics.Healthy() {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		writeJSON(w, metrics)
	}
}
