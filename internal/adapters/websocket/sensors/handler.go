package sensors

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	sensorIDQueryParam = "sensor_id"
	sensorTokenHeader  = "X-Auth-Sensor-Token" // #nosec G101
)

type Config struct {
	AllowedSensorsToToken map[string]string
	HandshakeTimeout      time.Duration
}

type Handler struct {
	cfg      *Config
	hub      *Hub
	logger   *slog.Logger
	upgrader websocket.Upgrader
}

func NewHandler(
	cfg *Config,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		cfg:    cfg,
		hub:    NewHub(logger),
		logger: logger,
		upgrader: websocket.Upgrader{
			HandshakeTimeout: cfg.HandshakeTimeout,
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
	}
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	sensorID := r.URL.Query().Get(sensorIDQueryParam)
	if sensorID == "" {
		http.Error(w, "sensor_id required", http.StatusBadRequest)
		return
	}

	expectedToken, exists := h.cfg.AllowedSensorsToToken[sensorID]
	if !exists {
		http.Error(w, "unknown sensor", http.StatusForbidden)
		return
	}

	providedToken := r.Header.Get(sensorTokenHeader)
	if providedToken != expectedToken {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("WebSocket upgrade failed", slog.Any("error", err))
		return
	}

	// Создаем логгер для этого конкретного соединения с метаданными
	connLogger := h.logger.With(
		"sensor_id", sensorID,
	)

	client := &Client{
		SensorID: sensorID,
		Conn:     conn,
		Logger:   connLogger,
	}

	h.hub.AddClient(client)

	connLogger.Info("Sensor connected", "total_clients", h.hub.GetClientCount())

	go h.listenClient(client)
}

// Слушаем конкретного клиента в горутине.
func (h *Handler) listenClient(client *Client) {
	defer func() {
		// Когда горутина завершается - удаляем клиента из хаба
		h.hub.RemoveClient(client.SensorID)

		err := client.Conn.Close()
		if err != nil {
			client.Logger.Error("Failed to close client connection", "error", err)
			return
		}

		client.Logger.Info("Sensor disconnected", "remaining_clients", h.hub.GetClientCount())
	}()

	for {
		_, messageBytes, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway, websocket.CloseAbnormalClosure,
			) {
				client.Logger.Error("WebSocket unexpectedly closed", "error", err)
			}
			break // Выходим из цикла - клиент отключился
		}

		var messageData MessageData
		if err = json.Unmarshal(messageBytes, &messageData); err != nil {
			client.Logger.Error(
				"Failed to parse JSON",
				"error", err,
				"raw_message", string(messageBytes),
			)
			continue // Пропускаем невалидные сообщения
		}

		client.Logger.Info("Received sensor data", "data", messageData)

		// TODO: Сохранить данные в базу или отправить дальше
	}
}
