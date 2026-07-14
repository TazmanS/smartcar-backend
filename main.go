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
	"github.com/TazmanS/smartcar-backend/internal/config"
	"github.com/TazmanS/smartcar-backend/internal/database"
	"github.com/TazmanS/smartcar-backend/internal/mqtt"
	"github.com/TazmanS/smartcar-backend/internal/routes"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

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

	router := routes.NewRouter(app)

	log.Println("Server started on", cfg.PORT)

	log.Fatal(
		http.ListenAndServe(cfg.PORT, router),
	)
}
