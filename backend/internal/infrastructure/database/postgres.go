package database

import (
	"fmt"
	"log"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/core/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewPostgresConnection(cfg *config.Config) (*Database, error) {
	var db *gorm.DB
	var err error

	// Configure GORM logger
	var gormLogger logger.Interface
	if cfg.IsDevelopment() {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Configure GORM config
	gormConfig := &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Attempt to connect with retries
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(cfg.Database.URL), gormConfig)
		if err == nil {
			break
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxConnections)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	return &Database{DB: db}, nil
}

func (d *Database) AutoMigrate() error {
	log.Println("Running database migrations...")

	err := d.DB.AutoMigrate(
		&domain.User{},
		&domain.ShortURL{},
		&domain.Click{},
	)

	if err != nil {
		return fmt.Errorf("failed to run auto migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func (d *Database) CreateIndexes() error {
	log.Println("Creating database indexes...")

	// Create indexes for better performance
	queries := []string{
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_short_urls_user_id_created_at ON short_urls(user_id, created_at DESC)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_short_urls_short_code_active ON short_urls(short_code, is_active) WHERE deleted_at IS NULL",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_clicks_short_url_id_clicked_at ON clicks(short_url_id, clicked_at DESC)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_clicks_ip_address ON clicks(ip_address)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_clicks_country_clicked_at ON clicks(country, clicked_at DESC)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email_active ON users(email) WHERE deleted_at IS NULL",
	}

	for _, query := range queries {
		if err := d.DB.Exec(query).Error; err != nil {
			// Log warning but don't fail if index already exists
			log.Printf("Warning: Failed to create index: %v", err)
		}
	}

	log.Println("Database indexes created successfully")
	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) Health() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (d *Database) GetStats() map[string]interface{} {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"open_connections":     stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"wait_count":          stats.WaitCount,
		"wait_duration":       stats.WaitDuration,
		"max_idle_closed":     stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}
}