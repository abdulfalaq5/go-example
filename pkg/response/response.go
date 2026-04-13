package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Meta is attached to every response and may carry pagination info.
type Meta struct {
	Timestamp string `json:"timestamp"`
	Page      int    `json:"page,omitempty"`
	PageSize  int    `json:"page_size,omitempty"`
	TotalRows int    `json:"total_rows,omitempty"`
}

// newMeta returns a Meta pre-filled with the current UTC timestamp.
func newMeta() *Meta {
	return &Meta{Timestamp: time.Now().UTC().Format(time.RFC3339)}
}

// envelope is the shared JSON wrapper for every API response.
type envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta"`
	Error   any    `json:"error,omitempty"`
}

// Success writes a 200 OK JSON response.
func Success(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, envelope{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    newMeta(),
	})
}

// SuccessWithMeta writes a 200 OK response with caller-supplied meta
// (pagination, total_rows, etc.). The timestamp is always injected.
func SuccessWithMeta(c *gin.Context, message string, data any, meta *Meta) {
	if meta == nil {
		meta = newMeta()
	} else {
		meta.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
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
		Meta:    newMeta(),
	})
}

// Error writes a JSON error response with the given HTTP status code.
func Error(c *gin.Context, statusCode int, message string, err any) {
	c.JSON(statusCode, envelope{
		Success: false,
		Message: message,
		Meta:    newMeta(),
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
