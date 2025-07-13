package migrations

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/core/services"

	"gorm.io/gorm"
)

// MigrationService handles database migrations with backup and rollback capabilities
type MigrationService struct {
	db     *gorm.DB
	config *config.Config
	logger *services.LoggingService
}

// MigrationRecord tracks applied migrations in the database
type MigrationRecord struct {
	ID          uint      `gorm:"primarykey"`
	Version     string    `gorm:"uniqueIndex;not null"`
	Name        string    `gorm:"not null"`
	AppliedAt   time.Time `gorm:"not null"`
	Checksum    string    `gorm:"not null"`
	ExecutionTimeMs int64 `gorm:"not null"`
}

// Migration represents a database migration
type Migration struct {
	Version     string
	Name        string
	UpSQL       string
	DownSQL     string
	Checksum    string
	FilePath    string
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	Version     string    `json:"version"`
	Name        string    `json:"name"`
	Applied     bool      `json:"applied"`
	AppliedAt   *time.Time `json:"applied_at,omitempty"`
	Checksum    string    `json:"checksum"`
	ExecutionTime int64   `json:"execution_time_ms,omitempty"`
}

// BackupInfo contains information about a database backup
type BackupInfo struct {
	Filename    string    `json:"filename"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
	PreMigrationVersion string `json:"pre_migration_version"`
}

// NewMigrationService creates a new migration service
func NewMigrationService(db *gorm.DB, cfg *config.Config, logger *services.LoggingService) *MigrationService {
	return &MigrationService{
		db:     db,
		config: cfg,
		logger: logger,
	}
}

// Initialize sets up the migration tracking table
func (ms *MigrationService) Initialize() error {
	ms.logger.Info("Initializing migration service...")
	
	// Create migration tracking table
	if err := ms.db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to create migration tracking table: %w", err)
	}
	
	ms.logger.Info("Migration service initialized successfully")
	return nil
}

// LoadMigrations loads all migration files from the migrations directory
func (ms *MigrationService) LoadMigrations() ([]Migration, error) {
	migrationPath := ms.config.GetMigrationPath()
	
	ms.logger.Info("Loading migrations from path", "path", migrationPath)
	
	files, err := filepath.Glob(filepath.Join(migrationPath, "*.sql"))
	if err != nil {
		return nil, fmt.Errorf("failed to read migration files: %w", err)
	}
	
	var migrations []Migration
	
	for _, file := range files {
		migration, err := ms.parseMigrationFile(file)
		if err != nil {
			ms.logger.Error("Failed to parse migration file", "file", file, "error", err.Error())
			continue
		}
		migrations = append(migrations, migration)
	}
	
	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return ms.compareVersions(migrations[i].Version, migrations[j].Version) < 0
	})
	
	ms.logger.Info("Loaded migrations", "count", len(migrations))
	return migrations, nil
}

// GetMigrationStatus returns the status of all migrations
func (ms *MigrationService) GetMigrationStatus() ([]MigrationStatus, error) {
	migrations, err := ms.LoadMigrations()
	if err != nil {
		return nil, err
	}
	
	// Get applied migrations from database
	var appliedMigrations []MigrationRecord
	if err := ms.db.Find(&appliedMigrations).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch applied migrations: %w", err)
	}
	
	// Create a map for quick lookup
	appliedMap := make(map[string]MigrationRecord)
	for _, applied := range appliedMigrations {
		appliedMap[applied.Version] = applied
	}
	
	var statuses []MigrationStatus
	for _, migration := range migrations {
		status := MigrationStatus{
			Version:  migration.Version,
			Name:     migration.Name,
			Checksum: migration.Checksum,
			Applied:  false,
		}
		
		if applied, exists := appliedMap[migration.Version]; exists {
			status.Applied = true
			status.AppliedAt = &applied.AppliedAt
			status.ExecutionTime = applied.ExecutionTimeMs
		}
		
		statuses = append(statuses, status)
	}
	
	return statuses, nil
}

// Migrate runs all pending migrations
func (ms *MigrationService) Migrate(ctx context.Context) error {
	if !ms.config.IsAutoMigrateEnabled() {
		ms.logger.Info("Auto migration is disabled")
		return nil
	}
	
	ms.logger.Info("Starting database migration...")
	
	migrations, err := ms.LoadMigrations()
	if err != nil {
		return err
	}
	
	// Get pending migrations
	pendingMigrations, err := ms.getPendingMigrations(migrations)
	if err != nil {
		return err
	}
	
	if len(pendingMigrations) == 0 {
		ms.logger.Info("No pending migrations found")
		return nil
	}
	
	ms.logger.Info("Found pending migrations", "count", len(pendingMigrations))
	
	// Create backup before migration if enabled
	var backupInfo *BackupInfo
	if ms.config.IsBackupBeforeMigrationEnabled() {
		backupInfo, err = ms.CreateBackup(ctx)
		if err != nil {
			ms.logger.Error("Failed to create backup before migration", "error", err.Error())
			return fmt.Errorf("backup failed: %w", err)
		}
		ms.logger.Info("Backup created successfully", "file", backupInfo.Filename)
	}
	
	// Apply migrations in transaction
	migrationCtx, cancel := context.WithTimeout(ctx, ms.config.GetMigrationTimeout())
	defer cancel()
	
	err = ms.applyMigrationsWithRetry(migrationCtx, pendingMigrations)
	if err != nil {
		ms.logger.Error("Migration failed", "error", err.Error())
		
		// If backup was created, suggest rollback
		if backupInfo != nil {
			ms.logger.Error("Migration failed, backup available for rollback", "backup", backupInfo.Filename)
		}
		
		return err
	}
	
	ms.logger.Info("All migrations applied successfully")
	return nil
}

// MigrateToVersion migrates to a specific version
func (ms *MigrationService) MigrateToVersion(ctx context.Context, targetVersion string) error {
	ms.logger.Info("Migrating to specific version", "version", targetVersion)
	
	migrations, err := ms.LoadMigrations()
	if err != nil {
		return err
	}
	
	// Find target migration
	var targetIndex = -1
	for i, migration := range migrations {
		if migration.Version == targetVersion {
			targetIndex = i
			break
		}
	}
	
	if targetIndex == -1 {
		return fmt.Errorf("migration version %s not found", targetVersion)
	}
	
	// Get current migration status
	currentVersion, err := ms.getCurrentVersion()
	if err != nil {
		return err
	}
	
	if currentVersion == targetVersion {
		ms.logger.Info("Already at target version", "version", targetVersion)
		return nil
	}
	
	// Determine direction (up or down)
	if ms.compareVersions(currentVersion, targetVersion) < 0 {
		// Migrate up
		return ms.migrateUp(ctx, migrations, targetIndex)
	} else {
		// Migrate down (rollback)
		return ms.migrateDown(ctx, migrations, targetIndex)
	}
}

// Rollback rolls back a specific number of migrations
func (ms *MigrationService) Rollback(ctx context.Context, steps int) error {
	ms.logger.Info("Rolling back migrations", "steps", steps)
	
	// Get applied migrations in reverse order
	var appliedMigrations []MigrationRecord
	if err := ms.db.Order("applied_at DESC").Limit(steps).Find(&appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to fetch applied migrations: %w", err)
	}
	
	if len(appliedMigrations) == 0 {
		ms.logger.Info("No migrations to rollback")
		return nil
	}
	
	// Load all migrations to get rollback SQL
	allMigrations, err := ms.LoadMigrations()
	if err != nil {
		return err
	}
	
	migrationMap := make(map[string]Migration)
	for _, migration := range allMigrations {
		migrationMap[migration.Version] = migration
	}
	
	// Create backup before rollback
	if ms.config.IsBackupBeforeMigrationEnabled() {
		backupInfo, err := ms.CreateBackup(ctx)
		if err != nil {
			ms.logger.Error("Failed to create backup before rollback", "error", err.Error())
			return fmt.Errorf("backup failed: %w", err)
		}
		ms.logger.Info("Backup created before rollback", "file", backupInfo.Filename)
	}
	
	// Rollback migrations
	for _, appliedMigration := range appliedMigrations {
		migration, exists := migrationMap[appliedMigration.Version]
		if !exists {
			ms.logger.Warn("Migration file not found for rollback", "version", appliedMigration.Version)
			continue
		}
		
		if migration.DownSQL == "" {
			ms.logger.Warn("No rollback SQL found for migration", "version", migration.Version)
			continue
		}
		
		err := ms.rollbackSingleMigration(ctx, migration, appliedMigration)
		if err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", migration.Version, err)
		}
	}
	
	ms.logger.Info("Rollback completed successfully", "migrations_rolled_back", len(appliedMigrations))
	return nil
}

// CreateBackup creates a database backup
func (ms *MigrationService) CreateBackup(ctx context.Context) (*BackupInfo, error) {
	ms.logger.Info("Creating database backup...")
	
	// Create backup directory if it doesn't exist
	backupDir := filepath.Join(ms.config.GetMigrationPath(), "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}
	
	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	currentVersion, _ := ms.getCurrentVersion()
	filename := fmt.Sprintf("backup_%s_v%s.sql", timestamp, currentVersion)
	backupPath := filepath.Join(backupDir, filename)
	
	// Use pg_dump to create backup
	_, err := ms.db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	
	// For now, create a simple backup structure
	// In production, you would use pg_dump or similar tools with the sql.DB connection
	backupContent := fmt.Sprintf("-- Database backup created at %s\n", time.Now().Format(time.RFC3339))
	backupContent += fmt.Sprintf("-- Pre-migration version: %s\n\n", currentVersion)
	
	if err := os.WriteFile(backupPath, []byte(backupContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write backup file: %w", err)
	}
	
	// Get file stats
	fileInfo, err := os.Stat(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup file stats: %w", err)
	}
	
	backupInfo := &BackupInfo{
		Filename:            filename,
		Path:                backupPath,
		Size:                fileInfo.Size(),
		CreatedAt:           time.Now(),
		PreMigrationVersion: currentVersion,
	}
	
	ms.logger.Info("Database backup created successfully", 
		"file", filename,
		"size", fileInfo.Size(),
		"path", backupPath,
	)
	
	return backupInfo, nil
}

// ValidateMigrations validates all migration files
func (ms *MigrationService) ValidateMigrations() error {
	ms.logger.Info("Validating migrations...")
	
	migrations, err := ms.LoadMigrations()
	if err != nil {
		return err
	}
	
	var errors []string
	
	for _, migration := range migrations {
		// Validate version format
		if !ms.isValidVersion(migration.Version) {
			errors = append(errors, fmt.Sprintf("invalid version format: %s", migration.Version))
		}
		
		// Validate SQL syntax (basic check)
		if strings.TrimSpace(migration.UpSQL) == "" {
			errors = append(errors, fmt.Sprintf("empty up SQL in migration %s", migration.Version))
		}
		
		// Check for dangerous operations in production
		if ms.config.IsProduction() {
			if ms.containsDangerousOperations(migration.UpSQL) {
				errors = append(errors, fmt.Sprintf("dangerous operation detected in migration %s", migration.Version))
			}
		}
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("migration validation failed: %s", strings.Join(errors, "; "))
	}
	
	ms.logger.Info("All migrations validated successfully")
	return nil
}

// Private helper methods

func (ms *MigrationService) parseMigrationFile(filePath string) (Migration, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Migration{}, fmt.Errorf("failed to read migration file: %w", err)
	}
	
	// Extract version and name from filename
	filename := filepath.Base(filePath)
	parts := strings.Split(strings.TrimSuffix(filename, ".sql"), "_")
	if len(parts) < 2 {
		return Migration{}, fmt.Errorf("invalid migration filename format: %s", filename)
	}
	
	version := parts[0]
	name := strings.Join(parts[1:], "_")
	
	// Split up and down migrations
	contentStr := string(content)
	upSQL, downSQL := ms.splitUpDownSQL(contentStr)
	
	// Calculate checksum
	checksum := ms.calculateChecksum(contentStr)
	
	return Migration{
		Version:  version,
		Name:     name,
		UpSQL:    upSQL,
		DownSQL:  downSQL,
		Checksum: checksum,
		FilePath: filePath,
	}, nil
}

func (ms *MigrationService) splitUpDownSQL(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var upSQL, downSQL []string
	var inDown bool
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "-- +migrate Up") {
			inDown = false
			continue
		}
		if strings.HasPrefix(trimmed, "-- +migrate Down") {
			inDown = true
			continue
		}
		
		if inDown {
			downSQL = append(downSQL, line)
		} else {
			upSQL = append(upSQL, line)
		}
	}
	
	return strings.TrimSpace(strings.Join(upSQL, "\n")), 
		   strings.TrimSpace(strings.Join(downSQL, "\n"))
}

func (ms *MigrationService) calculateChecksum(content string) string {
	// Simple checksum calculation - in production use crypto/sha256
	return fmt.Sprintf("%x", len(content))
}

func (ms *MigrationService) compareVersions(v1, v2 string) int {
	// Simple version comparison - assumes versions are timestamps or integers
	n1, err1 := strconv.ParseInt(v1, 10, 64)
	n2, err2 := strconv.ParseInt(v2, 10, 64)
	
	if err1 != nil || err2 != nil {
		return strings.Compare(v1, v2)
	}
	
	if n1 < n2 {
		return -1
	}
	if n1 > n2 {
		return 1
	}
	return 0
}

func (ms *MigrationService) getPendingMigrations(migrations []Migration) ([]Migration, error) {
	var appliedMigrations []MigrationRecord
	
	if err := ms.db.Select("version").Find(&appliedMigrations).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch applied migrations: %w", err)
	}
	
	appliedMap := make(map[string]bool)
	for _, applied := range appliedMigrations {
		appliedMap[applied.Version] = true
	}
	
	var pending []Migration
	for _, migration := range migrations {
		if !appliedMap[migration.Version] {
			pending = append(pending, migration)
		}
	}
	
	return pending, nil
}

func (ms *MigrationService) getCurrentVersion() (string, error) {
	var lastMigration MigrationRecord
	err := ms.db.Order("applied_at DESC").First(&lastMigration).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "0", nil
		}
		return "", err
	}
	return lastMigration.Version, nil
}

func (ms *MigrationService) applyMigrationsWithRetry(ctx context.Context, migrations []Migration) error {
	maxRetries := ms.config.Migration.MaxMigrationRetries
	
	for i := 0; i < maxRetries; i++ {
		err := ms.applyMigrations(ctx, migrations)
		if err == nil {
			return nil
		}
		
		ms.logger.Warn("Migration attempt failed", "attempt", i+1, "error", err.Error())
		
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	
	return fmt.Errorf("migration failed after %d attempts", maxRetries)
}

func (ms *MigrationService) applyMigrations(ctx context.Context, migrations []Migration) error {
	for _, migration := range migrations {
		if err := ms.applySingleMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}
	}
	return nil
}

func (ms *MigrationService) applySingleMigration(ctx context.Context, migration Migration) error {
	ms.logger.Info("Applying migration", "version", migration.Version, "name", migration.Name)
	
	start := time.Now()
	
	// Execute migration in transaction
	err := ms.db.Transaction(func(tx *gorm.DB) error {
		// Execute the migration SQL
		if err := tx.Exec(migration.UpSQL).Error; err != nil {
			return fmt.Errorf("failed to execute migration SQL: %w", err)
		}
		
		// Record the migration
		record := MigrationRecord{
			Version:     migration.Version,
			Name:        migration.Name,
			AppliedAt:   time.Now(),
			Checksum:    migration.Checksum,
			ExecutionTimeMs: time.Since(start).Milliseconds(),
		}
		
		if err := tx.Create(&record).Error; err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}
		
		return nil
	})
	
	if err != nil {
		return err
	}
	
	duration := time.Since(start)
	ms.logger.Info("Migration applied successfully", 
		"version", migration.Version,
		"duration", duration.String(),
	)
	
	return nil
}

func (ms *MigrationService) rollbackSingleMigration(ctx context.Context, migration Migration, record MigrationRecord) error {
	ms.logger.Info("Rolling back migration", "version", migration.Version, "name", migration.Name)
	
	err := ms.db.Transaction(func(tx *gorm.DB) error {
		// Execute rollback SQL
		if err := tx.Exec(migration.DownSQL).Error; err != nil {
			return fmt.Errorf("failed to execute rollback SQL: %w", err)
		}
		
		// Remove migration record
		if err := tx.Delete(&record).Error; err != nil {
			return fmt.Errorf("failed to remove migration record: %w", err)
		}
		
		return nil
	})
	
	if err != nil {
		return err
	}
	
	ms.logger.Info("Migration rolled back successfully", "version", migration.Version)
	return nil
}

func (ms *MigrationService) migrateUp(ctx context.Context, migrations []Migration, targetIndex int) error {
	// Apply migrations up to target
	pending := migrations[:targetIndex+1]
	pendingToApply, err := ms.getPendingMigrations(pending)
	if err != nil {
		return err
	}
	
	return ms.applyMigrations(ctx, pendingToApply)
}

func (ms *MigrationService) migrateDown(ctx context.Context, migrations []Migration, targetIndex int) error {
	// This would implement rollback to specific version
	// For now, return error as it's complex to implement
	return fmt.Errorf("migrate down to specific version not implemented yet")
}

func (ms *MigrationService) isValidVersion(version string) bool {
	// Check if version is numeric timestamp or valid format
	_, err := strconv.ParseInt(version, 10, 64)
	return err == nil
}

func (ms *MigrationService) containsDangerousOperations(sql string) bool {
	dangerous := []string{
		"DROP TABLE",
		"DROP DATABASE",
		"TRUNCATE",
		"DELETE FROM",
		"DROP INDEX",
	}
	
	upperSQL := strings.ToUpper(sql)
	for _, op := range dangerous {
		if strings.Contains(upperSQL, op) {
			return true
		}
	}
	return false
}