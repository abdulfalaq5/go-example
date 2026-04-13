package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Meta holds pagination or extra context metadata attached to a response.
type Meta struct {
	Page      int `json:"page,omitempty"`
	PageSize  int `json:"page_size,omitempty"`
	TotalRows int `json:"total_rows,omitempty"`
}

// envelope is the shared JSON wrapper for every API response.
type envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
	Error   any    `json:"error,omitempty"`
}

// Success writes a 200 OK (or custom code) JSON response.
func Success(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, envelope{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta writes a successful response that includes pagination metadata.
func SuccessWithMeta(c *gin.Context, message string, data any, meta *Meta) {
	c.JSON(http.StatusOK, envelope{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Created writes a 201 Created JSON response.
func Created(c *gin.Context, message string, data any) {
	c.JSON(http.StatusCreated, envelope{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error writes a JSON error response with the supplied HTTP status code.
func Error(c *gin.Context, statusCode int, message string, err any) {
	c.JSON(statusCode, envelope{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// BadRequest is a convenience wrapper for 400 responses.
func BadRequest(c *gin.Context, message string, err any) {
	Error(c, http.StatusBadRequest, message, err)
}

// Unauthorized is a convenience wrapper for 401 responses.
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

// Forbidden is a convenience wrapper for 403 responses.
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, nil)
}

// NotFound is a convenience wrapper for 404 responses.
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}

// InternalServerError is a convenience wrapper for 500 responses.
func InternalServerError(c *gin.Context, message string, err any) {
	Error(c, http.StatusInternalServerError, message, err)
}
