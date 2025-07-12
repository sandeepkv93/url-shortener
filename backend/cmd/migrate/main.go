package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/core/services"
	"url-shortener/internal/infrastructure/database"
	"url-shortener/internal/infrastructure/database/migrations"
)

const (
	ExitSuccess = 0
	ExitError   = 1
)

func main() {
	var (
		configFlag = flag.String("config", "", "Path to configuration file")
		envFlag    = flag.String("env", "", "Environment to use (development, staging, production)")
		timeoutFlag = flag.Int("timeout", 300, "Migration timeout in seconds")
	)
	flag.Parse()

	if len(flag.Args()) == 0 {
		printUsage()
		os.Exit(ExitError)
	}

	command := flag.Args()[0]

	// Load configuration
	cfg, err := loadConfig(*configFlag, *envFlag)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override timeout if provided
	if *timeoutFlag > 0 {
		cfg.Migration.Timeout = time.Duration(*timeoutFlag) * time.Second
	}

	// Initialize logger
	logger := services.NewLoggingService(&services.LoggingConfig{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
	})

	// Connect to database
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize migration service
	migrationService := migrations.NewMigrationService(db.DB, cfg, logger)
	if err := migrationService.Initialize(); err != nil {
		log.Fatalf("Failed to initialize migration service: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetMigrationTimeout())
	defer cancel()

	// Execute command
	switch command {
	case "up", "migrate":
		err = runMigrations(ctx, migrationService)
	case "down", "rollback":
		steps := 1
		if len(flag.Args()) > 1 {
			if steps, err = strconv.Atoi(flag.Args()[1]); err != nil {
				log.Fatalf("Invalid rollback steps: %v", err)
			}
		}
		err = rollbackMigrations(ctx, migrationService, steps)
	case "status":
		err = showMigrationStatus(migrationService)
	case "validate":
		err = validateMigrations(migrationService)
	case "backup":
		err = createBackup(ctx, migrationService)
	case "create":
		if len(flag.Args()) < 2 {
			log.Fatal("Migration name is required for create command")
		}
		migrationName := flag.Args()[1]
		err = createMigration(cfg, migrationName)
	case "version":
		err = showVersion(migrationService)
	case "reset":
		err = resetDatabase(ctx, migrationService)
	case "legacy":
		err = runLegacyMigration(db)
	case "help", "-h", "--help":
		printUsage()
		os.Exit(ExitSuccess)
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(ExitError)
	}

	if err != nil {
		log.Fatalf("Command failed: %v", err)
	}

	fmt.Println("✅ Command completed successfully")
}

func loadConfig(configPath, env string) (*config.Config, error) {
	// Set environment if provided
	if env != "" {
		os.Setenv("GO_ENV", env)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Load environment-specific config if available
	if err := cfg.LoadEnvironmentSpecificConfig(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func runMigrations(ctx context.Context, migrationService *migrations.MigrationService) error {
	fmt.Println("🚀 Running database migrations...")
	
	// Validate migrations first
	if err := migrationService.ValidateMigrations(); err != nil {
		return fmt.Errorf("migration validation failed: %w", err)
	}

	// Run migrations
	if err := migrationService.Migrate(ctx); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("✅ All migrations completed successfully")
	return nil
}

func rollbackMigrations(ctx context.Context, migrationService *migrations.MigrationService, steps int) error {
	fmt.Printf("🔄 Rolling back %d migration(s)...\n", steps)
	
	if err := migrationService.Rollback(ctx, steps); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	fmt.Printf("✅ Successfully rolled back %d migration(s)\n", steps)
	return nil
}

func showMigrationStatus(migrationService *migrations.MigrationService) error {
	fmt.Println("📊 Migration Status:")
	fmt.Println("==================")

	statuses, err := migrationService.GetMigrationStatus()
	if err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	if len(statuses) == 0 {
		fmt.Println("No migrations found")
		return nil
	}

	fmt.Printf("%-20s %-40s %-10s %-20s %-15s\n", "Version", "Name", "Applied", "Applied At", "Execution Time")
	fmt.Printf("%-20s %-40s %-10s %-20s %-15s\n", "-------", "----", "-------", "----------", "--------------")

	for _, status := range statuses {
		appliedStr := "❌ No"
		appliedAtStr := "-"
		executionTimeStr := "-"

		if status.Applied {
			appliedStr = "✅ Yes"
			if status.AppliedAt != nil {
				appliedAtStr = status.AppliedAt.Format("2006-01-02 15:04:05")
			}
			if status.ExecutionTime > 0 {
				executionTimeStr = fmt.Sprintf("%dms", status.ExecutionTime)
			}
		}

		fmt.Printf("%-20s %-40s %-10s %-20s %-15s\n",
			status.Version,
			truncateString(status.Name, 40),
			appliedStr,
			appliedAtStr,
			executionTimeStr,
		)
	}

	return nil
}

func validateMigrations(migrationService *migrations.MigrationService) error {
	fmt.Println("🔍 Validating migrations...")
	
	if err := migrationService.ValidateMigrations(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Println("✅ All migrations are valid")
	return nil
}

func createBackup(ctx context.Context, migrationService *migrations.MigrationService) error {
	fmt.Println("💾 Creating database backup...")
	
	backupInfo, err := migrationService.CreateBackup(ctx)
	if err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	fmt.Printf("✅ Backup created successfully:\n")
	fmt.Printf("   File: %s\n", backupInfo.Filename)
	fmt.Printf("   Path: %s\n", backupInfo.Path)
	fmt.Printf("   Size: %d bytes\n", backupInfo.Size)
	fmt.Printf("   Created: %s\n", backupInfo.CreatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

func createMigration(cfg *config.Config, name string) error {
	fmt.Printf("📝 Creating new migration: %s\n", name)
	
	// Generate timestamp-based version
	version := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.sql", version, name)
	
	migrationPath := cfg.GetMigrationPath()
	filePath := fmt.Sprintf("%s/%s", migrationPath, filename)
	
	// Create migration template
	template := fmt.Sprintf(`-- +migrate Up
-- %s

-- Add your UP migration SQL here


-- +migrate Down
-- Rollback for %s

-- Add your DOWN migration SQL here

`, name, name)

	if err := os.WriteFile(filePath, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	fmt.Printf("✅ Migration file created: %s\n", filePath)
	fmt.Printf("   Edit the file to add your migration SQL\n")
	
	return nil
}

func showVersion(migrationService *migrations.MigrationService) error {
	fmt.Println("📋 Database Version Information:")
	fmt.Println("===============================")

	statuses, err := migrationService.GetMigrationStatus()
	if err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	// Find latest applied migration
	var latestVersion string = "0"
	var latestAppliedAt *time.Time
	appliedCount := 0
	totalCount := len(statuses)

	for _, status := range statuses {
		if status.Applied {
			appliedCount++
			latestVersion = status.Version
			latestAppliedAt = status.AppliedAt
		}
	}

	fmt.Printf("Current Version: %s\n", latestVersion)
	if latestAppliedAt != nil {
		fmt.Printf("Last Migration: %s\n", latestAppliedAt.Format("2006-01-02 15:04:05"))
	}
	fmt.Printf("Applied Migrations: %d/%d\n", appliedCount, totalCount)
	
	if appliedCount < totalCount {
		fmt.Printf("⚠️  Pending Migrations: %d\n", totalCount-appliedCount)
	} else {
		fmt.Println("✅ Database is up to date")
	}

	return nil
}

func resetDatabase(ctx context.Context, migrationService *migrations.MigrationService) error {
	fmt.Println("⚠️  WARNING: This will reset the entire database!")
	fmt.Print("Are you sure you want to continue? (yes/no): ")
	
	var response string
	fmt.Scanln(&response)
	
	if response != "yes" {
		fmt.Println("❌ Database reset cancelled")
		return nil
	}

	fmt.Println("🔄 Resetting database...")
	
	// This would implement database reset functionality
	// For safety, this is not implemented in this example
	fmt.Println("❌ Database reset functionality is disabled for safety")
	fmt.Println("   To reset the database, manually drop and recreate it")
	
	return nil
}

func runLegacyMigration(db *database.Database) error {
	fmt.Println("🔄 Running legacy GORM auto-migration...")
	
	// Run the old GORM auto-migration for backward compatibility
	if err := db.AutoMigrate(); err != nil {
		return fmt.Errorf("failed to run auto migrations: %w", err)
	}

	// Create indexes
	if err := db.CreateIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	fmt.Println("✅ Legacy migration completed successfully")
	return nil
}

func printUsage() {
	fmt.Println("URL Shortener Migration Tool")
	fmt.Println("===========================")
	fmt.Println()
	fmt.Println("Usage: migrate [options] <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up, migrate           Run all pending migrations")
	fmt.Println("  down, rollback [n]    Rollback n migrations (default: 1)")
	fmt.Println("  status                Show migration status")
	fmt.Println("  validate              Validate all migration files")
	fmt.Println("  backup                Create database backup")
	fmt.Println("  create <name>         Create new migration file")
	fmt.Println("  version               Show current database version")
	fmt.Println("  reset                 Reset database (development only)")
	fmt.Println("  legacy                Run legacy GORM auto-migration")
	fmt.Println("  help                  Show this help message")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -config string        Path to configuration file")
	fmt.Println("  -env string          Environment (development, staging, production)")
	fmt.Println("  -timeout int         Migration timeout in seconds (default: 300)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  migrate up                           # Run all pending migrations")
	fmt.Println("  migrate -env production up           # Run migrations in production")
	fmt.Println("  migrate down 2                       # Rollback 2 migrations")
	fmt.Println("  migrate status                       # Show migration status")
	fmt.Println("  migrate create add_user_roles        # Create new migration")
	fmt.Println("  migrate backup                       # Create database backup")
	fmt.Println("  migrate validate                     # Validate migration files")
	fmt.Println("  migrate legacy                       # Run legacy auto-migration")
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}