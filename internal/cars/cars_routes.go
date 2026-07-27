package cars

import (
	"github.com/go-chi/chi/v5"
)

func RegisterCarsRoutes(r chi.Router, handler *CarHandler) {
	r.Post("/cars/list", handler.GetCarsList)
	r.Get("/cars/{id}/status", handler.GetCarStatus)

	r.Get("/car-stream", handler.CarStream)
	r.Post("/car-actions", handler.CarActions)
}
