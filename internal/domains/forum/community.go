package forum

import (
	"time"

	"github.com/google/uuid"
)

type CommunityId uuid.UUID

// Community represents a community, i.e. a subreddit.
type Community struct {
	Id        CommunityId `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt *time.Time  `json:"deleted_at"`

	Name        string `json:"name"`
	Description string `json:"description"`
	MemberCount int    `json:"member_count"`
}
