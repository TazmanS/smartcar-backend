package cars

import (
	"encoding/json"
	"net/http"

	"github.com/TazmanS/smartcar-backend/internal/app"
	"github.com/TazmanS/smartcar-backend/internal/cars/dto"
)

type CarHandler struct {
	service *CarService
}

func NewCarHandler(app *app.App) *CarHandler {
	return &CarHandler{
		service: NewCarService(app),
	}
}

// GetCarStatus godoc
//
//	@Summary		Get car status
//	@Description	Returns current car status
//	@Tags			Car
//	@Produce		json
//	@Success		200	{object}	models.CarStatusResponse
//	@Failure		500	{object}	models.CarStatusResponse
//	@Router			/api/car-status [get]
func (h *CarHandler) GetCarStatus(w http.ResponseWriter, r *http.Request) {
	response := h.service.GetCarStatus()

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}
}

// CarStream godoc
//
//	@Summary		Get car camera stream
//	@Description	Returns current car stream
//	@Tags			Car
//	@Produce		json
//	@Success		200
//	@Failure		500
//	@Router			/api/car-stream [get]
func (h *CarHandler) CarStream(w http.ResponseWriter, r *http.Request) {
	if err := h.service.CarStream(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CarActions godoc
//
//	@Summary		Send car action
//	@Description	Send action to the car via MQTT
//	@Tags			Car
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CarActionRequest	true	"Car action"
//	@Success		200
//	@Failure		400	{string}	string	"Bad Request"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/api/car-actions [post]
func (h *CarHandler) CarActions(w http.ResponseWriter, r *http.Request) {
	if err := h.service.CarActions(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CarGetSessionIdHandler(app *app.App, payload dto.CarsSessionRequest) {
	if payload.Key != app.Config.MQTTCarSessionKey {
		// save to logs the information
		return
	}

	CarGetSessionIdService()

}
