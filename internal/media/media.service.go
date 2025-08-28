package media

import (
	"ecom_be/common"
	"ecom_be/configs"
	"ecom_be/utils"
	"errors"
	"mime/multipart"
	"net/http"

	// Import image formats for decoding
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"gorm.io/gorm"
)

// MediaService provides media-related services.
type MediaService struct {
	db *gorm.DB
}

// NewMediaService creates a new MediaService.
func NewMediaService() *MediaService {
	return &MediaService{
		db: configs.GetDB(), // get DB Collection form configs
	}
}

func (ms *MediaService) FindAll(p *common.Pagination) ([]Media, int64, int, error) {
	var medias []Media
	var total int64

	if err := ms.db.Model(&Media{}).Count(&total).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, err
	}

	if total == 0 {
		return nil, 0, http.StatusNotFound, errors.New("no media found")
	}

	if p == nil {
		p = &common.Pagination{
			Page:  common.DEFAULT_PAGE,
			Limit: common.DEFAULT_LIMIT,
		}
	}

	offset := (p.Page - 1) * p.Limit
	if err := ms.db.
		Limit(p.Limit).
		Offset(offset).
		Order("created_at desc").
		Find(&medias).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, err
	}

	if len(medias) == 0 {
		return nil, total, http.StatusNotFound, errors.New("page has no data")
	}

	return medias, total, http.StatusOK, nil
}

func (ms *MediaService) FindByID(id uint) (*Media, error) {
	var media Media
	if err := ms.db.First(&media, id).Error; err != nil {
		return nil, err
	}
	return &media, nil
}

func (ms *MediaService) Create(media *Media) (*Media, error) {
	if err := ms.db.Create(media).Error; err != nil {
		return nil, err
	}
	return media, nil
}

func (ms *MediaService) Update(id uint, updatedData map[string]interface{}) (*Media, error) {
	var media Media
	if err := ms.db.First(&media, id).Error; err != nil {
		return nil, err
	}

	if err := ms.db.Model(&media).Updates(updatedData).Error; err != nil {
		return nil, err
	}

	return &media, nil
}

func (ms *MediaService) Delete(id uint) error {
	if err := ms.db.Delete(&Media{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (ms *MediaService) UploadFile(file *multipart.FileHeader, altText string) (*Media, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Process and save the image versions
	processedResult, err := utils.ProcessAndSaveImage(src, file.Filename)
	if err != nil {
		return nil, err
	}

	// Create a new Media object
	media := &Media{
		FileName:     file.Filename, // Store original filename
		MimeType:     "image/webp",  // MimeType is now webp
		Size:         file.Size,     // This is the original size, might want to store the new total size
		URL:          processedResult.URL,
		ThumbnailURL: processedResult.ThumbnailURL,
		BlurURL:      processedResult.BlurURL,
		Alt:          altText,
		Width:        processedResult.Width,
		Height:       processedResult.Height,
	}

	// Save the media record to the database
	if err := ms.db.Create(media).Error; err != nil {
		// Here you might want to add cleanup logic to delete the saved files if DB write fails
		return nil, err
	}

	return media, nil
}
