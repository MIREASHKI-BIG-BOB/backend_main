package sensors

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

type Client struct {
	SensorID string
	Conn     *websocket.Conn
	Logger   *slog.Logger
}
