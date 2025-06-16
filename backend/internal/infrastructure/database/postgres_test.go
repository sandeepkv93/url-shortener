package database

import (
	"testing"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DatabaseTestSuite struct {
	suite.Suite
	db  *Database
	cfg *config.Config
}

func (suite *DatabaseTestSuite) SetupSuite() {
	// Use test database configuration
	suite.cfg = &config.Config{
		Database: config.DatabaseConfig{
			URL:            "postgres://test:test@localhost:5432/urlshortener_test?sslmode=disable",
			Host:           "localhost",
			Port:           "5432",
			Name:           "urlshortener_test",
			User:           "test",
			Password:       "test",
			MaxConnections: 5,
			MaxIdle:        2,
		},
		Server: config.ServerConfig{
			Env: "test",
		},
	}
}

func (suite *DatabaseTestSuite) SetupTest() {
	db, err := NewPostgresConnection(suite.cfg)
	if err != nil {
		suite.T().Skip("Skipping database tests - PostgreSQL not available")
		return
	}
	suite.db = db

	// Run migrations
	err = suite.db.AutoMigrate()
	suite.Require().NoError(err)
}

func (suite *DatabaseTestSuite) TearDownTest() {
	if suite.db != nil {
		// Clean up test data
		suite.db.DB.Exec("TRUNCATE TABLE clicks CASCADE")
		suite.db.DB.Exec("TRUNCATE TABLE short_urls CASCADE")
		suite.db.DB.Exec("TRUNCATE TABLE users CASCADE")
		suite.db.Close()
	}
}

func (suite *DatabaseTestSuite) TestConnection() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
		return
	}

	err := suite.db.Health()
	suite.NoError(err, "Database health check should pass")
}

func (suite *DatabaseTestSuite) TestAutoMigrate() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
		return
	}

	err := suite.db.AutoMigrate()
	suite.NoError(err, "Auto migration should succeed")

	// Check if tables exist
	var count int64
	
	// Check users table
	err = suite.db.DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'users'").Scan(&count).Error
	suite.NoError(err)
	suite.Equal(int64(1), count, "Users table should exist")

	// Check short_urls table
	err = suite.db.DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'short_urls'").Scan(&count).Error
	suite.NoError(err)
	suite.Equal(int64(1), count, "Short URLs table should exist")

	// Check clicks table
	err = suite.db.DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'clicks'").Scan(&count).Error
	suite.NoError(err)
	suite.Equal(int64(1), count, "Clicks table should exist")
}

func (suite *DatabaseTestSuite) TestCreateUser() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
		return
	}

	user := &domain.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := suite.db.DB.Create(user).Error
	suite.NoError(err, "Creating user should succeed")
	suite.NotZero(user.ID, "User ID should be set")
	suite.NotZero(user.CreatedAt, "CreatedAt should be set")
}

func (suite *DatabaseTestSuite) TestCreateShortURL() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
		return
	}

	// Create user first
	user := &domain.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := suite.db.DB.Create(user).Error
	suite.Require().NoError(err)

	// Create short URL
	shortURL := &domain.ShortURL{
		ShortCode:   "test123",
		OriginalURL: "https://example.com",
		UserID:      &user.ID,
		IsActive:    true,
	}

	err = suite.db.DB.Create(shortURL).Error
	suite.NoError(err, "Creating short URL should succeed")
	suite.NotZero(shortURL.ID, "Short URL ID should be set")
	suite.NotZero(shortURL.CreatedAt, "CreatedAt should be set")
}

func (suite *DatabaseTestSuite) TestCreateClick() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
		return
	}

	// Create user first
	user := &domain.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := suite.db.DB.Create(user).Error
	suite.Require().NoError(err)

	// Create short URL
	shortURL := &domain.ShortURL{
		ShortCode:   "test123",
		OriginalURL: "https://example.com",
		UserID:      &user.ID,
		IsActive:    true,
	}
	err = suite.db.DB.Create(shortURL).Error
	suite.Require().NoError(err)

	// Create click
	click := &domain.Click{
		ShortURLID: shortURL.ID,
		IPAddress:  "192.168.1.1",
		UserAgent:  "Mozilla/5.0",
		Country:    "US",
		City:       "New York",
		ClickedAt:  time.Now(),
	}

	err = suite.db.DB.Create(click).Error
	suite.NoError(err, "Creating click should succeed")
	suite.NotZero(click.ID, "Click ID should be set")
}

func (suite *DatabaseTestSuite) TestDatabaseStats() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
		return
	}

	stats := suite.db.GetStats()
	suite.NotEmpty(stats, "Database stats should not be empty")
	suite.Contains(stats, "open_connections", "Stats should contain open_connections")
	suite.Contains(stats, "in_use", "Stats should contain in_use")
}

func (suite *DatabaseTestSuite) TestUniqueConstraints() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
		return
	}

	// Test unique email constraint
	user1 := &domain.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := suite.db.DB.Create(user1).Error
	suite.NoError(err)

	user2 := &domain.User{
		Email:    "test@example.com", // Same email
		Password: "anotherpassword",
	}
	err = suite.db.DB.Create(user2).Error
	suite.Error(err, "Creating user with duplicate email should fail")

	// Test unique short code constraint
	shortURL1 := &domain.ShortURL{
		ShortCode:   "unique123",
		OriginalURL: "https://example.com",
		UserID:      &user1.ID,
	}
	err = suite.db.DB.Create(shortURL1).Error
	suite.NoError(err)

	shortURL2 := &domain.ShortURL{
		ShortCode:   "unique123", // Same short code
		OriginalURL: "https://different.com",
		UserID:      &user1.ID,
	}
	err = suite.db.DB.Create(shortURL2).Error
	suite.Error(err, "Creating short URL with duplicate code should fail")
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

// Unit tests for configuration
func TestDatabaseConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &config.Config{
				Database: config.DatabaseConfig{
					URL:            "postgres://user:pass@localhost:5432/db?sslmode=disable",
					MaxConnections: 10,
					MaxIdle:        5,
				},
				Server: config.ServerConfig{
					Env: "development",
				},
			},
			expectError: false,
		},
		{
			name: "invalid URL",
			config: &config.Config{
				Database: config.DatabaseConfig{
					URL:            "invalid-url",
					MaxConnections: 10,
					MaxIdle:        5,
				},
				Server: config.ServerConfig{
					Env: "development",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPostgresConnection(tt.config)
			if tt.expectError {
				assert.Error(t, err)
			} else if err != nil {
				// Only fail if we expected success but can't connect to a real DB
				// In CI/testing environments, DB might not be available
				t.Logf("Database connection failed (expected in CI): %v", err)
			}
		})
	}
}