package cars

import (
	"github.com/TazmanS/smartcar-backend/internal/app"
	"github.com/go-chi/chi/v5"
)

func RegisterCarRoutes(r chi.Router, app *app.App) {
	handler := NewCarHandler(app)

	r.Get("/car-status", handler.GetCarStatus)
	r.Get("/car-stream", handler.CarStream)
	r.Post("/car-actions", handler.CarActions)
}
