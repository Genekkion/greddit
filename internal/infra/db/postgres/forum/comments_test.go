package forumdb

import (
	"testing"
	"time"

	"greddit/internal/domains/auth"
	"greddit/internal/domains/forum"
	"greddit/internal/infra/db/postgres"
	"greddit/internal/test"

	authdb "greddit/internal/infra/db/postgres/auth"
)

func TestCommentsRepo_CreateComment(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommentsRepo(pool)
	postsRepo := NewPostsRepo(pool)
	communitiesRepo := NewCommunitiesRepo(pool)
	usersRepo := authdb.NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("create single comment", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create a user
		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Create a community and post
		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "golang",
			Description: "Go programming",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		value := forum.CommentValue{
			Body: "This is my first comment",
		}

		comment, err := repo.CreateComment(ctx, post.Id, commenter.Id, value, nil)
		test.NilErr(t, err)
		test.AssertEqual(t, "Body not as expected", value.Body, comment.Body)
		test.AssertEqual(t, "PostId not as expected", post.Id, comment.PostId)
		test.AssertEqual(t, "CommenterId not as expected", commenter.Id, comment.CommenterId)
		test.Assert(t, "ParentId should be nil", comment.ParentId == nil)
		test.Assert(t, "CreatedAt should not be zero", !comment.CreatedAt.IsZero())
		test.Assert(t, "UpdatedAt should not be zero", !comment.UpdatedAt.IsZero())
	})

	t.Run("create comment with parent", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create a user
		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Create a community and post
		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "testing",
			Description: "Testing community",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		// Create parent comment
		parentComment, err := repo.CreateComment(ctx, post.Id, commenter.Id, forum.CommentValue{
			Body: "Parent comment",
		}, nil)
		test.NilErr(t, err)

		// Create child comment
		childValue := forum.CommentValue{
			Body: "Child comment reply",
		}

		childComment, err := repo.CreateComment(ctx, post.Id, commenter.Id, childValue, &parentComment.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Body not as expected", childValue.Body, childComment.Body)
		test.Assert(t, "ParentId should not be nil", childComment.ParentId != nil)
		test.AssertEqual(t, "ParentId not as expected", parentComment.Id, *childComment.ParentId)
	})

	t.Run("create multiple comments on same post", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "multiple",
			Description: "Multiple comments",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		value1 := forum.CommentValue{Body: "First comment"}
		value2 := forum.CommentValue{Body: "Second comment"}

		comment1, err := repo.CreateComment(ctx, post.Id, commenter.Id, value1, nil)
		test.NilErr(t, err)

		comment2, err := repo.CreateComment(ctx, post.Id, commenter.Id, value2, nil)
		test.NilErr(t, err)

		test.Assert(t, "Comment IDs should be different", comment1.Id != comment2.Id)
	})
}

func TestCommentsRepo_GetCommentById(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommentsRepo(pool)
	postsRepo := NewPostsRepo(pool)
	communitiesRepo := NewCommunitiesRepo(pool)
	usersRepo := authdb.NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("retrieve an existing comment", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "testing",
			Description: "Testing community",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		value := forum.CommentValue{
			Body: "Test comment",
		}

		createdComment, err := repo.CreateComment(ctx, post.Id, commenter.Id, value, nil)
		test.NilErr(t, err)

		// Retrieve the comment
		comment, err := repo.GetCommentById(ctx, createdComment.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Body not as expected", createdComment.Body, comment.Body)
		test.AssertEqual(t, "ID not as expected", createdComment.Id, comment.Id)
	})

	t.Run("non-existent comment", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		nonExistentId := forum.CommentId{}

		comment, err := repo.GetCommentById(ctx, nonExistentId)
		test.Assert(t, "Expected error for non-existent comment", err != nil)
		test.Assert(t, "Expected comment to be nil", comment == nil)
	})
}

