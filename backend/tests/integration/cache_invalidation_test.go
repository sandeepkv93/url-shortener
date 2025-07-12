package integration

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"url-shortener/internal/core/domain"
)

// CacheInvalidationTestSuite tests cache invalidation scenarios
type CacheInvalidationTestSuite struct {
	IntegrationTestSuite
}

// TestBasicCacheOperations tests basic cache hit/miss scenarios
func (s *CacheInvalidationTestSuite) TestBasicCacheOperations() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping cache tests")
		return
	}

	ctx := context.Background()
	key := "test:basic:key"
	value := "test-value"

	// Test cache miss
	result, err := s.cacheService.Get(ctx, key)
	s.NoError(err)
	s.Empty(result)

	// Set cache value
	err = s.cacheService.Set(ctx, key, value, 5*time.Minute)
	s.NoError(err)

	// Test cache hit
	result, err = s.cacheService.Get(ctx, key)
	s.NoError(err)
	s.Equal(value, result)

	// Delete cache value
	err = s.cacheService.Del(ctx, key)
	s.NoError(err)

	// Test cache miss after deletion
	result, err = s.cacheService.Get(ctx, key)
	s.NoError(err)
	s.Empty(result)
}

// TestURLCacheInvalidationOnUpdate tests cache invalidation when URLs are updated
func (s *CacheInvalidationTestSuite) TestURLCacheInvalidationOnUpdate() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping cache tests")
		return
	}

	// Create a URL
	urlData := s.createTestURL("https://example.com/cache-update-test")
	urlID := strconv.Itoa(int(urlData["id"].(float64)))
	shortCode := urlData["shortCode"].(string)

	// Access the URL to trigger caching
	resp := s.makeRequest("GET", "/"+shortCode, nil, nil)
	s.Equal(http.StatusFound, resp.StatusCode)
	resp.Body.Close()

	// Update the URL
	updateReq := map[string]interface{}{
		"title":       "Updated Cache Test",
		"description": "Cache invalidation test updated",
		"isPublic":    false,
	}

	updateResp := s.makeAuthenticatedRequest("PUT", "/api/v1/urls/"+urlID, updateReq)
	s.Equal(http.StatusOK, updateResp.StatusCode)
	updateResp.Body.Close()

	// Verify cache was invalidated by checking if updated data is returned
	time.Sleep(100 * time.Millisecond) // Allow time for cache invalidation

	detailResp := s.makeAuthenticatedRequest("GET", "/api/v1/urls/"+urlID, nil)
	var urlDetails struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsPublic    bool   `json:"isPublic"`
	}
	s.assertJSONResponse(detailResp, http.StatusOK, &urlDetails)

	// Verify updated data is returned (not cached old data)
	s.Equal("Updated Cache Test", urlDetails.Title)
	s.Equal("Cache invalidation test updated", urlDetails.Description)
	s.False(urlDetails.IsPublic)
}

// TestURLCacheInvalidationOnDeletion tests cache invalidation when URLs are deleted
func (s *CacheInvalidationTestSuite) TestURLCacheInvalidationOnDeletion() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping cache tests")
		return
	}

	// Create a URL
	urlData := s.createTestURL("https://example.com/cache-delete-test")
	urlID := strconv.Itoa(int(urlData["id"].(float64)))
	shortCode := urlData["shortCode"].(string)

	// Access the URL to trigger caching
	resp := s.makeRequest("GET", "/"+shortCode, nil, nil)
	s.Equal(http.StatusFound, resp.StatusCode)
	resp.Body.Close()

	// Delete the URL
	deleteResp := s.makeAuthenticatedRequest("DELETE", "/api/v1/urls/"+urlID, nil)
	s.Equal(http.StatusNoContent, deleteResp.StatusCode)

	// Verify cache was invalidated by checking that URL is no longer accessible
	time.Sleep(100 * time.Millisecond) // Allow time for cache invalidation

	accessResp := s.makeRequest("GET", "/"+shortCode, nil, nil)
	s.NotEqual(http.StatusFound, accessResp.StatusCode)
	accessResp.Body.Close()
}

