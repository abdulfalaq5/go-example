package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/falaqmsi/go-example/internal/model"
	"github.com/falaqmsi/go-example/internal/repository"
)

// UserService defines the business-logic contract for user operations.
type UserService interface {
	GetAll(ctx context.Context) ([]model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, input model.CreateUserInput) (*model.User, error)
	Update(ctx context.Context, id int64, input model.UpdateUserInput) (*model.User, error)
	Delete(ctx context.Context, id int64) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService wires the service to its repository dependency.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// GetAll returns every user ordered by ID ascending.
func (s *userService) GetAll(ctx context.Context) ([]model.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("userService.GetAll: %w", err)
	}
	// Always return a non-nil slice so the JSON encodes as [] not null.
	if users == nil {
		users = []model.User{}
	}
	return users, nil
}

// GetByID returns a single user or ErrUserNotFound if it does not exist.
func (s *userService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("userService.GetByID: %w", err)
	}
	return user, nil
}

// Create inserts a new user and returns the persisted record.
func (s *userService) Create(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	user, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("userService.Create: %w", err)
	}
	return user, nil
}

// Update modifies an existing user's fields and returns the updated record.
func (s *userService) Update(ctx context.Context, id int64, input model.UpdateUserInput) (*model.User, error) {
	user, err := s.repo.Update(ctx, id, input)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("userService.Update: %w", err)
	}
	return user, nil
}

// Delete removes a user by ID.
func (s *userService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return repository.ErrUserNotFound
		}
		return fmt.Errorf("userService.Delete: %w", err)
	}
	return nil
}
