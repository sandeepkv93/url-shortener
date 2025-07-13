package handlers

import (
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

type AnalyticsHandlerTestSuite struct {
	suite.Suite
	handler             *AnalyticsHandler
	mockAnalyticsService *MockAnalyticsServiceForTest
}

func TestAnalyticsHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyticsHandlerTestSuite))
}

func (suite *AnalyticsHandlerTestSuite) SetupTest() {
	suite.mockAnalyticsService = &MockAnalyticsServiceForTest{}
	suite.handler = NewAnalyticsHandler(suite.mockAnalyticsService)
}

func (suite *AnalyticsHandlerTestSuite) TestGetDashboardStats_Success() {
	// Setup
	stats := &domain.DashboardStats{
		TotalURLs:       10,
		ActiveURLs:      8,
		TotalClicks:     1000,
		ClickGrowthRate: 15.5,
		URLGrowthRate:   10.2,
		ClicksByDate: map[string]int64{
			"2024-01-01": 50,
			"2024-01-02": 45,
		},
		TopURLs:        []domain.TopURLStat{},
		RecentActivity: []domain.ActivityItem{},
	}

	user := &domain.User{ID: 1, Email: "test@example.com"}
	suite.mockAnalyticsService.On("GetDashboardStats", mock.Anything, uint(1)).Return(stats, nil)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/analytics/dashboard", nil)
	
	// Add user context (both user and user_id as middleware would)
	ctx := context.WithValue(httpReq.Context(), "user", user)
	ctx = context.WithValue(ctx, "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetDashboard(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response domain.DashboardStats
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), stats.TotalURLs, response.TotalURLs)
	assert.Equal(suite.T(), stats.TotalClicks, response.TotalClicks)
	
	suite.mockAnalyticsService.AssertExpectations(suite.T())
}

func (suite *AnalyticsHandlerTestSuite) TestGetDashboardStats_Unauthenticated() {
	// Create request without user context
	httpReq := httptest.NewRequest("GET", "/api/analytics/dashboard", nil)
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetDashboard(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
}

func (suite *AnalyticsHandlerTestSuite) TestGetURLAnalytics_Success() {
	// Setup
	analytics := &domain.URLAnalytics{
		ShortURLID:   1,
		ShortCode:    "test123",
		TotalClicks:  500,
		UniqueClicks: 200,
		ClicksByDate: map[string]int64{
			"2024-01-01": 25,
			"2024-01-02": 30,
		},
		ClicksByTime: map[int]int64{
			9:  10,
			10: 15,
		},
		CountryStats: map[string]int64{
			"US": 100,
			"UK": 50,
		},
		RegionStats: map[string]int64{
			"California": 60,
			"London":     40,
		},
		CityStats: map[string]int64{
			"Los Angeles": 30,
			"London":      40,
		},
		TopDevices: []domain.DeviceStat{
			{Device: "Desktop", Count: 300},
			{Device: "Mobile", Count: 200},
		},
		TopBrowsers: []domain.BrowserStat{
			{Browser: "Chrome", Count: 200},
			{Browser: "Firefox", Count: 100},
		},
		TopReferers: []domain.RefererStat{
			{Referer: "google.com", Count: 150},
			{Referer: "facebook.com", Count: 50},
		},
		RecentClicks: []domain.RecentClickStat{},
	}

	suite.mockAnalyticsService.On("GetURLAnalytics", mock.Anything, uint(1), uint(1)).Return(analytics, nil)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/analytics/url/1", nil)
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetURLAnalytics(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response domain.URLAnalytics
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), analytics.TotalClicks, response.TotalClicks)
	
	suite.mockAnalyticsService.AssertExpectations(suite.T())
}

func (suite *AnalyticsHandlerTestSuite) TestGetClickTimeline_Success() {
	// Setup
	timeline := &domain.TimelineStats{
		Period: "day",
		Data: map[string]int64{
			time.Now().Format("15:04"):                       10,
			time.Now().Add(-1 * time.Hour).Format("15:04"): 15,
		},
	}

	suite.mockAnalyticsService.On("GetClickTimeline", mock.Anything, uint(1), uint(1), "day").Return(timeline, nil)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/analytics/url/1/timeline?period=day", nil)
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetClickTimeline(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response domain.TimelineStats
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), timeline.Period, response.Period)
	assert.Len(suite.T(), response.Data, 2)
	
	suite.mockAnalyticsService.AssertExpectations(suite.T())
}

