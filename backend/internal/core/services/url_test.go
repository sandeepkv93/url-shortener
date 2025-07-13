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

type URLServiceTestSuite struct {
	suite.Suite
	urlService    *urlService
	mockURLRepo   *MockURLRepository
	mockClickRepo *MockClickRepository
	mockCacheRepo *MockCacheService
	mockConfigRepo *MockConfigService
}

func TestURLServiceSuite(t *testing.T) {
	suite.Run(t, new(URLServiceTestSuite))
}

func (suite *URLServiceTestSuite) SetupTest() {
	suite.mockURLRepo = &MockURLRepository{}
	suite.mockClickRepo = &MockClickRepository{}
	suite.mockCacheRepo = &MockCacheService{}
	suite.mockConfigRepo = &MockConfigService{}
	
	suite.urlService = &urlService{
		urlRepo:    suite.mockURLRepo,
		clickRepo:  suite.mockClickRepo,
		cacheRepo:  suite.mockCacheRepo,
		configRepo: suite.mockConfigRepo,
	}
}

func (suite *URLServiceTestSuite) TestShortenURL_Success() {
	ctx := context.Background()
	req := domain.ShortenURLRequest{
		OriginalURL: "https://example.com",
		UserID:      1,
		Title:       "Example Site",
	}

	// Mock expectations
	suite.mockURLRepo.On("ExistsByShortCode", ctx, mock.AnythingOfType("string")).Return(false, nil)
	suite.mockURLRepo.On("Create", ctx, mock.AnythingOfType("*domain.ShortURL")).Return(nil)
	suite.mockCacheRepo.On("CacheURL", ctx, mock.AnythingOfType("string"), req.OriginalURL, req.UserID, time.Hour*24).Return(nil)

	// Execute
	result, err := suite.urlService.ShortenURL(ctx, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), req.OriginalURL, result.OriginalURL)
	assert.Equal(suite.T(), req.UserID, result.UserID)
	assert.Equal(suite.T(), req.Title, result.Title)
	assert.True(suite.T(), result.IsActive)
	assert.NotEmpty(suite.T(), result.ShortCode)

	suite.mockURLRepo.AssertExpectations(suite.T())
	suite.mockCacheRepo.AssertExpectations(suite.T())
}

func (suite *URLServiceTestSuite) TestShortenURL_InvalidURL() {
	ctx := context.Background()
	req := domain.ShortenURLRequest{
		OriginalURL: "invalid-url",
		UserID:      1,
	}

	// Execute
	result, err := suite.urlService.ShortenURL(ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), domain.ErrInvalidURL, err)
}

func (suite *URLServiceTestSuite) TestShortenURL_CustomAlias() {
	ctx := context.Background()
	req := domain.ShortenURLRequest{
		OriginalURL: "https://example.com",
		UserID:      1,
		CustomAlias: "custom",
	}

	// Mock expectations
	suite.mockURLRepo.On("ExistsByShortCode", ctx, "custom").Return(false, nil)
	suite.mockURLRepo.On("Create", ctx, mock.AnythingOfType("*domain.ShortURL")).Return(nil)
	suite.mockCacheRepo.On("CacheURL", ctx, "custom", req.OriginalURL, req.UserID, time.Hour*24).Return(nil)

	// Execute
	result, err := suite.urlService.ShortenURL(ctx, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "custom", result.ShortCode)

	suite.mockURLRepo.AssertExpectations(suite.T())
	suite.mockCacheRepo.AssertExpectations(suite.T())
}

func (suite *URLServiceTestSuite) TestGetOriginalURL_FromCache() {
	ctx := context.Background()
	shortCode := "abc123"
	originalURL := "https://example.com"

	shortURL := &domain.ShortURL{
		ID:          1,
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		IsActive:    true,
	}

	// Mock expectations
	suite.mockCacheRepo.On("GetCachedURL", ctx, shortCode).Return(originalURL, uint(1), nil)
	suite.mockURLRepo.On("GetActiveByShortCode", ctx, shortCode).Return(shortURL, nil)

	// Execute
	result, err := suite.urlService.GetOriginalURL(ctx, shortCode)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), originalURL, result.OriginalURL)
	assert.Equal(suite.T(), shortCode, result.ShortCode)

	suite.mockCacheRepo.AssertExpectations(suite.T())
	suite.mockURLRepo.AssertExpectations(suite.T())
}

func (suite *URLServiceTestSuite) TestGetOriginalURL_FromDatabase() {
	ctx := context.Background()
	shortCode := "abc123"
	originalURL := "https://example.com"

	shortURL := &domain.ShortURL{
		ID:          1,
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		IsActive:    true,
	}

	// Mock expectations - cache miss, then database hit
	suite.mockCacheRepo.On("GetCachedURL", ctx, shortCode).Return("", uint(0), assert.AnError)
	suite.mockURLRepo.On("GetActiveByShortCode", ctx, shortCode).Return(shortURL, nil)
	suite.mockCacheRepo.On("CacheURL", ctx, shortCode, originalURL, shortURL.UserID, time.Hour*24).Return(nil)

	// Execute
	result, err := suite.urlService.GetOriginalURL(ctx, shortCode)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), originalURL, result.OriginalURL)

	suite.mockCacheRepo.AssertExpectations(suite.T())
	suite.mockURLRepo.AssertExpectations(suite.T())
}

