package authdb

import (
	"testing"
	"time"

	"greddit/internal/test"

	"greddit/internal/domains/auth"
	"greddit/internal/infra/db/postgres"
)

func TestUsersRepo_CreateUser(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("create single user", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		value := auth.UserValue{
			Username:    "testuser",
			DisplayName: "Test User",
			Role:        auth.RoleUser,
		}

		user, err := repo.CreateUser(ctx, value)
		test.NilErr(t, err)
		test.AssertEqual(t, "Username not as expected", value.Username, user.Username)
		test.AssertEqual(t, "DisplayName not as expected", value.DisplayName, user.DisplayName)
		test.AssertEqual(t, "Role not as expected", value.Role, user.Role)
	})

	t.Run("different usernames", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		value1 := auth.UserValue{
			Username:    "user1",
			DisplayName: "User One",
			Role:        auth.RoleUser,
		}
		value2 := auth.UserValue{
			Username:    "user2",
			DisplayName: "User Two",
			Role:        auth.RoleAdmin,
		}

		user1, err := repo.CreateUser(ctx, value1)
		test.NilErr(t, err)

		user2, err := repo.CreateUser(ctx, value2)
		test.NilErr(t, err)

		test.Assert(t, "User IDs should be different", user1.Id != user2.Id)
	})

	t.Run("duplicate username", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		value := auth.UserValue{
			Username:    "duplicate",
			DisplayName: "Duplicate User",
			Role:        "user",
		}

		_, err := repo.CreateUser(ctx, value)
		test.NilErr(t, err)

		// Try to create another user with the same username
		_, err = repo.CreateUser(ctx, value)
		test.Assert(t, "Expected error for duplicate username", err != nil)
	})
}

func TestUsersRepo_GetUserById(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("retrieve an existing user", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create a user first
		value := auth.UserValue{
			Username:    "getuser",
			DisplayName: "Get User",
			Role:        "user",
		}
		createdUser, err := repo.CreateUser(ctx, value)
		test.NilErr(t, err)

		// Retrieve the user
		user, err := repo.GetUserById(ctx, createdUser.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Username not as expected", createdUser.Username, user.Username)
		test.AssertEqual(t, "DisplayName not as expected", createdUser.DisplayName, user.DisplayName)
		test.AssertEqual(t, "Role not as expected", createdUser.Role, user.Role)
	})

	t.Run("non-existent user", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		nonExistentId := auth.UserId{}

		user, err := repo.GetUserById(ctx, nonExistentId)
		test.Assert(t, "Expected error for non-existent user", err != nil)
		test.Assert(t, "Expected user to be nil", user == nil)
	})
}

func TestUsersRepo_UpdateDisplayName(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("update display name", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create a user
		value := auth.UserValue{
			Username:    "updateuser",
			DisplayName: "Original Name",
			Role:        auth.RoleUser,
		}

		user, err := repo.CreateUser(ctx, value)
		test.NilErr(t, err)

		// Update display name
		newDisplayName := "Updated Name"

		updatedAt, err := repo.UpdateDisplayName(ctx, user.Id, newDisplayName)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil updated timestamp", updatedAt != nil)

		// Verify the update
		updatedUser, err := repo.GetUserById(ctx, user.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexpected display name", newDisplayName, updatedUser.DisplayName)
		test.AssertEqual(t, "Unexpected username", user.Username, updatedUser.Username)
		test.AssertEqual(t, "Unexpected role", user.Role, updatedUser.Role)
	})

	t.Run("non-existent user", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		nonExistentId := auth.UserId{}

		updatedAt, err := repo.UpdateDisplayName(ctx, nonExistentId, "New Name")
		test.Assert(t, "Expected error for non-existent user", err != nil)
		test.Assert(t, "Expected updated timestamp to be nil", updatedAt == nil)
	})

	t.Run("updates only the specified user", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create two users
		user1, err := repo.CreateUser(ctx, auth.UserValue{
			Username:    "upd1",
			DisplayName: "Name 1",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)
		user2, err := repo.CreateUser(ctx, auth.UserValue{
			Username:    "upd2",
			DisplayName: "Name 2",
			Role:        auth.RoleUser,
		})
		test.NilErr(t, err)

		// Update only user1
		updatedDisplayName := "Updated Name 1"
		updatedAt, err := repo.UpdateDisplayName(ctx, user1.Id, updatedDisplayName)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil updated timestamp", updatedAt != nil)

		// Verify user1 is updated but user2 is not
		retrieved1, err := repo.GetUserById(ctx, user1.Id)
		test.NilErr(t, err)

		retrieved2, _ := repo.GetUserById(ctx, user2.Id)
		test.NilErr(t, err)

		test.AssertEqual(t, "Unexpected name", updatedDisplayName, retrieved1.DisplayName)
		test.AssertEqual(t, "Unexpected name", user2.DisplayName, retrieved2.DisplayName)
	})
}

