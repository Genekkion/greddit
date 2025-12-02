package auth

import "greddit/internal/util/set"

// Role represents a user role.
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

var allowedRoles = set.New[Role](set.WithSlice([]Role{
	RoleAdmin,
	RoleUser,
}))

// InvalidRoleError is returned when a role is invalid.
type InvalidRoleError struct {
	value string
}

// Error implements the error interface.
func (e InvalidRoleError) Error() string {
	return "invalid role value: " + e.value
}

// Validate checks that the role is valid.
func (r Role) Validate() (err error) {
	if !allowedRoles.Contains(r) {
		return InvalidRoleError{
			value: string(r),
		}
	}
	return nil
}
