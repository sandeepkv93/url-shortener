package integration

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// EndToEndTestSuite tests complete end-to-end API workflows
type EndToEndTestSuite struct {
	IntegrationTestSuite
}

// TestCompleteUserJourney tests the complete user journey from registration to URL management
func (s *EndToEndTestSuite) TestCompleteUserJourney() {
	// Step 1: User Registration
	userEmail := "journey-user@example.com"
	userPassword := "securepassword123"
	
	registerReq := map[string]string{
		"email":    userEmail,
		"password": userPassword,
	}

	registerResp := s.makeRequest("POST", "/api/v1/auth/register", registerReq, nil)
	
	var registerData struct {
		User struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
		Token string `json:"token"`
	}
	s.assertJSONResponse(registerResp, http.StatusCreated, &registerData)
	
	userToken := registerData.Token
	
	s.Equal(userEmail, registerData.User.Email)
	s.NotEmpty(userToken)

	// Step 2: Create Multiple URLs
	urls := []struct {
		originalURL string
		title       string
		description string
		customAlias string
		password    string
		isPublic    bool
	}{
		{
			originalURL: "https://google.com",
			title:       "Google Search",
			description: "The world's most popular search engine",
			isPublic:    true,
		},
		{
			originalURL: "https://github.com",
			title:       "GitHub",
			description: "Developer platform",
			customAlias: "my-github",
			isPublic:    true,
		},
		{
			originalURL: "https://private-site.com",
			title:       "Private Site",
			description: "A private website",
			password:    "secret123",
			isPublic:    false,
		},
	}

	var createdURLs []map[string]interface{}
	
	for _, urlData := range urls {
		urlReq := map[string]interface{}{
			"originalUrl": urlData.originalURL,
			"title":       urlData.title,
			"description": urlData.description,
			"isPublic":    urlData.isPublic,
		}
		
		if urlData.customAlias != "" {
			urlReq["customAlias"] = urlData.customAlias
		}
		
		if urlData.password != "" {
			urlReq["password"] = urlData.password
		}

		headers := map[string]string{
			"Authorization": "Bearer " + userToken,
		}
		
		createResp := s.makeRequest("POST", "/api/v1/urls/", urlReq, headers)
		
		var urlResp map[string]interface{}
		s.assertJSONResponse(createResp, http.StatusCreated, &urlResp)
		
		createdURLs = append(createdURLs, urlResp)
		
		// Verify URL properties
		s.Equal(urlData.originalURL, urlResp["originalUrl"])
		s.Equal(urlData.title, urlResp["title"])
		s.Equal(urlData.description, urlResp["description"])
		s.Equal(urlData.isPublic, urlResp["isPublic"])
		
		if urlData.customAlias != "" {
			s.Equal(urlData.customAlias, urlResp["shortCode"])
		}
	}

	// Step 3: Test URL Redirection and Click Tracking
	for i, urlResp := range createdURLs {
		shortCode := urlResp["shortCode"].(string)
		originalURL := urls[i].originalURL
		
		// Simulate multiple clicks
		for j := 0; j < 3; j++ {
			redirectResp := s.makeRequest("GET", "/"+shortCode, nil, map[string]string{
				"User-Agent": fmt.Sprintf("TestAgent-%d", j),
				"Referer":    "https://test-site.com",
			})
			
			if urls[i].password == "" { // Only test redirection for non-password protected URLs
				s.Equal(http.StatusFound, redirectResp.StatusCode)
				s.Equal(originalURL, redirectResp.Header.Get("Location"))
			}
			redirectResp.Body.Close()
			
			time.Sleep(50 * time.Millisecond) // Delay between clicks
		}
	}

	// Allow time for async click processing
	time.Sleep(300 * time.Millisecond)

	// Step 4: Verify Analytics Data
	for _, urlResp := range createdURLs {
		urlID := strconv.Itoa(int(urlResp["id"].(float64)))
		
		headers := map[string]string{
			"Authorization": "Bearer " + userToken,
		}
		
		analyticsResp := s.makeRequest("GET", "/api/v1/analytics/urls/"+urlID, nil, headers)
		
		var analyticsData struct {
			Stats map[string]interface{} `json:"stats"`
		}
		s.assertJSONResponse(analyticsResp, http.StatusOK, &analyticsData)
		
		// Should have recorded clicks (except for password-protected URLs)
		totalClicks := int(analyticsData.Stats["totalClicks"].(float64))
		if urlResp["isProtected"].(bool) {
			// Password-protected URLs might have 0 clicks
			s.GreaterOrEqual(totalClicks, 0)
		} else {
			s.GreaterOrEqual(totalClicks, 3)
		}
	}

	// Step 5: Test Dashboard Analytics
	headers := map[string]string{
		"Authorization": "Bearer " + userToken,
	}
	
	dashboardResp := s.makeRequest("GET", "/api/v1/analytics/dashboard", nil, headers)
	
	var dashboardData struct {
		Summary struct {
			TotalURLs   int `json:"totalUrls"`
			TotalClicks int `json:"totalClicks"`
		} `json:"summary"`
	}
	s.assertJSONResponse(dashboardResp, http.StatusOK, &dashboardData)
	
	s.Equal(3, dashboardData.Summary.TotalURLs)
	s.GreaterOrEqual(dashboardData.Summary.TotalClicks, 6) // At least 6 clicks from non-protected URLs

	// Step 6: Update User Profile
	updateProfileReq := map[string]string{
		"email": "updated-journey-user@example.com",
	}
	
	updateResp := s.makeRequest("PUT", "/api/v1/auth/profile", updateProfileReq, headers)
	s.Equal(http.StatusOK, updateResp.StatusCode)
	updateResp.Body.Close()

	// Step 7: Update URL Properties
	firstURL := createdURLs[0]
	urlID := strconv.Itoa(int(firstURL["id"].(float64)))
	
	updateURLReq := map[string]interface{}{
		"title":       "Updated Google Search",
		"description": "Updated description for Google",
		"isPublic":    false,
	}
	
	updateURLResp := s.makeRequest("PUT", "/api/v1/urls/"+urlID, updateURLReq, headers)
	s.Equal(http.StatusOK, updateURLResp.StatusCode)
	updateURLResp.Body.Close()

	// Step 8: List All URLs
	listResp := s.makeRequest("GET", "/api/v1/urls/", nil, headers)
	
	var listData struct {
		URLs  []map[string]interface{} `json:"urls"`
		Total int                     `json:"total"`
	}
	s.assertJSONResponse(listResp, http.StatusOK, &listData)
	
	s.Equal(3, listData.Total)
	s.Len(listData.URLs, 3)

	// Step 9: Delete a URL
	secondURL := createdURLs[1]
	deleteURLID := strconv.Itoa(int(secondURL["id"].(float64)))
	
	deleteResp := s.makeRequest("DELETE", "/api/v1/urls/"+deleteURLID, nil, headers)
	s.Equal(http.StatusNoContent, deleteResp.StatusCode)

	// Verify deletion
	deletedResp := s.makeRequest("GET", "/api/v1/urls/"+deleteURLID, nil, headers)
	s.Equal(http.StatusNotFound, deletedResp.StatusCode)
	deletedResp.Body.Close()

	// Step 10: Final URL List Check
	finalListResp := s.makeRequest("GET", "/api/v1/urls/", nil, headers)
	s.assertJSONResponse(finalListResp, http.StatusOK, &listData)
	
	s.Equal(2, listData.Total) // One URL deleted
	s.Len(listData.URLs, 2)
}

