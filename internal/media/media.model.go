package media

import (
	"gorm.io/gorm"
)

type Media struct {
	*gorm.Model
	BlurURL      string `gorm:"column:blur_url" json:"blur_url"`
	FileName     string `gorm:"column:file_name" json:"file_name"`
	MimeType     string `gorm:"column:mime_type" json:"mime_type"`
	Size         int64  `gorm:"column:size" json:"size"`
	URL          string `gorm:"column:url" json:"url"`
	ThumbnailURL string `gorm:"column:thumbnail_url" json:"thumbnail_url"`
	Alt          string `gorm:"column:alt" json:"alt"`
	Width        int    `gorm:"column:width" json:"width"`
	Height       int    `gorm:"column:height" json:"height"`
}
