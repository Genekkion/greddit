package forumdb

import (
	"fmt"
	"greddit/internal/domains/forum"
	"greddit/internal/infra/db/postgres"
	"greddit/internal/test"
	"testing"
	"time"
)

func TestCommunitiesRepo_CreateCommunity(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommunitiesRepo(pool)
	ctx := t.Context()

	t.Run("create single community", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		value := forum.CommunityValue{
			Name:        "golang",
			Description: "Discussion about Go programming language",
		}

		community, err := repo.CreateCommunity(ctx, value)
		test.NilErr(t, err)
		test.AssertEqual(t, "Name not as expected", value.Name, community.Name)
		test.AssertEqual(t, "Description not as expected", value.Description, community.Description)
	})

	t.Run("different community names", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		value1 := forum.CommunityValue{
			Name:        "programming",
			Description: "General programming discussions",
		}
		value2 := forum.CommunityValue{
			Name:        "python",
			Description: "Python programming language",
		}

		community1, err := repo.CreateCommunity(ctx, value1)
		test.NilErr(t, err)

		community2, err := repo.CreateCommunity(ctx, value2)
		test.NilErr(t, err)

		test.Assert(t, "Community IDs should be different", community1.Id != community2.Id)
	})

	t.Run("duplicate community name", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		value := forum.CommunityValue{
			Name:        "duplicate",
			Description: "First community",
		}

		_, err := repo.CreateCommunity(ctx, value)
		test.NilErr(t, err)

		// Try to create another community with the same name
		_, err = repo.CreateCommunity(ctx, value)
		test.Assert(t, "Expected error for duplicate community name", err != nil)
	})
}

func TestCommunitiesRepo_GetCommunityById(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommunitiesRepo(pool)
	ctx := t.Context()

	t.Run("retrieve an existing community", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create a community first
		value := forum.CommunityValue{
			Name:        "testing",
			Description: "Testing community",
		}
		createdCommunity, err := repo.CreateCommunity(ctx, value)
		test.NilErr(t, err)

		// Retrieve the community
		community, err := repo.GetCommunityById(ctx, createdCommunity.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Name not as expected", createdCommunity.Name, community.Name)
		test.AssertEqual(t, "Description not as expected", createdCommunity.Description, community.Description)
	})

	t.Run("non-existent community", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		nonExistentId := forum.CommunityId{}

		community, err := repo.GetCommunityById(ctx, nonExistentId)
		test.Assert(t, "Expected error for non-existent community", err != nil)
		test.Assert(t, "Expected community to be nil", community == nil)
	})
}

func TestCommunitiesRepo_GetAllCommunitiesSortedByName(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommunitiesRepo(pool)
	ctx := t.Context()

	t.Run("retrieve communities sorted by name", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create communities in random order
		_, err := repo.CreateCommunity(ctx, forum.CommunityValue{Name: "zebra", Description: "Z"})
		test.NilErr(t, err)
		_, err = repo.CreateCommunity(ctx, forum.CommunityValue{Name: "alpha", Description: "A"})
		test.NilErr(t, err)
		_, err = repo.CreateCommunity(ctx, forum.CommunityValue{Name: "beta", Description: "B"})
		test.NilErr(t, err)

		communities, err := repo.GetAllCommunitiesSortedByName(ctx, 3, 0)
		test.NilErr(t, err)
		test.Assert(t, "Expected 3 communities", len(communities) == 3)

		test.AssertEqual(t, "First community should be 'alpha'", "alpha", communities[0].Name)
		test.AssertEqual(t, "Second community should be 'beta'", "beta", communities[1].Name)
		test.AssertEqual(t, "Third community should be 'zebra'", "zebra", communities[2].Name)
	})

	t.Run("pagination", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create multiple communities
		n := 5
		for i := range n {
			_, err := repo.CreateCommunity(ctx, forum.CommunityValue{
				Name:        string(rune('a' + i - 1)),
				Description: "Community",
			})
			test.NilErr(t, err)
		}

		// Get first page
		page1, err := repo.GetAllCommunitiesSortedByName(ctx, 2, 0)
		test.NilErr(t, err)
		fmt.Println(page1)
		test.AssertEqual(t, "Unexpected number of communities", 2, len(page1))

		// Get second page
		page2, err := repo.GetAllCommunitiesSortedByName(ctx, 2, 2)
		test.NilErr(t, err)
		fmt.Println(page2)
		test.AssertEqual(t, "Unexpected number of communities", 2, len(page2))
	})

	t.Run("soft-deleted communities are excluded", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create two communities
		c1, err := repo.CreateCommunity(ctx, forum.CommunityValue{Name: "active", Description: "Active"})
		test.NilErr(t, err)
		c2, err := repo.CreateCommunity(ctx, forum.CommunityValue{Name: "deleted", Description: "To be deleted"})
		test.NilErr(t, err)

		// Delete one community
		_, err = repo.DeleteCommunity(ctx, c2.Id)
		test.NilErr(t, err)

		// Retrieve all communities
		communities, err := repo.GetAllCommunitiesSortedByName(ctx, 2, 0)
		test.NilErr(t, err)

		test.Assert(t, "Expected 1 community", len(communities) == 1)
		test.AssertEqual(t, "Should be the active community", c1.Name, communities[0].Name)
	})
}