// TestAnalyticsCacheInvalidation tests cache invalidation for analytics data
func (s *CacheInvalidationTestSuite) TestAnalyticsCacheInvalidation() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping cache tests")
		return
	}

	// Create a URL
	urlData := s.createTestURL("https://example.com/analytics-cache-test")
	urlID := strconv.Itoa(int(urlData["id"].(float64)))
	shortCode := urlData["shortCode"].(string)

	// Get initial analytics (should cache the data)
	analyticsResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/urls/"+urlID, nil)
	var initialAnalytics struct {
		Stats map[string]interface{} `json:"stats"`
	}
	s.assertJSONResponse(analyticsResp, http.StatusOK, &initialAnalytics)
	initialClicks := int(initialAnalytics.Stats["totalClicks"].(float64))

	// Generate a click to change analytics data
	clickResp := s.makeRequest("GET", "/"+shortCode, nil, nil)
	s.Equal(http.StatusFound, clickResp.StatusCode)
	clickResp.Body.Close()

	// Allow time for async click processing and cache invalidation
	time.Sleep(200 * time.Millisecond)

	// Get analytics again (should return updated data, not cached)
	newAnalyticsResp := s.makeAuthenticatedRequest("GET", "/api/v1/analytics/urls/"+urlID, nil)
	var newAnalytics struct {
		Stats map[string]interface{} `json:"stats"`
	}
	s.assertJSONResponse(newAnalyticsResp, http.StatusOK, &newAnalytics)
	newClicks := int(newAnalytics.Stats["totalClicks"].(float64))

	// Verify analytics data was updated (cache was invalidated)
	s.Greater(newClicks, initialClicks)
}

// TestCacheExpiration tests cache expiration functionality
func (s *CacheInvalidationTestSuite) TestCacheExpiration() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping cache tests")
		return
	}

	ctx := context.Background()
	key := "test:expiration:key"
	value := "expiring-value"

	// Set cache value with short expiration
	err := s.cacheService.Set(ctx, key, value, 100*time.Millisecond)
	s.NoError(err)

	// Verify value exists
	result, err := s.cacheService.Get(ctx, key)
	s.NoError(err)
	s.Equal(value, result)

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Verify value has expired
	result, err = s.cacheService.Get(ctx, key)
	s.NoError(err)
	s.Empty(result)
}

// TestBulkCacheInvalidation tests bulk cache invalidation operations
func (s *CacheInvalidationTestSuite) TestBulkCacheInvalidation() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping cache tests")
		return
	}

	ctx := context.Background()

	// Set multiple cache values
	keys := []string{
		"test:bulk:key1",
		"test:bulk:key2",
		"test:bulk:key3",
		"test:bulk:key4",
	}

	for i, key := range keys {
		err := s.cacheService.Set(ctx, key, fmt.Sprintf("value-%d", i+1), 5*time.Minute)
		s.NoError(err)
	}

	// Verify all values exist
	for i, key := range keys {
		result, err := s.cacheService.Get(ctx, key)
		s.NoError(err)
		s.Equal(fmt.Sprintf("value-%d", i+1), result)
	}

	// Delete individually since bulk delete by pattern is not in the interface
	for _, key := range keys {
		err := s.cacheService.Del(ctx, key)
		s.NoError(err)
	}

	// Verify all values are deleted
	for _, key := range keys {
		result, err := s.cacheService.Get(ctx, key)
		s.NoError(err)
		s.Empty(result)
	}
}

