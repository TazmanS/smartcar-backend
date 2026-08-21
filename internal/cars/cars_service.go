package cars

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
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
	app     *app.App
	repo    *Repository
	mu      sync.Mutex
	streams map[uuid.UUID]http.ResponseWriter
}

func NewCarService(app *app.App, repo *Repository) *CarService {
	return &CarService{
		app:     app,
		repo:    repo,
		streams: make(map[uuid.UUID]http.ResponseWriter),
	}
}

func (s *CarService) GetCarInfo(ctx context.Context, id uuid.UUID) (models.Car, error) {
	return s.repo.GetCarInfo(ctx, id)
}

func (s *CarService) CarStream(carID uuid.UUID) error {
	topic := fmt.Sprintf("%s/%s", mqtt.MQTTTopicActions, carID)

	payload := fmt.Sprintf(
		`{"action":"start_stream","url":"%s/api/cars/%s/stream"}`,
		s.app.Config.BackendHost,
		carID,
	)

	return s.app.MQTT.Publish(topic, payload)
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

func (s *CarService) GetStreamWriter(carID uuid.UUID) (http.ResponseWriter, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w, ok := s.streams[carID]

	return w, ok
}

func (s *CarService) RegisterStream(carID uuid.UUID, w http.ResponseWriter) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.streams[carID]; exists {
		return fmt.Errorf("stream is already occupied")
	}

	s.streams[carID] = w

	return nil
}

func (s *CarService) UnregisterStream(carID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.streams, carID)
}

func (s *CarService) CarStreamStop(carID uuid.UUID) error {
	topic := fmt.Sprintf("%s/%s", mqtt.MQTTTopicActions, carID)

	return s.app.MQTT.Publish(
		topic,
		`{"action":"stop_stream"}`,
	)
}
