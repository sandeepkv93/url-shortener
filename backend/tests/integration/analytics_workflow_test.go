package integration

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// AnalyticsWorkflowTestSuite tests complete analytics workflows
type AnalyticsWorkflowTestSuite struct {
	IntegrationTestSuite
}

// TestClickTrackingWorkflow tests the complete click tracking process
func (s *AnalyticsWorkflowTestSuite) TestClickTrackingWorkflow() {
	// Create a test URL
	urlData := s.createTestURL("https://example.com/analytics-test")
	shortCode := urlData["shortCode"].(string)
	urlID := strconv.Itoa(int(urlData["id"].(float64)))
	
	// Simulate multiple clicks with different user agents and referrers
	clicks := []struct {
		userAgent string
		referer   string
	}{
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36", "https://google.com"},
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15", "https://facebook.com"},
		{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36", "https://twitter.com"},
		{"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X)", ""},
		{"Mozilla/5.0 (Android 11; Mobile; rv:68.0) Gecko/68.0", "https://reddit.com"},
	}
	
	// Perform clicks
	for _, click := range clicks {
		headers := map[string]string{}
		if click.userAgent != "" {
			headers["User-Agent"] = click.userAgent
		}
		if click.referer != "" {
			headers["Referer"] = click.referer
		}
		
		resp := s.makeRequest("GET", "/"+shortCode, nil, headers)
		s.Equal(http.StatusFound, resp.StatusCode)
		resp.Body.Close()
		
		// Small delay between clicks to simulate real usage
		time.Sleep(50 * time.Millisecond)
	}
	
	// Allow time for async click processing
	time.Sleep(200 * time.Millisecond)
	
	// Verify clicks were recorded
	urlDetails := s.getURLDetails(urlData["id"])
	s.Equal(len(clicks), int(urlDetails["clickCount"].(float64)))
	
	// Test URL-specific analytics
	analyticsResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/urls/"+urlID, nil)
	
	var analyticsData struct {
		URL         map[string]interface{} `json:"url"`
		Stats       map[string]interface{} `json:"stats"`
		TopCountries []map[string]interface{} `json:"topCountries"`
		TopDevices   []map[string]interface{} `json:"topDevices"`
		TopReferrers []map[string]interface{} `json:"topReferrers"`
		RecentClicks []map[string]interface{} `json:"recentClicks"`
	}
	s.assertJSONResponse(analyticsResp, http.StatusOK, &analyticsData)
	
	// Verify analytics data
	s.Equal(len(clicks), int(analyticsData.Stats["totalClicks"].(float64)))
	s.GreaterOrEqual(int(analyticsData.Stats["uniqueClicks"].(float64)), 1)
	s.GreaterOrEqual(len(analyticsData.RecentClicks), 1)
}

// TestDashboardAnalytics tests the analytics dashboard endpoint
func (s *AnalyticsWorkflowTestSuite) TestDashboardAnalytics() {
	// Create multiple URLs and simulate activity
	for i := 0; i < 3; i++ {
		urlData := s.createTestURL(fmt.Sprintf("https://example.com/dashboard-test-%d", i))
		shortCode := urlData["shortCode"].(string)
		
		// Simulate clicks for each URL
		for j := 0; j < (i+1)*2; j++ {
			resp := s.makeRequest("GET", "/"+shortCode, nil, nil)
			resp.Body.Close()
			time.Sleep(20 * time.Millisecond)
		}
	}
	
	// Allow time for processing
	time.Sleep(200 * time.Millisecond)
	
	// Test dashboard analytics
	resp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/dashboard", nil)
	
	var dashboardData struct {
		Summary struct {
			TotalURLs        int     `json:"totalUrls"`
			TotalClicks      int     `json:"totalClicks"`
			UniqueVisitors   int     `json:"uniqueVisitors"`
			AvgClicksPerURL  float64 `json:"avgClicksPerUrl"`
		} `json:"summary"`
		RecentActivity   []map[string]interface{} `json:"recentActivity"`
		TopPerforming    []map[string]interface{} `json:"topPerforming"`
		ClickTrend       []map[string]interface{} `json:"clickTrend"`
		DeviceBreakdown  []map[string]interface{} `json:"deviceBreakdown"`
		CountryBreakdown []map[string]interface{} `json:"countryBreakdown"`
	}
	s.assertJSONResponse(resp, http.StatusOK, &dashboardData)
	
	// Verify dashboard data
	s.GreaterOrEqual(dashboardData.Summary.TotalURLs, 3)
	s.GreaterOrEqual(dashboardData.Summary.TotalClicks, 6) // 2+4+6 clicks minimum
	s.Greater(dashboardData.Summary.AvgClicksPerURL, 0.0)
	s.GreaterOrEqual(len(dashboardData.TopPerforming), 1)
}

// TestGeographicAnalytics tests geographic analytics data
func (s *AnalyticsWorkflowTestSuite) TestGeographicAnalytics() {
	// Create a test URL
	urlData := s.createTestURL("https://example.com/geo-test")
	shortCode := urlData["shortCode"].(string)
	urlID := strconv.Itoa(int(urlData["id"].(float64)))
	
	// Simulate clicks from different IP addresses (simulated geographic locations)
	// In a real scenario, these would come from actual different locations
	ipAddresses := []string{
		"192.168.1.1",   // This would typically be mapped to a country
		"10.0.0.1",      
		"172.16.0.1",    
		"203.0.113.1",   
	}
	
	for _, ip := range ipAddresses {
		headers := map[string]string{
			"X-Forwarded-For": ip,
			"X-Real-IP":       ip,
		}
		
		resp := s.makeRequest("GET", "/"+shortCode, nil, headers)
		resp.Body.Close()
		time.Sleep(50 * time.Millisecond)
	}
	
	// Allow time for processing
	time.Sleep(200 * time.Millisecond)
	
	// Test geographic analytics endpoint
	geoResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/urls/"+urlID+"/geo", nil)
	
	var geoData struct {
		Countries []map[string]interface{} `json:"countries"`
		Cities    []map[string]interface{} `json:"cities"`
		Regions   []map[string]interface{} `json:"regions"`
		Total     int                      `json:"total"`
	}
	s.assertJSONResponse(geoResp, http.StatusOK, &geoData)
	
	// Verify geographic data structure
	s.Equal(len(ipAddresses), geoData.Total)
	// Note: In a real implementation with geolocation service, 
	// we would verify actual country/city mappings
}

// TestDeviceAnalytics tests device and browser analytics
func (s *AnalyticsWorkflowTestSuite) TestDeviceAnalytics() {
	// Create a test URL
	urlData := s.createTestURL("https://example.com/device-test")
	shortCode := urlData["shortCode"].(string)
	urlID := strconv.Itoa(int(urlData["id"].(float64)))
	
	// Simulate clicks from different devices and browsers
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Android 11; Mobile; rv:68.0) Gecko/68.0 Firefox/88.0",
	}
	
	for _, userAgent := range userAgents {
		headers := map[string]string{
			"User-Agent": userAgent,
		}
		
		resp := s.makeRequest("GET", "/"+shortCode, nil, headers)
		resp.Body.Close()
		time.Sleep(50 * time.Millisecond)
	}
	
	// Allow time for processing
	time.Sleep(200 * time.Millisecond)
	
	// Test device analytics endpoint
	deviceResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/urls/"+urlID+"/devices", nil)
	
	var deviceData struct {
		Browsers []map[string]interface{} `json:"browsers"`
		OS       []map[string]interface{} `json:"os"`
		Devices  []map[string]interface{} `json:"devices"`
		Total    int                      `json:"total"`
	}
	s.assertJSONResponse(deviceResp, http.StatusOK, &deviceData)
	
	// Verify device data structure
	s.Equal(len(userAgents), deviceData.Total)
	s.GreaterOrEqual(len(deviceData.Browsers), 1)
	s.GreaterOrEqual(len(deviceData.OS), 1)
	s.GreaterOrEqual(len(deviceData.Devices), 1)
}

