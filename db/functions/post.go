package functions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"segokuning/configs"
	"segokuning/db/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
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

func (p *Post) AddComment(ctx context.Context, postID int, comment entity.CommentPerPost) (entity.CommentPerPost, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return entity.CommentPerPost{}, err
	}
	defer conn.Release()

	commentJSON, err := json.Marshal(comment)
	if err != nil {
		return entity.CommentPerPost{}, err
	}

	sql := `UPDATE posts SET comments = comments || $1 WHERE id = $2 RETURNING created_at`
	fmt.Println(sql, commentJSON, postID)
	err = conn.QueryRow(ctx, sql, commentJSON, postID).Scan(&comment.CreatedAt)
	if err != nil {
		return entity.CommentPerPost{}, err
	}

	return comment, nil
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

	var (
		sql        = `SELECT id, post_in_html, tags, user_id, created_at, comments FROM posts where 1 = 1`
		arg        = 1
		args []any = []any{}
	)

	// only show post from friends join with friends table
	sql = fmt.Sprintf("%s AND user_id IN (SELECT friend_id FROM friends WHERE user_id = $%d)", sql, arg)
	args = append(args, query.UserId)
	arg++

	if query.Search != "" {
		sql = fmt.Sprintf("%s, AND post_in_html LIKE '%%%s%%'", sql, query.Search)
	}

	if len(query.SearchTags) > 0 {
		sql = fmt.Sprintf("%s AND $%v <@ tags", sql, arg)
		args = append(args, pq.Array(query.SearchTags))
		arg++
	}

	sql = fmt.Sprintf("%s LIMIT $%d", sql, arg)
	args = append(args, query.Limit)
	arg++

	sql = fmt.Sprintf("%s OFFSET $%d", sql, arg)
	args = append(args, query.Offset)
	arg++

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]entity.Post, 0)
	for rows.Next() {
		var post entity.Post
		err = rows.Scan(&post.Id, &post.PostInHtml, &post.Tags, &post.UserID, &post.CreatedAt, &post.Comments)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
