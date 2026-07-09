package main

// @title SmartCar API
// @version 1.0
// @description SmartCar Backend API
// @host localhost:8080
// @BasePath /

import (
	"log"
	"net/http"

	"github.com/TazmanS/smartcar-backend/internal/config"
	"github.com/TazmanS/smartcar-backend/internal/routes"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	cfg := config.Load()

	router := routes.NewRouter()

	log.Println("Server started on", cfg.PORT)

	log.Fatal(
		http.ListenAndServe(cfg.PORT, router),
	)
}