func (suite *URLServiceTestSuite) TestRecordClick_Success() {
	ctx := context.Background()
	shortURL := &domain.ShortURL{
		ID:        1,
		ShortCode: "abc123",
	}
	clickData := domain.ClickData{
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
		Country:   "US",
	}

	// Mock expectations
	suite.mockClickRepo.On("Create", ctx, mock.AnythingOfType("*domain.Click")).Return(nil)
	suite.mockURLRepo.On("IncrementClickCount", ctx, shortURL.ID).Return(nil)
	suite.mockCacheRepo.On("CacheUniqueClick", ctx, mock.AnythingOfType("string"), clickData.IPAddress).Return(false, nil)

	// Execute
	err := suite.urlService.RecordClick(ctx, shortURL, clickData)

	// Assert
	assert.NoError(suite.T(), err)

	suite.mockClickRepo.AssertExpectations(suite.T())
	suite.mockURLRepo.AssertExpectations(suite.T())
	suite.mockCacheRepo.AssertExpectations(suite.T())
}

func (suite *URLServiceTestSuite) TestUpdateURL_Success() {
	ctx := context.Background()
	urlID := uint(1)
	userID := uint(1)
	title := "Updated Title"
	isActive := false

	req := domain.UpdateURLRequest{
		Title:    &title,
		IsActive: &isActive,
	}

	existingURL := &domain.ShortURL{
		ID:        urlID,
		UserID:    userID,
		ShortCode: "abc123",
		Title:     "Original Title",
		IsActive:  true,
	}

	// Mock expectations
	suite.mockURLRepo.On("GetByID", ctx, urlID).Return(existingURL, nil)
	suite.mockURLRepo.On("Update", ctx, mock.AnythingOfType("*domain.ShortURL")).Return(nil)
	suite.mockCacheRepo.On("InvalidateURL", ctx, existingURL.ShortCode).Return(nil)

	// Execute
	result, err := suite.urlService.UpdateURL(ctx, urlID, userID, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), title, result.Title)
	assert.Equal(suite.T(), isActive, result.IsActive)

	suite.mockURLRepo.AssertExpectations(suite.T())
	suite.mockCacheRepo.AssertExpectations(suite.T())
}

func (suite *URLServiceTestSuite) TestUpdateURL_Unauthorized() {
	ctx := context.Background()
	urlID := uint(1)
	userID := uint(2) // Different user
	title := "Updated Title"

	req := domain.UpdateURLRequest{
		Title: &title,
	}

	existingURL := &domain.ShortURL{
		ID:     urlID,
		UserID: 1, // Different from userID
	}

	// Mock expectations
	suite.mockURLRepo.On("GetByID", ctx, urlID).Return(existingURL, nil)

	// Execute
	result, err := suite.urlService.UpdateURL(ctx, urlID, userID, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), domain.ErrUnauthorized, err)

	suite.mockURLRepo.AssertExpectations(suite.T())
}

func (suite *URLServiceTestSuite) TestDeleteURL_Success() {
	ctx := context.Background()
	urlID := uint(1)
	userID := uint(1)

	existingURL := &domain.ShortURL{
		ID:        urlID,
		UserID:    userID,
		ShortCode: "abc123",
	}

	// Mock expectations
	suite.mockURLRepo.On("GetByID", ctx, urlID).Return(existingURL, nil)
	suite.mockURLRepo.On("Delete", ctx, urlID).Return(nil)
	suite.mockCacheRepo.On("InvalidateURL", ctx, existingURL.ShortCode).Return(nil)

	// Execute
	err := suite.urlService.DeleteURL(ctx, urlID, userID)

	// Assert
	assert.NoError(suite.T(), err)

	suite.mockURLRepo.AssertExpectations(suite.T())
	suite.mockCacheRepo.AssertExpectations(suite.T())
}

func (suite *URLServiceTestSuite) TestValidatePassword_NoPassword() {
	ctx := context.Background()
	shortCode := "abc123"

	shortURL := &domain.ShortURL{
		ID:        1,
		ShortCode: shortCode,
		Password:  nil, // No password set
	}

	// Mock expectations
	suite.mockURLRepo.On("GetByShortCode", ctx, shortCode).Return(shortURL, nil)

	// Execute
	valid, err := suite.urlService.ValidatePassword(ctx, shortCode, "anypassword")

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), valid)

	suite.mockURLRepo.AssertExpectations(suite.T())
}
