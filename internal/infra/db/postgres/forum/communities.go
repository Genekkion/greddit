package forumdb

import (
	"context"
	"time"

	"greddit/internal/domains/forum"
	"greddit/internal/infra/db/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CommunitiesRepo implements the dbforumports.CommunitiesRepo interface.
type CommunitiesRepo struct {
	postgres.BaseRepo
}

func NewCommunitiesRepo(pool *pgxpool.Pool) CommunitiesRepo {
	return CommunitiesRepo{
		BaseRepo: postgres.NewBaseRepo(pool),
	}
}

func (r CommunitiesRepo) CreateCommunity(ctx context.Context, value forum.CommunityValue) (community *forum.Community, err error) {
	const stmt = "INSERT INTO forum_communities (name, description) VALUES ($1, $2) RETURNING id, created_at"
	args := []any{value.Name, value.Description}

	community = &forum.Community{
		CommunityValue: value,
	}
	err = r.QueryRow(ctx, stmt, args...).Scan(&community.Id, &community.CreatedAt)
	if err != nil {
		return nil, err
	}

	community.UpdatedAt = community.CreatedAt

	return community, nil
}

func (r CommunitiesRepo) GetCommunityById(ctx context.Context, id forum.CommunityId) (community *forum.Community, err error) {
	const stmt = "SELECT name, description, created_at, updated_at, deleted_at FROM forum_communities WHERE id = $1"
	args := []any{id}

	community = &forum.Community{}

	err = r.QueryRow(ctx, stmt, args...).Scan(
		&community.Name, &community.Description, &community.CreatedAt, &community.UpdatedAt, &community.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return community, nil
}

func (r CommunitiesRepo) GetAllCommunitiesSortedByName(ctx context.Context, limit int, offset int) (
	communities []forum.Community, err error,
) {
	const stmt = "SELECT id, name, description, created_at, updated_at FROM forum_communities WHERE deleted_at IS NULL ORDER BY name LIMIT $1 OFFSET $2"
	args := []any{limit, offset}

	return r.getAllCommunitiesAux(ctx, stmt, args, limit)
}

func (r CommunitiesRepo) GetAllCommunitiesSortedByCreatedAt(ctx context.Context, limit int, offset int) (
	communities []forum.Community, err error,
) {
	const stmt = "SELECT id, name, description, created_at, updated_at FROM forum_communities WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2"
	args := []any{limit, offset}

	return r.getAllCommunitiesAux(ctx, stmt, args, limit)
}

func (r CommunitiesRepo) GetAllCommunitiesSortedByUpdatedAt(ctx context.Context, limit int, offset int) (
	communities []forum.Community, err error,
) {
	const stmt = "SELECT id, name, description, created_at, updated_at FROM forum_communities WHERE deleted_at IS NULL ORDER BY updated_at DESC LIMIT $1 OFFSET $2"
	args := []any{limit, offset}

	return r.getAllCommunitiesAux(ctx, stmt, args, limit)
}

func (r CommunitiesRepo) getAllCommunitiesAux(ctx context.Context, stmt string, args []any, limit int) (communities []forum.Community, err error) {
	communities = make([]forum.Community, 0, limit)

	rows, err := r.Query(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := forum.Community{}
		err = rows.Scan(&c.Id, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		communities = append(communities, c)
	}
	return communities, nil
}

func (r CommunitiesRepo) UpdateCommunityDescription(ctx context.Context, id forum.CommunityId, description string) (updatedAt *time.Time, err error) {
	const stmt = "UPDATE forum_communities SET description = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at"
	args := []any{description, id}

	updatedAt = &time.Time{}
	err = r.QueryRow(ctx, stmt, args...).Scan(&updatedAt)
	if err != nil {
		return nil, err
	}

	return updatedAt, nil
}

func (r CommunitiesRepo) DeleteCommunity(ctx context.Context, id forum.CommunityId) (deletedAt *time.Time, err error) {
	const stmt = "UPDATE forum_communities SET deleted_at = NOW() WHERE id = $1 RETURNING deleted_at"
	args := []any{id}

	deletedAt = &time.Time{}
	err = r.QueryRow(ctx, stmt, args...).Scan(&deletedAt)
	if err != nil {
		return nil, err
	}
	return deletedAt, nil
}
