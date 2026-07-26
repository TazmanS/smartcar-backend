package cars

import (
	"github.com/go-chi/chi/v5"
)

func RegisterCarRoutes(r chi.Router, handler *CarHandler) {
	r.Get("/car-status", handler.GetCarStatus)
	r.Get("/car-stream", handler.CarStream)
	r.Post("/car-actions", handler.CarActions)
}
