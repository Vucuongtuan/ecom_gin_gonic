package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

// Define constants for image processing
const (
	UploadPath     = "public/uploads"
	MainWidth      = 1024 // Max width for the main image
	ThumbnailWidth = 300  // Width for the thumbnail
)

// ProcessedImageResult holds the data after processing an uploaded image.
type ProcessedImageResult struct {
	BlurURL      string
	URL          string
	ThumbnailURL string
	Width        int
	Height       int
}

// BlurImage generates a small, blurred base64-encoded data URI for a given image.
// It uses the provided width and height if specified, otherwise defaults to a small thumbnail.
func BlurImage(src image.Image, width int, height *int) (string, error) {
	if src == nil {
		return "", errors.New("source image is nil")
	}

	var blurWidth, blurHeight int
	if height != nil {
		blurWidth = width
		blurHeight = *height
	} else {
		bounds := src.Bounds()
		origW := bounds.Dx()
		origH := bounds.Dy()
		blurWidth = width
		blurHeight = int(float64(origH) * float64(blurWidth) / float64(origW))
		if blurHeight < 1 {
			blurHeight = 1
		}
	}

	thumbnail := imaging.Resize(src, blurWidth, blurHeight, imaging.Lanczos)
	blurImage := imaging.Blur(thumbnail, 30)

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, blurImage, &jpeg.Options{Quality: 30}); err != nil {
		return "", fmt.Errorf("failed to encode blurred image: %w", err)
	}

	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	if base64Str == "" {
		return "", errors.New("generated base64 string is empty")
	}

	return "data:image/jpeg;base64," + base64Str, nil
}

func ProcessAndSaveImage(src io.Reader, originalFilename string) (*ProcessedImageResult, error) {
	// Decode the image
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	blurURL, err := BlurImage(img, MainWidth, nil)
	if err != nil {
		// Log the error but don't fail the whole process
		fmt.Printf("Warning: failed to generate blur image for %s: %v\n", originalFilename, err)
	}

	ext := filepath.Ext(originalFilename)
	baseName := strings.TrimSuffix(filepath.Base(originalFilename), ext)
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate uuid: %w", err)
	}
	uniqueBaseName := fmt.Sprintf("%s-%s", baseName, strings.Replace(uuid.String(), "-", "", -1))

	if err := os.MkdirAll(UploadPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	mainImg := imaging.Resize(img, MainWidth, 0, imaging.Lanczos)
	mainFilename := uniqueBaseName + ".jpg"
	mainPath := filepath.Join(UploadPath, mainFilename)

	mainFile, err := os.Create(mainPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file for main image: %w", err)
	}
	defer mainFile.Close()

	if err := jpeg.Encode(mainFile, mainImg, &jpeg.Options{Quality: 80}); err != nil {
		return nil, fmt.Errorf("failed to encode main image to jpeg: %w", err)
	}

	thumbImg := imaging.Resize(img, ThumbnailWidth, 0, imaging.Lanczos)
	thumbFilename := uniqueBaseName + "_thumb.jpg"
	thumbPath := filepath.Join(UploadPath, thumbFilename)

	thumbFile, err := os.Create(thumbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file for thumbnail image: %w", err)
	}
	defer thumbFile.Close()

	if err := jpeg.Encode(thumbFile, thumbImg, &jpeg.Options{Quality: 80}); err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail to jpeg: %w", err)
	}

	bounds := mainImg.Bounds()

	return &ProcessedImageResult{
		BlurURL:      blurURL,
		URL:          "/" + strings.ReplaceAll(mainPath, "\\", "/"),
		ThumbnailURL: "/" + strings.ReplaceAll(thumbPath, "\\", "/"),
		Width:        bounds.Dx(),
		Height:       bounds.Dy(),
	}, nil
}
