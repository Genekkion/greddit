package auth

import (
	"fmt"
	"strings"
	"unicode"

	"greddit/internal/domains/shared"

	"github.com/google/uuid"
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
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Role        Role   `json:"role"`
}

const (
	nameMinLength = 3
	nameMaxLength = 32
)

var allowedUsernameChars = []*unicode.RangeTable{
	unicode.Letter,
	unicode.Digit,
}

func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return InvalidUserParamsError{
			reason: "username cannot be empty",
		}
	} else if len(username) < nameMinLength || len(username) > nameMaxLength {
		return InvalidUserParamsError{
			reason: fmt.Sprintf("username must be between %d and %d characters", nameMinLength, nameMaxLength),
		}
	}
	for _, r := range username {
		if !unicode.IsOneOf(allowedUsernameChars, r) {
			return InvalidUserParamsError{
				reason: "username contains invalid characters",
			}
		}
	}
	return nil
}

// Validate checks that the user value is valid.
func (v UserValue) Validate() (err error) {
	{
		err = ValidateUsername(v.Username)
		if err != nil {
			return err
		}
	}

	{
		displayName := v.DisplayName
		displayName = strings.TrimSpace(displayName)
		if displayName == "" {
			return InvalidUserParamsError{
				reason: "display name cannot be empty",
			}
		} else if len(displayName) > nameMaxLength {
			return InvalidUserParamsError{
				reason: fmt.Sprintf("display name must be less than %d characters", nameMaxLength),
			}
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
