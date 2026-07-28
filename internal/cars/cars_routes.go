package cars

import (
	"github.com/go-chi/chi/v5"
)

func RegisterCarsRoutes(r chi.Router, handler *CarHandler) {
	r.Post("/cars/list", handler.GetCarsList)
	r.Get("/cars/{id}/info", handler.GetCarInfo)

	r.Get("/cars/{id}/stream", handler.CarStream)
	r.Post("/car-actions", handler.CarActions)
}
