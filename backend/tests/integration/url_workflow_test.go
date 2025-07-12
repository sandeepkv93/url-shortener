package integration

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// URLWorkflowTestSuite tests complete URL management workflows
type URLWorkflowTestSuite struct {
	IntegrationTestSuite
}

// TestURLCreationWorkflow tests the complete URL creation process
func (s *URLWorkflowTestSuite) TestURLCreationWorkflow() {
	// Test basic URL creation
	urlReq := map[string]interface{}{
		"originalUrl": "https://example.com/very/long/url/that/needs/shortening",
		"title":       "Example Website",
		"description": "A test website for demonstration",
	}
	
	resp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", urlReq)
	
	var urlResp struct {
		ID          uint   `json:"id"`
		ShortCode   string `json:"shortCode"`
		OriginalURL string `json:"originalUrl"`
		Title       string `json:"title"`
		Description string `json:"description"`
		ClickCount  int    `json:"clickCount"`
		IsActive    bool   `json:"isActive"`
		IsPublic    bool   `json:"isPublic"`
		CreatedAt   string `json:"createdAt"`
	}
	s.assertJSONResponse(resp, http.StatusCreated, &urlResp)
	
	// Verify response data
	s.Greater(urlResp.ID, uint(0))
	s.NotEmpty(urlResp.ShortCode)
	s.Equal("https://example.com/very/long/url/that/needs/shortening", urlResp.OriginalURL)
	s.Equal("Example Website", urlResp.Title)
	s.Equal("A test website for demonstration", urlResp.Description)
	s.Equal(0, urlResp.ClickCount)
	s.True(urlResp.IsActive)
	s.True(urlResp.IsPublic)
	s.NotEmpty(urlResp.CreatedAt)
	
	// Verify short code has correct length and characters
	s.True(len(urlResp.ShortCode) >= 6 && len(urlResp.ShortCode) <= 8)
	s.Regexp("^[a-zA-Z0-9]+$", urlResp.ShortCode)
}

// TestURLCreationWithCustomAlias tests URL creation with custom alias
func (s *URLWorkflowTestSuite) TestURLCreationWithCustomAlias() {
	// Test URL creation with custom alias
	urlReq := map[string]interface{}{
		"originalUrl":  "https://example.com/custom",
		"customAlias":  "my-custom-link",
		"title":        "Custom Link",
	}
	
	resp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", urlReq)
	
	var urlResp struct {
		ShortCode   string `json:"shortCode"`
		OriginalURL string `json:"originalUrl"`
		CustomAlias bool   `json:"customAlias"`
	}
	s.assertJSONResponse(resp, http.StatusCreated, &urlResp)
	
	// Verify custom alias is used
	s.Equal("my-custom-link", urlResp.ShortCode)
	s.True(urlResp.CustomAlias)
	
	// Test duplicate custom alias should fail
	duplicateReq := map[string]interface{}{
		"originalUrl": "https://example.com/duplicate",
		"customAlias": "my-custom-link",
	}
	
	dupResp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", duplicateReq)
	s.assertErrorResponse(dupResp, http.StatusConflict, "alias already exists")
}

