package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"url-shortener/internal/core/domain"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	handler         *AuthHandler
	mockAuthService *MockAuthService
}

func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}

func (suite *AuthHandlerTestSuite) SetupTest() {
	suite.mockAuthService = &MockAuthService{}
	suite.handler = NewAuthHandler(suite.mockAuthService)
}

func (suite *AuthHandlerTestSuite) TestRegister_Success() {
	req := domain.RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	response := &domain.AuthResponse{
		User: &domain.UserResponse{
			ID:        1,
			Email:     req.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		},
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	suite.mockAuthService.On("Register", mock.Anything, req).Return(response, nil)

	// Create request
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Register(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, rr.Code)
	
	var result domain.AuthResponse
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), req.Email, result.User.Email)
	assert.Equal(suite.T(), "access_token", result.AccessToken)

	suite.mockAuthService.AssertExpectations(suite.T())
}

func (suite *AuthHandlerTestSuite) TestRegister_UserAlreadyExists() {
	req := domain.RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	suite.mockAuthService.On("Register", mock.Anything, req).Return(nil, domain.ErrUserAlreadyExists)

	// Create request
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Register(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusConflict, rr.Code)
	assert.Contains(suite.T(), rr.Body.String(), "User already exists")

	suite.mockAuthService.AssertExpectations(suite.T())
}

func (suite *AuthHandlerTestSuite) TestRegister_InvalidJSON() {
	httpReq := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Register(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
	assert.Contains(suite.T(), rr.Body.String(), "Invalid request body")
}

func (suite *AuthHandlerTestSuite) TestLogin_Success() {
	req := domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	response := &domain.AuthResponse{
		User: &domain.UserResponse{
			ID:    1,
			Email: req.Email,
		},
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	suite.mockAuthService.On("Login", mock.Anything, req).Return(response, nil)

	// Create request
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Login(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var result domain.AuthResponse
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), req.Email, result.User.Email)
	assert.Equal(suite.T(), "access_token", result.AccessToken)

	suite.mockAuthService.AssertExpectations(suite.T())
}

func (suite *AuthHandlerTestSuite) TestLogin_InvalidCredentials() {
	req := domain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	suite.mockAuthService.On("Login", mock.Anything, req).Return(nil, domain.ErrInvalidCredentials)

	// Create request
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Login(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
	assert.Contains(suite.T(), rr.Body.String(), "Invalid email or password")

	suite.mockAuthService.AssertExpectations(suite.T())
}

func (suite *AuthHandlerTestSuite) TestRefreshToken_Success() {
	reqBody := map[string]string{
		"refresh_token": "refresh_token_123",
	}

	response := &domain.AuthResponse{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	suite.mockAuthService.On("RefreshToken", mock.Anything, "refresh_token_123").Return(response, nil)

	// Create request
	body, _ := json.Marshal(reqBody)
	httpReq := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.RefreshToken(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var result domain.AuthResponse
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "new_access_token", result.AccessToken)

	suite.mockAuthService.AssertExpectations(suite.T())
}

func (suite *AuthHandlerTestSuite) TestRefreshToken_MissingToken() {
	reqBody := map[string]string{}

	// Create request
	body, _ := json.Marshal(reqBody)
	httpReq := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.RefreshToken(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
	assert.Contains(suite.T(), rr.Body.String(), "Refresh token is required")
}

func (suite *AuthHandlerTestSuite) TestLogout_Success() {
	userID := uint(1)
	suite.mockAuthService.On("Logout", mock.Anything, userID).Return(nil)

	// Create request with user context
	httpReq := httptest.NewRequest("POST", "/auth/logout", nil)
	ctx := context.WithValue(httpReq.Context(), "user_id", userID)
	httpReq = httpReq.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Logout(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	assert.Contains(suite.T(), rr.Body.String(), "Logged out successfully")

	suite.mockAuthService.AssertExpectations(suite.T())
}

func (suite *AuthHandlerTestSuite) TestLogout_Unauthenticated() {
	// Create request without user context
	httpReq := httptest.NewRequest("POST", "/auth/logout", nil)
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Logout(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
	assert.Contains(suite.T(), rr.Body.String(), "Authentication required")
}

// Mock AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}

func (m *MockAuthService) GetProfile(ctx context.Context, userID uint) (*domain.UserResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

func (m *MockAuthService) UpdateProfile(ctx context.Context, userID uint, req domain.UpdateProfileRequest) (*domain.UserResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

func (m *MockAuthService) ChangePassword(ctx context.Context, userID uint, req domain.ChangePasswordRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}