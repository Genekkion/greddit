package forum

import (
	"greddit/internal/domains/auth"
	"greddit/internal/domains/shared"

	"github.com/google/uuid"
)

type CommentId uuid.UUID

// Comment represents a comment on a post.
type Comment struct {
	shared.Base

	Id CommentId `json:"id"`

	Body string `json:"body"`

	UserId auth.UserId `json:"user_id"`
	PostId PostId      `json:"post_id"`
}
