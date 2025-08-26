package media

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type Media struct {
	*gorm.Model
	BlurURL      string    `gorm:"column:blur_url" json:"blur_url"`
	FileName     string    `gorm:"column:file_name" json:"file_name"`
	MimeType     string    `gorm:"column:mime_type" json:"mime_type"`
	Size         int64     `gorm:"column:size" json:"size"`
	URL          string    `gorm:"column:url" json:"url"`
	ThumbnailURL string    `gorm:"column:thumbnail_url" json:"thumbnail_url"`
	MediaSizes   MediaSize `gorm:"column:media_sizes;type:json" json:"media_sizes"`
	Alt          string    `gorm:"column:alt" json:"alt"`
}

type MediaSize struct {
	Small  *MediaSizeField `json:"small,omitempty"`
	Medium *MediaSizeField `json:"medium,omitempty"`
	Large  *MediaSizeField `json:"large,omitempty"`
}

// Scan - Implement the sql.Scanner interface for MediaSize
func (m *MediaSize) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, &m)
}

func (m MediaSize) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type MediaSizeField struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}