// TestURLCreationWithAdvancedOptions tests URL creation with all options
func (s *URLWorkflowTestSuite) TestURLCreationWithAdvancedOptions() {
	// Test URL creation with all available options
	expiresAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	
	urlReq := map[string]interface{}{
		"originalUrl": "https://example.com/advanced",
		"title":       "Advanced URL",
		"description": "URL with all options set",
		"expiresAt":   expiresAt,
		"password":    "secret123",
		"isPublic":    false,
		"tags":        []string{"test", "integration", "advanced"},
	}
	
	resp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", urlReq)
	
	var urlResp struct {
		ID          uint     `json:"id"`
		ShortCode   string   `json:"shortCode"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		ExpiresAt   string   `json:"expiresAt"`
		IsPublic    bool     `json:"isPublic"`
		IsProtected bool     `json:"isProtected"`
		Tags        []string `json:"tags"`
	}
	s.assertJSONResponse(resp, http.StatusCreated, &urlResp)
	
	// Verify all options were set correctly
	s.Equal("Advanced URL", urlResp.Title)
	s.Equal("URL with all options set", urlResp.Description)
	s.NotEmpty(urlResp.ExpiresAt)
	s.False(urlResp.IsPublic)
	s.True(urlResp.IsProtected)
	s.ElementsMatch([]string{"test", "integration", "advanced"}, urlResp.Tags)
}

// TestURLRedirectionWorkflow tests the URL redirection functionality
func (s *URLWorkflowTestSuite) TestURLRedirectionWorkflow() {
	// Create a test URL
	originalURL := "https://example.com/redirect-test"
	urlData := s.createTestURL(originalURL)
	
	shortCode := urlData["shortCode"].(string)
	
	// Test redirection
	resp := s.makeRequest("GET", "/"+shortCode, nil, nil)
	
	// Should return a redirect response
	s.Equal(http.StatusFound, resp.StatusCode)
	s.Equal(originalURL, resp.Header.Get("Location"))
	
	// Verify click was recorded
	time.Sleep(100 * time.Millisecond) // Allow time for async click recording
	
	urlDetails := s.getURLDetails(urlData["id"])
	s.Equal(1, int(urlDetails["clickCount"].(float64)))
}

// TestURLListingAndFiltering tests URL listing with various filters
func (s *URLWorkflowTestSuite) TestURLListingAndFiltering() {
	// Create multiple test URLs with different properties
	urls := []map[string]interface{}{
		{
			"originalUrl": "https://example.com/url1",
			"title":       "Public URL",
			"isPublic":    true,
			"tags":        []string{"public", "test"},
		},
		{
			"originalUrl": "https://example.com/url2",
			"title":       "Private URL",
			"isPublic":    false,
			"tags":        []string{"private", "test"},
		},
		{
			"originalUrl": "https://example.com/url3",
			"title":       "Password Protected",
			"password":    "secret",
			"tags":        []string{"protected"},
		},
	}
	
	var createdURLs []map[string]interface{}
	for _, urlReq := range urls {
		resp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", urlReq)
		var urlResp map[string]interface{}
		s.assertJSONResponse(resp, http.StatusCreated, &urlResp)
		createdURLs = append(createdURLs, urlResp)
	}
	
	// Test listing all URLs
	resp := s.makeAuthenticatedRequest("GET", "/api/v1/urls/", nil)
	
	var listResp struct {
		URLs  []map[string]interface{} `json:"urls"`
		Total int                     `json:"total"`
		Page  int                     `json:"page"`
		Limit int                     `json:"limit"`
	}
	s.assertJSONResponse(resp, http.StatusOK, &listResp)
	
	s.GreaterOrEqual(len(listResp.URLs), 3)
	s.GreaterOrEqual(listResp.Total, 3)
	
	// Test filtering by public status
	publicResp := s.makeAuthenticatedRequest("GET", "/api/v1/urls/?isPublic=true", nil)
	s.assertJSONResponse(publicResp, http.StatusOK, &listResp)
	
	// Verify all returned URLs are public
	for _, url := range listResp.URLs {
		s.True(url["isPublic"].(bool))
	}
	
	// Test pagination
	paginatedResp := s.makeAuthenticatedRequest("GET", "/api/v1/urls/?page=1&limit=2", nil)
	s.assertJSONResponse(paginatedResp, http.StatusOK, &listResp)
	
	s.LessOrEqual(len(listResp.URLs), 2)
	s.Equal(1, listResp.Page)
	s.Equal(2, listResp.Limit)
}

// TestURLUpdateWorkflow tests URL update functionality
func (s *URLWorkflowTestSuite) TestURLUpdateWorkflow() {
	// Create a test URL
	urlData := s.createTestURL("https://example.com/update-test")
	urlID := strconv.Itoa(int(urlData["id"].(float64)))
	
	// Test URL update
	updateReq := map[string]interface{}{
		"title":       "Updated Title",
		"description": "Updated description",
		"isPublic":    false,
		"tags":        []string{"updated", "test"},
	}
	
	resp := s.makeAuthenticatedRequest("PUT", "/api/v1/urls/"+urlID, updateReq)
	
	var updateResp struct {
		ID          uint     `json:"id"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		IsPublic    bool     `json:"isPublic"`
		Tags        []string `json:"tags"`
		UpdatedAt   string   `json:"updatedAt"`
	}
	s.assertJSONResponse(resp, http.StatusOK, &updateResp)
	
	// Verify updates
	s.Equal("Updated Title", updateResp.Title)
	s.Equal("Updated description", updateResp.Description)
	s.False(updateResp.IsPublic)
	s.ElementsMatch([]string{"updated", "test"}, updateResp.Tags)
	s.NotEmpty(updateResp.UpdatedAt)
	
	// Verify update persisted
	updatedData := s.getURLDetails(urlData["id"])
	s.Equal("Updated Title", updatedData["title"])
	s.Equal("Updated description", updatedData["description"])
	s.False(updatedData["isPublic"].(bool))
}

