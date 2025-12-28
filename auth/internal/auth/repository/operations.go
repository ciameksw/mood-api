package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ciameksw/mood-api/pkg/postgres"
)

type DBOperations struct {
	Postgres *postgres.PostgresDB
}

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

// CreateUser inserts a new user into the database
func (o *DBOperations) CreateUser(ctx context.Context, username, email, passwordHash string) (int, error) {
	var userID int
	query := "INSERT INTO users (username, email, password_hash, created_at) VALUES ($1, $2, $3, $4) RETURNING id"

	err := o.Postgres.DB.QueryRowContext(ctx, query, username, email, passwordHash, time.Now()).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUserByEmail retrieves a user by email
func (o *DBOperations) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	query := "SELECT id, username, email, password_hash, created_at FROM users WHERE email = $1"

	err := o.Postgres.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// UserExistsByEmail checks if a user with the given email exists
func (o *DBOperations) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"

	err := o.Postgres.DB.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// UserExistsByUsername checks if a user with the given username exists
func (o *DBOperations) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"

	err := o.Postgres.DB.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetUserByID retrieves a user by ID
func (o *DBOperations) GetUserByID(ctx context.Context, userID int) (*User, error) {
	user := &User{}
	query := "SELECT id, username, email, password_hash, created_at FROM users WHERE id = $1"

	err := o.Postgres.DB.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (o *DBOperations) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}
	query := "SELECT id, username, email, password_hash, created_at FROM users WHERE username = $1"

	err := o.Postgres.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// UpdateUser updates user profile data
func (o *DBOperations) UpdateUser(ctx context.Context, userID int, username, email string, passwordHash *string) error {
	query := "UPDATE users SET "
	args := []interface{}{}
	argIndex := 1

	updates := []string{}

	if username != "" {
		updates = append(updates, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, username)
		argIndex++
	}

	if email != "" {
		updates = append(updates, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, email)
		argIndex++
	}

	if passwordHash != nil {
		updates = append(updates, fmt.Sprintf("password_hash = $%d", argIndex))
		args = append(args, *passwordHash)
		argIndex++
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	query += buildUpdateQuery(updates, argIndex)
	args = append(args, userID)

	_, err := o.Postgres.DB.ExecContext(ctx, query, args...)
	return err
}

// DeleteUser deletes a user from the database
func (o *DBOperations) DeleteUser(ctx context.Context, userID int) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := o.Postgres.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Helper function to build UPDATE query
func buildUpdateQuery(updates []string, nextArgIndex int) string {
	query := ""
	for i, update := range updates {
		query += update
		if i < len(updates)-1 {
			query += ", "
		}
	}
	query += fmt.Sprintf(" WHERE id = $%d", nextArgIndex)
	return query
}
