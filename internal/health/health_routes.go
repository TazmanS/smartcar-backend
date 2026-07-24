package health

import (
	"github.com/TazmanS/smartcar-backend/internal/app"
	"github.com/go-chi/chi/v5"
)

func RegisterHealthRoutes(r chi.Router, app *app.App) {
	handler := NewHealthHandler(app)

	r.Get("/health", handler.Get)
}
