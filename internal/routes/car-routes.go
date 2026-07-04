package routes

import (
	"github.com/TazmanS/smartcar-backend/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterCarRoutes(r chi.Router) {
	r.Get("/car-status", handlers.CarStatusHandler)
}