// TestQRCodeWorkflow tests the complete QR code generation workflow
func (s *EndToEndTestSuite) TestQRCodeWorkflow() {
	// Create a URL
	urlData := s.createTestURL("https://example.com/qr-test")
	urlID := strconv.Itoa(int(urlData["id"].(float64)))

	// Test QR code generation in different formats
	formats := []string{"png", "svg"}
	sizes := []string{"256", "512"}

	for _, format := range formats {
		for _, size := range sizes {
			qrResp := s.makeAuthenticatedRequest("GET", 
				fmt.Sprintf("/api/v1/urls/%s/qr?format=%s&size=%s", urlID, format, size), nil)
			
			s.Equal(http.StatusOK, qrResp.StatusCode)
			
			// Check content type
			contentType := qrResp.Header.Get("Content-Type")
			if format == "png" {
				s.Contains(contentType, "image/png")
			} else if format == "svg" {
				s.Contains(contentType, "image/svg+xml")
			}
			
			// Check that we got some content
			s.Greater(qrResp.ContentLength, int64(0))
			qrResp.Body.Close()
		}
	}

	// Test bulk QR code generation
	bulkQRResp := s.makeAuthenticatedRequest("GET", "/api/v1/qr/bulk?format=png&size=256", nil)
	s.Equal(http.StatusOK, bulkQRResp.StatusCode)
	bulkQRResp.Body.Close()
}

