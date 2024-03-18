package functions

import (
	"context"
	"errors"
	"shopifyx/configs"
	"shopifyx/db/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	config configs.Config
	dbPool *pgxpool.Pool
}

func NewUser(dbPool *pgxpool.Pool, config configs.Config) *User {
	return &User{
		dbPool: dbPool,
		config: config,
	}
}

func (u *User) Register(ctx context.Context, usr entity.User) (entity.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return entity.User{}, err
	}
	defer conn.Release()

	// Hash the password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), u.config.BcryptSalt)
	if err != nil {
		return entity.User{}, err
	}

	var existingId string

	err = conn.QueryRow(ctx, `SELECT id FROM users WHERE credential_value = $1`, usr.CredentialValue).Scan(&existingId)
	if existingId != "" {
		return entity.User{}, errors.New("EXISTING_USERNAME")
	}

	sql := `
		INSERT INTO users (name, credential_type, credential_value, password) VALUES ($1, $2, $3, $4)
	`

	_, err = conn.Exec(ctx, sql, usr.Name, usr.CredentialType, usr.CredentialValue, string(hashedPassword))

	var result entity.User

	err = conn.QueryRow(ctx, `SELECT id, name, credential_type, credential_value FROM users WHERE credential_value = $1`, usr.CredentialValue).Scan(&result.Id, &result.Name, &result.CredentialType, &result.CredentialValue)

	if err != nil {
		return entity.User{}, err
	}

	return entity.User{
		Id:              result.Id,
		Name:            result.Name,
		CredentialValue: result.CredentialValue,
		CredentialType:  result.CredentialType,
	}, nil
}

func (u *User) Login(ctx context.Context, usr entity.User) (entity.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return entity.User{}, err
	}
	defer conn.Release()

	var result entity.User

	err = conn.QueryRow(ctx, `SELECT id, name, credential_type, credential_value, password FROM users WHERE credential_value = $1`, usr.CredentialValue).Scan(
		&result.Id, &result.Name, &result.CredentialType, &result.CredentialValue, &result.Password,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, errors.New("USER_NOT_FOUND")
	}
	if err != nil {
		return result, err
	}

	// Compare the provided password with the hashed password from the database
	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(usr.Password)); err != nil {
		return result, errors.New("INVALID_PASSWORD")
	}

	return result, nil
}

func (u *User) GetUserById(ctx context.Context, userID string) (entity.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return entity.User{}, err
	}
	defer conn.Release()

	var result entity.User

	err = conn.QueryRow(ctx, `SELECT id, name, credential_type, credential_value FROM users WHERE id = $1`, userID).Scan(&result.Id, &result.Name, &result.CredentialType, &result.CredentialValue)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, ErrNoRow
	}
	if err != nil {
		return result, err
	}

	return result, nil
}