func TestUsersRepo_DeleteUser(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("soft delete a user", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create a user
		value := auth.UserValue{
			Username:    "deleteuser",
			DisplayName: "Delete User",
			Role:        "user",
		}

		user, err := repo.CreateUser(ctx, value)
		test.NilErr(t, err)

		// Delete the user
		deletedAt, err := repo.DeleteUser(ctx, user.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		// Verify that the user is soft-deleted
		deletedUser, err := repo.GetUserById(ctx, user.Id)
		test.NilErr(t, err)

		test.AssertEqual(t, "Unexpected username", user.Username, deletedUser.Username)
		test.AssertEqual(t, "Unexpected display name", user.DisplayName, deletedUser.DisplayName)
		test.AssertEqual(t, "Unexpected role", user.Role, deletedUser.Role)
	})

	t.Run("non-existent user", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		nonExistentId := auth.UserId{}

		deletedAt, err := repo.DeleteUser(ctx, nonExistentId)
		test.Assert(t, "Expected error for non-existent user", err != nil)
		test.Assert(t, "Expected deleted timestamp to be nil", deletedAt == nil)
	})

	t.Run("Deletes only the specified user", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create two users
		user1, err := repo.CreateUser(ctx, auth.UserValue{
			Username:    "del1",
			DisplayName: "Delete 1",
			Role:        "user",
		})
		test.NilErr(t, err)

		user2, err := repo.CreateUser(ctx, auth.UserValue{
			Username:    "del2",
			DisplayName: "Delete 2",
			Role:        "user",
		})
		test.NilErr(t, err)

		// Delete only user1
		deletedAt, err := repo.DeleteUser(ctx, user1.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		// Verify user1 is deleted but user2 is not
		retrieved1, err := repo.GetUserById(ctx, user1.Id)
		test.NilErr(t, err)

		retrieved2, err := repo.GetUserById(ctx, user2.Id)
		test.NilErr(t, err)

		test.Assert(t, "user1 should be deleted", retrieved1.DeletedAt != nil)
		test.Assert(t, "user2 should not be deleted", retrieved2.DeletedAt == nil)
	})

	t.Run("user data still accessible after soft delete", func(t *testing.T) {
		postgres.ClearAllTables(t, pool)

		// Create and delete a user
		value := auth.UserValue{
			Username:    "softdel",
			DisplayName: "Soft Delete",
			Role:        auth.RoleUser,
		}

		user, err := repo.CreateUser(ctx, value)
		test.NilErr(t, err)

		deletedAt, err := repo.DeleteUser(ctx, user.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		// User should still be retrievable
		deletedUser, err := repo.GetUserById(ctx, user.Id)
		test.NilErr(t, err)

		test.AssertEqual(t, "Unexpected username", user.Username, deletedUser.Username)
		test.AssertEqual(t, "Unexpected display name", user.DisplayName, deletedUser.DisplayName)
		test.Assert(t, "Deleted timestamp should not be nil", deletedUser.DeletedAt != nil)
	})
}

func TestUsersRepo_Integration(t *testing.T) {
	t.Parallel()

	pool, cleanup := postgres.NewTestPool(t)
	defer cleanup()

	repo := NewUsersRepo(pool)
	ctx := t.Context()

	t.Run("full CRUD lifecycle", func(t *testing.T) {
		// Create
		value := auth.UserValue{
			Username:    "lifecycle",
			DisplayName: "Lifecycle User",
			Role:        auth.RoleUser,
		}
		user, err := repo.CreateUser(ctx, value)
		test.NilErr(t, err)

		// Read
		retrieved, err := repo.GetUserById(ctx, user.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "User not as expected", user, retrieved)

		// Update
		time.Sleep(time.Millisecond * 10) // Ensure time difference
		updatedDisplayName := "Updated Name"
		updatedAt, err := repo.UpdateDisplayName(ctx, user.Id, updatedDisplayName)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil updated timestamp", updatedAt != nil)

		updated, err := repo.GetUserById(ctx, user.Id)
		test.NilErr(t, err)
		test.AssertEqual(t, "Unexepected display name", updatedDisplayName, updated.DisplayName)

		// Delete
		deletedAt, err := repo.DeleteUser(ctx, user.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected a non-nil deleted timestamp", deletedAt != nil)

		deleted, err := repo.GetUserById(ctx, user.Id)
		test.NilErr(t, err)
		test.Assert(t, "Expected user to be soft-deleted", deleted.DeletedAt != nil)
	})
}
