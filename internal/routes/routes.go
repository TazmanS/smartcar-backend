package routes

import (
	"github.com/TazmanS/smartcar-backend/internal/handlers"
	"github.com/TazmanS/smartcar-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(frontendURL string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Cors(frontendURL))

	r.Route("/api", func(r chi.Router) {
		r.Get("/car-status", handlers.CarStatusHandler)
	})

	return r
}
