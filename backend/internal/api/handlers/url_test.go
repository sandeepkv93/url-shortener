package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"url-shortener/internal/core/domain"
)

type URLHandlerTestSuite struct {
	suite.Suite
	handler          *URLHandler
	mockURLService   *MockURLService
	mockAnalytics    *MockAnalyticsService
}

func TestURLHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(URLHandlerTestSuite))
}

func (suite *URLHandlerTestSuite) SetupTest() {
	suite.mockURLService = &MockURLService{}
	suite.mockAnalytics = &MockAnalyticsService{}
	suite.handler = NewURLHandler(suite.mockURLService, suite.mockAnalytics)
}

func (suite *URLHandlerTestSuite) TestCreateShortURL_Success() {
	// Setup
	req := domain.ShortenURLRequest{
		OriginalURL: "https://example.com",
		CustomAlias: "test",
	}
	shortURL := &domain.ShortURL{
		ID:          1,
		ShortCode:   "test123",
		OriginalURL: req.OriginalURL,
		UserID:      1,
		CreatedAt:   time.Now(),
	}

	suite.mockURLService.On("ShortenURL", mock.Anything, mock.AnythingOfType("domain.ShortenURLRequest")).Return(shortURL, nil)

	// Create request
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/urls", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.CreateShortURL(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, rr.Code)
	
	var response domain.ShortURL
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), shortURL.ShortCode, response.ShortCode)
	
	suite.mockURLService.AssertExpectations(suite.T())
}

func (suite *URLHandlerTestSuite) TestCreateShortURL_Unauthenticated() {
	// Create request without user context
	req := domain.ShortenURLRequest{
		OriginalURL: "https://example.com",
	}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/urls", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.CreateShortURL(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
}

func (suite *URLHandlerTestSuite) TestCreateShortURL_InvalidRequest() {
	// Create request with invalid JSON
	httpReq := httptest.NewRequest("POST", "/api/urls", bytes.NewBufferString("invalid json"))
	httpReq.Header.Set("Content-Type", "application/json")
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.CreateShortURL(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

func (suite *URLHandlerTestSuite) TestGetUserURLs_Success() {
	// Setup
	urls := []*domain.ShortURL{
		{
			ID:          1,
			ShortCode:   "test123",
			OriginalURL: "https://example.com",
			UserID:      1,
		},
		{
			ID:          2,
			ShortCode:   "test456",
			OriginalURL: "https://example.org",
			UserID:      1,
		},
	}

	suite.mockURLService.On("GetUserURLs", mock.Anything, uint(1), 0, 20).Return(urls, int64(2), nil)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/urls", nil)
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetUserURLs(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response struct {
		URLs  []*domain.ShortURL `json:"urls"`
		Total int64              `json:"total"`
		Page  int                `json:"page"`
		Limit int                `json:"limit"`
	}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response.URLs, 2)
	assert.Equal(suite.T(), int64(2), response.Total)
	
	suite.mockURLService.AssertExpectations(suite.T())
}

func (suite *URLHandlerTestSuite) TestUpdateURL_Success() {
	// Setup
	title := "Updated Title"
	isActive := true
	updateReq := domain.UpdateURLRequest{
		Title:    &title,
		IsActive: &isActive,
	}
	updatedURL := &domain.ShortURL{
		ID:        1,
		ShortCode: "test123",
		Title:     title,
		UserID:    1,
		IsActive:  isActive,
	}

	suite.mockURLService.On("UpdateURL", mock.Anything, uint(1), uint(1), updateReq).Return(updatedURL, nil)

	// Create request
	body, _ := json.Marshal(updateReq)
	httpReq := httptest.NewRequest("PUT", "/api/urls/1", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.UpdateURL(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response domain.ShortURL
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedURL.Title, response.Title)
	
	suite.mockURLService.AssertExpectations(suite.T())
}

func (suite *URLHandlerTestSuite) TestDeleteURL_Success() {
	// Setup
	suite.mockURLService.On("DeleteURL", mock.Anything, uint(1), uint(1)).Return(nil)

	// Create request
	httpReq := httptest.NewRequest("DELETE", "/api/urls/1", nil)
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.DeleteURL(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	suite.mockURLService.AssertExpectations(suite.T())
}

func (suite *URLHandlerTestSuite) TestRedirectURL_Success() {
	// Setup
	shortURL := &domain.ShortURL{
		ID:          1,
		ShortCode:   "test123",
		OriginalURL: "https://example.com",
		IsActive:    true,
	}

	suite.mockURLService.On("GetOriginalURL", mock.Anything, "test123").Return(shortURL, nil)
	suite.mockURLService.On("RecordClick", mock.Anything, shortURL, mock.AnythingOfType("domain.ClickData")).Return(nil)

	// Create request
	httpReq := httptest.NewRequest("GET", "/test123", nil)
	httpReq.RemoteAddr = "192.168.1.1:12345"
	httpReq.Header.Set("User-Agent", "test-agent")
	httpReq.Header.Set("Referer", "https://referrer.com")
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("shortCode", "test123")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.RedirectURL(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusMovedPermanently, rr.Code)
	assert.Equal(suite.T(), shortURL.OriginalURL, rr.Header().Get("Location"))
	
	suite.mockURLService.AssertExpectations(suite.T())
}

func (suite *URLHandlerTestSuite) TestRedirectURL_NotFound() {
	// Setup
	suite.mockURLService.On("GetOriginalURL", mock.Anything, "invalid").Return(nil, domain.ErrURLNotFound)

	// Create request
	httpReq := httptest.NewRequest("GET", "/invalid", nil)
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("shortCode", "invalid")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.RedirectURL(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, rr.Code)
	
	suite.mockURLService.AssertExpectations(suite.T())
}

func (suite *URLHandlerTestSuite) TestGetURL_Success() {
	// Setup
	urlStats := &domain.URLStats{
		ShortCode:    "test123",
		TotalClicks:  100,
		URL: &domain.ShortURL{
			ID:        1,
			ShortCode: "test123",
			UserID:    1,
		},
	}

	suite.mockURLService.On("GetURLStats", mock.Anything, uint(1), uint(1)).Return(urlStats, nil)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/urls/1", nil)
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetURL(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response domain.URLStats
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), urlStats.ShortCode, response.ShortCode)
	
	suite.mockURLService.AssertExpectations(suite.T())
}

// Mock implementations
type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) ShortenURL(ctx context.Context, req domain.ShortenURLRequest) (*domain.ShortURL, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ShortURL), args.Error(1)
}

func (m *MockURLService) GetOriginalURL(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	args := m.Called(ctx, shortCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ShortURL), args.Error(1)
}

func (m *MockURLService) GetUserURLs(ctx context.Context, userID uint, offset, limit int) ([]*domain.ShortURL, int64, error) {
	args := m.Called(ctx, userID, offset, limit)
	return args.Get(0).([]*domain.ShortURL), args.Get(1).(int64), args.Error(2)
}

func (m *MockURLService) UpdateURL(ctx context.Context, id uint, userID uint, req domain.UpdateURLRequest) (*domain.ShortURL, error) {
	args := m.Called(ctx, id, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ShortURL), args.Error(1)
}

func (m *MockURLService) DeleteURL(ctx context.Context, id uint, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockURLService) RecordClick(ctx context.Context, shortURL *domain.ShortURL, clickData domain.ClickData) error {
	args := m.Called(ctx, shortURL, clickData)
	return args.Error(0)
}

func (m *MockURLService) ValidatePassword(ctx context.Context, shortCode, password string) (bool, error) {
	args := m.Called(ctx, shortCode, password)
	return args.Bool(0), args.Error(1)
}

func (m *MockURLService) GetURLStats(ctx context.Context, id uint, userID uint) (*domain.URLStats, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.URLStats), args.Error(1)
}

