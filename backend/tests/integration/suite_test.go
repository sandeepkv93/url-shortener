package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"url-shortener/internal/api/handlers"
	"url-shortener/internal/api/middleware"
	"url-shortener/internal/api/routes"
	"url-shortener/internal/config"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
	"url-shortener/internal/core/services"
	"url-shortener/internal/infrastructure/database/repositories"
)

// TestQRCodeProvider is a simple QR code provider for testing
type TestQRCodeProvider struct{}

func (t *TestQRCodeProvider) GenerateQRCode(url string, options domain.QRGenerationOptions) ([]byte, error) {
	// Return a simple test QR code - in reality this would generate actual QR code
	return []byte("test-qr-code-data"), nil
}

// IntegrationTestSuite is the base test suite for all integration tests
type IntegrationTestSuite struct {
	suite.Suite
	
	// Test infrastructure
	db           *gorm.DB
	redisClient  *redis.Client
	cacheService ports.CacheService
	server       *httptest.Server
	router       http.Handler
	
	// Repositories
	userRepo  ports.UserRepository
	urlRepo   ports.URLRepository
	clickRepo ports.ClickRepository
	
	// Services
	authService      ports.AuthService
	urlService       ports.URLService
	analyticsService ports.AnalyticsService
	qrService        ports.QRService
	jwtService       ports.JWTService
	
	// Handlers
	authHandler      *handlers.AuthHandler
	urlHandler       *handlers.URLHandler
	analyticsHandler *handlers.AnalyticsHandler
	qrHandler        *handlers.QRHandler
	
	// Middleware
	authMiddleware *middleware.AuthMiddleware
	
	// Test data
	testUser *domain.User
	testJWT  string
}

// SetupSuite runs once before the entire test suite
func (s *IntegrationTestSuite) SetupSuite() {
	// Set up test environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("JWT_SECRET", "test-secret-key-for-integration-tests")
	os.Setenv("ENVIRONMENT", "test")
	
	// Initialize test database (SQLite for integration tests)
	s.setupTestDatabase()
	
	// Initialize test Redis
	s.setupTestRedis()
	
	// Initialize repositories
	s.setupRepositories()
	
	// Initialize services
	s.setupServices()
	
	// Initialize handlers
	s.setupHandlers()
	
	// Initialize middleware
	s.setupMiddleware()
	
	// Setup HTTP router and server
	s.setupServer()
	
	// Create test user
	s.createTestUser()
}

// TearDownSuite runs once after the entire test suite
func (s *IntegrationTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
	
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	
	if s.db != nil {
		sqlDB, _ := s.db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
}

// SetupTest runs before each test
func (s *IntegrationTestSuite) SetupTest() {
	// Clean database tables
	s.cleanDatabase()
	
	// Clean Redis cache
	s.cleanRedisCache()
	
	// Recreate test user for each test
	s.createTestUser()
}

// TearDownTest runs after each test
func (s *IntegrationTestSuite) TearDownTest() {
	// Additional cleanup if needed
}

// setupTestDatabase initializes an in-memory SQLite database for testing
func (s *IntegrationTestSuite) setupTestDatabase() {
	var err error
	
	// Use SQLite in-memory database for fast integration tests
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Reduce noise in tests
	})
	s.Require().NoError(err, "Failed to connect to test database")
	
	// Auto-migrate tables
	err = s.db.AutoMigrate(
		&domain.User{},
		&domain.ShortURL{},
		&domain.Click{},
	)
	s.Require().NoError(err, "Failed to migrate test database")
}

// setupTestRedis initializes a test Redis client
func (s *IntegrationTestSuite) setupTestRedis() {
	// Use Redis database 15 for testing (separate from main DB 0)
	s.redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // Use separate test database
	})
	
	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err := s.redisClient.Ping(ctx).Result()
	if err != nil {
		s.T().Skip("Redis not available for integration tests")
		return
	}
	
	// Initialize cache service (will be nil if Redis is not available)
	s.cacheService = nil
}

// setupRepositories initializes all repository instances
func (s *IntegrationTestSuite) setupRepositories() {
	s.userRepo = repositories.NewUserRepository(s.db)
	s.urlRepo = repositories.NewURLRepository(s.db)
	s.clickRepo = repositories.NewClickRepository(s.db)
}

// setupServices initializes all service instances
func (s *IntegrationTestSuite) setupServices() {
	// Load configuration
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret-key-for-integration-tests",
			Expiry: 24 * time.Hour,
		},
		App: config.AppConfig{
			BaseURL: "http://localhost:8080",
		},
		Server: config.ServerConfig{
			Env: "test",
		},
	}
	
	// Initialize services with nil cache service for now (since Redis might not be available)
	s.jwtService = services.NewJWTService(cfg.JWT.Secret)
	
	// Create a config service for tests
	configService := services.NewConfigService(cfg)
	
	s.authService = services.NewAuthService(s.userRepo, s.cacheService, s.jwtService, configService)
	s.urlService = services.NewURLService(s.urlRepo, s.clickRepo, s.cacheService, configService)
	s.analyticsService = services.NewAnalyticsService(s.urlRepo, s.clickRepo, s.userRepo, s.cacheService, configService)
	
	// Create a simple QR code provider for testing
	qrProvider := &TestQRCodeProvider{}
	s.qrService = services.NewQRService(s.urlRepo, configService, qrProvider)
}

// setupHandlers initializes all handler instances
func (s *IntegrationTestSuite) setupHandlers() {
	s.authHandler = handlers.NewAuthHandler(s.authService)
	s.urlHandler = handlers.NewURLHandler(s.urlService, s.analyticsService)
	s.analyticsHandler = handlers.NewAnalyticsHandler(s.analyticsService)
	s.qrHandler = handlers.NewQRHandler(s.qrService)
}

