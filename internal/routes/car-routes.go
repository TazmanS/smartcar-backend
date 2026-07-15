package routes

import (
	"github.com/TazmanS/smartcar-backend/internal/app"
	"github.com/TazmanS/smartcar-backend/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterCarRoutes(r chi.Router, app *app.App) {
	handler := handlers.NewCarHandler(app)

	r.Get("/car-status", handler.GetCarStatus)
	r.Get("/car-stream", handler.CarStream)
	r.Post("/car-actions", handler.CarActions)
}
