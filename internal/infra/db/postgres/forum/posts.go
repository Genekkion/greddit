package forumdb

import (
	"context"
	"time"

	"greddit/internal/domains/auth"

	"greddit/internal/domains/forum"
	"greddit/internal/infra/db/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostsRepo implements the dbportsforum.PostsRepo interface.
type PostsRepo struct {
	postgres.BaseRepo
}

// NewPostsRepo creates a new PostsRepo.
func NewPostsRepo(pool *pgxpool.Pool) PostsRepo {
	return PostsRepo{
		BaseRepo: postgres.NewBaseRepo(pool),
	}
}

func (p PostsRepo) CreatePost(ctx context.Context, communityId forum.CommunityId, posterId auth.UserId, value forum.PostValue) (
	post *forum.Post, err error,
) {
	const stmt = "INSERT INTO forum_posts (community_id, poster_id, title, body) VALUES ($1, $2, $3, $4) RETURNING id, created_at"
	args := []any{communityId, posterId, value.Title, value.Body}

	post = &forum.Post{
		PostValue: value,
	}

	err = p.QueryRow(ctx, stmt, args...).Scan(&post.Id, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	post.UpdatedAt = post.CreatedAt

	return post, nil
}

func (p PostsRepo) GetPostById(ctx context.Context, id forum.PostId) (post *forum.Post, err error) {
	const stmt = "SELECT title, body, created_at, updated_at, deleted_at FROM forum_posts WHERE id = $1"
	args := []any{id}

	post = &forum.Post{}

	err = p.QueryRow(ctx, stmt, args...).Scan(
		&post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (p PostsRepo) GetPostsByCommunitySortedCreatedAt(ctx context.Context, communityId forum.CommunityId, limit int,
	offset int,
) (posts []forum.Post, err error) {
	const stmt = "SELECT id, title, body, created_at, updated_at, deleted_at FROM forum_posts WHERE community_id = $1 ORDER BY created_at LIMIT $2 OFFSET $3"
	args := []any{communityId, limit, offset}

	rows, err := p.Query(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts = make([]forum.Post, 0, limit)
	for rows.Next() {
		post := forum.Post{}
		err = rows.Scan(
			&post.Id, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (p PostsRepo) UpdatePostContent(ctx context.Context, id forum.PostId, content string) (updatedAt *time.Time, err error) {
	const stmt = "UPDATE forum_posts SET body = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at"
	args := []any{content, id}

	updatedAt = &time.Time{}
	err = p.QueryRow(ctx, stmt, args...).Scan(&updatedAt)
	if err != nil {
		return nil, err
	}

	return updatedAt, nil
}

func (p PostsRepo) DeletePost(ctx context.Context, id forum.PostId) (deletedAt *time.Time, err error) {
	const stmt = "UPDATE forum_posts SET deleted_at = NOW() WHERE id = $1 RETURNING deleted_at"
	args := []any{id}

	deletedAt = &time.Time{}
	err = p.QueryRow(ctx, stmt, args...).Scan(&deletedAt)
	if err != nil {
		return nil, err
	}

	return deletedAt, nil
}
