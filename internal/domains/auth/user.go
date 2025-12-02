package auth

import (
	"greddit/internal/domains/shared"
	"greddit/internal/util/set"

	"github.com/google/uuid"
)

type UserId uuid.UUID

// User represents a user in the system. Should be reused across
// all sub applications.
type User struct {
	shared.Base

	Id UserId `json:"id"`

	Role Role `json:"role"`
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

var (
	allowedRoles = set.New[Role](set.WithSlice([]Role{
		RoleAdmin,
		RoleUser,
	}))
)

// InvalidRoleError is returned when a role is invalid.
type InvalidRoleError struct {
	value string
}

// Error implements the error interface.
func (e InvalidRoleError) Error() string {
	return "invalid role value: " + e.value
}

// NewUser creates a new user.
func NewUser(id UserId, role Role, base shared.Base) (user *User, err error) {
	if !allowedRoles.Contains(role) {
		return nil, InvalidRoleError{
			value: string(role),
		}
	}

	user = &User{
		Base: base,

		Id:   id,
		Role: role,
	}

	return user, nil
}
