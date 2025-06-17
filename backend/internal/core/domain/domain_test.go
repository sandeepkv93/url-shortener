package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DomainTestSuite struct {
	suite.Suite
}

func (suite *DomainTestSuite) TestUserToResponse() {
	user := &User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	response := user.ToResponse()

	suite.Equal(user.ID, response.ID)
	suite.Equal(user.Email, response.Email)
	suite.Equal(user.CreatedAt, response.CreatedAt)
	suite.Equal(user.UpdatedAt, response.UpdatedAt)
	// UserResponse doesn't have Password field, which is correct for security
}

func (suite *DomainTestSuite) TestShortURLIsExpired() {
	now := time.Now()
	
	// Test URL without expiration
	url1 := &ShortURL{
		ExpiresAt: nil,
	}
	suite.False(url1.IsExpired())

	// Test URL with future expiration
	futureTime := now.Add(time.Hour)
	url2 := &ShortURL{
		ExpiresAt: &futureTime,
	}
	suite.False(url2.IsExpired())

	// Test URL with past expiration
	pastTime := now.Add(-time.Hour)
	url3 := &ShortURL{
		ExpiresAt: &pastTime,
	}
	suite.True(url3.IsExpired())
}

func (suite *DomainTestSuite) TestShortURLIsAccessible() {
	now := time.Now()
	
	// Test active URL without expiration
	url1 := &ShortURL{
		IsActive:  true,
		ExpiresAt: nil,
	}
	suite.True(url1.IsAccessible())

	// Test inactive URL
	url2 := &ShortURL{
		IsActive:  false,
		ExpiresAt: nil,
	}
	suite.False(url2.IsAccessible())

	// Test active but expired URL
	pastTime := now.Add(-time.Hour)
	url3 := &ShortURL{
		IsActive:  true,
		ExpiresAt: &pastTime,
	}
	suite.False(url3.IsAccessible())

	// Test active URL with future expiration
	futureTime := now.Add(time.Hour)
	url4 := &ShortURL{
		IsActive:  true,
		ExpiresAt: &futureTime,
	}
	suite.True(url4.IsAccessible())
}

