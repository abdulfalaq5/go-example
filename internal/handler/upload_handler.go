package handler

import (
	"net/http"
	"time"

	"github.com/falaqmsi/go-example/internal/service"
	"github.com/falaqmsi/go-example/pkg/response"
	"github.com/gin-gonic/gin"
)

// UploadHandler encapsulates the dependencies for HTTP upload handling.
type UploadHandler struct {
	svc service.UploadService
}

// NewUploadHandler returns a configured UploadHandler instance.
func NewUploadHandler(svc service.UploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

// RegisterRoutes sets up the upload route onto the router group.
// It applies the authentication middleware so uploads remain secure.
func (h *UploadHandler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// Upload routes are usually protected
	group := rg.Group("/upload", authMiddleware)
	{
		group.POST("", h.Upload)
	}
}

// Upload godoc
//
//	@Summary		Upload a file
//	@Description	Accepts a multipart form containing a file and uploads it to configured storage (Local or MinIO)
//	@Tags			Upload
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"The file to upload"
//	@Success		201		{object}	object{success=bool,message=string,data=object{file_url=string},meta=response.Meta}	"Upload successful"
//	@Failure		400		{object}	response.ErrorResponse	"Bad request"
//	@Failure		401		{object}	response.ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/api/v1/upload [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	// Gin extracts the file from the multipart form
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "No file is received or invalid form", err.Error())
		return
	}

	url, err := h.svc.UploadFile(c.Request.Context(), file)
	if err != nil {
		response.InternalServerError(c, "Failed to upload file", err.Error())
		return
	}

	// Format matching requirement: { "file_url": "string" }
	data := gin.H{
		"file_url": url,
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "File uploaded successfully",
		"data":    data,
		"meta":    gin.H{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}
