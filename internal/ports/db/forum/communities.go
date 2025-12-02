package dbportsforum

import (
	"context"
	"time"

	"greddit/internal/domains/forum"
)

// CommunitiesRepo is a repository for communities.
type CommunitiesRepo interface {
	// CreateCommunity creates a community.
	CreateCommunity(ctx context.Context, value forum.CommunityValue) (community *forum.Community, err error)

	// GetCommunityById returns a community by its ID.
	GetCommunityById(ctx context.Context, id forum.CommunityId) (community *forum.Community, err error)

	// GetAllCommunitiesSortedByName returns all communities sorted by name. Note
	// that soft-deleted communities are not returned.
	GetAllCommunitiesSortedByName(ctx context.Context, limit int, offset int) (
		communities []forum.Community, err error)

	// GetAllCommunitiesSortedByCreatedAt returns all communities sorted by creation
	// date. Note that soft-deleted communities are not returned.
	GetAllCommunitiesSortedByCreatedAt(ctx context.Context, limit int, offset int) (
		communities []forum.Community, err error)

	// GetAllCommunitiesSortedByUpdatedAt returns all communities sorted by update
	// date. Note that soft-deleted communities are not returned.
	GetAllCommunitiesSortedByUpdatedAt(ctx context.Context, limit int, offset int) (
		communities []forum.Community, err error)

	// UpdateCommunityDescription updates the description of a community.
	UpdateCommunityDescription(ctx context.Context, id forum.CommunityId, description string) (updatedAt *time.Time, err error)

	// DeleteCommunity soft deletes a community.
	DeleteCommunity(ctx context.Context, id forum.CommunityId) (deletedAt *time.Time, err error)
}
