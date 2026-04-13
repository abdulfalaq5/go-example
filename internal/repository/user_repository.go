package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/falaqmsi/go-example/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrUserNotFound is returned when a requested user does not exist.
var ErrUserNotFound = errors.New("user not found")

// UserRepository defines the data-access contract for the users table.
type UserRepository interface {
	FindAll(ctx context.Context) ([]model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, input model.CreateUserInput) (*model.User, error)
	Update(ctx context.Context, id int64, input model.UpdateUserInput) (*model.User, error)
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db *pgxpool.Pool // DBMain
}

// NewUserRepository creates a UserRepository backed by the provided pool.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// ── FindAll ───────────────────────────────────────────────────────────────────

func (r *userRepository) FindAll(ctx context.Context) ([]model.User, error) {
	const q = `
		SELECT id, name, email, created_at, updated_at
		FROM   users
		ORDER  BY id ASC`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("userRepository.FindAll: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("userRepository.FindAll scan: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("userRepository.FindAll rows: %w", err)
	}

	return users, nil
}

// ── FindByID ──────────────────────────────────────────────────────────────────

func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	const q = `
		SELECT id, name, email, created_at, updated_at
		FROM   users
		WHERE  id = $1`

	var u model.User
	err := r.db.QueryRow(ctx, q, id).
		Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByID: %w", err)
	}

	return &u, nil
}

// ── Create ────────────────────────────────────────────────────────────────────

func (r *userRepository) Create(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	const q = `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at, updated_at`

	var u model.User
	err := r.db.QueryRow(ctx, q, input.Name, input.Email).
		Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("userRepository.Create: %w", err)
	}

	return &u, nil
}

// ── Update ────────────────────────────────────────────────────────────────────

func (r *userRepository) Update(ctx context.Context, id int64, input model.UpdateUserInput) (*model.User, error) {
	const q = `
		UPDATE users
		SET    name       = $1,
		       email      = $2,
		       updated_at = NOW()
		WHERE  id = $3
		RETURNING id, name, email, created_at, updated_at`

	var u model.User
	err := r.db.QueryRow(ctx, q, input.Name, input.Email, id).
		Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepository.Update: %w", err)
	}

	return &u, nil
}

// ── Delete ────────────────────────────────────────────────────────────────────

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM users WHERE id = $1`

	tag, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("userRepository.Delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
