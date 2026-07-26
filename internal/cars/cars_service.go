package cars

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/TazmanS/smartcar-backend/internal/app"
	"github.com/TazmanS/smartcar-backend/internal/cars/dto"
	"github.com/TazmanS/smartcar-backend/internal/cars/models"
	"github.com/google/uuid"
)

type CarService struct {
	app  *app.App
	repo *Repository
}

func NewCarService(repo *Repository) *CarService {
	return &CarService{
		repo: repo,
	}
}

func (s *CarService) GetCarStatus() dto.CarsStatusResponse {
	return dto.CarsStatusResponse{
		Status:  "ok",
		Message: "SmartCar backend is running!",
	}
}

func (s *CarService) CarStream(w http.ResponseWriter, r *http.Request) error {
	target, err := url.Parse(s.app.Config.ESP32URL)
	if err != nil {
		return err
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(target)
			pr.Out.URL.Path = "/stream"
		},
		FlushInterval: -1,
	}

	proxy.ServeHTTP(w, r)

	return nil
}

func (s *CarService) CarActions(w http.ResponseWriter, r *http.Request) error {
	var request dto.CarActionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return err
	}

	if !request.Action.IsValid() {
		return errors.New("invalid action")
	}

	return s.app.MQTT.Publish(
		s.app.Config.MQTTActions,
		string(request.Action),
	)
}

func (s *CarService) CarGetSessionIdService(ctx context.Context) (string, error) {
	//mqttClient.Publish("smartcar/session_id", "session_id: 123")
	sessionID := uuid.New().String()

	car := &models.Car{
		ID:   sessionID,
		Name: "Tesla",
	}

	if err := s.repo.Create(ctx, car); err != nil {
		return "", err
	}

	return sessionID, nil
}
