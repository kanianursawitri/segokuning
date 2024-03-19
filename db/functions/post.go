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

	sql := `SELECT id,user_id FROM posts WHERE id = $1`
	row := conn.QueryRow(ctx, sql, postID)
	post := entity.Post{}
	err = row.Scan(&post.Id, &post.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Post{}, nil
		}
		return entity.Post{}, err
	}

	return post, nil
}

func (p *Post) Get(ctx context.Context, query entity.QueryGetPosts) ([]entity.Post, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sql := `SELECT id, post_in_html, tags, user_id, created_at FROM posts LIMIT $1 OFFSET $2`
	rows, err := conn.Query(ctx, sql, query.Limit, query.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]entity.Post, 0)
	for rows.Next() {
		var post entity.Post
		err = rows.Scan(&post.Id, &post.PostInHtml, &post.Tags, &post.UserID, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