// TestCacheConsistencyAcrossOperations tests cache consistency during concurrent operations
func (s *CacheInvalidationTestSuite) TestCacheConsistencyAcrossOperations() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping cache tests")
		return
	}

	// Create multiple URLs for testing
	urlData1 := s.createTestURL("https://example.com/consistency-test-1")
	urlData2 := s.createTestURL("https://example.com/consistency-test-2")
	
	urlID1 := strconv.Itoa(int(urlData1["id"].(float64)))
	urlID2 := strconv.Itoa(int(urlData2["id"].(float64)))
	
	shortCode1 := urlData1["shortCode"].(string)
	shortCode2 := urlData2["shortCode"].(string)

	// Access both URLs to populate cache
	resp1 := s.makeRequest("GET", "/"+shortCode1, nil, nil)
	s.Equal(http.StatusFound, resp1.StatusCode)
	resp1.Body.Close()

	resp2 := s.makeRequest("GET", "/"+shortCode2, nil, nil)
	s.Equal(http.StatusFound, resp2.StatusCode)
	resp2.Body.Close()

	// Update one URL
	updateReq := map[string]interface{}{
		"title": "Updated for Consistency Test",
	}

	updateResp := s.makeAuthenticatedRequest("PUT", "/api/v1/urls/"+urlID1, updateReq)
	s.Equal(http.StatusOK, updateResp.StatusCode)
	updateResp.Body.Close()

	// Verify that only the updated URL's cache was invalidated
	time.Sleep(100 * time.Millisecond)

	// Check updated URL returns new data
	detail1Resp := s.makeAuthenticatedRequest("GET", "/api/v1/urls/"+urlID1, nil)
	var url1Details struct {
		Title string `json:"title"`
	}
	s.assertJSONResponse(detail1Resp, http.StatusOK, &url1Details)
	s.Equal("Updated for Consistency Test", url1Details.Title)

	// Check that the other URL's cache is still valid
	detail2Resp := s.makeAuthenticatedRequest("GET", "/api/v1/urls/"+urlID2, nil)
	var url2Details struct {
		Title string `json:"title"`
	}
	s.assertJSONResponse(detail2Resp, http.StatusOK, &url2Details)
	s.Equal("Test URL", url2Details.Title) // Should still be the original title
}

// TestCacheInvalidationOnUserOperations tests cache invalidation for user-related operations
func (s *CacheInvalidationTestSuite) TestCacheInvalidationOnUserOperations() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping cache tests")
		return
	}

	// Get user profile to potentially cache it
	profileResp := s.makeAuthenticatedRequest("GET", "/api/v1/auth/profile", nil)
	var initialProfile struct {
		User struct {
			Email string `json:"email"`
		} `json:"user"`
	}
	s.assertJSONResponse(profileResp, http.StatusOK, &initialProfile)
	initialEmail := initialProfile.User.Email

	// Update user profile
	updateReq := map[string]string{
		"email": "updated-cache-test@example.com",
	}

	updateResp := s.makeAuthenticatedRequest("PUT", "/api/v1/auth/profile", updateReq)
	s.Equal(http.StatusOK, updateResp.StatusCode)
	updateResp.Body.Close()

	// Verify cache was invalidated
	time.Sleep(100 * time.Millisecond)

	newProfileResp := s.makeAuthenticatedRequest("GET", "/api/v1/auth/profile", nil)
	var newProfile struct {
		User struct {
			Email string `json:"email"`
		} `json:"user"`
	}
	s.assertJSONResponse(newProfileResp, http.StatusOK, &newProfile)

	// Verify updated data is returned (not cached old data)
	s.NotEqual(initialEmail, newProfile.User.Email)
	s.Equal("updated-cache-test@example.com", newProfile.User.Email)
}

// TestCacheInvalidationOnRateLimitUpdates tests cache invalidation for rate limit data
func (s *CacheInvalidationTestSuite) TestCacheInvalidationOnRateLimitUpdates() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping cache tests")
		return
	}

	ctx := context.Background()
	rateLimitKey := "rate_limit:test_user"

	// Simulate rate limit data in cache
	err := s.cacheService.Set(ctx, rateLimitKey, "10", 1*time.Minute)
	s.NoError(err)

	// Verify rate limit data exists
	result, err := s.cacheService.Get(ctx, rateLimitKey)
	s.NoError(err)
	s.Equal("10", result)

	// Simulate rate limit reset/update
	err = s.cacheService.Del(ctx, rateLimitKey)
	s.NoError(err)

	// Verify rate limit data was cleared
	result, err = s.cacheService.Get(ctx, rateLimitKey)
	s.NoError(err)
	s.Empty(result)
}

