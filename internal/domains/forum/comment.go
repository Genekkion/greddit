package forum

import (
	"fmt"
	"greddit/internal/domains/auth"
	"greddit/internal/domains/shared"
	"strings"

	"github.com/google/uuid"
)

const (
	commentMaxBodyLength = 1024 * 1024
)

type CommentId = uuid.UUID

// Comment represents a comment on a post.
type Comment struct {
	shared.Base

	CommentValue
	CommentMetadata
}

// CommentValue represents the value of a comment.
type CommentValue struct {
	Body string `json:"body"`
}

// Validate checks that the comment value is valid.
func (v CommentValue) Validate() error {
	body := v.Body
	body = strings.TrimSpace(body)
	if body == "" {
		return InvalidCommentParamsError{
			reason: "body cannot be empty",
		}
	} else if len(body) > commentMaxBodyLength {
		return InvalidCommentParamsError{
			reason: fmt.Sprintf("body must be less than %d characters", commentMaxBodyLength),
		}
	}

	return nil
}

// CommentMetadata represents metadata about a comment.
type CommentMetadata struct {
	Id     CommentId   `json:"id"`
	UserId auth.UserId `json:"user_id"`
	PostId PostId      `json:"post_id"`
}

// InvalidCommentParamsError represents an error when creating a comment with invalid parameters.
type InvalidCommentParamsError struct {
	reason string
}

// Error implements the error interface.
func (e InvalidCommentParamsError) Error() string {
	return "invalid comment params: " + e.reason
}

// NewComment creates a new comment.
func NewComment(value CommentValue, metadata CommentMetadata, base shared.Base) (comment *Comment, err error) {
	err = value.Validate()
	if err != nil {
		return nil, err
	}

	comment = &Comment{
		Base: base,

		CommentValue:    value,
		CommentMetadata: metadata,
	}

	return comment, nil
}
