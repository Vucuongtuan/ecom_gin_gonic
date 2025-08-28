package configs

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	connStr := os.Getenv("STR_SQL")
	if connStr == "" {
		panic("STR_SQL environment variable is not set")
	}

	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
			PrepareStmt: true,
		})

		if err == nil {
			break
		}

		fmt.Printf("Failed to connect to database (attempt %d/%d): %v\n", i+1, maxRetries, err)
		if i < maxRetries-1 {
			fmt.Println("Retrying in 5 seconds...")
			time.Sleep(time.Second * 5)
		}
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database after %d retries: %v", maxRetries, err))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get database instance: %v", err))
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = sqlDB.Ping()
	if err != nil {
		panic(fmt.Sprintf("Failed to ping database: %v", err))
	}
	fmt.Println("Database connection established successfully🚀")
}

func GetDB() *gorm.DB {
	if DB == nil {
		panic("Database connection not initialized. Make sure to call configs.InitDB() at startup.")
	}
	return DB
}
