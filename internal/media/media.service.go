package media

import (
	"ecom_be/utils"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	// Import image formats for decoding
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// MediaService provides media-related services.
type MediaService struct {
	repo *MediaRepo
}

// NewMediaService creates a new MediaService.
func NewMediaService(repo *MediaRepo) *MediaService {
	return &MediaService{repo: repo}
}

// saveOriginalFile saves the uploaded file as is and returns its public URL.
func saveOriginalFile(file *multipart.FileHeader) (string, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Ensure upload directory exists
	if err := os.MkdirAll(utils.UploadPath, 0755); err != nil {
		return "", err
	}

	// Create destination file
	dstPath := filepath.Join(utils.UploadPath, file.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// Return the public URL
	urlPath := strings.ReplaceAll(filepath.ToSlash(filepath.Join(utils.URLPrefix, file.Filename)), "\"", "/")
	return urlPath, nil
}

func (s *MediaService) UploadFile(file *multipart.FileHeader, altText string) (*Media, error) {
	// 1. Save the original file first to get the main URL
	originalURL, err := saveOriginalFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to save original file: %w", err)
	}

	// 2. Open and decode the file for image processing
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file for processing: %w", err)
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image '%s': %w", file.Filename, err)
	}

	// 3. Generate blurred placeholder and resized (placeholder) AVIF versions
	blurData, err := utils.BlurImage(img)
	if err != nil {
		fmt.Printf("Warning: could not generate blur image for %s: %v\n", file.Filename, err)
	}

	utilMediaSizes, err := utils.ResizeImage(img, file.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to resize image '%s': %w", file.Filename, err)
	}

	// 4. Convert utils.ImageSizes to media.MediaSize
	mediaSize := MediaSize{}
	if utilMediaSizes.Small != nil {
		mediaSize.Small = &MediaSizeField{
			Width:  utilMediaSizes.Small.Width,
			Height: utilMediaSizes.Small.Height,
			URL:    utilMediaSizes.Small.URL,
		}
	}
	if utilMediaSizes.Medium != nil {
		mediaSize.Medium = &MediaSizeField{
			Width:  utilMediaSizes.Medium.Width,
			Height: utilMediaSizes.Medium.Height,
			URL:    utilMediaSizes.Medium.URL,
		}
	}
	if utilMediaSizes.Large != nil {
		mediaSize.Large = &MediaSizeField{
			Width:  utilMediaSizes.Large.Width,
			Height: utilMediaSizes.Large.Height,
			URL:    utilMediaSizes.Large.URL,
		}
	}

	// 5. Prepare the Media struct with the new URL structure
	thumbnailURL := ""
	if mediaSize.Small != nil {
		thumbnailURL = mediaSize.Small.URL // This is now a .avif url
	}

	media := &Media{
		FileName:     file.Filename,
		MimeType:     file.Header.Get("Content-Type"),
		Size:         file.Size,
		URL:          originalURL,  // URL of the original image
		ThumbnailURL: thumbnailURL, // URL of the small .avif image
		BlurURL:      blurData,
		MediaSizes:   mediaSize, // Contains .avif urls
		Alt:          altText,
	}

	// 6. Save the media metadata to the database
	createdMedia, err := s.repo.CreateMedia(media)
	if err != nil {
		return nil, fmt.Errorf("failed to save media to database: %w", err)
	}

	return createdMedia, nil
}
