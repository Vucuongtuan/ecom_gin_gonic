package media

import (
	"ecom_be/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MediaController struct {
	service *MediaService
}

func NewMediaController(service *MediaService) *MediaController {
	return &MediaController{service: service}
}

func (mc *MediaController) GetAllMedia(c *gin.Context) {
	p := common.GetPagination(c)
	listMedia, total, status, err := mc.service.FindAll(p)
	if err != nil || listMedia == nil {
		common.PaginatedResponse(c, status, []Media{}, 0, p)
		return
	}
	common.PaginatedResponse(c, status, listMedia, int(total), p)
	// c.JSON(http.StatusOK, gin.H{"message": "Get all media - Not implemented yet"})
}

func (c *MediaController) UploadSingle(ctx *gin.Context) {
	file, err := ctx.FormFile("file") // Key for single file
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file was uploaded or invalid key: " + err.Error()})
		return
	}

	altText := ctx.PostForm("alt") // Key for single alt text

	media, err := c.service.UploadFile(file, altText)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to upload file: " + file.Filename,
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, media)
}

func (c *MediaController) UploadMultiple(ctx *gin.Context) {
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
