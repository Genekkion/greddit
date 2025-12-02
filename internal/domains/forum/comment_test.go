package forum

import (
	"greddit/internal/domains/auth"
	"greddit/internal/domains/shared"
	"greddit/internal/test"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewComment(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		value := CommentValue{
			Body: "This is a great post! Thanks for sharing.",
		}
		metadata := CommentMetadata{
			Id:     CommentId(uuid.New()),
			UserId: auth.UserId(uuid.New()),
			PostId: PostId(uuid.New()),
		}
		tt := time.Date(2025, 10, 10, 12, 0, 0, 0, time.UTC)
		base := shared.Base{
			CreatedAt: tt,
			UpdatedAt: tt,
			DeletedAt: nil,
		}

		comment, err := NewComment(value, metadata, base)
		test.NilErr(t, err)

		expected := Comment{
			Base:            base,
			CommentValue:    value,
			CommentMetadata: metadata,
		}

		test.AssertEqual(t, "Comment not as expected", expected, *comment)
	})

	t.Run("invalid", func(t *testing.T) {
		data := []struct {
			name  string
			value CommentValue
		}{
			{
				name: "empty body",
				value: CommentValue{
					Body: "",
				},
			},
			{
				name: "body with only spaces",
				value: CommentValue{
					Body: "   ",
				},
			},
			{
				name: "body with only whitespace",
				value: CommentValue{
					Body: "  \n\t  ",
				},
			},
			{
				name: "body too long",
				value: CommentValue{
					Body: strings.Repeat("a", commentMaxBodyLength+1),
				},
			},
		}

		for _, d := range data {
			d := d
			t.Run(d.name, func(t *testing.T) {
				t.Parallel()

				metadata := CommentMetadata{
					Id:     CommentId(uuid.New()),
					UserId: auth.UserId(uuid.New()),
					PostId: PostId(uuid.New()),
				}
				tt := time.Date(2025, 10, 10, 12, 0, 0, 0, time.UTC)
				base := shared.Base{
					CreatedAt: tt,
					UpdatedAt: tt,
					DeletedAt: nil,
				}

				comment, err := NewComment(d.value, metadata, base)
				test.Assert(t, "Expected non-nil error", err != nil)
				test.Assert(t, "Expected comment to be nil", comment == nil)
			})
		}
	})
}
