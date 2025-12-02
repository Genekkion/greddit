package authdb

import (
	"context"
	"greddit/internal/domains/auth"
	"greddit/internal/infra/db/postgres"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UsersRepo implements the forumports.UsersRepo interface.
type UsersRepo struct {
	postgres.BaseRepo
}

func NewUsersRepo(pool *pgxpool.Pool) UsersRepo {
	return UsersRepo{
		BaseRepo: postgres.NewBaseRepo(pool),
	}
}

func (r UsersRepo) CreateUser(ctx context.Context, value auth.UserValue) (user *auth.User, err error) {
	const stmt = "INSERT INTO auth_users (username, display_name, role) VALUES ($1, $2, $3) RETURNING id, created_at"
	args := []any{value.Username, value.DisplayName, value.Role}

	user = &auth.User{
		UserValue: value,
	}

	err = r.QueryRow(ctx, stmt, args...).Scan(
		&user.UserMetadata.Id,
		&user.Base.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	user.Base.UpdatedAt = user.Base.CreatedAt

	return user, nil
}

func (r UsersRepo) GetUserById(ctx context.Context, id auth.UserId) (user *auth.User, err error) {
	const stmt = "SELECT username, display_name, role, created_at, updated_at, deleted_at FROM auth_users WHERE id = $1"
	args := []any{id}

	user = &auth.User{}
	user.Id = id

	err = r.QueryRow(ctx, stmt, args...).Scan(
		&user.Username,
		&user.DisplayName,
		&user.Role,
		&user.Base.CreatedAt,
		&user.Base.UpdatedAt,
		&user.Base.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r UsersRepo) UpdateDisplayName(ctx context.Context, id auth.UserId, displayName string) (updatedAt *time.Time, err error) {
	const stmt = "UPDATE auth_users SET display_name = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at"
	args := []any{displayName, id}

	updatedAt = &time.Time{}
	err = r.QueryRow(ctx, stmt, args...).Scan(&updatedAt)
	if err != nil {
		return nil, err
	}

	return updatedAt, nil
}

func (r UsersRepo) DeleteUser(ctx context.Context, id auth.UserId) (deletedAt *time.Time, err error) {
	const stmt = "UPDATE auth_users SET deleted_at = NOW() WHERE id = $1 RETURNING deleted_at"
	args := []any{id}

	deletedAt = &time.Time{}
	err = r.QueryRow(ctx, stmt, args...).Scan(&deletedAt)
	if err != nil {
		return nil, err
	}

	return deletedAt, nil
}
