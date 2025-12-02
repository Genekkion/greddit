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

	CommunityValue
	CommunityMetadata
}

// CommunityValue represents the value of a community.
type CommunityValue struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Validate checks that the community value is valid.
func (v CommunityValue) Validate() error {
	{
		name := v.Name
		name = strings.TrimSpace(name)
		if name == "" {
			return InvalidCommunityParamsError{
				reason: "name cannot be empty",
			}
		} else if len(name) > communityMaxNameLength {
			return InvalidCommunityParamsError{
				reason: fmt.Sprintf("name must be less than %d characters", communityMaxNameLength),
			}
		}

		r := []rune(name)[0]
		if !unicode.IsLetter(r) {
			return InvalidCommunityParamsError{
				reason: "name must start with a letter",
			}
		}
	}

	{

		description := v.Description
		description = strings.TrimSpace(description)
		if len(description) > communityMaxDescriptionLength {
			return InvalidCommunityParamsError{
				reason: fmt.Sprintf("description must be less than %d characters", communityMaxDescriptionLength),
			}
		}
	}

	return nil
}

// CommunityMetadata represents metadata about a community.
type CommunityMetadata struct {
	Id          CommunityId `json:"id"`
	MemberCount int         `json:"member_count"`
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
func NewCommunity(value CommunityValue, metadata CommunityMetadata, base shared.Base) (community *Community, err error) {
	err = value.Validate()
	if err != nil {
		return nil, err
	}

	metadata.MemberCount = 0

	community = &Community{
		Base: base,

		CommunityValue:    value,
		CommunityMetadata: metadata,
	}

	return community, nil
}
