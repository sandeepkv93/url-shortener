package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type RepositoryTestSuite struct {
	suite.Suite
	db              *gorm.DB
	userRepo        ports.UserRepository
	urlRepo         ports.URLRepository
	clickRepo       ports.ClickRepository
	ctx             context.Context
	testUser        *domain.User
	testURL         *domain.ShortURL
}

func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// Auto-migrate the schema
	err = db.AutoMigrate(&domain.User{}, &domain.ShortURL{}, &domain.Click{})
	suite.Require().NoError(err)

	suite.db = db

	// Initialize repositories
	suite.userRepo = NewUserRepository(db)
	suite.urlRepo = NewURLRepository(db)
	suite.clickRepo = NewClickRepository(db)
}

func (suite *RepositoryTestSuite) SetupTest() {
	// Clean up tables before each test
	suite.db.Exec("DELETE FROM clicks")
	suite.db.Exec("DELETE FROM short_urls")
	suite.db.Exec("DELETE FROM users")

	// Create test user
	suite.testUser = &domain.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := suite.userRepo.Create(suite.ctx, suite.testUser)
	suite.Require().NoError(err)

	// Create test URL
	suite.testURL = &domain.ShortURL{
		ShortCode:   "test123",
		OriginalURL: "https://example.com",
		UserID:      &suite.testUser.ID,
		IsActive:    true,
	}
	err = suite.urlRepo.Create(suite.ctx, suite.testURL)
	suite.Require().NoError(err)
}

func (suite *RepositoryTestSuite) TearDownSuite() {
	// Clean up database connection
	db, err := suite.db.DB()
	if err == nil {
		db.Close()
	}
}

// User Repository Tests
func (suite *RepositoryTestSuite) TestUserRepository_Create() {
	user := &domain.User{
		Email:    "new@example.com",
		Password: "hashedpassword",
	}

	err := suite.userRepo.Create(suite.ctx, user)
	suite.NoError(err)
	suite.NotZero(user.ID)

	// Test duplicate email
	duplicateUser := &domain.User{
		Email:    "new@example.com",
		Password: "anotherpassword",
	}
	err = suite.userRepo.Create(suite.ctx, duplicateUser)
	suite.Error(err)
	suite.Equal(domain.ErrUserAlreadyExists, err)
}

func (suite *RepositoryTestSuite) TestUserRepository_GetByID() {
	user, err := suite.userRepo.GetByID(suite.ctx, suite.testUser.ID)
	suite.NoError(err)
	suite.Equal(suite.testUser.Email, user.Email)

	// Test non-existent user
	_, err = suite.userRepo.GetByID(suite.ctx, 999)
	suite.Error(err)
	suite.Equal(domain.ErrUserNotFound, err)
}

func (suite *RepositoryTestSuite) TestUserRepository_GetByEmail() {
	user, err := suite.userRepo.GetByEmail(suite.ctx, suite.testUser.Email)
	suite.NoError(err)
	suite.Equal(suite.testUser.ID, user.ID)

	// Test non-existent email
	_, err = suite.userRepo.GetByEmail(suite.ctx, "nonexistent@example.com")
	suite.Error(err)
	suite.Equal(domain.ErrUserNotFound, err)
}

func (suite *RepositoryTestSuite) TestUserRepository_Update() {
	suite.testUser.Email = "updated@example.com"
	err := suite.userRepo.Update(suite.ctx, suite.testUser)
	suite.NoError(err)

	// Verify update
	user, err := suite.userRepo.GetByID(suite.ctx, suite.testUser.ID)
	suite.NoError(err)
	suite.Equal("updated@example.com", user.Email)
}

func (suite *RepositoryTestSuite) TestUserRepository_Delete() {
	err := suite.userRepo.Delete(suite.ctx, suite.testUser.ID)
	suite.NoError(err)

	// Verify deletion
	_, err = suite.userRepo.GetByID(suite.ctx, suite.testUser.ID)
	suite.Error(err)
	suite.Equal(domain.ErrUserNotFound, err)

	// Test deleting non-existent user
	err = suite.userRepo.Delete(suite.ctx, 999)
	suite.Error(err)
	suite.Equal(domain.ErrUserNotFound, err)
}

