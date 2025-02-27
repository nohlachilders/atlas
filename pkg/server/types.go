package server

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	IsAdmin        bool      `json:"-"`
	Token          string    `json:"token,omitempty"`
	RefreshToken   string    `json:"refresh_token,omitempty"`
}
