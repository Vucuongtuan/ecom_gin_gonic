package media

import "gorm.io/gorm"

// MediaRepo handles database operations for Media entities.
type MediaRepo struct {
	db *gorm.DB
}

// NewMediaRepo creates a new instance of MediaRepo.
func NewMediaRepo(db *gorm.DB) *MediaRepo {
	return &MediaRepo{db: db}
}

// CreateMedia inserts a new Media record into the database.
func (r *MediaRepo) CreateMedia(media *Media) (*Media, error) {
	if err := r.db.Create(media).Error; err != nil {
		return nil, err
	}
	return media, nil
}