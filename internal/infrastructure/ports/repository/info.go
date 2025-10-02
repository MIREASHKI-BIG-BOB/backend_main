package repository

import (
	"context"

	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/domain/entities"
)

type InfoRepository interface {
	GetDoctorByID(ctx context.Context, id int) (*entities.Doctor, error)
	GetMedicalByID(ctx context.Context, id int) (*entities.Medical, error)
}

