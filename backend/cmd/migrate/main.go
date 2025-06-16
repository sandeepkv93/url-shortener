package main

import (
	"log"
	"os"

	"url-shortener/internal/config"
	"url-shortener/internal/infrastructure/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create indexes
	if err := db.CreateIndexes(); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	log.Println("Database migration completed successfully")
}

func init() {
	// Ensure we're in the correct directory
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		log.Fatal("go.mod not found. Please run this command from the backend directory.")
	}
}