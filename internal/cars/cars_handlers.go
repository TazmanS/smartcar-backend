package cars

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/TazmanS/smartcar-backend/internal/app"
	"github.com/TazmanS/smartcar-backend/internal/cars/dto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CarHandler struct {
	service *CarService
}

func NewCarHandler(app *app.App, repo *Repository) *CarHandler {
	return &CarHandler{
		service: NewCarService(app, repo),
	}
}

// GetCarInfo godoc
//
//	@Summary		Get car information
//	@Description	Returns information about a specific car by its ID
//	@Tags			Car
//	@Produce		json
//	@Param			id	path		string	true	"Car ID (UUID)"
//	@Success		200	{object}	models.Car
//	@Failure		400	{string}	string	"Invalid car ID"
//	@Failure		404	{string}	string	"Car not found"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/api/cars/{id}/info [get]
func (h *CarHandler) GetCarInfo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	carID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid car id", http.StatusBadRequest)
		return
	}

	car, err := h.service.GetCarInfo(r.Context(), carID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(car)
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
//	@Param			request	body		dto.CarActionRequest	true	"Car action"
//	@Success		200
//	@Failure		400	{string}	string	"Bad Request"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/api/car-actions [post]
func (h *CarHandler) CarActions(w http.ResponseWriter, r *http.Request) {
	if err := h.service.CarActions(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetCarsList godoc
//
//	@Summary		Get paginated list of cars
//	@Description	Returns a paginated list of registered cars
//	@Tags			Car
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.GetCarsListRequest	true	"Pagination parameters"
//	@Success		200		{object}	dto.GetCarsListResponse
//	@Failure		400		{string}	string	"Bad Request"
//	@Failure		500		{string}	string	"Internal Server Error"
//	@Router			/api/cars/list [post]
func (h *CarHandler) GetCarsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.GetCarsListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cars, err := h.service.GetCarsList(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(cars); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CarHandler) CarGetSessionId(ctx context.Context, payload dto.CarsSessionRequest) uuid.UUID {
	if payload.Key != h.service.app.Config.MQTTCarSessionKey {
		return uuid.Nil
	}

	sessionID, err := h.service.CarGetSessionId(ctx)
	if err != nil {
		return uuid.Nil
	}

	return sessionID
}

func (h *CarHandler) CarHeartbeat(ctx context.Context, req dto.CarsHeartbeatRequest) error {
	return h.service.CarHeartbeat(ctx, req)
}
