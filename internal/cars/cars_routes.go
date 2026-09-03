package cars

import (
	"github.com/go-chi/chi/v5"
)

func RegisterCarsRoutes(r chi.Router, handler *CarHandler) {
	r.Post("/cars/list", handler.GetCarsList)
	r.Get("/cars/{id}/info", handler.GetCarInfo)

	r.Get("/cars/{id}/stream", handler.CarStream)
	r.Post("/cars/{id}/stream", handler.CarStreamUpload)
	r.Post("/cars/{id}/stream/stop", handler.CarStreamStop)
	r.Post("/cars/{id}/actions", handler.CarActions)

	r.Get("/cars/{id}/socket", handler.CarSocket)
}
