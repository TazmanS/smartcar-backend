package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/TazmanS/smartcar-backend/internal/app"
	"github.com/TazmanS/smartcar-backend/internal/models"
)

type CarService struct {
	app *app.App
}

func NewCarService(app *app.App) *CarService {
	return &CarService{
		app: app,
	}
}

func (s *CarService) GetCarStatus() models.CarStatusResponse {
	return models.CarStatusResponse{
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
	var request models.CarActionRequest

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
