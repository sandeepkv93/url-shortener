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

