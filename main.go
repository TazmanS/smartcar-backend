package main

// @title SmartCar API
// @version 1.0
// @description SmartCar Backend API
// @host localhost:8080
// @BasePath /

import (
	"log"
	"net/http"

	"github.com/TazmanS/smartcar-backend/internal/app"
	"github.com/TazmanS/smartcar-backend/internal/cars"
	"github.com/TazmanS/smartcar-backend/internal/config"
	"github.com/TazmanS/smartcar-backend/internal/database"
	"github.com/TazmanS/smartcar-backend/internal/health"
	"github.com/TazmanS/smartcar-backend/internal/mqtt"
	"github.com/TazmanS/smartcar-backend/internal/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mqttClient, err := mqtt.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer mqttClient.Close()

	app := &app.App{
		DB:     db,
		MQTT:   mqttClient,
		Config: cfg,
	}

	carsRepo := cars.NewRepository(app.DB)
	carsHandler := cars.NewCarHandler(app, carsRepo)
	carsMQTTHandler := cars.NewMQTTHandler(app, carsHandler)

	carService := cars.NewCarService(app, carsRepo)
	go carService.StartCleanupTask()

	if err := mqttClient.Subscribe(mqtt.MQTTSubSessionId, carsMQTTHandler.HandleMQTTSessionMessage); err != nil {
		log.Fatal(err)
	}

	if err := mqttClient.Subscribe(mqtt.MQTTTopicHeartbeat, carsMQTTHandler.HandleMQTTHeartbeat); err != nil {
		log.Fatal(err)
	}

	if err := mqttClient.Subscribe(mqtt.MQTTMessage, carsMQTTHandler.HandleMQTTMessage); err != nil {
		log.Fatal(err)
	}

	router := routes.NewRouter()

	router.Route("/api", func(api chi.Router) {
		health.RegisterHealthRoutes(api)
		cars.RegisterCarsRoutes(api, carsHandler)
	})

	log.Println("Server started on", cfg.PORT)

	log.Fatal(
		http.ListenAndServe(cfg.PORT, router),
	)
}
