package cars

import (
	"context"
	"encoding/json"
	"log"

	"github.com/TazmanS/smartcar-backend/internal/app"
	cars_dto "github.com/TazmanS/smartcar-backend/internal/cars/dto"
	paho "github.com/eclipse/paho.mqtt.golang"
)

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
	log.Printf("Topic: %s", msg.Topic())
	log.Printf("Payload: %s", string(msg.Payload()))

	if msg.Topic() != h.app.Config.MQTTCarSessionSub {
		return
	}

	ctx := context.Background()
	var req cars_dto.CarsSessionRequest

	if err := json.Unmarshal(msg.Payload(), &req); err != nil {
		log.Printf("Invalid JSON: %v", err)
		return
	}

	sessionID := h.handler.CarGetSessionIdHandler(ctx, req)

	_ = sessionID
}
