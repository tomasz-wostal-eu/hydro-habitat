package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/tomasz-wostal-eu/hydro-habitat/internal/model"
)

// PostgresUserRepository implements the UserRepository interface for PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

// GetAll fetches all users from the database
func (r *PostgresUserRepository) GetAll(ctx context.Context) ([]model.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return users, nil
}

// GetByID fetches a user by ID from the database
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int) (model.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("user with id %d not found", id)
		}
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Create inserts a new user into the database
func (r *PostgresUserRepository) Create(ctx context.Context, user model.UserCreate) (int, error) {
	query := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(ctx, query, user.Name, user.Email).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

// Update updates an existing user in the database
func (r *PostgresUserRepository) Update(ctx context.Context, id int, updates model.UserUpdate) error {
	// First, get the current user to know what fields to update
	currentUser, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Use current values if updates doesn't provide new ones
	name := currentUser.Name
	if updates.Name != "" {
		name = updates.Name
	}

	email := currentUser.Email
	if updates.Email != "" {
		email = updates.Email
	}

	query := `
		UPDATE users
		SET name = $1, email = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, name, email, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

// Delete removes a user from the database
func (r *PostgresUserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}
