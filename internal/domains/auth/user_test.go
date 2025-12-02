package auth

import (
	"testing"
	"time"

	"greddit/internal/domains/shared"
	"greddit/internal/test"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		value := UserValue{
			Username:    "gene",
			DisplayName: "Gene",
			Role:        RoleAdmin,
		}
		metadata := UserMetadata{
			Id: uuid.New(),
		}
		tt := time.Date(2025, 10, 10, 12, 0, 0, 0, time.UTC)
		base := shared.Base{
			CreatedAt: tt,
			UpdatedAt: tt,
			DeletedAt: nil,
		}

		user, err := NewUser(value, metadata, base)
		test.NilErr(t, err)

		expected := User{
			Base:         base,
			UserValue:    value,
			UserMetadata: metadata,
		}

		test.AssertEqual(t, "User not as expected", expected, *user)
	})

	t.Run("invalid", func(t *testing.T) {
		data := []struct {
			name  string
			value UserValue
		}{
			{
				name: "empty username",
				value: UserValue{
					Username:    "",
					DisplayName: "Gene",
					Role:        RoleAdmin,
				},
			},
			{
				name: "username with spaces",
				value: UserValue{
					Username: "gene user",
				},
			},
			{
				name: "invalid role",
				value: UserValue{
					Username:    "gene",
					DisplayName: "Gene",
					Role:        Role("invalid"),
				},
			},
			{
				name: "username with only spaces",
				value: UserValue{
					Username:    "   ",
					DisplayName: "Gene",
					Role:        RoleAdmin,
				},
			},
			{
				name: "username too short",
				value: UserValue{
					Username:    "ab",
					DisplayName: "Gene",
					Role:        RoleAdmin,
				},
			},
			{
				name: "username too long",
				value: UserValue{
					Username:    "thisusernameiswaytoolongandexceedsthemaximumlength",
					DisplayName: "Gene",
					Role:        RoleAdmin,
				},
			},
			{
				name: "username with invalid characters",
				value: UserValue{
					Username:    "gene@123",
					DisplayName: "Gene",
					Role:        RoleAdmin,
				},
			},
			{
				name: "empty display name",
				value: UserValue{
					Username:    "gene",
					DisplayName: "",
					Role:        RoleAdmin,
				},
			},
			{
				name: "display name with only spaces",
				value: UserValue{
					Username:    "gene",
					DisplayName: "   ",
					Role:        RoleAdmin,
				},
			},
			{
				name: "display name too long",
				value: UserValue{
					Username:    "gene",
					DisplayName: "This display name is way too long and exceeds the maximum length allowed",
					Role:        RoleAdmin,
				},
			},
		}

		for _, d := range data {
			d := d
			t.Run(d.name, func(t *testing.T) {
				t.Parallel()

				metadata := UserMetadata{
					Id: uuid.New(),
				}
				tt := time.Date(2025, 10, 10, 12, 0, 0, 0, time.UTC)
				base := shared.Base{
					CreatedAt: tt,
					UpdatedAt: tt,
					DeletedAt: nil,
				}

				user, err := NewUser(d.value, metadata, base)
				test.Assert(t, "Expected non-nil error", err != nil)
				test.Assert(t, "Expected user to be nil", user == nil)
			})
		}
	})
}
