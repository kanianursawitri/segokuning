package functions

import (
	"context"
	"errors"
	"shopifyx/configs"
	"shopifyx/db/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	config configs.Config
	dbPool *pgxpool.Pool
}

func NewPost(dbPool *pgxpool.Pool, config configs.Config) *Post {
	return &Post{
		dbPool: dbPool,
		config: config,
	}
}

func (p *Post) Add(ctx context.Context, post entity.Post) (entity.Post, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return entity.Post{}, err
	}
	defer conn.Release()

	sql := `INSERT INTO posts (post_in_html, tags, user_id) VALUES ($1, $2, $3) RETURNING id,created_at`
	err = conn.QueryRow(ctx, sql, post.PostInHtml, post.Tags, post.UserID).Scan(&post.Id, &post.CreatedAt)
	if err != nil {
		return entity.Post{}, err
	}

	return post, nil
}

func (p *Post) GetByID(ctx context.Context, postID int) (entity.Post, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return entity.Post{}, err
	}
	defer conn.Release()

	sql := `SELECT id WHERE id = $1`
	row := conn.QueryRow(ctx, sql, postID)
	post := entity.Post{}
	err = row.Scan(&post.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Post{}, nil
		}
		return entity.Post{}, err
	}

	return post, nil
}
