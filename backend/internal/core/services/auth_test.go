package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"url-shortener/internal/core/domain"
)

type AuthServiceTestSuite struct {
	suite.Suite
	authService   *authService
	mockUserRepo  *MockUserRepository
	mockCacheRepo *MockCacheService
	mockJWTRepo   *MockJWTService
	mockConfigRepo *MockConfigService
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (suite *AuthServiceTestSuite) SetupTest() {
	suite.mockUserRepo = &MockUserRepository{}
	suite.mockCacheRepo = &MockCacheService{}
	suite.mockJWTRepo = &MockJWTService{}
	suite.mockConfigRepo = &MockConfigService{}
	
	suite.authService = &authService{
		userRepo:   suite.mockUserRepo,
		cacheRepo:  suite.mockCacheRepo,
		jwtService: suite.mockJWTRepo,
		configRepo: suite.mockConfigRepo,
	}
}

func (suite *AuthServiceTestSuite) TestRegister_Success() {
	ctx := context.Background()
	req := domain.RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Mock expectations
	suite.mockUserRepo.On("Exists", ctx, req.Email).Return(false, nil)
	suite.mockUserRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
	suite.mockJWTRepo.On("GenerateAccessToken", mock.AnythingOfType("uint"), req.Email).Return("access_token", nil)
	suite.mockJWTRepo.On("GenerateRefreshToken", mock.AnythingOfType("uint")).Return("refresh_token", nil)
	suite.mockCacheRepo.On("SetSession", ctx, "refresh_token", mock.AnythingOfType("uint"), time.Hour*24*7).Return(nil)

	// Execute
	response, err := suite.authService.Register(ctx, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), "access_token", response.AccessToken)
	assert.Equal(suite.T(), "refresh_token", response.RefreshToken)
	assert.Equal(suite.T(), "Bearer", response.TokenType)
	assert.Equal(suite.T(), 3600, response.ExpiresIn)
	assert.Equal(suite.T(), req.Email, response.User.Email)

	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockCacheRepo.AssertExpectations(suite.T())
	suite.mockJWTRepo.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestRegister_UserExists() {
	ctx := context.Background()
	req := domain.RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Mock expectations
	suite.mockUserRepo.On("Exists", ctx, req.Email).Return(true, nil)

	// Execute
	response, err := suite.authService.Register(ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Equal(suite.T(), domain.ErrUserAlreadyExists, err)

	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestLogin_Success() {
	ctx := context.Background()
	req := domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword, _ := suite.authService.hashPassword("password123")
	user := &domain.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  hashedPassword,
		FirstName: "John",
		LastName:  "Doe",
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	// Mock expectations
	suite.mockUserRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)
	suite.mockUserRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
	suite.mockJWTRepo.On("GenerateAccessToken", user.ID, user.Email).Return("access_token", nil)
	suite.mockJWTRepo.On("GenerateRefreshToken", user.ID).Return("refresh_token", nil)
	suite.mockCacheRepo.On("SetSession", ctx, "refresh_token", user.ID, time.Hour*24*7).Return(nil)

	// Execute
	response, err := suite.authService.Login(ctx, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), "access_token", response.AccessToken)
	assert.Equal(suite.T(), "refresh_token", response.RefreshToken)

	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockCacheRepo.AssertExpectations(suite.T())
	suite.mockJWTRepo.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestLogin_InvalidCredentials() {
	ctx := context.Background()
	req := domain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	hashedPassword, _ := suite.authService.hashPassword("password123")
	user := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashedPassword,
		IsActive: true,
	}

	// Mock expectations
	suite.mockUserRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)

	// Execute
	response, err := suite.authService.Login(ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Equal(suite.T(), domain.ErrInvalidCredentials, err)

	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestRefreshToken_Success() {
	ctx := context.Background()
	refreshToken := "refresh_token"
	userID := uint(1)

	claims := &domain.TokenClaims{
		UserID: userID,
	}

	user := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	// Mock expectations
	suite.mockJWTRepo.On("ValidateRefreshToken", refreshToken).Return(claims, nil)
	suite.mockCacheRepo.On("GetSession", ctx, refreshToken).Return(userID, nil)
	suite.mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)
	suite.mockJWTRepo.On("GenerateAccessToken", userID, user.Email).Return("new_access_token", nil)
	suite.mockJWTRepo.On("GenerateRefreshToken", userID).Return("new_refresh_token", nil)
	suite.mockCacheRepo.On("SetSession", ctx, "new_refresh_token", userID, time.Hour*24*7).Return(nil)

	// Execute
	response, err := suite.authService.RefreshToken(ctx, refreshToken)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), "new_access_token", response.AccessToken)
	assert.Equal(suite.T(), "new_refresh_token", response.RefreshToken)

	suite.mockJWTRepo.AssertExpectations(suite.T())
	suite.mockCacheRepo.AssertExpectations(suite.T())
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestLogout_Success() {
	ctx := context.Background()
	userID := uint(1)

	// Execute
	err := suite.authService.Logout(ctx, userID)

	// Assert
	assert.NoError(suite.T(), err)
}

