package main

// @title SmartCar API
// @version 1.0
// @description SmartCar Backend API
// @host localhost:8080
// @BasePath /

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TazmanS/smartcar-backend/internal/app"
	"github.com/TazmanS/smartcar-backend/internal/cars"
	cars_dto "github.com/TazmanS/smartcar-backend/internal/cars/dto"
	"github.com/TazmanS/smartcar-backend/internal/config"
	"github.com/TazmanS/smartcar-backend/internal/database"
	"github.com/TazmanS/smartcar-backend/internal/mqtt"
	"github.com/TazmanS/smartcar-backend/internal/routes"
	paho "github.com/eclipse/paho.mqtt.golang"
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

	err = mqttClient.Subscribe(cfg.MQTTCarSessionSub, func(client paho.Client, msg paho.Message) {
		handleSessionMessage(mqttClient, app, client, msg)
	})

	if err != nil {
		log.Fatal(err)
	}

	router := routes.NewRouter(app)

	log.Println("Server started on", cfg.PORT)

	log.Fatal(
		http.ListenAndServe(cfg.PORT, router),
	)
}

func handleSessionMessage(mqttClient *mqtt.Client, app *app.App, client paho.Client, msg paho.Message) {
	log.Printf("Topic: %s", msg.Topic())
	log.Printf("Payload: %s", string(msg.Payload()))

	if msg.Topic() == app.Config.MQTTCarSessionSub {
		var req cars_dto.CarsSessionRequest

		err := json.Unmarshal(msg.Payload(), &req)
		if err != nil {
			log.Printf("Invalid JSON: %v", err)
			return
		}
		cars.CarGetSessionIdHandler(app, req)
	}
}
