package forumdb

import (
	"testing"
	"time"

	"greddit/internal/domains/auth"

	authdb "greddit/internal/infra/db/postgres/auth"

	"greddit/internal/domains/forum"
	"greddit/internal/infra/db/postgres"
	"greddit/internal/test"
)

func TestPostsRepo_CreatePost(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewPostsRepo(pool)
	communitiesRepo := NewCommunitiesRepo(pool)
	usersRepo := authdb.NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("create single post", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create a user first
		poster, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "testuser",
			DisplayName: "testuser",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Create a community
		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "golang",
			Description: "Go programming",
		})
		test.NilErr(t, err)

		value := forum.PostValue{
			Title: "My first post",
			Body:  "This is the body of my first post",
		}

		post, err := repo.CreatePost(ctx, community.Id, poster.Id, value)
		test.NilErr(t, err)
		test.AssertEqual(t, "Title not as expected", value.Title, post.Title)
		test.AssertEqual(t, "Body not as expected", value.Body, post.Body)
		test.Assert(t, "CreatedAt should not be zero", !post.CreatedAt.IsZero())
		test.Assert(t, "UpdatedAt should not be zero", !post.UpdatedAt.IsZero())
	})

	t.Run("create multiple posts in same community", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		poster, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "testuser",
			DisplayName: "testuser",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "programming",
			Description: "General programming",
		})
		test.NilErr(t, err)

		value1 := forum.PostValue{
			Title: "First Post",
			Body:  "Body of first post",
		}
		value2 := forum.PostValue{
			Title: "Second Post",
			Body:  "Body of second post",
		}

		post1, err := repo.CreatePost(ctx, community.Id, poster.Id, value1)
		test.NilErr(t, err)

		post2, err := repo.CreatePost(ctx, community.Id, poster.Id, value2)
		test.NilErr(t, err)

		test.Assert(t, "Post IDs should be different", post1.Id != post2.Id)
	})
}

func TestPostsRepo_GetPostById(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewPostsRepo(pool)
	communitiesRepo := NewCommunitiesRepo(pool)
	usersRepo := authdb.NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("retrieve an existing post", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		poster, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "testuser",
			DisplayName: "testuser",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Create a community and post
		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "testing",
			Description: "Testing community",
		})
		test.NilErr(t, err)

		value := forum.PostValue{
			Title: "Test Post",
			Body:  "Test Body",
		}

		createdPost, err := repo.CreatePost(ctx, community.Id, poster.Id, value)
		test.NilErr(t, err)

		// Retrieve the post
		post, err := repo.GetPostById(ctx, createdPost.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Title not as expected", createdPost.Title, post.Title)
		test.AssertEqual(t, "Body not as expected", createdPost.Body, post.Body)
	})

	t.Run("non-existent post", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		nonExistentId := forum.PostId{}

		post, err := repo.GetPostById(ctx, nonExistentId)
		test.Assert(t, "Expected error for non-existent post", err != nil)
		test.Assert(t, "Expected post to be nil", post == nil)
	})
}

