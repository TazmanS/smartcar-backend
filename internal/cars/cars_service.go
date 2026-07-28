package cars

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/TazmanS/smartcar-backend/internal/app"
	"github.com/TazmanS/smartcar-backend/internal/cars/dto"
	"github.com/TazmanS/smartcar-backend/internal/cars/models"
	"github.com/TazmanS/smartcar-backend/internal/mqtt"
	"github.com/google/uuid"
)

var carNames = []string{
	"Tesla",
	"Mustang",
	"Camaro",
	"Corvette",
	"Supra",
	"Skyline",
	"Civic",
	"Impreza",
	"Eclipse",
	"RX-7",
	"Charger",
	"Challenger",
	"Viper",
	"GT-R",
	"Focus RS",
	"Lancer",
	"Miata",
	"Silvia",
	"NSX",
	"911",
}

type CarService struct {
	app  *app.App
	repo *Repository
}

func NewCarService(app *app.App, repo *Repository) *CarService {
	return &CarService{
		app:  app,
		repo: repo,
	}
}

func (s *CarService) GetCarInfo(ctx context.Context, id uuid.UUID) (models.Car, error) {
	return s.repo.GetCarInfo(ctx, id)
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
		mqtt.MQTTTopicActions,
		string(request.Action),
	)
}

func (s *CarService) CarGetSessionId(ctx context.Context) (uuid.UUID, error) {
	sessionID := uuid.New()

	car := &models.Car{
		ID:   sessionID,
		Name: carNames[rand.Intn(len(carNames))],
	}

	if err := s.repo.Create(ctx, car); err != nil {
		return uuid.Nil, err
	}

	return sessionID, nil
}

func (s *CarService) CarHeartbeat(ctx context.Context, req dto.CarsHeartbeatRequest) error {
	return s.repo.CarHeartbeat(ctx, req)
}

func (s *CarService) GetCarsList(ctx context.Context, req *dto.GetCarsListRequest) (*dto.GetCarsListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PerPage <= 0 {
		req.PerPage = 20
	}

	if req.PerPage > 100 {
		req.PerPage = 100
	}

	cars, totalItems, err := s.repo.GetCarsList(ctx, req)
	if err != nil {
		return nil, err
	}

	totalPages := (totalItems + req.PerPage - 1) / req.PerPage

	return &dto.GetCarsListResponse{
		Data: cars,
		Page: dto.PageInfo{
			Page:       req.Page,
			PerPage:    req.PerPage,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *CarService) StartCleanupTask() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := s.repo.DeleteInactiveCars(context.Background()); err != nil {
			log.Printf("cleanup error: %v", err)
		}
	}
}
