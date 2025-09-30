package examinations

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/adapters/websocket/sensors"
	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/infrastructure/database"
	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/infrastructure/ports/repository"
)

type examRepository struct {
	db *database.DB
}

func (e *examRepository) CreateExamination(ctx context.Context) error {
	query := `
		INSERT INTO examinations (
			client_id, med_id, doctor_id, notes, status, 
			start_time, created_by, updated_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := e.db.ExecContext(ctx, query,
		1,          // client_id
		1,          // med_id
		1,          // doctor_id
		"",         // notes
		0,          // status
		time.Now(), // start_time
		1,          // created_by
		1,          // updated_by
	)

	if err != nil {
		return fmt.Errorf("failed to create examination: %w", err)
	}

	return nil
}

func (e *examRepository) AddCtgRow(ctx context.Context, data sensors.MessageData) error {
	var examinationID int
	err := e.db.GetContext(ctx, &examinationID, "SELECT id FROM examinations ORDER BY id DESC LIMIT 1")
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no examinations found, create examination first")
		}
		return fmt.Errorf("failed to get examination ID: %w", err)
	}

	query := `
		INSERT INTO ctg (
			examination_id, uuid, bpm, uterus, spasms, created_at
		) VALUES (?, ?, ?, ?, ?, ?)`

	_, err = e.db.ExecContext(ctx, query,
		examinationID,
		data.SensorID,
		data.Data.BPMChild,
		data.Data.Uterus,
		data.Data.Spasms,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert CTG data: %w", err)
	}

	return nil
}

func NewExamRepository(db *database.DB) repository.ExamRepository {
	return &examRepository{
		db: db,
	}
}
