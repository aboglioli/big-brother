package models

import (
	"time"

	"github.com/google/uuid"
)

// ID
func NewID() string {
	return uuid.New().String()
}

// Base
type Base struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

func NewBase() Base {
	return Base{
		ID:      NewID(),
		Enabled: true,
	}
}

// Timestamp
type Timestamp struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func NewTimestamp() Timestamp {
	return Timestamp{
		CreatedAt: time.Now().UTC(),
	}
}