func TestCommunitiesRepo_GetAllCommunitiesSortedByCreatedAt(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommunitiesRepo(pool)
	ctx := t.Context()

	t.Run("retrieve communities sorted by creation date", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create communities with slight delays
		c1, err := repo.CreateCommunity(ctx, forum.CommunityValue{Name: "first", Description: "First"})
		test.NilErr(t, err)
		time.Sleep(time.Millisecond * 10)

		c2, err := repo.CreateCommunity(ctx, forum.CommunityValue{Name: "second", Description: "Second"})
		test.NilErr(t, err)
		time.Sleep(time.Millisecond * 10)

		c3, err := repo.CreateCommunity(ctx, forum.CommunityValue{Name: "third", Description: "Third"})
		test.NilErr(t, err)

		communities, err := repo.GetAllCommunitiesSortedByCreatedAt(ctx, 10, 0)
		test.NilErr(t, err)

		test.Assert(t, "Expected 3 communities", len(communities) == 3)

		// Should be in reverse chronological order (newest first)
		test.AssertEqual(t, "First should be 'third'", c3.Id, communities[0].Id)
		test.AssertEqual(t, "Second should be 'second'", c2.Id, communities[1].Id)
		test.AssertEqual(t, "Third should be 'first'", c1.Id, communities[2].Id)
	})
}

func TestCommunitiesRepo_GetAllCommunitiesSortedByUpdatedAt(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommunitiesRepo(pool)
	ctx := t.Context()

	t.Run("retrieve communities sorted by update date", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create three communities
		c1, err := repo.CreateCommunity(ctx, forum.CommunityValue{Name: "c1", Description: "First"})
		test.NilErr(t, err)
		c2, err := repo.CreateCommunity(ctx, forum.CommunityValue{Name: "c2", Description: "Second"})
		test.NilErr(t, err)
		c3, err := repo.CreateCommunity(ctx, forum.CommunityValue{Name: "c3", Description: "Third"})
		test.NilErr(t, err)

		time.Sleep(time.Millisecond * 10)

		// Update the first community (making it most recently updated)
		_, err = repo.UpdateCommunityDescription(ctx, c1.Id, "Updated first")
		test.NilErr(t, err)

		communities, err := repo.GetAllCommunitiesSortedByUpdatedAt(ctx, 3, 0)
		test.NilErr(t, err)

		test.Assert(t, "Expected 3 communities", len(communities) == 3)

		// c1 should be first (most recently updated)
		test.AssertEqual(t, "First should be c1", c1.Id, communities[0].Id)
		test.AssertEqual(t, "Second should be c3", c3.Id, communities[1].Id)
		test.AssertEqual(t, "Third should be c2", c2.Id, communities[2].Id)
	})
}

func TestCommunitiesRepo_UpdateCommunityDescription(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommunitiesRepo(pool)
	ctx := t.Context()

	t.Run("update description", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create a community
		value := forum.CommunityValue{
			Name:        "updatetest",
			Description: "Original description",
		}

		community, err := repo.CreateCommunity(ctx, value)
		test.NilErr(t, err)

		// Update description
		newDescription := "Updated description"

		updatedAt, err := repo.UpdateCommunityDescription(ctx, community.Id, newDescription)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil updated timestamp", updatedAt != nil)

		// Verify the update
		updatedCommunity, err := repo.GetCommunityById(ctx, community.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexpected description", newDescription, updatedCommunity.Description)
		test.AssertEqual(t, "Unexpected name", community.Name, updatedCommunity.Name)
	})

	t.Run("non-existent community", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		nonExistentId := forum.CommunityId{}

		updatedAt, err := repo.UpdateCommunityDescription(ctx, nonExistentId, "New Description")
		test.Assert(t, "Expected error for non-existent community", err != nil)
		test.Assert(t, "Expected updated timestamp to be nil", updatedAt == nil)
	})

	t.Run("updates only the specified community", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create two communities
		community1, err := repo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "comm1",
			Description: "Description 1",
		})
		test.NilErr(t, err)
		community2, err := repo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "comm2",
			Description: "Description 2",
		})
		test.NilErr(t, err)

		// Update only community1
		updatedDescription := "Updated Description 1"
		updatedAt, err := repo.UpdateCommunityDescription(ctx, community1.Id, updatedDescription)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil updated timestamp", updatedAt != nil)

		// Verify community1 is updated but community2 is not
		retrieved1, err := repo.GetCommunityById(ctx, community1.Id)
		test.NilErr(t, err)

		retrieved2, err := repo.GetCommunityById(ctx, community2.Id)
		test.NilErr(t, err)

		test.AssertEqual(t, "Unexpected description", updatedDescription, retrieved1.Description)
		test.AssertEqual(t, "Unexpected description", community2.Description, retrieved2.Description)
	})
}

