package forumports

import (
	"context"
	"greddit/internal/domains/auth"
	"time"
)

// UsersRepo is a repository for users.
type UsersRepo interface {
	// CreateUser creates a user.
	CreateUser(ctx context.Context, value auth.UserValue) (user *auth.User, err error)

	// GetUserById returns a user by its ID.
	GetUserById(ctx context.Context, id auth.UserId) (user *auth.User, err error)

	// UpdateDisplayName updates the display name of a user.
	UpdateDisplayName(ctx context.Context, id auth.UserId, displayName string) (updatedAt *time.Time, err error)

	// DeleteUser soft deletes a user.
	DeleteUser(ctx context.Context, id auth.UserId) (deletedAt *time.Time, err error)
}