// TestReferrerAnalytics tests referrer analytics data
func (s *AnalyticsWorkflowTestSuite) TestReferrerAnalytics() {
	// Create a test URL
	urlData := s.createTestURL("https://example.com/referrer-test")
	shortCode := urlData["shortCode"].(string)
	urlID := strconv.Itoa(int(urlData["id"].(float64)))
	
	// Simulate clicks from different referrers
	referrers := []string{
		"https://google.com/search?q=test",
		"https://facebook.com/posts/123",
		"https://twitter.com/user/status/456",
		"https://reddit.com/r/programming",
		"", // Direct visit
	}
	
	for _, referrer := range referrers {
		headers := map[string]string{}
		if referrer != "" {
			headers["Referer"] = referrer
		}
		
		resp := s.makeRequest("GET", "/"+shortCode, nil, headers)
		resp.Body.Close()
		time.Sleep(50 * time.Millisecond)
	}
	
	// Allow time for processing
	time.Sleep(200 * time.Millisecond)
	
	// Test referrer analytics endpoint
	referrerResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/urls/"+urlID+"/referrers", nil)
	
	var referrerData struct {
		Referrers []map[string]interface{} `json:"referrers"`
		Categories []map[string]interface{} `json:"categories"`
		Total     int                      `json:"total"`
	}
	s.assertJSONResponse(referrerResp, http.StatusOK, &referrerData)
	
	// Verify referrer data structure
	s.Equal(len(referrers), referrerData.Total)
	s.GreaterOrEqual(len(referrerData.Referrers), 1)
	s.GreaterOrEqual(len(referrerData.Categories), 1)
}

