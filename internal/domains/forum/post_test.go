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

func TestNewPost(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		value := PostValue{
			Title: "Why Go is awesome",
			Body:  "Go is a great programming language for building scalable systems. Here's why...",
		}
		metadata := PostMetadata{
			Id:          PostId(uuid.New()),
			PosterId:    auth.UserId(uuid.New()),
			CommunityId: CommunityId(uuid.New()),
			UpvoteCount: 0,
		}
		tt := time.Date(2025, 10, 10, 12, 0, 0, 0, time.UTC)
		base := shared.Base{
			CreatedAt: tt,
			UpdatedAt: tt,
			DeletedAt: nil,
		}

		post, err := NewPost(value, metadata, base)
		test.NilErr(t, err)

		expected := Post{
			Base:         base,
			PostValue:    value,
			PostMetadata: metadata,
		}

		test.AssertEqual(t, "Post not as expected", expected, *post)
	})

	t.Run("invalid", func(t *testing.T) {
		data := []struct {
			name  string
			value PostValue
		}{
			{
				name: "empty title",
				value: PostValue{
					Title: "",
					Body:  "This is a valid body",
				},
			},
			{
				name: "title with only spaces",
				value: PostValue{
					Title: "   ",
					Body:  "This is a valid body",
				},
			},
			{
				name: "title too long",
				value: PostValue{
					Title: strings.Repeat("a", postMaxTitleLength+1),
					Body:  "This is a valid body",
				},
			},
			{
				name: "empty body",
				value: PostValue{
					Title: "Valid title",
					Body:  "",
				},
			},
			{
				name: "body with only spaces",
				value: PostValue{
					Title: "Valid title",
					Body:  "   ",
				},
			},
			{
				name: "body with only whitespace",
				value: PostValue{
					Title: "Valid title",
					Body:  "  \n\t  ",
				},
			},
			{
				name: "body too long",
				value: PostValue{
					Title: "Valid title",
					Body:  strings.Repeat("a", postMaxBodyLength+1),
				},
			},
			{
				name: "both title and body empty",
				value: PostValue{
					Title: "",
					Body:  "",
				},
			},
			{
				name: "both title and body only spaces",
				value: PostValue{
					Title: "   ",
					Body:  "   ",
				},
			},
		}

		for _, d := range data {
			d := d
			t.Run(d.name, func(t *testing.T) {
				t.Parallel()

				metadata := PostMetadata{
					Id:          PostId(uuid.New()),
					PosterId:    uuid.New(),
					CommunityId: CommunityId(uuid.New()),
					UpvoteCount: 0,
				}
				tt := time.Date(2025, 10, 10, 12, 0, 0, 0, time.UTC)
				base := shared.Base{
					CreatedAt: tt,
					UpdatedAt: tt,
					DeletedAt: nil,
				}

				post, err := NewPost(d.value, metadata, base)
				test.Assert(t, "Expected non-nil error", err != nil)
				test.Assert(t, "Expected post to be nil", post == nil)
			})
		}
	})
}