func TestCommentsRepo_GetCommentsByPostSortedCreatedAt(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommentsRepo(pool)
	postsRepo := NewPostsRepo(pool)
	communitiesRepo := NewCommunitiesRepo(pool)
	usersRepo := authdb.NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("retrieve comments sorted by creation date", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "test",
			Description: "Test community",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		// Create comments with slight delays
		c1, err := repo.CreateComment(ctx, post.Id, commenter.Id, forum.CommentValue{
			Body: "First comment",
		}, nil)
		test.NilErr(t, err)
		time.Sleep(time.Millisecond * 10)

		c2, err := repo.CreateComment(ctx, post.Id, commenter.Id, forum.CommentValue{
			Body: "Second comment",
		}, nil)
		test.NilErr(t, err)
		time.Sleep(time.Millisecond * 10)

		c3, err := repo.CreateComment(ctx, post.Id, commenter.Id, forum.CommentValue{
			Body: "Third comment",
		}, nil)
		test.NilErr(t, err)

		comments, err := repo.GetCommentsByPostSortedCreatedAt(ctx, post.Id, 10, 0)
		test.NilErr(t, err)

		test.AssertEqual(t, "Expected 3 comments", 3, len(comments))

		// Should be in chronological order (oldest first)
		test.AssertEqual(t, "First should be c1", c1.Id, comments[0].Id)
		test.AssertEqual(t, "Second should be c2", c2.Id, comments[1].Id)
		test.AssertEqual(t, "Third should be c3", c3.Id, comments[2].Id)
	})

	t.Run("pagination", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "pagination",
			Description: "Pagination test",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		// Create multiple comments
		n := 5
		for i := range n {
			_, err := repo.CreateComment(ctx, post.Id, commenter.Id, forum.CommentValue{
				Body: string(rune('A' + i)),
			}, nil)
			test.NilErr(t, err)
		}

		// Get first page
		page1, err := repo.GetCommentsByPostSortedCreatedAt(ctx, post.Id, 2, 0)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexpected number of comments", 2, len(page1))

		// Get second page
		page2, err := repo.GetCommentsByPostSortedCreatedAt(ctx, post.Id, 2, 2)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexpected number of comments", 2, len(page2))
	})

	t.Run("only returns comments from specified post", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "multipost",
			Description: "Multi post test",
		})
		test.NilErr(t, err)

		// Create two posts
		post1, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Post 1",
			Body:  "First post",
		})
		test.NilErr(t, err)

		post2, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Post 2",
			Body:  "Second post",
		})
		test.NilErr(t, err)

		// Create comments on both posts
		_, err = repo.CreateComment(ctx, post1.Id, commenter.Id, forum.CommentValue{
			Body: "Comment on post 1",
		}, nil)
		test.NilErr(t, err)

		_, err = repo.CreateComment(ctx, post2.Id, commenter.Id, forum.CommentValue{
			Body: "Comment on post 2",
		}, nil)
		test.NilErr(t, err)

		// Get comments from post1
		comments, err := repo.GetCommentsByPostSortedCreatedAt(ctx, post1.Id, 10, 0)
		test.NilErr(t, err)

		test.Assert(t, "Expected 1 comment", len(comments) == 1)
		test.AssertEqual(t, "Comment should be from post 1", "Comment on post 1", comments[0].Body)
	})
}

func TestCommentsRepo_GetCommentsByCommenterSortedCreatedAt(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommentsRepo(pool)
	postsRepo := NewPostsRepo(pool)
	communitiesRepo := NewCommunitiesRepo(pool)
	usersRepo := authdb.NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("retrieve comments by commenter", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter1, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter1",
			DisplayName: "commenter1",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		commenter2, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter2",
			DisplayName: "commenter2",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "test",
			Description: "Test community",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter1.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		// Create comments by different users
		c1, err := repo.CreateComment(ctx, post.Id, commenter1.Id, forum.CommentValue{
			Body: "Comment by user 1",
		}, nil)
		test.NilErr(t, err)

		_, err = repo.CreateComment(ctx, post.Id, commenter2.Id, forum.CommentValue{
			Body: "Comment by user 2",
		}, nil)
		test.NilErr(t, err)

		c3, err := repo.CreateComment(ctx, post.Id, commenter1.Id, forum.CommentValue{
			Body: "Another comment by user 1",
		}, nil)
		test.NilErr(t, err)

		// Get comments by commenter1
		comments, err := repo.GetCommentsByCommenterSortedCreatedAt(ctx, commenter1.Id, 10, 0)
		test.NilErr(t, err)

		test.AssertEqual(t, "Expected 2 comments", 2, len(comments))
		test.AssertEqual(t, "First comment should be c1", c1.Id, comments[0].Id)
		test.AssertEqual(t, "Second comment should be c3", c3.Id, comments[1].Id)
	})

	t.Run("pagination by commenter", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "pagination",
			Description: "Pagination test",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		// Create multiple comments
		n := 5
		for i := range n {
			_, err := repo.CreateComment(ctx, post.Id, commenter.Id, forum.CommentValue{
				Body: string(rune('A' + i)),
			}, nil)
			test.NilErr(t, err)
		}

		// Get first page
		page1, err := repo.GetCommentsByCommenterSortedCreatedAt(ctx, commenter.Id, 2, 0)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexpected number of comments", 2, len(page1))

		// Get second page
		page2, err := repo.GetCommentsByCommenterSortedCreatedAt(ctx, commenter.Id, 2, 2)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexpected number of comments", 2, len(page2))
	})
}

