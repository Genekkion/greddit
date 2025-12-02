package forum

import (
	"testing"
	"time"

	"greddit/internal/domains/shared"
	"greddit/internal/test"

	"github.com/google/uuid"
)

func TestNewCommunity(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		value := CommunityValue{
			Name:        "golang",
			Description: "A community for Go programming language enthusiasts",
		}
		metadata := CommunityMetadata{
			Id: uuid.New(),
		}
		tt := time.Date(2025, 10, 10, 12, 0, 0, 0, time.UTC)
		base := shared.Base{
			CreatedAt: tt,
			UpdatedAt: tt,
			DeletedAt: nil,
		}

		community, err := NewCommunity(value, metadata, base)
		test.NilErr(t, err)

		expected := Community{
			Base:              base,
			CommunityValue:    value,
			CommunityMetadata: metadata,
		}

		test.AssertEqual(t, "Community not as expected", expected, *community)
	})

	t.Run("invalid", func(t *testing.T) {
		data := []struct {
			name  string
			value CommunityValue
		}{
			{
				name: "empty name",
				value: CommunityValue{
					Name:        "",
					Description: "A valid description",
				},
			},
			{
				name: "name with only spaces",
				value: CommunityValue{
					Name:        "   ",
					Description: "A valid description",
				},
			},
			{
				name: "name too long",
				value: CommunityValue{
					Name:        "thisnameiswaytoolongandexceedsthemaximumlengthallowedforcommunitynameandshouldnotbevalid",
					Description: "A valid description",
				},
			},
			{
				name: "name starting with non-letter",
				value: CommunityValue{
					Name:        "123golang",
					Description: "A valid description",
				},
			},
			{
				name: "name starting with special character",
				value: CommunityValue{
					Name:        "_golang",
					Description: "A valid description",
				},
			},
			{
				name: "description too long",
				value: CommunityValue{
					Name: "golang",
					Description: "This is a very long description that exceeds the maximum allowed length for a community description. " +
						"It contains way too many characters and should fail validation. " +
						"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
						"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. " +
						"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. " +
						"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. " +
						"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
						"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. " +
						"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. " +
						"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. " +
						"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
						"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. " +
						"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. " +
						"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. " +
						"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
						"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. " +
						"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. " +
						"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. " +
						"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
						"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. " +
						"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. " +
						"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. " +
						"This extra text ensures we definitely exceed the 2048 character limit for the description validation to work properly and trigger an error.",
				},
			},
		}

		for _, d := range data {
			t.Run(d.name, func(t *testing.T) {
				t.Parallel()

				metadata := CommunityMetadata{
					Id: uuid.New(),
				}
				tt := time.Date(2025, 10, 10, 12, 0, 0, 0, time.UTC)
				base := shared.Base{
					CreatedAt: tt,
					UpdatedAt: tt,
					DeletedAt: nil,
				}

				c, err := NewCommunity(d.value, metadata, base)
				test.Assert(t, "Expected non-nil error", err != nil)
				test.Assert(t, "Expected community to be nil", c == nil)
			})
		}
	})
}