func TestCommunitiesRepo_DeleteCommunity(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommunitiesRepo(pool)
	ctx := t.Context()

	t.Run("soft delete a community", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create a community
		value := forum.CommunityValue{
			Name:        "deletetest",
			Description: "To be deleted",
		}

		community, err := repo.CreateCommunity(ctx, value)
		test.NilErr(t, err)

		// Delete the community
		deletedAt, err := repo.DeleteCommunity(ctx, community.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		// Verify that the community is soft-deleted
		deletedCommunity, err := repo.GetCommunityById(ctx, community.Id)
		test.NilErr(t, err)

		test.AssertEqual(t, "Unexpected name", community.Name, deletedCommunity.Name)
		test.AssertEqual(t, "Unexpected description", community.Description, deletedCommunity.Description)
		test.Assert(t, "DeletedAt should not be nil", deletedCommunity.DeletedAt != nil)
	})

	t.Run("non-existent community", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		nonExistentId := forum.CommunityId{}

		deletedAt, err := repo.DeleteCommunity(ctx, nonExistentId)
		test.Assert(t, "Expected error for non-existent community", err != nil)
		test.Assert(t, "Expected deleted timestamp to be nil", deletedAt == nil)
	})

	t.Run("deletes only the specified community", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create two communities
		community1, err := repo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "del1",
			Description: "Delete 1",
		})
		test.NilErr(t, err)

		community2, err := repo.CreateCommunity(ctx, forum.CommunityValue{
			Name:        "del2",
			Description: "Delete 2",
		})
		test.NilErr(t, err)

		// Delete only community1
		deletedAt, err := repo.DeleteCommunity(ctx, community1.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		// Verify community1 is deleted but community2 is not
		retrieved1, err := repo.GetCommunityById(ctx, community1.Id)
		test.NilErr(t, err)

		retrieved2, err := repo.GetCommunityById(ctx, community2.Id)
		test.NilErr(t, err)

		test.Assert(t, "community1 should be deleted", retrieved1.DeletedAt != nil)
		test.Assert(t, "community2 should not be deleted", retrieved2.DeletedAt == nil)
	})

	t.Run("community data still accessible after soft delete", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create and delete a community
		value := forum.CommunityValue{
			Name:        "softdel",
			Description: "Soft Delete",
		}

		community, err := repo.CreateCommunity(ctx, value)
		test.NilErr(t, err)

		deletedAt, err := repo.DeleteCommunity(ctx, community.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		// Community should still be retrievable
		deletedCommunity, err := repo.GetCommunityById(ctx, community.Id)
		test.NilErr(t, err)

		test.AssertEqual(t, "Unexpected name", community.Name, deletedCommunity.Name)
		test.AssertEqual(t, "Unexpected description", community.Description, deletedCommunity.Description)
		test.Assert(t, "Deleted timestamp should not be nil", deletedCommunity.DeletedAt != nil)
	})
}

func TestCommunitiesRepo_Integration(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewCommunitiesRepo(pool)
	ctx := t.Context()

	t.Run("full CRUD lifecycle", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create
		value := forum.CommunityValue{
			Name:        "lifecycle",
			Description: "Lifecycle Community",
		}
		community, err := repo.CreateCommunity(ctx, value)
		test.NilErr(t, err)

		// Read
		retrieved, err := repo.GetCommunityById(ctx, community.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Name not as expected", community.Name, retrieved.Name)
		test.AssertEqual(t, "Description not as expected", community.Description, retrieved.Description)

		// Update
		time.Sleep(time.Millisecond * 10) // Ensure time difference
		updatedDescription := "Updated Description"
		updatedAt, err := repo.UpdateCommunityDescription(ctx, community.Id, updatedDescription)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil updated timestamp", updatedAt != nil)

		updated, err := repo.GetCommunityById(ctx, community.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexpected description", updatedDescription, updated.Description)

		// Delete
		deletedAt, err := repo.DeleteCommunity(ctx, community.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		deleted, err := repo.GetCommunityById(ctx, community.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected community to be soft-deleted", deleted.DeletedAt != nil)
	})
}
