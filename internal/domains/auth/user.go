package auth

import (
	"time"

	"github.com/google/uuid"
)

type UserId uuid.UUID

// User represents a user in the system. Should be reused across
// all sub applications.
type User struct {
	Id        UserId     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