func (suite *AnalyticsHandlerTestSuite) TestGetGeographicStats_Success() {
	// Setup
	geoStats := &domain.GeoStats{
		CountryStats: map[string]int64{
			"US":     500,
			"UK":     200,
			"Canada": 100,
		},
		RegionStats: map[string]int64{
			"California": 300,
			"London":     200,
			"Ontario":    100,
		},
		CityStats: map[string]int64{
			"New York": 200,
			"London":   150,
			"Toronto":  80,
		},
	}

	suite.mockAnalyticsService.On("GetGeographicStats", mock.Anything, uint(1), uint(1)).Return(geoStats, nil)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/analytics/url/1/geo", nil)
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetGeographicStats(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response domain.GeoStats
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response.CountryStats, 3)
	assert.Equal(suite.T(), geoStats.CountryStats["US"], response.CountryStats["US"])
	
	suite.mockAnalyticsService.AssertExpectations(suite.T())
}

func (suite *AnalyticsHandlerTestSuite) TestGetDeviceStats_Success() {
	// Setup
	deviceStats := &domain.DeviceStats{
		TopDevices: []map[string]interface{}{
			{"device": "Desktop", "count": 300},
			{"device": "Mobile", "count": 200},
			{"device": "Tablet", "count": 50},
		},
		TopBrowsers: []map[string]interface{}{
			{"browser": "Chrome", "count": 250},
			{"browser": "Firefox", "count": 150},
			{"browser": "Safari", "count": 100},
		},
	}

	suite.mockAnalyticsService.On("GetDeviceStats", mock.Anything, uint(1), uint(1)).Return(deviceStats, nil)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/analytics/url/1/devices", nil)
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetDeviceStats(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response domain.DeviceStats
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response.TopDevices, 3)
	assert.Equal(suite.T(), deviceStats.TopDevices[0]["device"], response.TopDevices[0]["device"])
	
	suite.mockAnalyticsService.AssertExpectations(suite.T())
}

func (suite *AnalyticsHandlerTestSuite) TestGetReferrerStats_Success() {
	// Setup
	referrerStats := []domain.RefererStat{
		{Referer: "google.com", Count: 500},
		{Referer: "facebook.com", Count: 300},
		{Referer: "twitter.com", Count: 200},
	}

	suite.mockAnalyticsService.On("GetReferrerStats", mock.Anything, uint(1), uint(1)).Return(referrerStats, nil)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/analytics/url/1/referrers", nil)
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetReferrerStats(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response struct {
		Referrers []domain.RefererStat `json:"referrers"`
	}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response.Referrers, 3)
	assert.Equal(suite.T(), referrerStats[0].Referer, response.Referrers[0].Referer)
	
	suite.mockAnalyticsService.AssertExpectations(suite.T())
}

func (suite *AnalyticsHandlerTestSuite) TestGetTopPerformingURLs_Success() {
	// Setup
	topURLs := []*domain.URLPerformance{
		{
			ShortURL:     "test123",
			OriginalURL:  "https://example.com",
			TotalClicks:  1000,
			UniqueClicks: 500,
			ClickRate:    50.0,
		},
		{
			ShortURL:     "test456",
			OriginalURL:  "https://example.org",
			TotalClicks:  800,
			UniqueClicks: 400,
			ClickRate:    40.0,
		},
	}

	user := &domain.User{ID: 1, Email: "test@example.com"}
	suite.mockAnalyticsService.On("GetTopPerformingURLs", mock.Anything, uint(1), 10).Return(topURLs, nil)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/analytics/top-urls?limit=10", nil)
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user", user)
	ctx = context.WithValue(ctx, "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetTopPerformingURLs(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response struct {
		URLs  []*domain.URLPerformance `json:"urls"`
		Limit int                       `json:"limit"`
	}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response.URLs, 2)
	assert.Equal(suite.T(), 10, response.Limit)
	assert.Equal(suite.T(), topURLs[0].ShortURL, response.URLs[0].ShortURL)
	
	suite.mockAnalyticsService.AssertExpectations(suite.T())
}

