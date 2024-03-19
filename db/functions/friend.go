package functions

import (
	"context"
	"errors"
	"shopifyx/configs"

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

	sql := `SELECT id FROM friends WHERE (first_user_id = $1 AND second_friend_id = $2) OR (first_user_id = $2 AND second_friend_id = $1)`
	row := conn.QueryRow(ctx, sql, userID, friendID)
	var id int
	err = row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
