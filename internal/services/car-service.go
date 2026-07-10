package services

import "github.com/TazmanS/smartcar-backend/internal/models"

func GetCarStatusService() models.CarStatusResponse {
	return models.CarStatusResponse{
		Status:  "ok",
		Message: "SmartCar backend is running!",
	}
}
