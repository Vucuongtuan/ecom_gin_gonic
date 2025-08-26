package internal

import (
	"ecom_be/configs"
	"ecom_be/internal/media"
)

// Func auto migrate database
func Migrate() {
	configs.DB.AutoMigrate(
		&media.Media{},
	)
}
