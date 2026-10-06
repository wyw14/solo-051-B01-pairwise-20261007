package api

import (
	"github.com/wyw14/cry-161/internal/service"
	"net/http"
)

func recipes(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { writeJSON(w, svc.RecipesList()) }
}
func recipeEdit(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { _ = svc; w.WriteHeader(http.StatusNoContent) }
}
