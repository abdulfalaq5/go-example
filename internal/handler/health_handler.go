package handler

import (
	"github.com/falaqmsi/go-example/internal/service"
	"github.com/falaqmsi/go-example/pkg/response"
	"github.com/gin-gonic/gin"
)

// HealthHandler handles HTTP requests for the health-check endpoint.
type HealthHandler struct {
	svc service.HealthService
	env string
}

// NewHealthHandler creates a HealthHandler with its dependencies.
func NewHealthHandler(svc service.HealthService, env string) *HealthHandler {
	return &HealthHandler{svc: svc, env: env}
}

// Check godoc
//
//	@Summary		Health check
//	@Description	Returns the current operational status of the API.
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	model.HealthStatus
//	@Router			/health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	status, err := h.svc.Check(h.env)
	if err != nil {
		response.InternalServerError(c, "service unavailable", err.Error())
		return
	}
	response.Success(c, "server is healthy", status)
}
