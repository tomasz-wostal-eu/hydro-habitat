package service

import (
	"context"
	"errors"
	"log"

	"github.com/tomasz-wostal-eu/hydro-habitat/internal/model"
	"github.com/tomasz-wostal-eu/hydro-habitat/internal/repository"
)

// UserService defines the user service operations
type UserService struct {
	repo   repository.UserRepository
	logger *log.Logger
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository, logger *log.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

// GetAllUsers fetches all users
func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Printf("Error fetching all users: %v", err)
		return nil, err
	}
	return users, nil
}

// GetUserByID fetches a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id int) (model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Printf("Error fetching user with ID %d: %v", id, err)
		return model.User{}, err
	}
	return user, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, userData model.UserCreate) (int, error) {
	// Basic validation
	if userData.Name == "" {
		return 0, errors.New("name is required")
	}
	if userData.Email == "" {
		return 0, errors.New("email is required")
	}

	id, err := s.repo.Create(ctx, userData)
	if err != nil {
		s.logger.Printf("Error creating user: %v", err)
		return 0, err
	}
	return id, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, id int, userData model.UserUpdate) error {
	err := s.repo.Update(ctx, id, userData)
	if err != nil {
		s.logger.Printf("Error updating user with ID %d: %v", id, err)
		return err
	}
	return nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Printf("Error deleting user with ID %d: %v", id, err)
		return err
	}
	return nil
}
