package forum

import (
	"fmt"
	"greddit/internal/domains/shared"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

const (
	communityMaxNameLength        = 64
	communityMaxDescriptionLength = 2048
)

type CommunityId uuid.UUID

// Community represents a community, i.e. a subreddit.
type Community struct {
	shared.Base

	Id CommunityId `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`
	MemberCount int    `json:"member_count"`
}

// InvalidCommunityParamsError represents an error when creating a community with invalid parameters.
type InvalidCommunityParamsError struct {
	reason string
}

// Error implements the error interface.
func (e InvalidCommunityParamsError) Error() string {
	return "invalid community params: " + e.reason
}

// NewCommunity creates a new community.
func NewCommunity(name string, description string, base shared.Base) (community *Community, err error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, InvalidCommunityParamsError{
			reason: "name cannot be empty",
		}
	} else if len(name) > communityMaxNameLength {
		return nil, InvalidCommunityParamsError{
			reason: fmt.Sprintf("name must be less than %d characters", communityMaxNameLength),
		}
	}

	r := []rune(name)[0]
	if !unicode.IsLetter(r) {
		return nil, InvalidCommunityParamsError{
			reason: "name must start with a letter",
		}
	}

	description = strings.TrimSpace(description)
	if len(description) > communityMaxDescriptionLength {
		return nil, InvalidCommunityParamsError{
			reason: fmt.Sprintf("description must be less than %d characters", communityMaxDescriptionLength),
		}
	}

	community = &Community{
		Base: base,

		Name:        name,
		Description: description,
	}

	return community, nil
}
