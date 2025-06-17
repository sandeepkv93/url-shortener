package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"url-shortener/internal/core/domain"
)

type AuthMiddlewareTestSuite struct {
	suite.Suite
	middleware   *AuthMiddleware
	mockJWT      *MockJWTService
	mockUserRepo *MockUserRepository
}

func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}

func (suite *AuthMiddlewareTestSuite) SetupTest() {
	suite.mockJWT = &MockJWTService{}
	suite.mockUserRepo = &MockUserRepository{}
	suite.middleware = NewAuthMiddleware(suite.mockJWT, suite.mockUserRepo)
}

func (suite *AuthMiddlewareTestSuite) TestRequireAuth_Success() {
	// Setup
	token := "valid-token"
	userID := uint(1)
	user := &domain.User{
		ID:       userID,
		Email:    "test@example.com",
		IsActive: true,
	}
	
	claims := &domain.TokenClaims{
		UserID: userID,
	}

	suite.mockJWT.On("ValidateAccessToken", token).Return(claims, nil)
	suite.mockUserRepo.On("GetByID", mock.Anything, userID).Return(user, nil)

	// Create handler
	handler := suite.middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify user is in context
		contextUser := GetUserFromContext(r.Context())
		assert.NotNil(suite.T(), contextUser)
		assert.Equal(suite.T(), userID, contextUser.ID)
		
		contextUserID := GetUserIDFromContext(r.Context())
		assert.Equal(suite.T(), userID, contextUserID)
		
		w.WriteHeader(http.StatusOK)
	}))

	// Create request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	suite.mockJWT.AssertExpectations(suite.T())
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *AuthMiddlewareTestSuite) TestRequireAuth_MissingToken() {
	// Create handler
	handler := suite.middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.T().Error("Handler should not be called")
	}))

	// Create request without token
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
	assert.Contains(suite.T(), rr.Body.String(), "Missing authorization token")
}

func (suite *AuthMiddlewareTestSuite) TestRequireAuth_InvalidToken() {
	// Setup
	token := "invalid-token"
	
	suite.mockJWT.On("ValidateAccessToken", token).Return(nil, assert.AnError)

	// Create handler
	handler := suite.middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.T().Error("Handler should not be called")
	}))

	// Create request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
	assert.Contains(suite.T(), rr.Body.String(), "Invalid token")
	suite.mockJWT.AssertExpectations(suite.T())
}

func (suite *AuthMiddlewareTestSuite) TestRequireAuth_UserNotFound() {
	// Setup
	token := "valid-token"
	userID := uint(1)
	
	claims := &domain.TokenClaims{
		UserID: userID,
	}

	suite.mockJWT.On("ValidateAccessToken", token).Return(claims, nil)
	suite.mockUserRepo.On("GetByID", mock.Anything, userID).Return(nil, assert.AnError)

	// Create handler
	handler := suite.middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.T().Error("Handler should not be called")
	}))

	// Create request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
	assert.Contains(suite.T(), rr.Body.String(), "User not found")
	suite.mockJWT.AssertExpectations(suite.T())
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *AuthMiddlewareTestSuite) TestRequireAuth_InactiveUser() {
	// Setup
	token := "valid-token"
	userID := uint(1)
	user := &domain.User{
		ID:       userID,
		Email:    "test@example.com",
		IsActive: false, // Inactive user
	}
	
	claims := &domain.TokenClaims{
		UserID: userID,
	}

	suite.mockJWT.On("ValidateAccessToken", token).Return(claims, nil)
	suite.mockUserRepo.On("GetByID", mock.Anything, userID).Return(user, nil)

	// Create handler
	handler := suite.middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.T().Error("Handler should not be called")
	}))

	// Create request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
	assert.Contains(suite.T(), rr.Body.String(), "User account is inactive")
	suite.mockJWT.AssertExpectations(suite.T())
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *AuthMiddlewareTestSuite) TestOptionalAuth_WithValidToken() {
	// Setup
	token := "valid-token"
	userID := uint(1)
	user := &domain.User{
		ID:       userID,
		Email:    "test@example.com",
		IsActive: true,
	}
	
	claims := &domain.TokenClaims{
		UserID: userID,
	}

	suite.mockJWT.On("ValidateAccessToken", token).Return(claims, nil)
	suite.mockUserRepo.On("GetByID", mock.Anything, userID).Return(user, nil)

	// Create handler
	handler := suite.middleware.OptionalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify user is in context
		contextUser := GetUserFromContext(r.Context())
		assert.NotNil(suite.T(), contextUser)
		assert.Equal(suite.T(), userID, contextUser.ID)
		
		w.WriteHeader(http.StatusOK)
	}))

	// Create request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	suite.mockJWT.AssertExpectations(suite.T())
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *AuthMiddlewareTestSuite) TestOptionalAuth_WithoutToken() {
	// Create handler
	handler := suite.middleware.OptionalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify no user in context
		contextUser := GetUserFromContext(r.Context())
		assert.Nil(suite.T(), contextUser)
		
		contextUserID := GetUserIDFromContext(r.Context())
		assert.Equal(suite.T(), uint(0), contextUserID)
		
		w.WriteHeader(http.StatusOK)
	}))

	// Create request without token
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
}

func (suite *AuthMiddlewareTestSuite) TestOptionalAuth_WithInvalidToken() {
	// Setup
	token := "invalid-token"
	
	suite.mockJWT.On("ValidateAccessToken", token).Return(nil, assert.AnError)

	// Create handler
	handler := suite.middleware.OptionalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify no user in context
		contextUser := GetUserFromContext(r.Context())
		assert.Nil(suite.T(), contextUser)
		
		w.WriteHeader(http.StatusOK)
	}))

	// Create request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	suite.mockJWT.AssertExpectations(suite.T())
}

func (suite *AuthMiddlewareTestSuite) TestExtractToken_FromHeader() {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	
	token := suite.middleware.extractToken(req)
	assert.Equal(suite.T(), "test-token", token)
}

func (suite *AuthMiddlewareTestSuite) TestExtractToken_FromQuery() {
	req := httptest.NewRequest("GET", "/test?token=test-token", nil)
	
	token := suite.middleware.extractToken(req)
	assert.Equal(suite.T(), "test-token", token)
}

func (suite *AuthMiddlewareTestSuite) TestExtractToken_InvalidHeader() {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Invalid test-token")
	
	token := suite.middleware.extractToken(req)
	assert.Equal(suite.T(), "", token)
}

func (suite *AuthMiddlewareTestSuite) TestHelperFunctions() {
	// Test with no context
	ctx := context.Background()
	assert.Nil(suite.T(), GetUserFromContext(ctx))
	assert.Equal(suite.T(), uint(0), GetUserIDFromContext(ctx))
	assert.False(suite.T(), IsAuthenticated(ctx))
	
	// Test with user context
	user := &domain.User{ID: 1, Email: "test@example.com"}
	ctx = context.WithValue(ctx, "user", user)
	ctx = context.WithValue(ctx, "user_id", user.ID)
	
	assert.Equal(suite.T(), user, GetUserFromContext(ctx))
	assert.Equal(suite.T(), user.ID, GetUserIDFromContext(ctx))
	assert.True(suite.T(), IsAuthenticated(ctx))
}

// Mock implementations
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}

func (m *MockJWTService) ValidateRefreshToken(token string) (*domain.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserStats), args.Error(1)
}