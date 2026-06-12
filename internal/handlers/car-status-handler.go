package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/TazmanS/smartcar-backend/internal/services"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func CarStatusHandler(w http.ResponseWriter, r *http.Request) {
	response := services.GetCarStatusService()

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}
}
