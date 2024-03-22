package functions

import (
	"context"
	"errors"
	"fmt"
	"segokuning/configs"
	"segokuning/db/entity"

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
	var sql string

	sql = fmt.Sprintf(`SELECT id FROM users WHERE %s = $1`, usr.CredentialType)

	err = conn.QueryRow(ctx, sql, usr.CredentialValue).Scan(&existingId)
	if existingId != "" {
		return entity.User{}, errors.New("EXISTING_USERNAME")
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		// Handle error
		return entity.User{}, err
	}
	defer tx.Rollback(ctx)

	sql = fmt.Sprintf(`INSERT INTO users (name, %s, password) VALUES ($1, $2, $3) RETURNING id, name, phone, email`, usr.CredentialType)

	err = tx.QueryRow(ctx, sql, usr.Name, usr.CredentialValue, string(hashedPassword)).Scan(&usr.Id, &usr.Name, &usr.Phone, &usr.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, errors.New("USER_NOT_FOUND")
		}
		return entity.User{}, err
	}
	// Insert into friends_counter with friend_count = 0
	_, err = tx.Exec(ctx, `INSERT INTO friends_counter (user_id, friend_count) VALUES ($1, $2)`, usr.Id, 0)
	if err != nil {
		// Handle error
		return entity.User{}, err
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		// Handle error
		return entity.User{}, err
	}

	return usr, nil
}

func (u *User) Login(ctx context.Context, usr entity.User) (entity.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return entity.User{}, err
	}
	defer conn.Release()

	var result entity.User
	var sql string

	sql = fmt.Sprintf(`SELECT id, name, phone, email name, password FROM users WHERE %s = $1`, usr.CredentialType)
	err = conn.QueryRow(ctx, sql, usr.CredentialValue).Scan(
		&result.Id, &result.Name, &result.Phone, &result.Email, &result.Password,
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

	err = conn.QueryRow(ctx, `SELECT id, name, phone, email name, password FROM users WHERE id = $1`, userID).Scan(&result.Id, &result.Name, &result.CredentialType, &result.CredentialValue)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, ErrNoRow
	}
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *User) UpdateEmail(ctx context.Context, userID string, email string) (entity.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return entity.User{}, err
	}
	defer conn.Release()

	var result entity.User

	// Check if the email already exists
	var existingEmail *string
	err = conn.QueryRow(ctx, `SELECT email FROM users WHERE email = $1`, email).Scan(&existingEmail)
	if existingEmail != nil {
		return result, errors.New("EMAIL_EXISTS")
	}
	// Check if the user already has an email
	err = conn.QueryRow(ctx, `SELECT email FROM users WHERE id = $1`, userID).Scan(&existingEmail)
	if existingEmail != nil {
		return result, errors.New("EMAIL_ALREADY_SET")
	}

	// If no errors, proceed to update the email
	err = conn.QueryRow(ctx, `UPDATE users SET email = $1 WHERE id = $2 RETURNING id, name, phone, email`, email, userID).Scan(&result.Id, &result.Name, &result.CredentialType, &result.CredentialValue)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, errors.New("USER_NOT_FOUND")
	}
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *User) UpdatePhone(ctx context.Context, userID string, phone string) (entity.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return entity.User{}, err
	}
	defer conn.Release()

	var result entity.User

	// Check if the phone already exists
	var existingPhone string
	err = conn.QueryRow(ctx, `SELECT phone FROM users WHERE phone = $1`, phone).Scan(&existingPhone)
	if err == nil {
		return result, errors.New("PHONE_EXISTS") // Returning 409 error
	}

	// Check if the user already has a phone
	err = conn.QueryRow(ctx, `SELECT phone FROM users WHERE id = $1`, userID).Scan(&existingPhone)
	if err == nil {
		return result, errors.New("PHONE_ALREADY_SET") // Returning 400 error
	}

	// If no errors, proceed to update the phone
	err = conn.QueryRow(ctx, `UPDATE users SET phone = $1 WHERE id = $2 RETURNING id, name, phone, email`, phone, userID).Scan(&result.Id, &result.Name, &result.CredentialType, &result.CredentialValue)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, errors.New("USER_NOT_FOUND")
	}
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *User) UpdateAccount(ctx context.Context, userID, name, imageURL string) (entity.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return entity.User{}, err
	}
	defer conn.Release()

	var result entity.User

	err = conn.QueryRow(ctx, `UPDATE users SET name = $1, image_url = $2 WHERE id = $3 RETURNING id, phone, email`, name, imageURL, userID).Scan(&result.Id, &result.Phone, &result.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, errors.New("USER_NOT_FOUND")
	}
	if err != nil {
		return result, err
	}
	return result, nil
}