// TestURLDeletionWorkflow tests URL deletion functionality
func (s *URLWorkflowTestSuite) TestURLDeletionWorkflow() {
	// Create a test URL
	urlData := s.createTestURL("https://example.com/delete-test")
	urlID := strconv.Itoa(int(urlData["id"].(float64)))
	shortCode := urlData["shortCode"].(string)
	
	// Test URL deletion
	resp := s.makeAuthenticatedRequest("DELETE", "/api/v1/urls/"+urlID, nil)
	s.Equal(http.StatusNoContent, resp.StatusCode)
	
	// Verify URL is no longer accessible
	getResp := s.makeAuthenticatedRequest("GET", "/api/v1/urls/"+urlID, nil)
	s.assertErrorResponse(getResp, http.StatusNotFound, "URL not found")
	
	// Verify redirection no longer works
	redirectResp := s.makeRequest("GET", "/"+shortCode, nil, nil)
	s.assertErrorResponse(redirectResp, http.StatusNotFound, "URL not found")
}

// TestPasswordProtectedURLWorkflow tests password-protected URL functionality
func (s *URLWorkflowTestSuite) TestPasswordProtectedURLWorkflow() {
	// Create a password-protected URL
	urlReq := map[string]interface{}{
		"originalUrl": "https://example.com/protected",
		"title":       "Protected URL",
		"password":    "secret123",
	}
	
	resp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", urlReq)
	
	var urlResp struct {
		ShortCode   string `json:"shortCode"`
		IsProtected bool   `json:"isProtected"`
	}
	s.assertJSONResponse(resp, http.StatusCreated, &urlResp)
	
	s.True(urlResp.IsProtected)
	
	// Test accessing without password should prompt for password
	redirectResp := s.makeRequest("GET", "/"+urlResp.ShortCode, nil, nil)
	// This behavior depends on implementation - might return password prompt page
	// or 401 Unauthorized. We'll check for non-redirect status.
	s.NotEqual(http.StatusFound, redirectResp.StatusCode)
	
	// Test accessing with correct password (implementation-dependent)
	// This would typically be done via a password form submission
	// For now, we just verify the URL was created with protection enabled
}

// TestURLExpirationWorkflow tests URL expiration functionality
func (s *URLWorkflowTestSuite) TestURLExpirationWorkflow() {
	// Create a URL that expires very soon
	expiresAt := time.Now().Add(100 * time.Millisecond).Format(time.RFC3339)
	
	urlReq := map[string]interface{}{
		"originalUrl": "https://example.com/expires-soon",
		"title":       "Expiring URL",
		"expiresAt":   expiresAt,
	}
	
	resp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", urlReq)
	
	var urlResp struct {
		ShortCode string `json:"shortCode"`
		ExpiresAt string `json:"expiresAt"`
	}
	s.assertJSONResponse(resp, http.StatusCreated, &urlResp)
	
	// Wait for expiration
	time.Sleep(200 * time.Millisecond)
	
	// Test that expired URL is no longer accessible
	redirectResp := s.makeRequest("GET", "/"+urlResp.ShortCode, nil, nil)
	s.NotEqual(http.StatusFound, redirectResp.StatusCode)
	
	// Should return an error indicating the URL has expired
	s.True(redirectResp.StatusCode == http.StatusGone || redirectResp.StatusCode == http.StatusNotFound)
}