func (suite *DomainTestSuite) TestShortURLToResponse() {
	baseURL := "https://short.ly"
	shortCode := "abc123"
	originalURL := "https://example.com"
	
	url := &ShortURL{
		ID:          1,
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CustomAlias: false,
		IsActive:    true,
		ClickCount:  100,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	response := url.ToResponse(baseURL)

	suite.Equal(url.ID, response.ID)
	suite.Equal(url.ShortCode, response.ShortCode)
	suite.Equal(url.OriginalURL, response.OriginalURL)
	suite.Equal(baseURL+"/"+shortCode, response.ShortURL)
	suite.Equal(url.CustomAlias, response.CustomAlias)
	suite.Equal(url.IsActive, response.IsActive)
	suite.Equal(url.ClickCount, response.ClickCount)
	suite.Equal(url.CreatedAt, response.CreatedAt)
	suite.Equal(url.UpdatedAt, response.UpdatedAt)
}

func (suite *DomainTestSuite) TestDomainErrorCreation() {
	// Test basic domain error
	err := NewDomainError("test_error", "Test error message", 400)
	suite.Equal("test_error", err.Type)
	suite.Equal("Test error message", err.Message)
	suite.Equal(400, err.Code)
	suite.Equal("[test_error] Test error message", err.Error())

	// Test validation error
	validationErr := NewValidationError("email", "is required")
	suite.Equal("validation_error", validationErr.Type)
	suite.Equal("email: is required", validationErr.Message)
	suite.Equal(400, validationErr.Code)

	// Test not found error
	notFoundErr := NewNotFoundError("user")
	suite.Equal("not_found", notFoundErr.Type)
	suite.Equal("user not found", notFoundErr.Message)
	suite.Equal(404, notFoundErr.Code)

	// Test conflict error
	conflictErr := NewConflictError("email")
	suite.Equal("conflict", conflictErr.Type)
	suite.Equal("email already exists", conflictErr.Message)
	suite.Equal(409, conflictErr.Code)

	// Test unauthorized error
	unauthorizedErr := NewUnauthorizedError("invalid token")
	suite.Equal("unauthorized", unauthorizedErr.Type)
	suite.Equal("invalid token", unauthorizedErr.Message)
	suite.Equal(401, unauthorizedErr.Code)

	// Test forbidden error
	forbiddenErr := NewForbiddenError("access denied")
	suite.Equal("forbidden", forbiddenErr.Type)
	suite.Equal("access denied", forbiddenErr.Message)
	suite.Equal(403, forbiddenErr.Code)

	// Test internal error
	internalErr := NewInternalError("database connection failed")
	suite.Equal("internal_error", internalErr.Type)
	suite.Equal("database connection failed", internalErr.Message)
	suite.Equal(500, internalErr.Code)
}

func (suite *DomainTestSuite) TestQROptionsValidation() {
	// Test valid QR options
	options := &QROptions{
		Size:       300,
		Format:     "png",
		ErrorLevel: "M",
		BorderSize: 4,
		Foreground: "#000000",
		Background: "#FFFFFF",
		LogoSize:   20,
	}
	
	// These would normally be validated by the validator package
	suite.Equal(300, options.Size)
	suite.Equal("png", options.Format)
	suite.Equal("M", options.ErrorLevel)
	suite.Equal(4, options.BorderSize)
	suite.Equal("#000000", options.Foreground)
	suite.Equal("#FFFFFF", options.Background)
	suite.Equal(20, options.LogoSize)
}

func (suite *DomainTestSuite) TestGeoLocationData() {
	location := &GeoLocation{
		IPAddress:   "192.168.1.1",
		Country:     "United States",
		CountryCode: "US",
		Region:      "California",
		RegionCode:  "CA",
		City:        "San Francisco",
		Latitude:    37.7749,
		Longitude:   -122.4194,
		Timezone:    "America/Los_Angeles",
		ISP:         "Example ISP",
		Source:      "ipapi",
		CachedAt:    time.Now(),
	}

	suite.Equal("192.168.1.1", location.IPAddress)
	suite.Equal("United States", location.Country)
	suite.Equal("US", location.CountryCode)
	suite.Equal("California", location.Region)
	suite.Equal("CA", location.RegionCode)
	suite.Equal("San Francisco", location.City)
	suite.Equal(37.7749, location.Latitude)
	suite.Equal(-122.4194, location.Longitude)
	suite.Equal("America/Los_Angeles", location.Timezone)
	suite.Equal("Example ISP", location.ISP)
	suite.Equal("ipapi", location.Source)
}

func (suite *DomainTestSuite) TestHealthStatusComponents() {
	health := &HealthStatus{
		Status:    "healthy",
		Version:   "1.0.0",
		Uptime:    time.Hour * 24,
		Timestamp: time.Now(),
		Components: map[string]*ComponentHealth{
			"database": {
				Status:      "up",
				Message:     "Connected successfully",
				LastChecked: time.Now(),
				Duration:    time.Millisecond * 50,
			},
			"cache": {
				Status:      "up",
				Message:     "Redis responding",
				LastChecked: time.Now(),
				Duration:    time.Millisecond * 10,
			},
		},
	}

	suite.Equal("healthy", health.Status)
	suite.Equal("1.0.0", health.Version)
	suite.Equal(time.Hour*24, health.Uptime)
	suite.Contains(health.Components, "database")
	suite.Contains(health.Components, "cache")
	suite.Equal("up", health.Components["database"].Status)
	suite.Equal("up", health.Components["cache"].Status)
}

func (suite *DomainTestSuite) TestAnalyticsStructures() {
	stats := &ClickStats{
		TotalClicks:  1000,
		UniqueClicks: 750,
		ClicksByDate: map[string]int64{
			"2024-01-01": 100,
			"2024-01-02": 150,
			"2024-01-03": 200,
		},
		ClicksByTime: map[int]int64{
			9:  50,
			10: 75,
			11: 100,
		},
		TopCountries: []CountryStat{
			{Country: "US", Count: 500},
			{Country: "CA", Count: 200},
			{Country: "UK", Count: 150},
		},
		TopDevices: []DeviceStat{
			{Device: "desktop", Count: 600},
			{Device: "mobile", Count: 300},
			{Device: "tablet", Count: 100},
		},
		TopBrowsers: []BrowserStat{
			{Browser: "chrome", Count: 700},
			{Browser: "firefox", Count: 200},
			{Browser: "safari", Count: 100},
		},
	}

	suite.Equal(int64(1000), stats.TotalClicks)
	suite.Equal(int64(750), stats.UniqueClicks)
	suite.Len(stats.ClicksByDate, 3)
	suite.Len(stats.ClicksByTime, 3)
	suite.Len(stats.TopCountries, 3)
	suite.Len(stats.TopDevices, 3)
	suite.Len(stats.TopBrowsers, 3)
	suite.Equal(int64(500), stats.TopCountries[0].Count)
	suite.Equal("US", stats.TopCountries[0].Country)
}

func TestDomainTestSuite(t *testing.T) {
	suite.Run(t, new(DomainTestSuite))
}

// Additional specific tests
func TestCreateUserRequest(t *testing.T) {
	// Test valid create user request
	req := &CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)

	// Test login request
	loginReq := &LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	assert.Equal(t, "test@example.com", loginReq.Email)
	assert.Equal(t, "password123", loginReq.Password)
}

func TestURLFilter(t *testing.T) {
	filter := &URLFilter{
		Status:    "active",
		DateFrom:  "2024-01-01",
		DateTo:    "2024-01-31",
		Search:    "example",
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	assert.Equal(t, "active", filter.Status)
	assert.Equal(t, "2024-01-01", filter.DateFrom)
	assert.Equal(t, "2024-01-31", filter.DateTo)
	assert.Equal(t, "example", filter.Search)
	assert.Equal(t, "created_at", filter.SortBy)
	assert.Equal(t, "desc", filter.SortOrder)
}

func TestRecordClickRequest(t *testing.T) {
	req := &RecordClickRequest{
		ShortCode: "abc123",
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		Referer:   "https://google.com",
	}

	assert.Equal(t, "abc123", req.ShortCode)
	assert.Equal(t, "192.168.1.1", req.IPAddress)
	assert.Contains(t, req.UserAgent, "Mozilla")
	assert.Equal(t, "https://google.com", req.Referer)
}