// TestTimelineAnalytics tests timeline analytics with different periods
func (s *AnalyticsWorkflowTestSuite) TestTimelineAnalytics() {
	// Create a test URL
	urlData := s.createTestURL("https://example.com/timeline-test")
	shortCode := urlData["shortCode"].(string)
	urlID := strconv.Itoa(int(urlData["id"].(float64)))
	
	// Simulate clicks spread over time
	for i := 0; i < 5; i++ {
		resp := s.makeRequest("GET", "/"+shortCode, nil, nil)
		resp.Body.Close()
		time.Sleep(100 * time.Millisecond)
	}
	
	// Allow time for processing
	time.Sleep(200 * time.Millisecond)
	
	// Test different timeline periods
	periods := []string{"1h", "24h", "7d", "30d"}
	
	for _, period := range periods {
		timelineResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/urls/"+urlID+"/timeline?period="+period, nil)
		
		var timelineData struct {
			Timeline []map[string]interface{} `json:"timeline"`
			Period   string                   `json:"period"`
			Total    int                      `json:"total"`
		}
		s.assertJSONResponse(timelineResp, http.StatusOK, &timelineData)
		
		// Verify timeline data structure
		s.Equal(period, timelineData.Period)
		s.Equal(5, timelineData.Total)
		s.GreaterOrEqual(len(timelineData.Timeline), 1)
	}
}

// TestGlobalAnalytics tests global analytics across all URLs
func (s *AnalyticsWorkflowTestSuite) TestGlobalAnalytics() {
	// Create multiple URLs and simulate activity
	for i := 0; i < 3; i++ {
		urlData := s.createTestURL(fmt.Sprintf("https://example.com/global-test-%d", i))
		shortCode := urlData["shortCode"].(string)
		
		// Simulate different amounts of activity for each URL
		for j := 0; j < (i+1)*3; j++ {
			resp := s.makeRequest("GET", "/"+shortCode, nil, nil)
			resp.Body.Close()
			time.Sleep(20 * time.Millisecond)
		}
	}
	
	// Allow time for processing
	time.Sleep(200 * time.Millisecond)
	
	// Test global analytics endpoint
	globalResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/global", nil)
	
	var globalData struct {
		Overview struct {
			TotalURLs       int     `json:"totalUrls"`
			TotalClicks     int     `json:"totalClicks"`
			UniqueVisitors  int     `json:"uniqueVisitors"`
			AvgClicksPerURL float64 `json:"avgClicksPerUrl"`
		} `json:"overview"`
		TopCountries []map[string]interface{} `json:"topCountries"`
		TopDevices   []map[string]interface{} `json:"topDevices"`
		TopReferrers []map[string]interface{} `json:"topReferrers"`
		ClickTrend   []map[string]interface{} `json:"clickTrend"`
	}
	s.assertJSONResponse(globalResp, http.StatusOK, &globalData)
	
	// Verify global analytics data
	s.GreaterOrEqual(globalData.Overview.TotalURLs, 3)
	s.GreaterOrEqual(globalData.Overview.TotalClicks, 18) // 3+6+9 clicks minimum
	s.Greater(globalData.Overview.AvgClicksPerURL, 0.0)
}

// TestTopPerformingURLs tests the top performing URLs analytics
func (s *AnalyticsWorkflowTestSuite) TestTopPerformingURLs() {
	// Create URLs with different performance levels
	urlPerformance := []struct {
		url    string
		clicks int
	}{
		{"https://example.com/top-performer", 10},
		{"https://example.com/medium-performer", 5},
		{"https://example.com/low-performer", 2},
	}
	
	var urlCodes []string
	for _, perf := range urlPerformance {
		urlData := s.createTestURL(perf.url)
		shortCode := urlData["shortCode"].(string)
		urlCodes = append(urlCodes, shortCode)
		
		// Simulate the specified number of clicks
		for i := 0; i < perf.clicks; i++ {
			resp := s.makeRequest("GET", "/"+shortCode, nil, nil)
			resp.Body.Close()
			time.Sleep(20 * time.Millisecond)
		}
	}
	
	// Allow time for processing
	time.Sleep(200 * time.Millisecond)
	
	// Test top performing URLs endpoint
	topResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/top-urls?limit=10", nil)
	
	var topData struct {
		URLs  []map[string]interface{} `json:"urls"`
		Total int                      `json:"total"`
		Limit int                      `json:"limit"`
	}
	s.assertJSONResponse(topResp, http.StatusOK, &topData)
	
	// Verify URLs are ordered by performance (click count)
	s.GreaterOrEqual(len(topData.URLs), 3)
	s.Equal(10, topData.Limit)
	
	// Check that URLs are sorted by click count (descending)
	for i := 0; i < len(topData.URLs)-1; i++ {
		currentClicks := int(topData.URLs[i]["clickCount"].(float64))
		nextClicks := int(topData.URLs[i+1]["clickCount"].(float64))
		s.GreaterOrEqual(currentClicks, nextClicks)
	}
}

