package repository

import (
	"context"
	"time"

	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/domain/entities"
)

type Examination struct {
	ID        int        `db:"id"`
	ClientID  int        `db:"client_id"`
	MedID     int        `db:"med_id"`
	DoctorID  int        `db:"doctor_id"`
	Notes     string     `db:"notes"`
	Status    int        `db:"status"`
	CloudID   *int       `db:"cloud_id"`
	StartTime time.Time  `db:"start_time"`
	EndTime   *time.Time `db:"end_time"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	CreatedBy int        `db:"created_by"`
	UpdatedBy int        `db:"updated_by"`
}

type ExamRepository interface {
	CreateExamination(ctx context.Context) error
	AddCtgRow(ctx context.Context, data entities.CTGData) error
	GetLastExamination(ctx context.Context) (*Examination, error)
	NeedsNewExamination(ctx context.Context) (bool, error)
	CloseLastExamination(ctx context.Context) error
}
