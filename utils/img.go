package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// ImageSizeField holds the data for a single resized image.
type ImageSizeField struct {
	Width  int
	Height int
	URL    string
}

// ImageSizes holds the data for all resized versions of an image.
type ImageSizes struct {
	Small  *ImageSizeField
	Medium *ImageSizeField
	Large  *ImageSizeField
}

// Define constants for image resizing
const (
	UploadPath  = "public/uploads"
	URLPrefix   = "/uploads"
	SmallWidth  = 320
	MediumWidth = 640
	LargeWidth  = 1024
)

// BlurImage generates a small, blurred base64-encoded data URI for a given image.
func BlurImage(src image.Image) (string, error) {
	if src == nil {
		return "", errors.New("source image is nil")
	}

	thumbnail := imaging.Thumbnail(src, 100, 100, imaging.Lanczos)
	blurImage := imaging.Blur(thumbnail, 5)

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

// saveResizedImage saves a resized image with a placeholder .avif extension.
func saveResizedImage(img image.Image, originalFilename, sizeName string, width int) (*ImageSizeField, error) {
	ext := filepath.Ext(originalFilename)
	baseName := strings.TrimSuffix(originalFilename, ext)
	// New filename with .avif extension
	newFilename := fmt.Sprintf("%s-%s.avif", baseName, sizeName)

	if err := os.MkdirAll(UploadPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	dstPath := filepath.Join(UploadPath, newFilename)

	// Resize image
	resizedImg := imaging.Resize(img, width, 0, imaging.Lanczos)

	// NOTE: This saves the image in JPEG format but with an .avif extension.
	// A real AVIF encoder (e.g., from a library like github.com/Kagami/go-avif)
	// should be used here for proper conversion.
	if err := imaging.Save(resizedImg, dstPath); err != nil {
		return nil, fmt.Errorf("failed to save (placeholder) AVIF image %s: %w", dstPath, err)
	}

	bounds := resizedImg.Bounds()
	urlPath := strings.ReplaceAll(filepath.ToSlash(filepath.Join(URLPrefix, newFilename)), "\\", "/")

	return &ImageSizeField{
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
		URL:    urlPath,
	}, nil
}

// ResizeImage takes an image, resizes it to different sizes, and saves them as placeholder AVIF files.
func ResizeImage(src image.Image, originalFilename string) (*ImageSizes, error) {
	if src == nil {
		return nil, errors.New("source image is nil")
	}
	if originalFilename == "" {
		return nil, errors.New("original filename is empty")
	}

	imageSizes := &ImageSizes{}
	var errorStrings []string

	// Small
	smallField, err := saveResizedImage(src, originalFilename, "small", SmallWidth)
	if err != nil {
		errorStrings = append(errorStrings, err.Error())
	}
	imageSizes.Small = smallField

	// Medium
	mediumField, err := saveResizedImage(src, originalFilename, "medium", MediumWidth)
	if err != nil {
		errorStrings = append(errorStrings, err.Error())
	}
	imageSizes.Medium = mediumField

	// Large
	largeField, err := saveResizedImage(src, originalFilename, "large", LargeWidth)
	if err != nil {
		errorStrings = append(errorStrings, err.Error())
	}
	imageSizes.Large = largeField

	if len(errorStrings) > 0 {
		return imageSizes, errors.New(strings.Join(errorStrings, "; "))
	}

	return imageSizes, nil
}
