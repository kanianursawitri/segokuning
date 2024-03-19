package functions

import (
	"context"
	"shopifyx/configs"
	"shopifyx/db/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Comment struct {
	Config configs.Config
	DBPool *pgxpool.Pool
}

func NewComment(dbPool *pgxpool.Pool, config configs.Config) *Comment {
	return &Comment{
		DBPool: dbPool,
		Config: config,
	}
}

func (c *Comment) Add(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	conn, err := c.DBPool.Acquire(ctx)
	if err != nil {
		return entity.Comment{}, err
	}
	defer conn.Release()

	sql := `INSERT INTO comments (comment, user_id, post_id) VALUES ($1, $2, $3) RETURNING id,created_at`
	err = conn.QueryRow(ctx, sql, comment.Comment, comment.UserID, comment.PostID).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return entity.Comment{}, err
	}

	return comment, nil
}
