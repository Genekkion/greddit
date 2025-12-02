package dbportsforum

import (
	"context"
	"time"

	"greddit/internal/domains/auth"

	"greddit/internal/domains/forum"
)

// CommentsRepo is a repository for comments.
type CommentsRepo interface {
	// CreateComment creates a comment.
	CreateComment(ctx context.Context, postId forum.PostId, commenterId auth.UserId, value forum.CommentValue) (
		comment *forum.Comment, err error)

	// GetCommentById returns a comment by its ID.
	GetCommentById(ctx context.Context, id forum.CommentId) (comment *forum.Comment, err error)

	// GetCommentsByPostSortedCreatedAt returns all comments in a post sorted by creation date.
	GetCommentsByPostSortedCreatedAt(ctx context.Context, postId forum.PostId, limit int, offset int) (
		comments []forum.Comment, err error)

	// GetCommentsByCommenterSortedCreatedAt returns all comments by a commenter sorted by creation date.
	GetCommentsByCommenterSortedCreatedAt(ctx context.Context, commenterId auth.UserId, limit int, offset int) (
		comments []forum.Comment, err error)

	// UpdateCommentBody updates the body of a comment.
	UpdateCommentBody(ctx context.Context, id forum.CommentId, body string) (updatedAt *time.Time, err error)

	// DeleteComment soft deletes a comment.
	DeleteComment(ctx context.Context, id forum.CommentId) (deletedAt *time.Time, err error)
}
