package forum

import (
	"time"

	"github.com/google/uuid"
)

type CommentId uuid.UUID

// Comment represents a comment on a post.
type Comment struct {
	Id        CommentId  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	Body string `json:"body"`

	PostId PostId `json:"post_id"`
}