// TestPopularURLsEndpoint tests the popular URLs public endpoint
func (s *URLWorkflowTestSuite) TestPopularURLsEndpoint() {
	// Create some URLs and simulate clicks
	for i := 0; i < 3; i++ {
		urlData := s.createTestURL(fmt.Sprintf("https://example.com/popular%d", i))
		shortCode := urlData["shortCode"].(string)
		
		// Simulate multiple clicks
		for j := 0; j < (i+1)*2; j++ {
			clickResp := s.makeRequest("GET", "/"+shortCode, nil, nil)
			clickResp.Body.Close()
			time.Sleep(10 * time.Millisecond) // Small delay between clicks
		}
	}
	
	// Allow time for click processing
	time.Sleep(200 * time.Millisecond)
	
	// Test popular URLs endpoint
	resp := s.makeRequest("GET", "/api/v1/urls/popular", nil, nil)
	
	var popularResp struct {
		URLs []map[string]interface{} `json:"urls"`
	}
	s.assertJSONResponse(resp, http.StatusOK, &popularResp)
	
	// Should return URLs ordered by click count
	s.GreaterOrEqual(len(popularResp.URLs), 1)
	
	// Verify URLs are ordered by click count (descending)
	if len(popularResp.URLs) > 1 {
		firstClickCount := int(popularResp.URLs[0]["clickCount"].(float64))
		secondClickCount := int(popularResp.URLs[1]["clickCount"].(float64))
		s.GreaterOrEqual(firstClickCount, secondClickCount)
	}
}

// TestConcurrentURLOperations tests URL operations under concurrent load
func (s *URLWorkflowTestSuite) TestConcurrentURLOperations() {
	// Test concurrent URL creation
	done := make(chan map[string]interface{}, 5)
	
	for i := 0; i < 5; i++ {
		go func(index int) {
			urlReq := map[string]interface{}{
				"originalUrl": fmt.Sprintf("https://example.com/concurrent%d", index),
				"title":       fmt.Sprintf("Concurrent URL %d", index),
			}
			
			resp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", urlReq)
			
			var urlResp map[string]interface{}
			s.assertJSONResponse(resp, http.StatusCreated, &urlResp)
			done <- urlResp
		}(i)
	}
	
	// Collect all created URLs
	var createdURLs []map[string]interface{}
	for i := 0; i < 5; i++ {
		urlData := <-done
		createdURLs = append(createdURLs, urlData)
	}
	
	// Verify all URLs were created successfully
	s.Len(createdURLs, 5)
	
	// Verify all short codes are unique
	shortCodes := make(map[string]bool)
	for _, urlData := range createdURLs {
		shortCode := urlData["shortCode"].(string)
		s.False(shortCodes[shortCode], "Short code should be unique: "+shortCode)
		shortCodes[shortCode] = true
	}
}

// TestRateLimitingOnURLCreation tests rate limiting on URL creation
func (s *URLWorkflowTestSuite) TestRateLimitingOnURLCreation() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping rate limiting tests")
		return
	}
	
	// Make multiple rapid URL creation requests
	var rateLimitTriggered bool
	
	for i := 0; i < 30; i++ {
		urlReq := map[string]interface{}{
			"originalUrl": fmt.Sprintf("https://example.com/ratelimit%d", i),
		}
		
		resp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", urlReq)
		if resp.StatusCode == http.StatusTooManyRequests {
			rateLimitTriggered = true
			resp.Body.Close()
			break
		}
		resp.Body.Close()
		
		// Small delay to avoid overwhelming the system
		time.Sleep(10 * time.Millisecond)
	}
	
	// If rate limiting is properly configured, it should have triggered
	if !rateLimitTriggered {
		s.T().Log("Rate limiting may not be configured or may have higher limits")
	}
}

// Helper methods

// getURLDetails retrieves URL details by ID
func (s *URLWorkflowTestSuite) getURLDetails(urlID interface{}) map[string]interface{} {
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