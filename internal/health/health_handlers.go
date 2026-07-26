package health

import (
	"encoding/json"
	"net/http"

	"github.com/TazmanS/smartcar-backend/internal/health/dto"
)

func Get(w http.ResponseWriter, r *http.Request) {
	response := dto.HealthResponse{
		Status: "ok",
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