// setupMiddleware initializes middleware instances
func (s *IntegrationTestSuite) setupMiddleware() {
	s.authMiddleware = middleware.NewAuthMiddleware(s.jwtService, s.userRepo)
}

// setupServer initializes the HTTP router and test server
func (s *IntegrationTestSuite) setupServer() {
	// Create router configuration
	routeConfig := &routes.Config{
		AuthHandler:      s.authHandler,
		URLHandler:       s.urlHandler,
		AnalyticsHandler: s.analyticsHandler,
		QRHandler:        s.qrHandler,
		AuthMiddleware:   s.authMiddleware,
		CacheService:     s.cacheService,
		EnableCORS:       true,
		EnableLogging:    false,
		AllowedOrigins:   []string{"*"},
	}
	
	// Create router and setup routes
	router := routes.NewRouter(routeConfig)
	s.router = router.SetupRoutes()
	
	// Create test server
	s.server = httptest.NewServer(s.router)
}

// createTestUser creates a test user for use in tests
func (s *IntegrationTestSuite) createTestUser() {
	testEmail := "test@example.com"
	testPassword := "password123"
	
	// Create user directly in database for testing
	hashedPassword, err := bcryptHashPassword(testPassword, 12)
	s.Require().NoError(err, "Failed to hash password")
	
	user := &domain.User{
		Email:     testEmail,
		Password:  hashedPassword,
		FirstName: "Test",
		LastName:  "User",
	}
	
	err = s.userRepo.Create(context.Background(), user)
	s.Require().NoError(err, "Failed to create test user")
	
	s.testUser = user
	
	// Generate JWT token for the test user
	token, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	s.Require().NoError(err, "Failed to generate test JWT")
	
	s.testJWT = token
}

// bcryptHashPassword hashes a password using bcrypt
func bcryptHashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

// cleanDatabase removes all data from database tables
func (s *IntegrationTestSuite) cleanDatabase() {
	// Delete in reverse order of dependencies
	s.db.Exec("DELETE FROM clicks")
	s.db.Exec("DELETE FROM short_urls")
	s.db.Exec("DELETE FROM users")
}

// cleanRedisCache clears all keys from the test Redis database
func (s *IntegrationTestSuite) cleanRedisCache() {
	if s.redisClient != nil {
		ctx := context.Background()
		s.redisClient.FlushDB(ctx)
	}
}

// makeRequest is a helper method to make HTTP requests to the test server
func (s *IntegrationTestSuite) makeRequest(method, path string, body interface{}, headers map[string]string) *http.Response {
	var reqBody *bytes.Buffer
	
	if body != nil {
		jsonData, err := json.Marshal(body)
		s.Require().NoError(err)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}
	
	req, err := http.NewRequest(method, s.server.URL+path, reqBody)
	s.Require().NoError(err)
	
	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	
	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}
	
	resp, err := client.Do(req)
	s.Require().NoError(err)
	
	return resp
}

// makeAuthenticatedRequest makes a request with the test user's JWT token
func (s *IntegrationTestSuite) makeAuthenticatedRequest(method, path string, body interface{}) *http.Response {
	headers := map[string]string{
		"Authorization": "Bearer " + s.testJWT,
	}
	return s.makeRequest(method, path, body, headers)
}

// parseJSONResponse parses the JSON response body into the provided interface
func (s *IntegrationTestSuite) parseJSONResponse(resp *http.Response, dest interface{}) {
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(dest)
	s.Require().NoError(err)
}

// assertJSONResponse checks the response status and parses JSON
func (s *IntegrationTestSuite) assertJSONResponse(resp *http.Response, expectedStatus int, dest interface{}) {
	s.Equal(expectedStatus, resp.StatusCode, "Unexpected response status")
	s.parseJSONResponse(resp, dest)
}

// assertErrorResponse checks for an error response with expected status and message
func (s *IntegrationTestSuite) assertErrorResponse(resp *http.Response, expectedStatus int, expectedMessage string) {
	s.Equal(expectedStatus, resp.StatusCode)
	
	var errorResp map[string]interface{}
	s.parseJSONResponse(resp, &errorResp)
	
	if expectedMessage != "" {
		s.Contains(fmt.Sprintf("%v", errorResp["error"]), expectedMessage)
	}
}

// Helper methods for common operations

// loginUser logs in the test user and returns the JWT token
func (s *IntegrationTestSuite) loginUser() string {
	loginReq := map[string]string{
		"email":    s.testUser.Email,
		"password": "password123",
	}
	
	resp := s.makeRequest("POST", "/api/v1/auth/login", loginReq, nil)
	
	var loginResp struct {
		Token string `json:"token"`
	}
	s.assertJSONResponse(resp, http.StatusOK, &loginResp)
	
	return loginResp.Token
}

// createTestURL creates a test URL and returns the response
func (s *IntegrationTestSuite) createTestURL(originalURL string) map[string]interface{} {
	urlReq := map[string]interface{}{
		"originalUrl": originalURL,
		"title":       "Test URL",
		"description": "A test URL for integration testing",
	}
	
	resp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", urlReq)
	
	var urlResp map[string]interface{}
	s.assertJSONResponse(resp, http.StatusCreated, &urlResp)
	
	return urlResp
}

// waitForAsyncOperation waits for an asynchronous operation to complete
func (s *IntegrationTestSuite) waitForAsyncOperation(timeout time.Duration, checkFunc func() bool) bool {
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		if checkFunc() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	
	return false
}

// Test runner function for the integration test suite
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}