package dto

import "github.com/google/uuid"

type CarsHeartbeatRequest struct {
	SessionID uuid.UUID `json:"session_id"`
}