func TestPostsRepo_GetPostsByCommunitySortedCreatedAt(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewPostsRepo(pool)
	communitiesRepo := NewCommunitiesRepo(pool)
	usersRepo := authdb.NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("retrieve posts sorted by creation date", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		poster, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "testuser",
			DisplayName: "testuser",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Create a community
		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "test",
			Description: "Test community",
		})
		test.NilErr(t, err)

		// Create posts with slight delays
		p1, err := repo.CreatePost(ctx, community.Id, poster.Id, forum.PostValue{
			Title: "First Post",
			Body:  "First",
		})
		test.NilErr(t, err)
		time.Sleep(time.Millisecond * 10)

		p2, err := repo.CreatePost(ctx, community.Id, poster.Id, forum.PostValue{
			Title: "Second Post",
			Body:  "Second",
		})
		test.NilErr(t, err)
		time.Sleep(time.Millisecond * 10)

		p3, err := repo.CreatePost(ctx, community.Id, poster.Id, forum.PostValue{
			Title: "Third Post",
			Body:  "Third",
		})
		test.NilErr(t, err)

		posts, err := repo.GetPostsByCommunitySortedCreatedAt(ctx, community.Id, 10, 0)
		test.NilErr(t, err)

		test.AssertEqual(t, "Expected 3 posts", 3, len(posts))

		// Should be in chronological order (oldest first)
		test.AssertEqual(t, "First should be p1", p1.Id, posts[0].Id)
		test.AssertEqual(t, "Second should be p2", p2.Id, posts[1].Id)
		test.AssertEqual(t, "Third should be p3", p3.Id, posts[2].Id)
	})

	t.Run("pagination", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		poster, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "testuser",
			DisplayName: "testuser",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Create a community
		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "pagination",
			Description: "Pagination test",
		})
		test.NilErr(t, err)

		// Create multiple posts
		n := 5
		for i := range n {
			_, err := repo.CreatePost(ctx, community.Id, poster.Id, forum.PostValue{
				Title: string(rune('A' + i)),
				Body:  "Post body",
			})
			test.NilErr(t, err)
		}

		// Get first page
		page1, err := repo.GetPostsByCommunitySortedCreatedAt(ctx, community.Id, 2, 0)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexpected number of posts", 2, len(page1))

		// Get second page
		page2, err := repo.GetPostsByCommunitySortedCreatedAt(ctx, community.Id, 2, 2)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexpected number of posts", 2, len(page2))
	})

	t.Run("only returns posts from specified community", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		poster, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "testuser",
			DisplayName: "testuser",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Create two communities
		community1, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "community1",
			Description: "First community",
		})
		test.NilErr(t, err)

		community2, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "community2",
			Description: "Second community",
		})
		test.NilErr(t, err)

		// Create posts in both communities
		_, err = repo.CreatePost(ctx, community1.Id, poster.Id, forum.PostValue{
			Title: "Community 1 Post",
			Body:  "Post in community 1",
		})
		test.NilErr(t, err)

		_, err = repo.CreatePost(ctx, community2.Id, poster.Id, forum.PostValue{
			Title: "Community 2 Post",
			Body:  "Post in community 2",
		})
		test.NilErr(t, err)

		// Get posts from community1
		posts, err := repo.GetPostsByCommunitySortedCreatedAt(ctx, community1.Id, 10, 0)
		test.NilErr(t, err)

		test.Assert(t, "Expected 1 post", len(posts) == 1)
		test.AssertEqual(t, "Post should be from community 1", "Community 1 Post", posts[0].Title)
	})
}

func TestPostsRepo_UpdatePostContent(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewPostsRepo(pool)
	communitiesRepo := NewCommunitiesRepo(pool)
	usersRepo := authdb.NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("update post content", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		poster, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "testuser",
			DisplayName: "testuser",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Create a community and post
		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "updatetest",
			Description: "Update test",
		})
		test.NilErr(t, err)

		value := forum.PostValue{
			Title: "Original Title",
			Body:  "Original content",
		}

		post, err := repo.CreatePost(ctx, community.Id, poster.Id, value)
		test.NilErr(t, err)

		// Update content
		newContent := "Updated content"

		updatedAt, err := repo.UpdatePostContent(ctx, post.Id, newContent)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil updated timestamp", updatedAt != nil)

		// Verify the update
		updatedPost, err := repo.GetPostById(ctx, post.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexpected body", newContent, updatedPost.Body)
		test.AssertEqual(t, "Title should remain unchanged", post.Title, updatedPost.Title)
	})

	t.Run("non-existent post", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		nonExistentId := forum.PostId{}

		updatedAt, err := repo.UpdatePostContent(ctx, nonExistentId, "New Content")
		test.Assert(t, "Expected error for non-existent post", err != nil)
		test.Assert(t, "Expected updated timestamp to be nil", updatedAt == nil)
	})

	t.Run("updates only the specified post", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		poster, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "testuser",
			DisplayName: "testuser",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Create a community and two posts
		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "multipost",
			Description: "Multiple posts",
		})
		test.NilErr(t, err)

		post1, err := repo.CreatePost(ctx, community.Id, poster.Id, forum.PostValue{
			Title: "Post 1",
			Body:  "Content 1",
		})
		test.NilErr(t, err)

		post2, err := repo.CreatePost(ctx, community.Id, poster.Id, forum.PostValue{
			Title: "Post 2",
			Body:  "Content 2",
		})
		test.NilErr(t, err)

		// Update only post1
		updatedContent := "Updated Content 1"
		updatedAt, err := repo.UpdatePostContent(ctx, post1.Id, updatedContent)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil updated timestamp", updatedAt != nil)

		// Verify post1 is updated but post2 is not
		retrieved1, err := repo.GetPostById(ctx, post1.Id)
		test.NilErr(t, err)

		retrieved2, err := repo.GetPostById(ctx, post2.Id)
		test.NilErr(t, err)

		test.AssertEqual(t, "Unexpected body for post1", updatedContent, retrieved1.Body)
		test.AssertEqual(t, "Unexpected body for post2", post2.Body, retrieved2.Body)
	})
}

