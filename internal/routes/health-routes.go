package routes

import (
	"github.com/TazmanS/smartcar-backend/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterHealthRoutes(r chi.Router) {
	r.Get("/health", handlers.HealthHandler)
}
