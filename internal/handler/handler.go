package handler

import (
	"encoding/json"
	"errors"
	"github.com/Royal17x/subscription-service/internal/model"
	"github.com/Royal17x/subscription-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Handler struct {
	service  service.SubscriptionService
	validate *validator.Validate
	log      *slog.Logger
}

func New(svc service.SubscriptionService, log *slog.Logger) *Handler {
	return &Handler{
		service:  svc,
		validate: validator.New(),
		log:      log,
	}
}

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/subscriptions", func(r chi.Router) {
			r.Post("/", h.Create)
			r.Get("/", h.List)
			r.Get("/total-cost", h.TotalCost)
			r.Get("/{id}", h.GetByID)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})
	})
	return r
}
func writeJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, code int, message string) {
	writeJson(w, code, map[string]string{"error": message})
}

func errorToStatus(err error) int {
	switch {
	case errors.Is(err, model.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, model.ErrInvalidUUID), errors.Is(err, model.ErrInvalidDateRange):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}

}
