package forum

import (
	"greddit/internal/domains/auth"
	"greddit/internal/domains/shared"
	"strings"

	"github.com/google/uuid"
)

const (
	postMaxTitleLength = 128
	postMaxBodyLength  = 1024 * 1024
)

type PostId uuid.UUID

// Post represents a post in a community.
type Post struct {
	shared.Base

	PostMetadata
	PostValue
}

// PostMetadata represents metadata about a post.
type PostMetadata struct {
	Id PostId `json:"id"`

	PosterId    auth.UserId `json:"poster_id"`
	CommunityId CommunityId `json:"community_id"`

	UpvoteCount int `json:"upvote_count"`
}

// PostValue represents the value of a post.
type PostValue struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Validate checks that the post value is valid.
func (v PostValue) Validate() error {
	title := v.Title
	title = strings.TrimSpace(title)
	if title == "" {
		return InvalidPostParamsError{
			reason: "title cannot be empty",
		}
	} else if len(title) > postMaxTitleLength {
		return InvalidPostParamsError{
			reason: "title must be less than " + string(postMaxTitleLength) + " characters",
		}
	}

	body := v.Body
	body = strings.TrimSpace(body)
	if body == "" {
		return InvalidPostParamsError{
			reason: "body cannot be empty",
		}
	} else if len(body) > postMaxBodyLength {
		return InvalidPostParamsError{
			reason: "body must be less than " + string(postMaxBodyLength) + " characters",
		}
	}

	return nil
}

// InvalidPostParamsError is returned when a post is created with invalid parameters.
type InvalidPostParamsError struct {
	reason string
}

// Error implements the error interface.
func (e InvalidPostParamsError) Error() string {
	return "invalid post params: " + e.reason
}

// NewPost creates a new post.
func NewPost(value PostValue, metadata PostMetadata, base shared.Base) (
	post *Post, err error) {

	err = value.Validate()
	if err != nil {
		return nil, err
	}

	post = &Post{
		Base: base,

		PostMetadata: metadata,
		PostValue:    value,
	}

	return post, nil
}