func (m *MockURLService) GetPopularURLs(ctx context.Context, limit int) ([]*domain.ShortURL, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*domain.ShortURL), args.Error(1)
}

func (m *MockURLService) CleanupExpiredURLs(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockAnalyticsService implements the analytics service interface
type MockAnalyticsService struct {
	mock.Mock
}

func (m *MockAnalyticsService) GetDashboardStats(ctx context.Context, userID uint) (*domain.DashboardStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DashboardStats), args.Error(1)
}

func (m *MockAnalyticsService) GetGlobalStats(ctx context.Context) (*domain.GlobalStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GlobalStats), args.Error(1)
}

func (m *MockAnalyticsService) GetURLAnalytics(ctx context.Context, shortURLID uint, userID uint) (*domain.URLAnalytics, error) {
	args := m.Called(ctx, shortURLID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.URLAnalytics), args.Error(1)
}

func (m *MockAnalyticsService) GetTopPerformingURLs(ctx context.Context, userID uint, limit int) ([]*domain.URLPerformance, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]*domain.URLPerformance), args.Error(1)
}

func (m *MockAnalyticsService) GetClickTimeline(ctx context.Context, shortURLID uint, userID uint, period string) (*domain.TimelineStats, error) {
	args := m.Called(ctx, shortURLID, userID, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TimelineStats), args.Error(1)
}

func (m *MockAnalyticsService) GetGeographicStats(ctx context.Context, shortURLID uint, userID uint) (*domain.GeoStats, error) {
	args := m.Called(ctx, shortURLID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GeoStats), args.Error(1)
}

func (m *MockAnalyticsService) GetDeviceStats(ctx context.Context, shortURLID uint, userID uint) (*domain.DeviceStats, error) {
	args := m.Called(ctx, shortURLID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DeviceStats), args.Error(1)
}

func (m *MockAnalyticsService) GetReferrerStats(ctx context.Context, shortURLID uint, userID uint) ([]domain.RefererStat, error) {
	args := m.Called(ctx, shortURLID, userID)
	return args.Get(0).([]domain.RefererStat), args.Error(1)
}

func (m *MockAnalyticsService) ExportAnalytics(ctx context.Context, userID uint, format string, dateRange domain.DateRange) ([]byte, error) {
	args := m.Called(ctx, userID, format, dateRange)
	return args.Get(0).([]byte), args.Error(1)
}