// TestAnalyticsExportWorkflow tests the complete analytics export workflow
func (s *EndToEndTestSuite) TestAnalyticsExportWorkflow() {
	// Create multiple URLs and generate activity
	numURLs := 5
	var urlIDs []string

	for i := 0; i < numURLs; i++ {
		urlData := s.createTestURL(fmt.Sprintf("https://example.com/export-test-%d", i))
		urlID := strconv.Itoa(int(urlData["id"].(float64)))
		urlIDs = append(urlIDs, urlID)
		
		shortCode := urlData["shortCode"].(string)
		
		// Generate different amounts of clicks for each URL
		for j := 0; j < (i+1)*2; j++ {
			clickResp := s.makeRequest("GET", "/"+shortCode, nil, map[string]string{
				"User-Agent": fmt.Sprintf("ExportTestAgent-%d-%d", i, j),
			})
			clickResp.Body.Close()
			time.Sleep(20 * time.Millisecond)
		}
	}

	// Allow time for processing
	time.Sleep(300 * time.Millisecond)

	// Test different export formats
	exportFormats := []string{"json", "csv"}
	
	for _, format := range exportFormats {
		// Test global analytics export
		globalExportResp := s.makeAuthenticatedRequest("GET", 
			"/api/v1/analytics/export?format="+format, nil)
		
		s.Equal(http.StatusOK, globalExportResp.StatusCode)
		
		contentType := globalExportResp.Header.Get("Content-Type")
		if format == "json" {
			s.Contains(contentType, "application/json")
		} else if format == "csv" {
			s.Contains(contentType, "text/csv")
		}
		
		globalExportResp.Body.Close()

		// Test URL-specific analytics export
		urlExportResp := s.makeAuthenticatedRequest("GET", 
			fmt.Sprintf("/api/v1/analytics/urls/%s/export?format=%s", urlIDs[0], format), nil)
		
		s.Equal(http.StatusOK, urlExportResp.StatusCode)
		urlExportResp.Body.Close()
	}
}

