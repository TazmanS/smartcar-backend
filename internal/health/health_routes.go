package health

import (
	"github.com/go-chi/chi/v5"
)

func RegisterHealthRoutes(r chi.Router) {
	r.Get("/health", Get)
}
