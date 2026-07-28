package cars

import (
	"context"
	"encoding/json"
	"log"

	"github.com/TazmanS/smartcar-backend/internal/app"
	cars_dto "github.com/TazmanS/smartcar-backend/internal/cars/dto"
	"github.com/TazmanS/smartcar-backend/internal/logger"
	"github.com/TazmanS/smartcar-backend/internal/mqtt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type SessionIDResponse struct {
	SessionID uuid.UUID `json:"session_id"`
	RandomID  uint32    `json:"random_id"`
}

type MQTTHandler struct {
	app     *app.App
	handler *CarHandler
}

func NewMQTTHandler(app *app.App, handler *CarHandler) *MQTTHandler {
	return &MQTTHandler{
		app:     app,
		handler: handler,
	}
}

func (h *MQTTHandler) HandleMQTTSessionMessage(client paho.Client, msg paho.Message) {
	logger.Info("Session received",
		"topic", msg.Topic(),
		"payload", string(msg.Payload()))

	ctx := context.Background()
	var req cars_dto.CarsSessionRequest

	if err := json.Unmarshal(msg.Payload(), &req); err != nil {
		log.Printf("Invalid JSON: %v", err)
		return
	}

	sessionID := h.handler.CarGetSessionId(ctx, req)

	msgPublish, _ := json.Marshal(SessionIDResponse{
		SessionID: sessionID,
		RandomID:  req.RandomID,
	})
	h.app.MQTT.Publish(mqtt.MQTTTopicSession, string(msgPublish))
}

func (h *MQTTHandler) HandleMQTTHeartbeat(client paho.Client, msg paho.Message) {
	logger.Info("Heartbeat received",
		"topic", msg.Topic(),
		"payload", string(msg.Payload()))

	ctx := context.Background()
	var req cars_dto.CarsHeartbeatRequest

	if err := json.Unmarshal(msg.Payload(), &req); err != nil {
		log.Printf("Invalid JSON: %v", err)
		return
	}

	if err := h.handler.CarHeartbeat(ctx, req); err != nil {
		log.Printf("Failed to process heartbeat: %v", err)
	}
}
