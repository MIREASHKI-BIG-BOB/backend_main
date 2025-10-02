package entities

import "time"

type Doctor struct {
	ID             int       `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Phone          string    `json:"phone" db:"phone"`
	Specialization string    `json:"specialization" db:"specialization"`
	LicenseNumber  int       `json:"license_number" db:"license_number"`
	MedID          int       `json:"med_id" db:"med_id"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

