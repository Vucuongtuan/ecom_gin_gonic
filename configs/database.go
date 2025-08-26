package configs

import (
	"log"
	"os"

	"gorm.io/gorm"
)

// DB is the database connection instance.
var DB *gorm.DB

// ConnectDB initializes the database connection.
func ConnectDB() error {
	dsn := os.Getenv("DATABASE_URI")
	if dsn == "" {
		log.Fatal("DATABASE_URI environment variable not set")
	}

	// var err error
	var db *gorm.DB

	// Retry connection
	// maxRetries := 5
	// retryInterval := 5 * time.Second
	// for i := 0; i < maxRetries; i++ {
	// 	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 		PrepareStmt: true,
	// 	})
	// 	if err == nil {
	// 		break
	// 	}
	// 	log.Printf("Failed to connect to database. Retrying in %v... (%d/%d)", retryInterval, i+1, maxRetries)
	// 	time.Sleep(retryInterval)
	// }

	// if err != nil {
	// 	log.Fatalf("Failed to connect to database after %d retries: %v", maxRetries, err)
	// }

	// sqlDB, err := db.DB()
	// if err != nil {
	// 	log.Fatalf("Failed to get database instance: %v", err)
	// }

	// // Set connection pool settings
	// sqlDB.SetMaxIdleConns(10)
	// sqlDB.SetMaxOpenConns(100)
	// sqlDB.SetConnMaxLifetime(time.Hour)

	// // Ping the database to verify the connection
	// if err := sqlDB.Ping(); err != nil {
	// 	log.Fatalf("Failed to ping database: %v", err)
	// }

	DB = db
	log.Println("Database connection successful.")
	return nil
}
