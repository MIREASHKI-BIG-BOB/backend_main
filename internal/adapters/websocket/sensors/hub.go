package sensors

import (
	"log/slog"
	"sync"
)

type Hub struct {
	clients map[string]*Client
	logger  *slog.Logger
	mutex   sync.RWMutex
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		logger:  logger,
	}
}

func (h *Hub) AddClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client.SensorID] = client
}

func (h *Hub) RemoveClient(clientID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.clients, clientID)
}

func (h *Hub) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return len(h.clients)
}

func (h *Hub) BroadcastToAll(message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for clientID, client := range h.clients {
		if err := client.Conn.WriteMessage(1, message); err != nil {
			h.logger.Error("Failed to send to client", "client_id", clientID, "error", err)
		}
	}
}
