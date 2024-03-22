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

func (f *Friend) Get(ctx context.Context, userID int) ([]entity.Friend, error) {
	conn, err := f.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	return []entity.Friend{}, nil
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

func (f *Friend) DeleteFriend(ctx context.Context, userID, friendID int) error {
	conn, err := f.DBPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

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
