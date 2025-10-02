package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/domain/entities"
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
		1,
		1,
		1,
		"",
		0,
		time.Now(),
		1,
		1,
	)

	if err != nil {
		return fmt.Errorf("failed to create examination: %w", err)
	}

	return nil
}

func (e *examRepository) AddCtgRow(ctx context.Context, data entities.CTGData) error {
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
			examination_id, sec_from_start, uuid, bpm, uterus, spasms, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = e.db.ExecContext(ctx, query,
		examinationID,
		data.SecFromStart,
		data.SensorID,
		data.BPMChild,
		data.Uterus,
		data.Spasms,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert CTG data: %w", err)
	}

	return nil
}

func (e *examRepository) GetLastExamination(ctx context.Context) (*repository.Examination, error) {
	var exam repository.Examination
	query := `SELECT id, client_id, med_id, doctor_id, notes, status, cloud_id, 
	                 start_time, end_time, created_at, updated_at, created_by, updated_by 
	          FROM examinations 
	          ORDER BY id DESC 
	          LIMIT 1`

	err := e.db.GetContext(ctx, &exam, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last examination: %w", err)
	}

	return &exam, nil
}

func (e *examRepository) NeedsNewExamination(ctx context.Context) (bool, error) {
	lastExam, err := e.GetLastExamination(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check last examination: %w", err)
	}

	// Если обследований нет - нужно создать новое
	if lastExam == nil {
		return true, nil
	}

	// Если у последнего обследования заполнен end_time - нужно создать новое
	if lastExam.EndTime != nil {
		return true, nil
	}

	// Иначе используем существующее
	return false, nil
}

func (e *examRepository) CloseLastExamination(ctx context.Context) error {
	query := `
		UPDATE examinations 
		SET end_time = ?, updated_at = ?
		WHERE id = (SELECT id FROM examinations ORDER BY id DESC LIMIT 1)
		AND end_time IS NULL`

	result, err := e.db.ExecContext(ctx, query, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to close examination: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil
	}

	return nil
}

func NewExamRepository(db *database.DB) repository.ExamRepository {
	return &examRepository{
		db: db,
	}
}
