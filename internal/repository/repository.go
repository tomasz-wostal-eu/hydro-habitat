package repository

import (
	"context"

	"github.com/tomasz-wostal-eu/hydro-habitat/internal/model"
)

// UserRepository defines operations for working with users
type UserRepository interface {
	// GetAll returns all users
	GetAll(ctx context.Context) ([]model.User, error)

	// GetByID returns a user by ID
	GetByID(ctx context.Context, id int) (model.User, error)

	// Create creates a new user
	Create(ctx context.Context, user model.UserCreate) (int, error)

	// Update updates an existing user
	Update(ctx context.Context, id int, user model.UserUpdate) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id int) error
}