// Mock implementations

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, offset, limit int) ([]*domain.User, int64, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]*domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) GetUserStats(ctx context.Context, userID uint) (*domain.UserStats, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*domain.UserStats), args.Error(1)
}

type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheService) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}



func (m *MockCacheService) CacheAnalyticsData(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, data, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetAnalyticsData(ctx context.Context, key string) (interface{}, error) {
	args := m.Called(ctx, key)
	return args.Get(0), args.Error(1)
}

// Additional methods to satisfy the ports.CacheService interface
func (m *MockCacheService) Del(ctx context.Context, keys ...string) error {
	return nil
}

func (m *MockCacheService) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (m *MockCacheService) TTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, nil
}

func (m *MockCacheService) Incr(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *MockCacheService) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return 0, nil
}

func (m *MockCacheService) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return nil
}

func (m *MockCacheService) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return false, nil
}

func (m *MockCacheService) SCard(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *MockCacheService) HSet(ctx context.Context, key string, values ...interface{}) error {
	return nil
}

func (m *MockCacheService) HGet(ctx context.Context, key, field string) (string, error) {
	return "", nil
}

func (m *MockCacheService) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return nil, nil
}

func (m *MockCacheService) HDel(ctx context.Context, key string, fields ...string) error {
	return nil
}

func (m *MockCacheService) CacheURL(ctx context.Context, shortCode, originalURL string, userID uint, expiration time.Duration) error {
	args := m.Called(ctx, shortCode, originalURL, userID, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetCachedURL(ctx context.Context, shortCode string) (string, uint, error) {
	args := m.Called(ctx, shortCode)
	return args.String(0), args.Get(1).(uint), args.Error(2)
}

func (m *MockCacheService) InvalidateURL(ctx context.Context, shortCode string) error {
	args := m.Called(ctx, shortCode)
	return args.Error(0)
}

func (m *MockCacheService) IsRateLimited(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	return false, nil
}

func (m *MockCacheService) IncrementRateLimit(ctx context.Context, key string, window time.Duration) (int64, error) {
	return 0, nil
}

func (m *MockCacheService) SetSession(ctx context.Context, token string, userID uint, expiration time.Duration) error {
	args := m.Called(ctx, token, userID, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetSession(ctx context.Context, token string) (uint, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockCacheService) InvalidateSession(ctx context.Context, token string) error {
	return nil
}

func (m *MockCacheService) CacheClickCount(ctx context.Context, shortCode string, count int64) error {
	return nil
}

func (m *MockCacheService) GetClickCount(ctx context.Context, shortCode string) (int64, error) {
	return 0, nil
}

func (m *MockCacheService) IncrementClickCount(ctx context.Context, shortCode string) (int64, error) {
	return 0, nil
}

func (m *MockCacheService) CacheUniqueClick(ctx context.Context, shortCode, ipAddress string) (bool, error) {
	args := m.Called(ctx, shortCode, ipAddress)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheService) GetUniqueClickCount(ctx context.Context, shortCode string) (int64, error) {
	return 0, nil
}

func (m *MockCacheService) Ping(ctx context.Context) error {
	return nil
}

func (m *MockCacheService) FlushDB(ctx context.Context) error {
	return nil
}

func (m *MockCacheService) Info(ctx context.Context) (string, error) {
	return "", nil
}

func (m *MockCacheService) Close() error {
	return nil
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateAccessToken(userID uint, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(userID uint) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateAccessToken(token string) (*domain.TokenClaims, error) {
	args := m.Called(token)
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}

func (m *MockJWTService) ValidateRefreshToken(token string) (*domain.TokenClaims, error) {
	args := m.Called(token)
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}

type MockConfigService struct {
	mock.Mock
}

func (m *MockConfigService) GetBaseURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetJWTSecret() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetDatabaseURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetRedisURL() string {
	args := m.Called()
	return args.String(0)
}