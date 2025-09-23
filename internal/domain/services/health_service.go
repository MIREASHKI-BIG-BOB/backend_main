package services

import "time"

const ServiceName = "backend_main"

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

type HealthInfo struct {
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

func (h *HealthService) GetHealthStatus() *HealthInfo {
	return &HealthInfo{
		Timestamp: time.Now(),
		Service:   ServiceName,
	}
}