func TestCommentsRepo_UpdateCommentBody(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommentsRepo(pool)
	postsRepo := NewPostsRepo(pool)
	communitiesRepo := NewCommunitiesRepo(pool)
	usersRepo := authdb.NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("update comment body", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "updatetest",
			Description: "Update test",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		value := forum.CommentValue{
			Body: "Original comment",
		}

		comment, err := repo.CreateComment(ctx, post.Id, commenter.Id, value, nil)
		test.NilErr(t, err)

		// Update body
		newBody := "Updated comment body"

		updatedAt, err := repo.UpdateCommentBody(ctx, comment.Id, newBody)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil updated timestamp", updatedAt != nil)

		// Verify the update
		updatedComment, err := repo.GetCommentById(ctx, comment.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexpected body", newBody, updatedComment.Body)
	})

	t.Run("non-existent comment", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		nonExistentId := forum.CommentId{}

		updatedAt, err := repo.UpdateCommentBody(ctx, nonExistentId, "New Body")
		test.Assert(t, "Expected error for non-existent comment", err != nil)
		test.Assert(t, "Expected updated timestamp to be nil", updatedAt == nil)
	})

	t.Run("updates only the specified comment", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "multicomment",
			Description: "Multiple comments",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		comment1, err := repo.CreateComment(ctx, post.Id, commenter.Id, forum.CommentValue{
			Body: "Comment 1",
		}, nil)
		test.NilErr(t, err)

		comment2, err := repo.CreateComment(ctx, post.Id, commenter.Id, forum.CommentValue{
			Body: "Comment 2",
		}, nil)
		test.NilErr(t, err)

		// Update only comment1
		updatedBody := "Updated Comment 1"
		updatedAt, err := repo.UpdateCommentBody(ctx, comment1.Id, updatedBody)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil updated timestamp", updatedAt != nil)

		// Verify comment1 is updated but comment2 is not
		retrieved1, err := repo.GetCommentById(ctx, comment1.Id)
		test.NilErr(t, err)

		retrieved2, err := repo.GetCommentById(ctx, comment2.Id)
		test.NilErr(t, err)

		test.AssertEqual(t, "Unexpected body for comment1", updatedBody, retrieved1.Body)
		test.AssertEqual(t, "Unexpected body for comment2", comment2.Body, retrieved2.Body)
	})
}

func TestCommentsRepo_DeleteComment(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommentsRepo(pool)
	postsRepo := NewPostsRepo(pool)
	communitiesRepo := NewCommunitiesRepo(pool)
	usersRepo := authdb.NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("soft delete a comment", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "deletetest",
			Description: "Delete test",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		value := forum.CommentValue{
			Body: "To be deleted",
		}

		comment, err := repo.CreateComment(ctx, post.Id, commenter.Id, value, nil)
		test.NilErr(t, err)

		// Delete the comment
		deletedAt, err := repo.DeleteComment(ctx, comment.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		// Verify that the comment is soft-deleted
		deletedComment, err := repo.GetCommentById(ctx, comment.Id)
		test.NilErr(t, err)

		test.AssertEqual(t, "Unexpected body", comment.Body, deletedComment.Body)
		test.Assert(t, "DeletedAt should not be nil", deletedComment.DeletedAt != nil)
	})

	t.Run("non-existent comment", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		nonExistentId := forum.CommentId{}

		deletedAt, err := repo.DeleteComment(ctx, nonExistentId)
		test.Assert(t, "Expected error for non-existent comment", err != nil)
		test.Assert(t, "Expected deleted timestamp to be nil", deletedAt == nil)
	})

	t.Run("deletes only the specified comment", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "multidel",
			Description: "Multiple deletes",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		comment1, err := repo.CreateComment(ctx, post.Id, commenter.Id, forum.CommentValue{
			Body: "Comment 1",
		}, nil)
		test.NilErr(t, err)

		comment2, err := repo.CreateComment(ctx, post.Id, commenter.Id, forum.CommentValue{
			Body: "Comment 2",
		}, nil)
		test.NilErr(t, err)

		// Delete only comment1
		deletedAt, err := repo.DeleteComment(ctx, comment1.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		// Verify comment1 is deleted but comment2 is not
		retrieved1, err := repo.GetCommentById(ctx, comment1.Id)
		test.NilErr(t, err)

		retrieved2, err := repo.GetCommentById(ctx, comment2.Id)
		test.NilErr(t, err)

		test.Assert(t, "comment1 should be deleted", retrieved1.DeletedAt != nil)
		test.Assert(t, "comment2 should not be deleted", retrieved2.DeletedAt == nil)
	})

	t.Run("comment data still accessible after soft delete", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "softdel",
			Description: "Soft Delete",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		value := forum.CommentValue{
			Body: "Soft Delete Comment",
		}

		comment, err := repo.CreateComment(ctx, post.Id, commenter.Id, value, nil)
		test.NilErr(t, err)

		deletedAt, err := repo.DeleteComment(ctx, comment.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		// Comment should still be retrievable
		deletedComment, err := repo.GetCommentById(ctx, comment.Id)
		test.NilErr(t, err)

		test.AssertEqual(t, "Unexpected body", comment.Body, deletedComment.Body)
		test.Assert(t, "Deleted timestamp should not be nil", deletedComment.DeletedAt != nil)
	})
}

