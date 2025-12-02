package dbportsforum

import (
	"context"
	"time"

	"greddit/internal/domains/auth"

	"greddit/internal/domains/forum"
)

// PostsRepo is a repository for posts.
type PostsRepo interface {
	// CreatePost creates a post.
	CreatePost(ctx context.Context, communityId forum.CommunityId, posterId auth.UserId, value forum.PostValue) (post *forum.Post, err error)

	// GetPostById returns a post by its ID.
	GetPostById(ctx context.Context, id forum.PostId) (post *forum.Post, err error)

	// GetPostsByCommunitySortedCreatedAt returns all posts in a community sorted by creation date.
	GetPostsByCommunitySortedCreatedAt(ctx context.Context, communityId forum.CommunityId, limit int, offset int) (posts []forum.Post, err error)

	// UpdatePostContent updates the content of a post.
	UpdatePostContent(ctx context.Context, id forum.PostId, content string) (updatedAt *time.Time, err error)

	// DeletePost soft deletes a post.
	DeletePost(ctx context.Context, id forum.PostId) (deletedAt *time.Time, err error)
}
