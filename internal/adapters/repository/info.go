package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/domain/entities"
	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/infrastructure/database"
	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/infrastructure/ports/repository"
)

type infoRepository struct {
	db *database.DB
}

func NewInfoRepository(db *database.DB) repository.InfoRepository {
	return &infoRepository{
		db: db,
	}
}

func (r *infoRepository) GetDoctorByID(ctx context.Context, id int) (*entities.Doctor, error) {
	var doctor entities.Doctor
	query := `
		SELECT id, name, phone, specialization, license_number, med_id, is_active, created_at
		FROM doctors
		WHERE id = ?`

	err := r.db.GetContext(ctx, &doctor, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("doctor with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}

	return &doctor, nil
}

func (r *infoRepository) GetMedicalByID(ctx context.Context, id int) (*entities.Medical, error) {
	var medical entities.Medical
	query := `
		SELECT id, name, address, phone, email, license_number, created_at
		FROM medicals
		WHERE id = ?`

	err := r.db.GetContext(ctx, &medical, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("medical with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get medical: %w", err)
	}

	return &medical, nil
}

