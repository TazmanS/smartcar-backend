package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/TazmanS/smartcar-backend/internal/app"
	"github.com/TazmanS/smartcar-backend/internal/services"
)

type CarHandler struct {
	service *services.CarService
}

func NewCarHandler(app *app.App) *CarHandler {
	return &CarHandler{
		service: services.NewCarService(app),
	}
}

// CarStatusHandler godoc
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

func (h *CarHandler) CarStream(w http.ResponseWriter, r *http.Request) {
	if err := h.service.CarStream(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