// TestErrorHandlingWorkflow tests comprehensive error handling scenarios
func (s *EndToEndTestSuite) TestErrorHandlingWorkflow() {
	// Test authentication errors
	
	// 1. Access protected endpoint without token
	noAuthResp := s.makeRequest("GET", "/api/v1/auth/profile", nil, nil)
	s.assertErrorResponse(noAuthResp, http.StatusUnauthorized, "authorization header required")

	// 2. Access with invalid token
	invalidAuthResp := s.makeRequest("GET", "/api/v1/auth/profile", nil, map[string]string{
		"Authorization": "Bearer invalid.token.here",
	})
	s.assertErrorResponse(invalidAuthResp, http.StatusUnauthorized, "invalid token")

	// Test validation errors
	
	// 3. Invalid email format in registration
	invalidEmailReq := map[string]string{
		"email":    "invalid-email",
		"password": "password123",
	}
	invalidEmailResp := s.makeRequest("POST", "/api/v1/auth/register", invalidEmailReq, nil)
	s.Equal(http.StatusBadRequest, invalidEmailResp.StatusCode)
	invalidEmailResp.Body.Close()

	// 4. Short password
	shortPasswordReq := map[string]string{
		"email":    "valid@example.com",
		"password": "123",
	}
	shortPasswordResp := s.makeRequest("POST", "/api/v1/auth/register", shortPasswordReq, nil)
	s.Equal(http.StatusBadRequest, shortPasswordResp.StatusCode)
	shortPasswordResp.Body.Close()

	// 5. Invalid URL format
	invalidURLReq := map[string]interface{}{
		"originalUrl": "not-a-valid-url",
		"title":       "Invalid URL Test",
	}
	invalidURLResp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", invalidURLReq)
	s.Equal(http.StatusBadRequest, invalidURLResp.StatusCode)
	invalidURLResp.Body.Close()

	// Test resource not found errors
	
	// 6. Access non-existent URL
	notFoundResp := s.makeAuthenticatedRequest("GET", "/api/v1/urls/99999", nil)
	s.assertErrorResponse(notFoundResp, http.StatusNotFound, "URL not found")

	// 7. Access non-existent short code
	notFoundRedirectResp := s.makeRequest("GET", "/nonexistent", nil, nil)
	s.assertErrorResponse(notFoundRedirectResp, http.StatusNotFound, "URL not found")

	// Test conflict errors
	
	// 8. Duplicate email registration
	validUserReq := map[string]string{
		"email":    "conflict-test@example.com",
		"password": "password123",
	}
	
	// First registration should succeed
	firstRegResp := s.makeRequest("POST", "/api/v1/auth/register", validUserReq, nil)
	s.Equal(http.StatusCreated, firstRegResp.StatusCode)
	firstRegResp.Body.Close()
	
	// Second registration should fail
	secondRegResp := s.makeRequest("POST", "/api/v1/auth/register", validUserReq, nil)
	s.assertErrorResponse(secondRegResp, http.StatusConflict, "email already exists")

	// 9. Duplicate custom alias
	aliasURL1 := map[string]interface{}{
		"originalUrl":  "https://example.com/alias1",
		"customAlias": "duplicate-alias",
	}
	aliasURL2 := map[string]interface{}{
		"originalUrl":  "https://example.com/alias2",
		"customAlias": "duplicate-alias",
	}
	
	firstAliasResp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", aliasURL1)
	s.Equal(http.StatusCreated, firstAliasResp.StatusCode)
	firstAliasResp.Body.Close()
	
	secondAliasResp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", aliasURL2)
	s.assertErrorResponse(secondAliasResp, http.StatusConflict, "alias already exists")

	// Test method not allowed
	
	// 10. Wrong HTTP method
	methodNotAllowedResp := s.makeRequest("PATCH", "/api/v1/auth/login", nil, nil)
	s.Equal(http.StatusMethodNotAllowed, methodNotAllowedResp.StatusCode)
	methodNotAllowedResp.Body.Close()
}

// TestPerformanceUnderLoad tests API performance under load
func (s *EndToEndTestSuite) TestPerformanceUnderLoad() {
	numConcurrentUsers := 10
	operationsPerUser := 5
	
	var wg sync.WaitGroup
	results := make(chan struct {
		success bool
		latency time.Duration
		operation string
	}, numConcurrentUsers*operationsPerUser)

	// Create concurrent users performing various operations
	for i := 0; i < numConcurrentUsers; i++ {
		wg.Add(1)
		
		go func(userIndex int) {
			defer wg.Done()
			
			// Register user
			userEmail := fmt.Sprintf("load-test-user-%d@example.com", userIndex)
			registerReq := map[string]string{
				"email":    userEmail,
				"password": "loadtest123",
			}
			
			start := time.Now()
			registerResp := s.makeRequest("POST", "/api/v1/auth/register", registerReq, nil)
			latency := time.Since(start)
			
			results <- struct {
				success   bool
				latency   time.Duration
				operation string
			}{
				success:   registerResp.StatusCode == http.StatusCreated,
				latency:   latency,
				operation: "register",
			}
			
			if registerResp.StatusCode != http.StatusCreated {
				registerResp.Body.Close()
				return
			}
			
			var registerData struct {
				Token string `json:"token"`
			}
			s.parseJSONResponse(registerResp, &registerData)
			
			headers := map[string]string{
				"Authorization": "Bearer " + registerData.Token,
			}

			// Create URLs
			for j := 0; j < operationsPerUser; j++ {
				urlReq := map[string]interface{}{
					"originalUrl": fmt.Sprintf("https://example.com/load-test-%d-%d", userIndex, j),
					"title":       fmt.Sprintf("Load Test URL %d-%d", userIndex, j),
				}
				
				start = time.Now()
				createResp := s.makeRequest("POST", "/api/v1/urls/", urlReq, headers)
				latency = time.Since(start)
				
				results <- struct {
					success   bool
					latency   time.Duration
					operation string
				}{
					success:   createResp.StatusCode == http.StatusCreated,
					latency:   latency,
					operation: "create_url",
				}
				
				createResp.Body.Close()
			}
		}(i)
	}

	// Wait for all operations to complete
	wg.Wait()
	close(results)

	// Analyze results
	stats := make(map[string]struct {
		count       int
		successCount int
		totalLatency time.Duration
		maxLatency   time.Duration
	})

	for result := range results {
		stat := stats[result.operation]
		stat.count++
		if result.success {
			stat.successCount++
		}
		stat.totalLatency += result.latency
		if result.latency > stat.maxLatency {
			stat.maxLatency = result.latency
		}
		stats[result.operation] = stat
	}

	// Assert performance requirements
	for operation, stat := range stats {
		successRate := float64(stat.successCount) / float64(stat.count) * 100
		avgLatency := stat.totalLatency / time.Duration(stat.count)
		
		s.T().Logf("Operation: %s, Success Rate: %.2f%%, Avg Latency: %v, Max Latency: %v", 
			operation, successRate, avgLatency, stat.maxLatency)
		
		// Performance assertions
		s.GreaterOrEqual(successRate, 95.0, "Success rate should be at least 95%")
		s.Less(avgLatency, 1*time.Second, "Average latency should be less than 1 second")
		s.Less(stat.maxLatency, 5*time.Second, "Max latency should be less than 5 seconds")
	}
}

