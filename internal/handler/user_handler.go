package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/falaqmsi/go-example/internal/model"
	"github.com/falaqmsi/go-example/internal/repository"
	"github.com/falaqmsi/go-example/internal/service"
	"github.com/falaqmsi/go-example/pkg/response"
	"github.com/gin-gonic/gin"
)

// UserHandler handles all HTTP requests for the /users resource.
type UserHandler struct {
	svc service.UserService
}

// NewUserHandler creates a UserHandler with its service dependency.
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// RegisterRoutes attaches all user routes onto the given router group and secures them.
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	users := rg.Group("/users", authMiddleware)
	{
		users.GET("", h.GetAll)
		users.POST("", h.Create)
		users.GET("/:id", h.GetByID)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
	}
}

// ── GetAll ────────────────────────────────────────────────────────────────────

// GetAll godoc
//
//	@Summary		List all users
//	@Description	Returns a list of all users ordered by ID ascending.
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	object{success=bool,message=string,data=[]model.User,meta=response.Meta}	"List of users"
//	@Failure		500	{object}	response.ErrorResponse	"Internal server error"
//	@Router			/api/v1/users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "failed to retrieve users", err.Error())
		return
	}
	response.Success(c, "users retrieved successfully", users)
}

// ── GetByID ───────────────────────────────────────────────────────────────────

// GetByID godoc
//
//	@Summary		Get a user by ID
//	@Description	Returns a single user record identified by its numeric ID.
//	@Tags			Users
//	@Produce		json
//	@Param			id	path		int	true	"User ID"	minimum(1)
//	@Success		200	{object}	object{success=bool,message=string,data=model.User,meta=response.Meta}	"User found"
//	@Failure		400	{object}	response.ErrorResponse	"Invalid ID"
//	@Failure		404	{object}	response.ErrorResponse	"User not found"
//	@Failure		500	{object}	response.ErrorResponse	"Internal server error"
//	@Router			/api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	user, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalServerError(c, "failed to retrieve user", err.Error())
		return
	}

	response.Success(c, "user retrieved successfully", user)
}

// ── Create ────────────────────────────────────────────────────────────────────

// Create godoc
//
//	@Summary		Create a new user
//	@Description	Inserts a new user record. Email must be unique.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		model.CreateUserInput	true	"User payload"
//	@Success		201		{object}	object{success=bool,message=string,data=model.User,meta=response.Meta}	"User created"
//	@Failure		400		{object}	response.ErrorResponse	"Validation error"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var input model.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	user, err := h.svc.Create(c.Request.Context(), input)
	if err != nil {
		response.InternalServerError(c, "failed to create user", err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "user created successfully",
		"data":    user,
		"meta":    gin.H{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

// ── Update ────────────────────────────────────────────────────────────────────

// Update godoc
//
//	@Summary		Update an existing user
//	@Description	Updates the name and email of an existing user by ID.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"User ID"	minimum(1)
//	@Param			body	body		model.UpdateUserInput	true	"Update payload"
//	@Success		200		{object}	object{success=bool,message=string,data=model.User,meta=response.Meta}	"User updated"
//	@Failure		400		{object}	response.ErrorResponse	"Validation error or invalid ID"
//	@Failure		404		{object}	response.ErrorResponse	"User not found"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	var input model.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	user, err := h.svc.Update(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalServerError(c, "failed to update user", err.Error())
		return
	}

	response.Success(c, "user updated successfully", user)
}

// ── Delete ────────────────────────────────────────────────────────────────────

// Delete godoc
//
//	@Summary		Delete a user
//	@Description	Permanently removes a user record by ID.
//	@Tags			Users
//	@Produce		json
//	@Param			id	path		int	true	"User ID"	minimum(1)
//	@Success		200	{object}	object{success=bool,message=string,meta=response.Meta}	"User deleted"
//	@Failure		400	{object}	response.ErrorResponse	"Invalid ID"
//	@Failure		404	{object}	response.ErrorResponse	"User not found"
//	@Failure		500	{object}	response.ErrorResponse	"Internal server error"
//	@Router			/api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalServerError(c, "failed to delete user", err.Error())
		return
	}

	response.Success(c, "user deleted successfully", nil)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// parseID reads and validates the :id path param.
// On invalid input it writes a 400 and returns a non-nil error so the caller
// can do an early return.
func parseID(c *gin.Context) (int64, error) {
	raw := c.Param("id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "invalid id: must be a positive integer", nil)
		return 0, errors.New("invalid id")
	}
	return id, nil
}
