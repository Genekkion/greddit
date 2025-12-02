package auth

import (
	"greddit/internal/domains/shared"
	"greddit/internal/test"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	t.Run("Valid user", func(t *testing.T) {
		t.Parallel()

		value := UserValue{
			Username: "gene",
			Role:     RoleAdmin,
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
}