func (suite *AnalyticsHandlerTestSuite) TestExportAnalytics_Success() {
	// Setup
	exportData := []byte("date,url,clicks\n2024-01-01,test123,100\n")

	suite.mockAnalyticsService.On("ExportAnalytics", mock.Anything, uint(1), "csv", mock.AnythingOfType("domain.DateRange")).Return(exportData, nil)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/analytics/export?format=csv&start_date=2024-01-01&end_date=2024-01-07", nil)
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.ExportAnalytics(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	assert.Equal(suite.T(), "text/csv", rr.Header().Get("Content-Type"))
	assert.Contains(suite.T(), rr.Header().Get("Content-Disposition"), "analytics.csv")
	assert.Equal(suite.T(), string(exportData), rr.Body.String())
	
	suite.mockAnalyticsService.AssertExpectations(suite.T())
}

func (suite *AnalyticsHandlerTestSuite) TestGetGlobalStats_Success() {
	// Setup
	globalStats := &domain.GlobalStats{
		TotalUsers:       1000,
		TotalURLs:        5000,
		TotalClicks:      100000,
		ActiveURLs:       4500,
		ClicksToday:      5000,
		URLsCreatedToday: 50,
		NewUsersToday:    10,
	}

	suite.mockAnalyticsService.On("GetGlobalStats", mock.Anything).Return(globalStats, nil)

	// Create request (admin endpoint, but we'll test as authenticated for now)
	httpReq := httptest.NewRequest("GET", "/api/analytics/global", nil)
	
	// Add user context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(1))
	httpReq = httpReq.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetGlobalStats(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response domain.GlobalStats
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), globalStats.TotalUsers, response.TotalUsers)
	
	suite.mockAnalyticsService.AssertExpectations(suite.T())
}

// Mock implementation for analytics service
type MockAnalyticsServiceForTest struct {
	mock.Mock
}

func (m *MockAnalyticsServiceForTest) GetDashboardStats(ctx context.Context, userID uint) (*domain.DashboardStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DashboardStats), args.Error(1)
}

func (m *MockAnalyticsServiceForTest) GetGlobalStats(ctx context.Context) (*domain.GlobalStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GlobalStats), args.Error(1)
}

func (m *MockAnalyticsServiceForTest) GetURLAnalytics(ctx context.Context, shortURLID uint, userID uint) (*domain.URLAnalytics, error) {
	args := m.Called(ctx, shortURLID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.URLAnalytics), args.Error(1)
}

func (m *MockAnalyticsServiceForTest) GetTopPerformingURLs(ctx context.Context, userID uint, limit int) ([]*domain.URLPerformance, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.URLPerformance), args.Error(1)
}

func (m *MockAnalyticsServiceForTest) GetClickTimeline(ctx context.Context, shortURLID uint, userID uint, period string) (*domain.TimelineStats, error) {
	args := m.Called(ctx, shortURLID, userID, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TimelineStats), args.Error(1)
}

func (m *MockAnalyticsServiceForTest) GetGeographicStats(ctx context.Context, shortURLID uint, userID uint) (*domain.GeoStats, error) {
	args := m.Called(ctx, shortURLID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GeoStats), args.Error(1)
}

func (m *MockAnalyticsServiceForTest) GetDeviceStats(ctx context.Context, shortURLID uint, userID uint) (*domain.DeviceStats, error) {
	args := m.Called(ctx, shortURLID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DeviceStats), args.Error(1)
}

func (m *MockAnalyticsServiceForTest) GetReferrerStats(ctx context.Context, shortURLID uint, userID uint) ([]domain.RefererStat, error) {
	args := m.Called(ctx, shortURLID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.RefererStat), args.Error(1)
}

func (m *MockAnalyticsServiceForTest) ExportAnalytics(ctx context.Context, userID uint, format string, dateRange domain.DateRange) ([]byte, error) {
	args := m.Called(ctx, userID, format, dateRange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}