func TestCommentsRepo_Integration(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommentsRepo(pool)
	postsRepo := NewPostsRepo(pool)
	communitiesRepo := NewCommunitiesRepo(pool)
	usersRepo := authdb.NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("full CRUD lifecycle", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "lifecycle",
			Description: "Lifecycle Community",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		// Create
		value := forum.CommentValue{
			Body: "Lifecycle Comment",
		}

		comment, err := repo.CreateComment(ctx, post.Id, commenter.Id, value, nil)
		test.NilErr(t, err)

		// Read
		retrieved, err := repo.GetCommentById(ctx, comment.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Body not as expected", comment.Body, retrieved.Body)

		// Update
		time.Sleep(time.Millisecond * 10) // Ensure time difference
		updatedBody := "Updated Body"
		updatedAt, err := repo.UpdateCommentBody(ctx, comment.Id, updatedBody)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil updated timestamp", updatedAt != nil)

		updated, err := repo.GetCommentById(ctx, comment.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexpected body", updatedBody, updated.Body)

		// Delete
		deletedAt, err := repo.DeleteComment(ctx, comment.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		deleted, err := repo.GetCommentById(ctx, comment.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected comment to be soft-deleted", deleted.DeletedAt != nil)
	})

	t.Run("nested comments hierarchy", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		commenter, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "commenter",
			DisplayName: "commenter",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "nested",
			Description: "Nested comments",
		})
		test.NilErr(t, err)

		post, err := postsRepo.CreatePost(ctx, community.Id, commenter.Id, forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		})
		test.NilErr(t, err)

		// Create parent comment
		parent, err := repo.CreateComment(ctx, post.Id, commenter.Id, forum.CommentValue{
			Body: "Parent comment",
		}, nil)
		test.NilErr(t, err)

		// Create child comment
		child, err := repo.CreateComment(ctx, post.Id, commenter.Id, forum.CommentValue{
			Body: "Child comment",
		}, &parent.Id)
		test.NilErr(t, err)

		// Create grandchild comment
		grandchild, err := repo.CreateComment(ctx, post.Id, commenter.Id, forum.CommentValue{
			Body: "Grandchild comment",
		}, &child.Id)
		test.NilErr(t, err)

		// Verify hierarchy
		test.Assert(t, "Parent should have no parent", parent.ParentId == nil)
		test.Assert(t, "Child should have parent", child.ParentId != nil)
		test.AssertEqual(t, "Child parent should be parent", parent.Id, *child.ParentId)
		test.Assert(t, "Grandchild should have parent", grandchild.ParentId != nil)
		test.AssertEqual(t, "Grandchild parent should be child", child.Id, *grandchild.ParentId)

		// All comments should be retrievable by post
		comments, err := repo.GetCommentsByPostSortedCreatedAt(ctx, post.Id, 10, 0)
		test.NilErr(t, err)
		test.AssertEqual(t, "Expected 3 comments", 3, len(comments))
	})
}
