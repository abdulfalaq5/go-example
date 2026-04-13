package model

import "time"

// User represents a system user stored in the main database.
// @Description Full user record returned by the API.
type User struct {
	ID        int64     `json:"id"         example:"1"`
	Name      string    `json:"name"       example:"Alice"`
	Email     string    `json:"email"      example:"alice@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2026-04-13T08:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-04-13T08:00:00Z"`
}

// CreateUserInput is the validated payload used when inserting a new user.
// @Description Request body for creating a user.
type CreateUserInput struct {
	Name  string `json:"name"  binding:"required,min=2,max=100" example:"Alice"`
	Email string `json:"email" binding:"required,email"          example:"alice@example.com"`
}

// UpdateUserInput carries the fields that may be updated on an existing user.
// @Description Request body for updating a user.
type UpdateUserInput struct {
	Name  string `json:"name"  binding:"required,min=2,max=100" example:"Alice Updated"`
	Email string `json:"email" binding:"required,email"          example:"alice.updated@example.com"`
}
