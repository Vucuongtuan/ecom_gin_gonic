package media

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MediaController handles HTTP requests for media.
type MediaController struct {
	service *MediaService
}

// NewMediaController creates a new MediaController.
func NewMediaController(service *MediaService) *MediaController {
	return &MediaController{service: service}
}

// RegisterRoutes registers the media routes to a gin router group.
// Example usage in main.go:
//
//	mediaController.RegisterRoutes(router.Group("/api/v1"))
func (c *MediaController) RegisterRoutes(router *gin.RouterGroup) {
	mediaRoutes := router.Group("/media")
	{
		mediaRoutes.POST("/upload", c.Upload)
	}
}

// Upload handles multipart file uploads.

func (c *MediaController) Upload(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form: " + err.Error()})
		return
	}

	files := form.File["files"] // Key for files
	alts := form.Value["alts"]  // Key for alt texts

	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No files were uploaded"})
		return
	}

	// Optional: Check if the number of alts matches the number of files
	if len(alts) > 0 && len(alts) != len(files) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "The number of alt texts does not match the number of files"})
		return
	}

	var uploadedMedia []*Media
	for i, file := range files {
		var altText string
		if i < len(alts) {
			altText = alts[i]
		}

		media, err := c.service.UploadFile(file, altText)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to upload file: " + file.Filename,
				"details": err.Error(),
			})
			return
		}
		uploadedMedia = append(uploadedMedia, media)
	}

	ctx.JSON(http.StatusCreated, uploadedMedia)
}
