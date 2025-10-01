package sensors

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/domain/entities"
	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/infrastructure/ports/repository"
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
	cfg             *Config
	hub             *Hub
	logger          *slog.Logger
	upgrader        websocket.Upgrader
	examRepo        repository.ExamRepository
	mlWSClient      *websocket.Conn
	mlWSMutex       sync.Mutex
	frontendHandler FrontendBroadcaster
	mlAddr          string
	mlPort          string
}

type FrontendBroadcaster interface {
	BroadcastToFrontend(message []byte)
}

func NewHandler(
	cfg *Config,
	logger *slog.Logger,
	examRepo repository.ExamRepository,
	frontendHandler FrontendBroadcaster,
	mlAddr string,
	mlPort string,
) *Handler {
	h := &Handler{
		cfg:             cfg,
		hub:             NewHub(logger),
		logger:          logger,
		examRepo:        examRepo,
		frontendHandler: frontendHandler,
		mlAddr:          mlAddr,
		mlPort:          mlPort,
		upgrader: websocket.Upgrader{
			HandshakeTimeout: cfg.HandshakeTimeout,
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
	}

	go h.connectToMLService()

	return h
}

func (h *Handler) connectToMLService() {
	mlURL := fmt.Sprintf("ws://%s:%s/ws/ctg", h.mlAddr, h.mlPort)

	for {
		conn, _, err := websocket.DefaultDialer.Dial(mlURL, nil)
		if err != nil {
			h.logger.Error("Failed to connect to ML service", "error", err)
			time.Sleep(5 * time.Second)
			continue
		}

		h.mlWSMutex.Lock()
		h.mlWSClient = conn
		h.mlWSMutex.Unlock()

		h.logger.Info("Connected to ML service", "url", mlURL)

		h.listenMLService(conn)

		h.logger.Warn("ML service connection lost, reconnecting...")
		time.Sleep(5 * time.Second)
	}
}

func (h *Handler) listenMLService(conn *websocket.Conn) {
	defer func() {
		h.mlWSMutex.Lock()
		h.mlWSClient = nil
		h.mlWSMutex.Unlock()
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			h.logger.Error("Error reading from ML service", "error", err)
			return
		}

		h.logger.Info("Received response from ML service", "response", string(message))

		if h.frontendHandler != nil {
			h.frontendHandler.BroadcastToFrontend(message)
			h.logger.Debug("ML response forwarded to frontend clients")
		}
	}
}

func (h *Handler) sendToMLService(data entities.CTGData) {
	h.mlWSMutex.Lock()
	defer h.mlWSMutex.Unlock()

	if h.mlWSClient == nil {
		h.logger.Warn("ML service not connected, skipping")
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("Failed to marshal CTG data", "error", err)
		return
	}

	if err := h.mlWSClient.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		h.logger.Error("Failed to send data to ML service", "error", err)
		return
	}

	h.logger.Debug("Sent data to ML service", "sensor_id", data.SensorID)
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

	// Проверяем нужно ли создать новое обследование
	ctx := context.Background()
	needsNew, err := h.examRepo.NeedsNewExamination(ctx)
	if err != nil {
		h.logger.Error("Failed to check if new examination is needed", slog.Any("error", err))
		conn.Close()
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Если нужно - создаем новое обследование
	if needsNew {
		if err := h.examRepo.CreateExamination(ctx); err != nil {
			h.logger.Error("Failed to create examination", slog.Any("error", err))
			conn.Close()
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		h.logger.Info("Created new examination")
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

func (h *Handler) listenClient(client *Client) {
	defer func() {
		h.hub.RemoveClient(client.SensorID)

		err := client.Conn.Close()
		if err != nil {
			client.Logger.Error("Failed to close client connection", "error", err)
			return
		}

		remainingClients := h.hub.GetClientCount()
		client.Logger.Info("Sensor disconnected", "remaining_clients", remainingClients)

		if remainingClients == 0 {
			ctx := context.Background()
			if err := h.examRepo.CloseLastExamination(ctx); err != nil {
				client.Logger.Error("Failed to close examination", "error", err)
			} else {
				client.Logger.Info("Examination closed - no clients remaining")
			}
		}
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
			break
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

		client.Logger.Info("Received sensor data",
			"sec_from_start", messageData.SecFromStart,
			"bpm", messageData.Data.BPMChild,
			"uterus", messageData.Data.Uterus,
			"spasms", messageData.Data.Spasms,
		)

		ctgData := entities.CTGData{
			SensorID:     client.SensorID,
			SecFromStart: messageData.SecFromStart,
			BPMChild:     messageData.Data.BPMChild,
			Uterus:       messageData.Data.Uterus,
			Spasms:       messageData.Data.Spasms,
		}

		ctx := context.Background()
		if err := h.examRepo.AddCtgRow(ctx, ctgData); err != nil {
			client.Logger.Error("Failed to save CTG data", "error", err)
			continue
		}

		client.Logger.Debug("CTG data saved successfully")
		h.sendToMLService(ctgData)
	}
}
