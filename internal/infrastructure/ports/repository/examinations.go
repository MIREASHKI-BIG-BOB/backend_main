package repository

import (
	"context"

	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/adapters/websocket/sensors"
)

type ExamRepository interface {
	CreateExamination(ctx context.Context) error
	AddCtgRow(ctx context.Context, data sensors.MessageData) error
}
