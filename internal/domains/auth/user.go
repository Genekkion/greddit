package auth

import (
	"greddit/internal/domains/shared"
	"strings"

	"github.com/google/uuid"
)

const (
	usernameMinLength = 3
	usernameMaxLength = 32
)

type UserId = uuid.UUID

// User represents a user in the system. Should be reused across
// all sub applications.
type User struct {
	shared.Base

	UserValue
	UserMetadata
}

// UserValue represents the value of a user.
type UserValue struct {
	Username string `json:"username"`
	Role     Role   `json:"role"`
}

// Validate checks that the user value is valid.
func (v UserValue) Validate() (err error) {
	{
		username := v.Username
		username = strings.TrimSpace(username)
		if username == "" {

		}
	}

	{
		role := v.Role
		err = role.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// UserMetadata represents metadata about a user.
type UserMetadata struct {
	Id UserId `json:"id"`
}

// InvalidUserParamsError represents an error when creating a user with invalid parameters.
type InvalidUserParamsError struct {
	reason string
}

// Error implements the error interface.
func (e InvalidUserParamsError) Error() string {
	return "invalid user params: " + e.reason
}

// NewUser creates a new user.
func NewUser(value UserValue, metadata UserMetadata, base shared.Base) (user *User, err error) {
	err = value.Validate()
	if err != nil {
		return nil, err
	}

	user = &User{
		Base: base,

		UserValue:    value,
		UserMetadata: metadata,
	}

	return user, nil
}