func (suite *RepositoryTestSuite) TestUserRepository_Exists() {
	exists, err := suite.userRepo.Exists(suite.ctx, suite.testUser.Email)
	suite.NoError(err)
	suite.True(exists)

	exists, err = suite.userRepo.Exists(suite.ctx, "nonexistent@example.com")
	suite.NoError(err)
	suite.False(exists)
}

func (suite *RepositoryTestSuite) TestUserRepository_List() {
	// Create additional users
	for i := 0; i < 5; i++ {
		user := &domain.User{
			Email:    fmt.Sprintf("user%d@example.com", i),
			Password: "password",
		}
		err := suite.userRepo.Create(suite.ctx, user)
		suite.NoError(err)
	}

	users, total, err := suite.userRepo.List(suite.ctx, 0, 10)
	suite.NoError(err)
	suite.Equal(int64(6), total) // 5 new + 1 test user
	suite.Len(users, 6)

	// Test pagination
	users, total, err = suite.userRepo.List(suite.ctx, 2, 2)
	suite.NoError(err)
	suite.Equal(int64(6), total)
	suite.Len(users, 2)
}

// URL Repository Tests
func (suite *RepositoryTestSuite) TestURLRepository_Create() {
	url := &domain.ShortURL{
		ShortCode:   "new123",
		OriginalURL: "https://newexample.com",
		UserID:      &suite.testUser.ID,
		IsActive:    true,
	}

	err := suite.urlRepo.Create(suite.ctx, url)
	suite.NoError(err)
	suite.NotZero(url.ID)

	// Test duplicate short code
	duplicateURL := &domain.ShortURL{
		ShortCode:   "new123",
		OriginalURL: "https://another.com",
		UserID:      &suite.testUser.ID,
	}
	err = suite.urlRepo.Create(suite.ctx, duplicateURL)
	suite.Error(err)
	suite.Equal(domain.ErrShortCodeExists, err)
}

func (suite *RepositoryTestSuite) TestURLRepository_GetByShortCode() {
	url, err := suite.urlRepo.GetByShortCode(suite.ctx, suite.testURL.ShortCode)
	suite.NoError(err)
	suite.Equal(suite.testURL.OriginalURL, url.OriginalURL)

	// Test non-existent short code
	_, err = suite.urlRepo.GetByShortCode(suite.ctx, "nonexistent")
	suite.Error(err)
	suite.Equal(domain.ErrShortURLNotFound, err)
}

func (suite *RepositoryTestSuite) TestURLRepository_GetActiveByShortCode() {
	// Test active URL
	url, err := suite.urlRepo.GetActiveByShortCode(suite.ctx, suite.testURL.ShortCode)
	suite.NoError(err)
	suite.Equal(suite.testURL.ID, url.ID)

	// Deactivate URL
	suite.testURL.IsActive = false
	err = suite.urlRepo.Update(suite.ctx, suite.testURL)
	suite.NoError(err)

	// Test inactive URL
	_, err = suite.urlRepo.GetActiveByShortCode(suite.ctx, suite.testURL.ShortCode)
	suite.Error(err)
	suite.Equal(domain.ErrShortURLNotFound, err)
}

func (suite *RepositoryTestSuite) TestURLRepository_IncrementClickCount() {
	originalCount := suite.testURL.ClickCount

	err := suite.urlRepo.IncrementClickCount(suite.ctx, suite.testURL.ID)
	suite.NoError(err)

	// Verify increment
	url, err := suite.urlRepo.GetByID(suite.ctx, suite.testURL.ID)
	suite.NoError(err)
	suite.Equal(originalCount+1, url.ClickCount)
}

func (suite *RepositoryTestSuite) TestURLRepository_GetExpiredURLs() {
	// Create expired URL
	pastTime := time.Now().Add(-time.Hour)
	expiredURL := &domain.ShortURL{
		ShortCode:   "expired",
		OriginalURL: "https://expired.com",
		UserID:      &suite.testUser.ID,
		ExpiresAt:   &pastTime,
		IsActive:    true,
	}
	err := suite.urlRepo.Create(suite.ctx, expiredURL)
	suite.NoError(err)

	expiredURLs, err := suite.urlRepo.GetExpiredURLs(suite.ctx, 10)
	suite.NoError(err)
	suite.Len(expiredURLs, 1)
	suite.Equal("expired", expiredURLs[0].ShortCode)
}

