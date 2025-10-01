package sensors

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/MIREASHKI-BIG-BOB/backend_main/config"
)

type SensorStatus struct {
	UUID      string `json:"uuid"`
	IP        string `json:"ip"`
	Connected bool   `json:"connected"`
}

type SensorsUseCase interface {
	ConnectSensor(ctx context.Context) (*SensorStatus, error)
	DisconnectSensors(ctx context.Context) error
}

type sensorsUseCase struct {
	cfg              *config.Config
	logger           *slog.Logger
	httpClient       *http.Client
	mu               sync.RWMutex
	currentSensorIdx int
}

func NewSensorsUseCase(cfg *config.Config, logger *slog.Logger) SensorsUseCase {
	return &sensorsUseCase{
		cfg:              cfg,
		logger:           logger,
		httpClient:       &http.Client{},
		currentSensorIdx: -1,
	}
}

func (uc *sensorsUseCase) ConnectSensor(ctx context.Context) (*SensorStatus, error) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.currentSensorIdx >= len(uc.cfg.Sensors.Entities)-1 {
		return nil, fmt.Errorf("all sensors are already started")
	}

	uc.currentSensorIdx++
	sensor := uc.cfg.Sensors.Entities[uc.currentSensorIdx]

	url := fmt.Sprintf("http://%s/api/on", sensor.IP)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		uc.currentSensorIdx--
		return nil, fmt.Errorf("failed to create request for sensor %s: %w", sensor.UUID, err)
	}

	resp, err := uc.httpClient.Do(req)
	if err != nil {
		uc.currentSensorIdx--
		return nil, fmt.Errorf("failed to start sensor %s: %w", sensor.UUID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		uc.currentSensorIdx--
		return nil, fmt.Errorf("sensor %s returned status %d", sensor.UUID, resp.StatusCode)
	}

	uc.logger.Info("Sensor started",
		"uuid", sensor.UUID,
		"ip", sensor.IP,
		"index", uc.currentSensorIdx+1,
	)

	return &SensorStatus{
		UUID:      sensor.UUID,
		IP:        sensor.IP,
		Connected: true,
	}, nil
}

func (uc *sensorsUseCase) DisconnectSensors(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	if uc.currentSensorIdx < 0 {
		uc.logger.Info("No active sensors to stop")
		return nil
	}

	stoppedCount := 0
	for i := 0; i <= uc.currentSensorIdx; i++ {
		sensor := uc.cfg.Sensors.Entities[i]

		url := fmt.Sprintf("http://%s/api/off", sensor.IP)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			uc.logger.Warn("Failed to create stop request for sensor", "uuid", sensor.UUID, "error", err)
			continue
		}

		resp, err := uc.httpClient.Do(req)
		if err != nil {
			uc.logger.Warn("Failed to stop sensor", "uuid", sensor.UUID, "error", err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			uc.logger.Warn("Sensor returned error on stop", "uuid", sensor.UUID, "status", resp.StatusCode)
		} else {
			stoppedCount++
		}
	}
	uc.currentSensorIdx = -1

	uc.logger.Info("All sensors stopped", "stopped_count", stoppedCount)

	return nil
}
