package forumdb

import (
	"context"
	"time"

	"greddit/internal/domains/auth"
	"greddit/internal/domains/forum"
	"greddit/internal/infra/db/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CommentsRepo implements the dbportsforum.CommentsRepo interface.
type CommentsRepo struct {
	postgres.BaseRepo
}

// NewCommentsRepo creates a new CommentsRepo.
func NewCommentsRepo(pool *pgxpool.Pool) CommentsRepo {
	return CommentsRepo{
		BaseRepo: postgres.NewBaseRepo(pool),
	}
}

func (c CommentsRepo) CreateComment(ctx context.Context, postId forum.PostId, commenterId auth.UserId,
	value forum.CommentValue, parentId *forum.CommentId,
) (comment *forum.Comment, err error) {
	const stmt = "INSERT INTO forum_comments (post_id, commenter_id, parent_id, body) VALUES ($1, $2, $3, $4) RETURNING id, created_at"
	args := []any{postId, commenterId, parentId, value.Body}

	comment = &forum.Comment{
		CommentValue: value,
		CommentMetadata: forum.CommentMetadata{
			PostId:      postId,
			CommenterId: commenterId,
			ParentId:    parentId,
		},
	}
	err = c.QueryRow(ctx, stmt, args...).Scan(&comment.Id, &comment.CreatedAt)
	if err != nil {
		return nil, err
	}

	comment.UpdatedAt = comment.CreatedAt

	return comment, nil
}

func (c CommentsRepo) GetCommentById(ctx context.Context, id forum.CommentId) (comment *forum.Comment, err error) {
	const stmt = "SELECT body, created_at, updated_at, deleted_at, parent_id FROM forum_comments WHERE id = $1"
	args := []any{id}

	comment = &forum.Comment{}
	comment.Id = id
	err = c.QueryRow(ctx, stmt, args...).Scan(
		&comment.Body, &comment.CreatedAt, &comment.UpdatedAt, &comment.DeletedAt, &comment.ParentId,
	)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (c CommentsRepo) GetCommentsByPostSortedCreatedAt(ctx context.Context, postId forum.PostId, limit int,
	offset int,
) (comments []forum.Comment, err error) {
	const stmt = "SELECT id, body, created_at, updated_at, deleted_at, post_id, commenter_id, parent_id FROM forum_comments WHERE post_id = $1 ORDER BY created_at LIMIT $2 OFFSET $3"
	args := []any{postId, limit, offset}

	return c.getCommentsAux(ctx, stmt, args, limit)
}

func (c CommentsRepo) GetCommentsByCommenterSortedCreatedAt(ctx context.Context, commenterId auth.UserId, limit int,
	offset int,
) (comments []forum.Comment, err error) {
	const stmt = "SELECT id, body, created_at, updated_at, deleted_at, post_id, commenter_id, parent_id FROM forum_comments WHERE commenter_id = $1 ORDER BY created_at LIMIT $2 OFFSET $3"
	args := []any{commenterId, limit, offset}

	return c.getCommentsAux(ctx, stmt, args, limit)
}

func (c CommentsRepo) getCommentsAux(ctx context.Context, stmt string, args []any, limit int) (
	comments []forum.Comment, err error,
) {
	comments = make([]forum.Comment, 0, limit)
	rows, err := c.Query(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		comment := forum.Comment{}
		err = rows.Scan(
			&comment.Id, &comment.Body,
			&comment.CreatedAt, &comment.UpdatedAt, &comment.DeletedAt,
			&comment.PostId, &comment.CommenterId, &comment.ParentId,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (c CommentsRepo) UpdateCommentBody(ctx context.Context, id forum.CommentId, body string) (
	updatedAt *time.Time, err error,
) {
	const stmt = "UPDATE forum_comments SET body = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at"
	args := []any{body, id}

	updatedAt = &time.Time{}
	err = c.QueryRow(ctx, stmt, args...).Scan(&updatedAt)
	if err != nil {
		return nil, err
	}

	return updatedAt, nil
}

func (c CommentsRepo) DeleteComment(ctx context.Context, id forum.CommentId) (deletedAt *time.Time, err error) {
	const stmt = "UPDATE forum_comments SET deleted_at = NOW() WHERE id = $1 RETURNING deleted_at"
	args := []any{id}

	deletedAt = &time.Time{}
	err = c.QueryRow(ctx, stmt, args...).Scan(&deletedAt)
	if err != nil {
		return nil, err
	}

	return deletedAt, nil
}
