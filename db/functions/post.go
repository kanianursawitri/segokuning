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

func (f *Post) GetCreator(ctx context.Context, userId int) (entity.Creator, error) {
	conn, err := f.dbPool.Acquire(ctx)
	if err != nil {
		return entity.Creator{}, err
	}
	defer conn.Release()
	var creator entity.Creator

	query := `SELECT u.id AS userId, u.name, u.image_url, fc.friend_count
						  FROM users u
						  LEFT JOIN friends_counter fc ON fc.user_id = u.id
						  WHERE u.id = $1`

	rows := conn.QueryRow(ctx, query, userId)
	err = rows.Scan(&creator.UserId, &creator.Name, &creator.ImageUrl, &creator.FriendCount)
	if err != nil {
		return entity.Creator{}, err
	}

	return creator, nil
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

	creator, err := p.GetCreator(ctx, comment.Creator.UserId)
	if err != nil {
		return entity.CommentPerPost{}, err
	}
	comment.Creator = creator

	commentJSON, err := json.Marshal(comment)
	if err != nil {
		return entity.CommentPerPost{}, err
	}

	sql := `UPDATE posts SET comments = comments || $1 WHERE id = $2 RETURNING created_at`
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
	sql = fmt.Sprintf("%s AND user_id IN (SELECT friend_id FROM friends WHERE user_id = $%d UNION SELECT $%d)", sql, arg, arg+1)
	args = append(args, query.UserId, query.UserId)
	arg += 2

	if query.Search != "" {
		sql = fmt.Sprintf("%s, AND post_in_html LIKE '%%%s%%'", sql, query.Search)
	}

	if len(query.SearchTags) > 0 {
		sql = fmt.Sprintf("%s AND $%v <@ tags", sql, arg)
		args = append(args, pq.Array(query.SearchTags))
		arg++
	}

	sql = fmt.Sprintf("%s ORDER BY created_at DESC", sql)

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

		//sort comments by created_at desc
		for i, j := 0, len(post.Comments)-1; i < j; i, j = i+1, j-1 {
			post.Comments[i], post.Comments[j] = post.Comments[j], post.Comments[i]
		}

		creator, err := p.GetCreator(ctx, post.UserID)
		if err != nil {
			return nil, err
		}
		post.Creator = creator
		posts = append(posts, post)
	}

	return posts, nil
}

func (p *Post) Count(ctx context.Context, query entity.QueryGetPosts) (int, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	var (
		sql        = `SELECT COUNT(*) FROM posts where 1 = 1`
		arg        = 1
		args []any = []any{}
	)

	// only show post from friends join with friends table
	sql = fmt.Sprintf("%s AND user_id IN (SELECT friend_id FROM friends WHERE user_id = $%d UNION SELECT $%d)", sql, arg, arg+1)
	args = append(args, query.UserId, query.UserId)
	arg += 2

	if query.Search != "" {
		sql = fmt.Sprintf("%s, AND post_in_html LIKE '%%%s%%'", sql, query.Search)
	}

	if len(query.SearchTags) > 0 {
		sql = fmt.Sprintf("%s AND $%v <@ tags", sql, arg)
		args = append(args, pq.Array(query.SearchTags))
		arg++
	}

	row := conn.QueryRow(ctx, sql, args...)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
