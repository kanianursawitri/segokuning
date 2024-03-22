package functions

import (
	"context"
	"errors"
	"segokuning/configs"
	"segokuning/db/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Friend struct {
	Config configs.Config
	DBPool *pgxpool.Pool
}

func NewFriend(dbPool *pgxpool.Pool, config configs.Config) *Friend {
	return &Friend{
		DBPool: dbPool,
		Config: config,
	}
}

func (f *Friend) IsFriend(ctx context.Context, userID, friendID int) (bool, error) {
	conn, err := f.DBPool.Acquire(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	sql := `SELECT id FROM friends WHERE user_id = $1 AND friend_id = $2`
	row := conn.QueryRow(ctx, sql, userID, friendID)
	var friend entity.Friend
	err = row.Scan(&friend.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (f *Friend) GetByID(ctx context.Context, friendID int) (entity.Friend, error) {
	conn, err := f.DBPool.Acquire(ctx)
	if err != nil {
		return entity.Friend{}, err
	}
	defer conn.Release()

	sql := `SELECT id FROM friends WHERE id = $1`
	row := conn.QueryRow(ctx, sql, friendID)
	var friend entity.Friend
	err = row.Scan(&friend.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Friend{}, nil
		}
		return entity.Friend{}, err
	}
	return friend, nil
}

func (f *Friend) Get(ctx context.Context, q entity.QueryGetFriends) ([]entity.Friend, error) {
	conn, err := f.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var (
		sql = `SELECT fs.friend_id AS id, fs.user_id AS userId, u.name, u.image_url, u.created_at 
                      FROM friends fs 
                      LEFT JOIN users u ON fs.friend_id = u.id 
                      WHERE 1 = 1`
		args []interface{}
	)

	if q.OnlyFriends {
		sql += " AND fs.user_id = $1"
		args = append(args, q.UserID)
	}

	if q.Search != "" {
		sql += " AND (u.name ILIKE '%' || $3 || '%' OR u.image_url ILIKE '%' || $3 || '%')"
		args = append(args, q.Search)
	}

	sql += " ORDER BY u.created_at"
	if q.OrderBy != "" {
		sql += " " + q.OrderBy
	}

	if q.Limit != 0 {
		sql += " LIMIT $4"
		args = append(args, q.Limit)
	} else {
		sql += " LIMIT 5"
		args = append(args, q.Limit, q.Offset)
	}

	if q.Offset != 0 {
		sql += " OFFSET $5"
		args = append(args, q.Offset)
	} else {
		sql += " OFFSET 0"
		args = append(args, q.Offset)
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	friends := make([]entity.Friend, 0)
	for rows.Next() {
		var friend entity.Friend
		err := rows.Scan(&friend.ID, &friend.UserID, &friend.Name, &friend.ImageUrl, &friend.CreatedAt)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

func (f *Friend) AddFriend(ctx context.Context, userID, friendID int) error {
	conn, err := f.DBPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if userID == friendID {
		return errors.New("NO_ADD_SELF")
	}

	getFriend, err := f.GetByID(ctx, friendID)
	if err != nil {
		return err
	}

	if getFriend.ID == 0 {
		return errors.New("FRIEND_NOT_FOUND")
	}

	// Check if the friendship already exists
	isFriend, err := f.IsFriend(ctx, userID, friendID)
	if err != nil {
		return err
	}
	if isFriend {
		return errors.New("FRIENDSHIP_EXISTS")
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		// Handle error
		return err
	}
	defer tx.Rollback(ctx)

	sql := `INSERT INTO friends (user_id, friend_id) VALUES ($1, $2)`
	_, err = tx.Exec(ctx, sql, userID, friendID)
	if err != nil {
		return err
	}

	// Update friends_counter
	sql = `UPDATE friends_counter AS fc SET friend_count=(SELECT COUNT (friend_id) FROM friends AS f WHERE f.user_id=fc.user_id AND user_id IN ($1,$2)) WHERE user_id IN ($1,$2)`
	_, err = tx.Exec(ctx, sql, userID, friendID)
	if err != nil {
		return err
	}
	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		// Handle error
		return err
	}

	return nil
}

func (f *Friend) DeleteFriend(ctx context.Context, userID, friendID int) error {
	conn, err := f.DBPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	getFriend, err := f.GetByID(ctx, friendID)
	if err != nil {
		return err
	}

	if getFriend.ID == 0 {
		return errors.New("FRIEND_NOT_FOUND")
	}

	isFriend, err := f.IsFriend(ctx, userID, friendID)
	if !isFriend {
		return errors.New("FRIENDSHIP_NOT_EXISTS")
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		// Handle error
		return err
	}
	defer tx.Rollback(ctx)

	sql := `DELETE FROM friends WHERE user_id = $1 AND friend_id = $2`
	_, err = tx.Exec(ctx, sql, userID, friendID)
	if err != nil {
		return err
	}

	// Update friends_counter
	sql = `UPDATE friends_counter SET friend_count=(SELECT COUNT (friend_id) FROM friends AS f WHERE f.user_id=fc.user_id AND user_id IN ($1,$2)) WHERE user_id IN ($1,$2)`
	_, err = tx.Exec(ctx, sql, userID, friendID)
	if err != nil {
		return err
	}
	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		// Handle error
		return err
	}

	return nil
}