func TestPostsRepo_DeletePost(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewPostsRepo(pool)
	communitiesRepo := NewCommunitiesRepo(pool)
	usersRepo := authdb.NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("soft delete a post", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		poster, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "testuser",
			DisplayName: "testuser",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Create a community and post
		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "deletetest",
			Description: "Delete test",
		})
		test.NilErr(t, err)

		value := forum.PostValue{
			Title: "To be deleted",
			Body:  "This post will be deleted",
		}

		post, err := repo.CreatePost(ctx, community.Id, poster.Id, value)
		test.NilErr(t, err)

		// Delete the post
		deletedAt, err := repo.DeletePost(ctx, post.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		// Verify that the post is soft-deleted
		deletedPost, err := repo.GetPostById(ctx, post.Id)
		test.NilErr(t, err)

		test.AssertEqual(t, "Unexpected title", post.Title, deletedPost.Title)
		test.AssertEqual(t, "Unexpected body", post.Body, deletedPost.Body)
		test.Assert(t, "DeletedAt should not be nil", deletedPost.DeletedAt != nil)
	})

	t.Run("non-existent post", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		nonExistentId := forum.PostId{}

		deletedAt, err := repo.DeletePost(ctx, nonExistentId)
		test.Assert(t, "Expected error for non-existent post", err != nil)
		test.Assert(t, "Expected deleted timestamp to be nil", deletedAt == nil)
	})

	t.Run("deletes only the specified post", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		poster, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "testuser",
			DisplayName: "testuser",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Create a community and two posts
		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "multidel",
			Description: "Multiple deletes",
		})
		test.NilErr(t, err)

		post1, err := repo.CreatePost(ctx, community.Id, poster.Id, forum.PostValue{
			Title: "Post 1",
			Body:  "Delete 1",
		})
		test.NilErr(t, err)

		post2, err := repo.CreatePost(ctx, community.Id, poster.Id, forum.PostValue{
			Title: "Post 2",
			Body:  "Delete 2",
		})
		test.NilErr(t, err)

		// Delete only post1
		deletedAt, err := repo.DeletePost(ctx, post1.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		// Verify post1 is deleted but post2 is not
		retrieved1, err := repo.GetPostById(ctx, post1.Id)
		test.NilErr(t, err)

		retrieved2, err := repo.GetPostById(ctx, post2.Id)
		test.NilErr(t, err)

		test.Assert(t, "post1 should be deleted", retrieved1.DeletedAt != nil)
		test.Assert(t, "post2 should not be deleted", retrieved2.DeletedAt == nil)
	})

	t.Run("post data still accessible after soft delete", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		poster, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "testuser",
			DisplayName: "testuser",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Create and delete a post
		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "softdel",
			Description: "Soft Delete",
		})
		test.NilErr(t, err)

		value := forum.PostValue{
			Title: "Soft Delete Post",
			Body:  "This will be soft deleted",
		}

		post, err := repo.CreatePost(ctx, community.Id, poster.Id, value)
		test.NilErr(t, err)

		deletedAt, err := repo.DeletePost(ctx, post.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		// Post should still be retrievable
		deletedPost, err := repo.GetPostById(ctx, post.Id)
		test.NilErr(t, err)

		test.AssertEqual(t, "Unexpected title", post.Title, deletedPost.Title)
		test.AssertEqual(t, "Unexpected body", post.Body, deletedPost.Body)
		test.Assert(t, "Deleted timestamp should not be nil", deletedPost.DeletedAt != nil)
	})
}

func TestPostsRepo_Integration(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewPostsRepo(pool)
	communitiesRepo := NewCommunitiesRepo(pool)
	usersRepo := authdb.NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("full CRUD lifecycle", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		poster, err := usersRepo.CreateUser(ctx, auth.UserValue{
			Username:    "testuser",
			DisplayName: "testuser",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Create a community
		community, err := communitiesRepo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "lifecycle",
			Description: "Lifecycle Community",
		})
		test.NilErr(t, err)

		// Create
		value := forum.PostValue{
			Title: "Lifecycle Post",
			Body:  "Lifecycle Body",
		}

		post, err := repo.CreatePost(ctx, community.Id, poster.Id, value)
		test.NilErr(t, err)

		// Read
		retrieved, err := repo.GetPostById(ctx, post.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Title not as expected", post.Title, retrieved.Title)
		test.AssertEqual(t, "Body not as expected", post.Body, retrieved.Body)

		// Update
		time.Sleep(time.Millisecond * 10) // Ensure time difference
		updatedContent := "Updated Body"
		updatedAt, err := repo.UpdatePostContent(ctx, post.Id, updatedContent)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil updated timestamp", updatedAt != nil)

		updated, err := repo.GetPostById(ctx, post.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexpected body", updatedContent, updated.Body)

		// Delete
		deletedAt, err := repo.DeletePost(ctx, post.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		deleted, err := repo.GetPostById(ctx, post.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected post to be soft-deleted", deleted.DeletedAt != nil)
	})
}
