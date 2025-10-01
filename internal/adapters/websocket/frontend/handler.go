package frontend

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Handler struct {
	logger   *slog.Logger
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	mu       sync.RWMutex
	broadcast chan []byte
}

func NewHandler(logger *slog.Logger) *Handler {
	h := &Handler{
		logger: logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte, 256),
	}

	go h.handleBroadcast()

	return h
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("WebSocket upgrade failed", "error", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	clientCount := len(h.clients)
	h.mu.Unlock()

	h.logger.Info("Frontend client connected", "total_clients", clientCount)

	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		remainingClients := len(h.clients)
		h.mu.Unlock()

		conn.Close()
		h.logger.Info("Frontend client disconnected", "remaining_clients", remainingClients)
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error("WebSocket unexpectedly closed", "error", err)
			}
			break
		}
	}
}

func (h *Handler) handleBroadcast() {
	for message := range h.broadcast {
		h.mu.RLock()
		for client := range h.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				h.logger.Error("Failed to send message to frontend client", "error", err)
				client.Close()
				h.mu.Lock()
				delete(h.clients, client)
				h.mu.Unlock()
			}
		}
		h.mu.RUnlock()
	}
}

func (h *Handler) BroadcastToFrontend(message []byte) {
	select {
	case h.broadcast <- message:
		h.logger.Debug("Message queued for broadcast to frontend")
	default:
		h.logger.Warn("Broadcast channel full, dropping message")
	}
}

func (h *Handler) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

