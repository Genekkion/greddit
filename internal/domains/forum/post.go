package forum

import (
	"greddit/internal/domains/auth"
	"time"

	"github.com/google/uuid"
)

type PostId uuid.UUID

// Post represents a post in a community.
type Post struct {
	Id        PostId     `json:"id"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	UpvoteCount int `json:"upvote_count"`

	PosterId    auth.UserId `json:"poster_id"`
	CommunityId CommunityId `json:"community_id"`
}
