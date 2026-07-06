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

	router := routes.NewRouter(cfg.FrontendURL)

	log.Println("Server started on :8080")

	log.Fatal(
		http.ListenAndServe(":8080", router),
	)
}