// Click Repository Tests
func (suite *RepositoryTestSuite) TestClickRepository_Create() {
	click := &domain.Click{
		ShortURLID: suite.testURL.ID,
		IPAddress:  "192.168.1.1",
		UserAgent:  "Mozilla/5.0",
		Country:    "US",
		City:       "New York",
		Device:     "desktop",
		Browser:    "chrome",
		ClickedAt:  time.Now(),
	}

	err := suite.clickRepo.Create(suite.ctx, click)
	suite.NoError(err)
	suite.NotZero(click.ID)
}

func (suite *RepositoryTestSuite) TestClickRepository_GetTotalClicks() {
	// Create test clicks
	for i := 0; i < 5; i++ {
		click := &domain.Click{
			ShortURLID: suite.testURL.ID,
			IPAddress:  fmt.Sprintf("192.168.1.%d", i+1),
			ClickedAt:  time.Now(),
		}
		err := suite.clickRepo.Create(suite.ctx, click)
		suite.NoError(err)
	}

	count, err := suite.clickRepo.GetTotalClicks(suite.ctx, suite.testURL.ID)
	suite.NoError(err)
	suite.Equal(int64(5), count)
}

func (suite *RepositoryTestSuite) TestClickRepository_GetUniqueClicks() {
	// Create clicks with some duplicate IPs
	ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.1", "192.168.1.3", "192.168.1.2"}
	for _, ip := range ips {
		click := &domain.Click{
			ShortURLID: suite.testURL.ID,
			IPAddress:  ip,
			ClickedAt:  time.Now(),
		}
		err := suite.clickRepo.Create(suite.ctx, click)
		suite.NoError(err)
	}

	count, err := suite.clickRepo.GetUniqueClicks(suite.ctx, suite.testURL.ID)
	suite.NoError(err)
	suite.Equal(int64(3), count) // 3 unique IPs
}

func (suite *RepositoryTestSuite) TestClickRepository_GetTopCountries() {
	// Create clicks from different countries
	countries := []string{"US", "CA", "US", "UK", "US"}
	for _, country := range countries {
		click := &domain.Click{
			ShortURLID: suite.testURL.ID,
			IPAddress:  "192.168.1.1",
			Country:    country,
			ClickedAt:  time.Now(),
		}
		err := suite.clickRepo.Create(suite.ctx, click)
		suite.NoError(err)
	}

	stats, err := suite.clickRepo.GetTopCountries(suite.ctx, suite.testURL.ID, 10)
	suite.NoError(err)
	suite.Len(stats, 3)
	suite.Equal("US", stats[0].Country)
	suite.Equal(int64(3), stats[0].Count)
}

func (suite *RepositoryTestSuite) TestClickRepository_GetGlobalStats() {
	// Create additional test data
	user2 := &domain.User{
		Email:    "user2@example.com",
		Password: "password",
	}
	err := suite.userRepo.Create(suite.ctx, user2)
	suite.NoError(err)

	// Create some clicks
	for i := 0; i < 3; i++ {
		click := &domain.Click{
			ShortURLID: suite.testURL.ID,
			IPAddress:  fmt.Sprintf("192.168.1.%d", i+1),
			ClickedAt:  time.Now(),
		}
		err := suite.clickRepo.Create(suite.ctx, click)
		suite.NoError(err)
	}

	stats, err := suite.clickRepo.GetGlobalStats(suite.ctx)
	suite.NoError(err)
	suite.Equal(int64(2), stats.TotalUsers)
	suite.Equal(int64(1), stats.TotalURLs)
	suite.Equal(int64(3), stats.TotalClicks)
	suite.Equal(int64(1), stats.ActiveURLs)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

// Additional helper tests
func TestIsDuplicateKeyError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "duplicate key error",
			err:      fmt.Errorf("duplicate key"),
			expected: true,
		},
		{
			name:     "unique constraint error",
			err:      fmt.Errorf("UNIQUE constraint failed"),
			expected: true,
		},
		{
			name:     "other error",
			err:      fmt.Errorf("connection failed"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDuplicateKeyError(tt.err)
			if result != tt.expected {
				t.Errorf("isDuplicateKeyError() = %v, want %v", result, tt.expected)
			}
		})
	}
}