// TestCompleteDataConsistency tests data consistency across all operations
func (s *EndToEndTestSuite) TestCompleteDataConsistency() {
	// Create user and URLs
	urlData := s.createTestURL("https://example.com/consistency-test")
	urlID := strconv.Itoa(int(urlData["id"].(float64)))
	shortCode := urlData["shortCode"].(string)

	// Generate clicks
	numClicks := 10
	for i := 0; i < numClicks; i++ {
		clickResp := s.makeRequest("GET", "/"+shortCode, nil, map[string]string{
			"User-Agent": fmt.Sprintf("ConsistencyAgent-%d", i),
		})
		clickResp.Body.Close()
		time.Sleep(10 * time.Millisecond)
	}

	// Allow time for processing
	time.Sleep(200 * time.Millisecond)

	// Verify click count consistency across different endpoints
	
	// 1. URL details endpoint
	urlDetailResp := s.makeAuthenticatedRequest("GET", "/api/v1/urls/"+urlID, nil)
	var urlDetails struct {
		ClickCount int `json:"clickCount"`
	}
	s.assertJSONResponse(urlDetailResp, http.StatusOK, &urlDetails)

	// 2. Analytics endpoint
	analyticsResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/urls/"+urlID, nil)
	var analyticsData struct {
		Stats map[string]interface{} `json:"stats"`
	}
	s.assertJSONResponse(analyticsResp, http.StatusOK, &analyticsData)
	
	analyticsClicks := int(analyticsData.Stats["totalClicks"].(float64))

	// 3. Dashboard endpoint
	dashboardResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/dashboard", nil)
	var dashboardData struct {
		Summary struct {
			TotalClicks int `json:"totalClicks"`
		} `json:"summary"`
	}
	s.assertJSONResponse(dashboardResp, http.StatusOK, &dashboardData)

	// Verify consistency
	s.Equal(numClicks, urlDetails.ClickCount, "URL details click count should match actual clicks")
	s.Equal(numClicks, analyticsClicks, "Analytics click count should match actual clicks")
	s.GreaterOrEqual(dashboardData.Summary.TotalClicks, numClicks, "Dashboard should include all clicks")

	// Test consistency after update
	updateReq := map[string]interface{}{
		"title": "Updated Consistency Test",
	}
	
	updateResp := s.makeAuthenticatedRequest("PUT", "/api/v1/urls/"+urlID, updateReq)
	s.Equal(http.StatusOK, updateResp.StatusCode)
	updateResp.Body.Close()

	// Verify click count persists after update
	time.Sleep(100 * time.Millisecond)
	
	updatedDetailResp := s.makeAuthenticatedRequest("GET", "/api/v1/urls/"+urlID, nil)
	var updatedDetails struct {
		ClickCount int    `json:"clickCount"`
		Title      string `json:"title"`
	}
	s.assertJSONResponse(updatedDetailResp, http.StatusOK, &updatedDetails)
	
	s.Equal(numClicks, updatedDetails.ClickCount, "Click count should persist after update")
	s.Equal("Updated Consistency Test", updatedDetails.Title, "Title should be updated")
}