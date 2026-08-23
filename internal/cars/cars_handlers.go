package cars

import (
	"context"
	"encoding/json"
	"io"
	"log"
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
//	@Summary		Start car camera stream
//	@Description	Sends a command to the specified car to start the camera stream
//	@Tags			Car
//	@Produce		json
//	@Param			id	path		string	true	"Car ID (UUID)"
//	@Success		200
//	@Failure		400	{string}	string	"Invalid car ID"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/api/cars/{id}/stream [get]
func (h *CarHandler) CarStream(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	carID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid car id", http.StatusBadRequest)
		return
	}

	if err := h.service.RegisterStream(carID, w); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	defer h.service.UnregisterStream(carID)

	w.Header().Set(
		"Content-Type",
		"multipart/x-mixed-replace; boundary=frame",
	)

	w.WriteHeader(http.StatusOK)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	flusher.Flush()

	if err := h.service.CarStream(carID); err != nil {
		return
	}

	<-r.Context().Done()

	if err := h.service.CarStreamStop(carID); err != nil {
		log.Printf(
			"failed to stop stream car=%s: %v",
			carID,
			err,
		)
	}
}

func (h *CarHandler) CarStreamUpload(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	carID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid car id", http.StatusBadRequest)
		return
	}

	frontendWriter, ok := h.service.GetStreamWriter(carID)
	if !ok {
		http.Error(w, "no active stream", http.StatusNotFound)
		return
	}

	flusher, ok := frontendWriter.(http.Flusher)
	if !ok {
		http.Error(
			w,
			"frontend stream does not support flushing",
			http.StatusInternalServerError,
		)
		return
	}

	buf := make([]byte, 4096)

	for {
		n, err := r.Body.Read(buf)

		if n > 0 {
			if _, writeErr := frontendWriter.Write(buf[:n]); writeErr != nil {
				log.Printf(
					"stream write failed car=%s: %v",
					carID,
					writeErr,
				)
				return
			}

			flusher.Flush()
		}

		if err != nil {
			if err != io.EOF {
				log.Printf(
					"stream read failed car=%s: %v",
					carID,
					err,
				)
			}

			break
		}
	}

	log.Printf("stream ended car=%s", carID)

	w.WriteHeader(http.StatusOK)
}

func (h *CarHandler) CarStreamStop(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	carID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid car id", http.StatusBadRequest)
		return
	}

	if err := h.service.CarStreamStop(carID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.service.UnregisterStream(carID)

	w.WriteHeader(http.StatusNoContent)
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
	id := chi.URLParam(r, "id")

	carID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid car id", http.StatusBadRequest)
		return
	}

	if err := h.service.CarActions(carID, r); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
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