// TestAnalyticsExport tests analytics data export functionality
func (s *AnalyticsWorkflowTestSuite) TestAnalyticsExport() {
	// Create test data
	urlData := s.createTestURL("https://example.com/export-test")
	shortCode := urlData["shortCode"].(string)
	
	// Simulate some clicks
	for i := 0; i < 5; i++ {
		resp := s.makeRequest("GET", "/"+shortCode, nil, nil)
		resp.Body.Close()
		time.Sleep(50 * time.Millisecond)
	}
	
	// Allow time for processing
	time.Sleep(200 * time.Millisecond)
	
	// Test analytics export endpoint
	exportResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/export?format=json", nil)
	s.Equal(http.StatusOK, exportResp.StatusCode)
	
	// Verify content type for JSON export
	contentType := exportResp.Header.Get("Content-Type")
	s.Contains(contentType, "application/json")
	
	// Test CSV export
	csvResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/export?format=csv", nil)
	s.Equal(http.StatusOK, csvResp.StatusCode)
	
	// Verify content type for CSV export
	csvContentType := csvResp.Header.Get("Content-Type")
	s.Contains(csvContentType, "text/csv")
	
	csvResp.Body.Close()
	exportResp.Body.Close()
}

// TestAnalyticsWithDateFilters tests analytics with date range filters
func (s *AnalyticsWorkflowTestSuite) TestAnalyticsWithDateFilters() {
	// Create a test URL
	urlData := s.createTestURL("https://example.com/date-filter-test")
	urlID := strconv.Itoa(int(urlData["id"].(float64)))
	
	// Test analytics with date filters
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)
	
	// Test with date range filters
	dateParams := fmt.Sprintf("?startDate=%s&endDate=%s", 
		yesterday.Format("2006-01-02"), 
		tomorrow.Format("2006-01-02"))
	
	analyticsResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/urls/"+urlID+dateParams, nil)
	
	var analyticsData struct {
		DateRange struct {
			StartDate string `json:"startDate"`
			EndDate   string `json:"endDate"`
		} `json:"dateRange"`
		Stats map[string]interface{} `json:"stats"`
	}
	s.assertJSONResponse(analyticsResp, http.StatusOK, &analyticsData)
	
	// Verify date range was applied
	s.NotEmpty(analyticsData.DateRange.StartDate)
	s.NotEmpty(analyticsData.DateRange.EndDate)
}

// TestConcurrentAnalyticsQueries tests analytics endpoints under concurrent load
func (s *AnalyticsWorkflowTestSuite) TestConcurrentAnalyticsQueries() {
	// Create test data first
	urlData := s.createTestURL("https://example.com/concurrent-analytics")
	shortCode := urlData["shortCode"].(string)
	urlID := strconv.Itoa(int(urlData["id"].(float64)))
	
	// Simulate some clicks
	for i := 0; i < 10; i++ {
		resp := s.makeRequest("GET", "/"+shortCode, nil, nil)
		resp.Body.Close()
		time.Sleep(20 * time.Millisecond)
	}
	
	// Allow time for processing
	time.Sleep(200 * time.Millisecond)
	
	// Test concurrent analytics queries
	done := make(chan bool, 5)
	
	endpoints := []string{
		"/api/v1/analytics/dashboard",
		"/api/v1/analytics/urls/" + urlID,
		"/api/v1/analytics/urls/" + urlID + "/timeline",
		"/api/v1/analytics/urls/" + urlID + "/geo",
		"/api/v1/analytics/urls/" + urlID + "/devices",
	}
	
	for _, endpoint := range endpoints {
		go func(url string) {
			defer func() { done <- true }()
			
			resp := s.makeAuthenticatedRequest("GET", url, nil)
			s.Equal(http.StatusOK, resp.StatusCode)
			resp.Body.Close()
		}(endpoint)
	}
	
	// Wait for all queries to complete
	for i := 0; i < len(endpoints); i++ {
		<-done
	}
}

// Helper method to get URL details (inherited from URLWorkflowTestSuite)
func (s *AnalyticsWorkflowTestSuite) getURLDetails(urlID interface{}) map[string]interface{} {
	var id string
	switch v := urlID.(type) {
	case float64:
		id = strconv.Itoa(int(v))
	case uint:
		id = strconv.Itoa(int(v))
	case string:
		id = v
	default:
		s.Fail("Invalid URL ID type")
	}
	
	resp := s.makeAuthenticatedRequest("GET", "/api/v1/urls/"+id, nil)
	
	var urlData map[string]interface{}
	s.assertJSONResponse(resp, http.StatusOK, &urlData)
	
	return urlData
}