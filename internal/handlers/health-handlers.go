package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/TazmanS/smartcar-backend/internal/app"
	"github.com/TazmanS/smartcar-backend/internal/models"
)

type HealthHandler struct {
	app *app.App
}

func NewHealthHandler(app *app.App) *HealthHandler {
	return &HealthHandler{
		app: app,
	}
}

func (h *HealthHandler) Get(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status: "ok",
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
