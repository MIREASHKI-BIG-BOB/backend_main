package entities

import "time"

type Medical struct {
	ID            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Address       string    `json:"address" db:"address"`
	Phone         string    `json:"phone" db:"phone"`
	Email         string    `json:"email" db:"email"`
	LicenseNumber int       `json:"license_number" db:"license_number"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

