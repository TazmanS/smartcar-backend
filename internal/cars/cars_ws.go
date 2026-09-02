package cars

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *CarHandler) CarSocket(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Printf("New connection %s", id)

	carID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid car id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.service.mu.Lock()
	h.service.sockets[carID] = conn
	h.service.mu.Unlock()
	defer func() {
		h.service.mu.Lock()
		delete(h.service.sockets, carID)
		h.service.mu.Unlock()

		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}

}