// TestCachePerformanceWithInvalidation tests cache performance during invalidation operations
func (s *CacheInvalidationTestSuite) TestCachePerformanceWithInvalidation() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping cache tests")
		return
	}

	ctx := context.Background()
	
	// Set up a large number of cache entries
	numEntries := 100
	keys := make([]string, numEntries)
	
	start := time.Now()
	
	// Bulk set operations
	for i := 0; i < numEntries; i++ {
		key := fmt.Sprintf("test:performance:%d", i)
		keys[i] = key
		err := s.cacheService.Set(ctx, key, fmt.Sprintf("value-%d", i), 5*time.Minute)
		s.NoError(err)
	}
	
	setBulkDuration := time.Since(start)
	s.T().Logf("Bulk set of %d entries took: %v", numEntries, setBulkDuration)

	// Verify all entries exist (test cache hit performance)
	start = time.Now()
	for _, key := range keys {
		result, err := s.cacheService.Get(ctx, key)
		s.NoError(err)
		s.NotEmpty(result)
	}
	getBulkDuration := time.Since(start)
	s.T().Logf("Bulk get of %d entries took: %v", numEntries, getBulkDuration)

	// Bulk delete operations (test invalidation performance)
	start = time.Now()
	for _, key := range keys {
		err := s.cacheService.Del(ctx, key)
		s.NoError(err)
	}
	deleteBulkDuration := time.Since(start)
	s.T().Logf("Bulk delete of %d entries took: %v", numEntries, deleteBulkDuration)

	// Verify all entries are deleted
	for _, key := range keys {
		result, err := s.cacheService.Get(ctx, key)
		s.NoError(err)
		s.Empty(result)
	}

	// Performance assertions (these are rough guidelines)
	s.Less(setBulkDuration, 5*time.Second, "Bulk set operation took too long")
	s.Less(getBulkDuration, 2*time.Second, "Bulk get operation took too long")
	s.Less(deleteBulkDuration, 3*time.Second, "Bulk delete operation took too long")
}

// TestCacheInvalidationWithDatabaseRollback tests cache behavior during database rollbacks
func (s *CacheInvalidationTestSuite) TestCacheInvalidationWithDatabaseRollback() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping cache tests")
		return
	}

	// Create a URL first
	urlData := s.createTestURL("https://example.com/rollback-cache-test")
	urlID := strconv.Itoa(int(urlData["id"].(float64)))

	// Cache the URL data by accessing it
	detailResp := s.makeAuthenticatedRequest("GET", "/api/v1/urls/"+urlID, nil)
	s.Equal(http.StatusOK, detailResp.StatusCode)
	detailResp.Body.Close()

	// Start a transaction and update the URL
	tx := s.db.Begin()
	s.Require().NoError(tx.Error)

	// Update URL in transaction
	err := tx.Model(&domain.ShortURL{}).Where("id = ?", urlData["id"]).Update("title", "Rollback Test Title").Error
	s.NoError(err)

	// At this point, if cache invalidation happened within transaction,
	// the cache might be inconsistent if we rollback

	// Rollback the transaction
	result := tx.Rollback()
	s.Require().NoError(result.Error)

	// Verify that cache is consistent with database state
	time.Sleep(100 * time.Millisecond)

	finalDetailResp := s.makeAuthenticatedRequest("GET", "/api/v1/urls/"+urlID, nil)
	var urlDetails struct {
		Title string `json:"title"`
	}
	s.assertJSONResponse(finalDetailResp, http.StatusOK, &urlDetails)

	// Should return original title, not the rolled-back update
	s.Equal("Test URL", urlDetails.Title)
}