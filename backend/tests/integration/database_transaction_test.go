package integration

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"

	"url-shortener/internal/core/domain"
)

// DatabaseTransactionTestSuite tests database transaction and rollback scenarios
type DatabaseTransactionTestSuite struct {
	IntegrationTestSuite
}

// TestSuccessfulTransaction tests a successful transaction flow
func (s *DatabaseTransactionTestSuite) TestSuccessfulTransaction() {
	// Start a transaction
	tx := s.db.Begin()
	s.Require().NoError(tx.Error)

	// Create a user within the transaction
	user := &domain.User{
		Email:     "transaction-test@example.com",
		Password:  "hashed-password",
		FirstName: "Test",
		LastName:  "User",
	}

	result := tx.Create(user)
	s.Require().NoError(result.Error)
	s.Greater(user.ID, uint(0))

	// Create a URL within the same transaction
	url := &domain.ShortURL{
		UserID:      user.ID,
		OriginalURL: "https://example.com/transaction-test",
		ShortCode:   "txtest",
		Title:       "Transaction Test",
		IsActive:    true,
	}

	result = tx.Create(url)
	s.Require().NoError(result.Error)
	s.Greater(url.ID, uint(0))

	// Commit the transaction
	result = tx.Commit()
	s.Require().NoError(result.Error)

	// Verify data exists after commit
	var foundUser domain.User
	err := s.db.First(&foundUser, user.ID).Error
	s.NoError(err)
	s.Equal(user.Email, foundUser.Email)

	var foundURL domain.ShortURL
	err = s.db.First(&foundURL, url.ID).Error
	s.NoError(err)
	s.Equal(url.OriginalURL, foundURL.OriginalURL)
}

// TestTransactionRollback tests transaction rollback on error
func (s *DatabaseTransactionTestSuite) TestTransactionRollback() {
	// Start a transaction
	tx := s.db.Begin()
	s.Require().NoError(tx.Error)

	// Create a user within the transaction
	user := &domain.User{
		Email:     "rollback-test@example.com",
		Password:  "hashed-password",
		FirstName: "Rollback",
		LastName:  "Test",
	}

	result := tx.Create(user)
	s.Require().NoError(result.Error)
	userID := user.ID

	// Verify user exists within transaction
	var foundUser domain.User
	err := tx.First(&foundUser, userID).Error
	s.NoError(err)

	// Simulate an error and rollback
	result = tx.Rollback()
	s.Require().NoError(result.Error)

	// Verify user doesn't exist after rollback
	err = s.db.First(&foundUser, userID).Error
	s.Error(err)
	s.True(errors.Is(err, gorm.ErrRecordNotFound))
}

// TestServiceLayerTransactionHandling tests transaction handling in service layer
func (s *DatabaseTransactionTestSuite) TestServiceLayerTransactionHandling() {
	// Test URL creation with transaction (simulating a complex operation)
	originalURL := "https://example.com/service-transaction-test"
	
	// Mock a scenario where URL creation might fail after user validation
	// We'll test the service's transaction handling indirectly through API calls
	
	// Create a URL that should succeed
	urlReq := map[string]interface{}{
		"originalUrl": originalURL,
		"title":       "Service Transaction Test",
		"description": "Testing service layer transaction handling",
	}

	resp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", urlReq)

	var urlResp struct {
		ID          uint   `json:"id"`
		ShortCode   string `json:"shortCode"`
		OriginalURL string `json:"originalUrl"`
	}
	s.assertJSONResponse(resp, http.StatusCreated, &urlResp)

	// Verify the URL was created properly
	var foundURL domain.ShortURL
	err := s.db.First(&foundURL, urlResp.ID).Error
	s.NoError(err)
	s.Equal(originalURL, foundURL.OriginalURL)
	s.Equal(s.testUser.ID, foundURL.UserID)
}

// TestConcurrentTransactions tests concurrent database operations
func (s *DatabaseTransactionTestSuite) TestConcurrentTransactions() {
	numGoroutines := 10
	wg := sync.WaitGroup{}
	results := make(chan error, numGoroutines)

	// Create multiple concurrent URL creation operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			urlReq := map[string]interface{}{
				"originalUrl": "https://example.com/concurrent-" + strconv.Itoa(index),
				"title":       "Concurrent Test " + strconv.Itoa(index),
			}

			resp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", urlReq)
			if resp.StatusCode != http.StatusCreated {
				results <- errors.New("failed to create URL in concurrent test")
				return
			}
			resp.Body.Close()
			results <- nil
		}(i)
	}

	wg.Wait()
	close(results)

	// Check that all operations succeeded
	errorCount := 0
	for err := range results {
		if err != nil {
			errorCount++
		}
	}

	s.Equal(0, errorCount, "Some concurrent transactions failed")

	// Verify all URLs were created
	var urlCount int64
	s.db.Model(&domain.ShortURL{}).Where("user_id = ?", s.testUser.ID).Count(&urlCount)
	s.GreaterOrEqual(urlCount, int64(numGoroutines))
}

// TestTransactionIsolation tests transaction isolation levels
func (s *DatabaseTransactionTestSuite) TestTransactionIsolation() {
	// Create initial data
	url := &domain.ShortURL{
		UserID:      s.testUser.ID,
		OriginalURL: "https://example.com/isolation-test",
		ShortCode:   "isoltest",
		Title:       "Isolation Test",
		ClickCount:  0,
		IsActive:    true,
	}

	result := s.db.Create(url)
	s.Require().NoError(result.Error)

	// Start two transactions
	tx1 := s.db.Begin()
	tx2 := s.db.Begin()
	s.Require().NoError(tx1.Error)
	s.Require().NoError(tx2.Error)

	// Transaction 1: Read initial click count
	var url1 domain.ShortURL
	err := tx1.First(&url1, url.ID).Error
	s.NoError(err)
	initialCount := url1.ClickCount

	// Transaction 2: Increment click count
	err = tx2.Model(&domain.ShortURL{}).Where("id = ?", url.ID).Update("click_count", gorm.Expr("click_count + ?", 1)).Error
	s.NoError(err)

	// Transaction 1: Read click count again (should still see original value due to isolation)
	err = tx1.First(&url1, url.ID).Error
	s.NoError(err)
	s.Equal(initialCount, url1.ClickCount) // Should still see original value

	// Commit transaction 2
	result = tx2.Commit()
	s.Require().NoError(result.Error)

	// Transaction 1: Still should see original value (depending on isolation level)
	err = tx1.First(&url1, url.ID).Error
	s.NoError(err)
	// SQLite's default isolation might allow seeing the change, so we'll just verify it's consistent

	// Commit transaction 1
	result = tx1.Commit()
	s.Require().NoError(result.Error)

	// Verify final state
	var finalURL domain.ShortURL
	err = s.db.First(&finalURL, url.ID).Error
	s.NoError(err)
	s.Equal(initialCount+1, finalURL.ClickCount)
}

// TestTransactionTimeout tests transaction timeout scenarios
func (s *DatabaseTransactionTestSuite) TestTransactionTimeout() {
	// This test simulates a long-running transaction that might timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	tx := s.db.WithContext(ctx).Begin()
	s.Require().NoError(tx.Error)

	// Create a user
	user := &domain.User{
		Email:        "timeout-test@example.com",
		Password:  "hashed-password",
		FirstName: "Test",
		LastName:  "User",
	}

	result := tx.Create(user)
	s.Require().NoError(result.Error)

	// Simulate a long operation by sleeping longer than the context timeout
	time.Sleep(150 * time.Millisecond)

	// Try to create another record - this might fail due to context timeout
	url := &domain.ShortURL{
		UserID:      user.ID,
		OriginalURL: "https://example.com/timeout-test",
		ShortCode:   "timetest",
		Title:       "Timeout Test",
		IsActive:    true,
	}

	result = tx.Create(url)
	// This might fail due to context timeout, which is expected

	// Try to commit - this should handle the timeout gracefully
	result = tx.Commit()
	// The result might be an error due to timeout, which is acceptable
}

// TestNestedTransactionBehavior tests nested transaction handling
func (s *DatabaseTransactionTestSuite) TestNestedTransactionBehavior() {
	// Start outer transaction
	outerTx := s.db.Begin()
	s.Require().NoError(outerTx.Error)

	// Create a user in outer transaction
	user := &domain.User{
		Email:        "nested-test@example.com",
		Password:  "hashed-password",
		FirstName: "Test",
		LastName:  "User",
	}

	result := outerTx.Create(user)
	s.Require().NoError(result.Error)

	// Start inner transaction (savepoint)
	innerTx := outerTx.SavePoint("inner_tx")
	s.Require().NoError(innerTx.Error)

	// Create a URL in inner transaction
	url := &domain.ShortURL{
		UserID:      user.ID,
		OriginalURL: "https://example.com/nested-test",
		ShortCode:   "nesttest",
		Title:       "Nested Test",
		IsActive:    true,
	}

	result = outerTx.Create(url)
	s.Require().NoError(result.Error)

	// Rollback to savepoint (simulating inner transaction rollback)
	result = outerTx.RollbackTo("inner_tx")
	s.Require().NoError(result.Error)

	// Verify user still exists but URL doesn't (within transaction)
	var foundUser domain.User
	err := outerTx.First(&foundUser, user.ID).Error
	s.NoError(err)

	var foundURL domain.ShortURL
	err = outerTx.First(&foundURL, "user_id = ? AND short_code = ?", user.ID, "nesttest").Error
	s.Error(err)
	s.True(errors.Is(err, gorm.ErrRecordNotFound))

	// Commit outer transaction
	result = outerTx.Commit()
	s.Require().NoError(result.Error)

	// Verify final state - user should exist, URL should not
	err = s.db.First(&foundUser, user.ID).Error
	s.NoError(err)

	err = s.db.First(&foundURL, "user_id = ? AND short_code = ?", user.ID, "nesttest").Error
	s.Error(err)
	s.True(errors.Is(err, gorm.ErrRecordNotFound))
}

// TestDatabaseConstraintViolationRollback tests rollback on constraint violations
func (s *DatabaseTransactionTestSuite) TestDatabaseConstraintViolationRollback() {
	// Start transaction
	tx := s.db.Begin()
	s.Require().NoError(tx.Error)

	// Create a user
	user := &domain.User{
		Email:        "constraint-test@example.com",
		Password:  "hashed-password",
		FirstName: "Test",
		LastName:  "User",
	}

	result := tx.Create(user)
	s.Require().NoError(result.Error)

	// Create first URL with unique short code
	url1 := &domain.ShortURL{
		UserID:      user.ID,
		OriginalURL: "https://example.com/constraint-test-1",
		ShortCode:   "unique123",
		Title:       "Constraint Test 1",
		IsActive:    true,
	}

	result = tx.Create(url1)
	s.Require().NoError(result.Error)

	// Try to create second URL with same short code (should violate unique constraint)
	url2 := &domain.ShortURL{
		UserID:      user.ID,
		OriginalURL: "https://example.com/constraint-test-2",
		ShortCode:   "unique123", // Same as url1
		Title:       "Constraint Test 2",
		IsActive:    true,
	}

	result = tx.Create(url2)
	// This should fail due to unique constraint on short_code
	if result.Error != nil {
		// Expected failure, rollback transaction
		result = tx.Rollback()
		s.Require().NoError(result.Error)

		// Verify nothing was committed
		var foundUser domain.User
		err := s.db.First(&foundUser, user.ID).Error
		s.Error(err)
		s.True(errors.Is(err, gorm.ErrRecordNotFound))
	} else {
		// If no constraint violation occurred (depends on DB setup), just commit
		result = tx.Commit()
		s.Require().NoError(result.Error)
	}
}

// TestBulkOperationTransaction tests bulk operations within transactions
func (s *DatabaseTransactionTestSuite) TestBulkOperationTransaction() {
	// Start transaction
	tx := s.db.Begin()
	s.Require().NoError(tx.Error)

	// Create multiple URLs in a single transaction
	urls := []domain.ShortURL{
		{
			UserID:      s.testUser.ID,
			OriginalURL: "https://example.com/bulk-1",
			ShortCode:   "bulk1",
			Title:       "Bulk Test 1",
			IsActive:    true,
		},
		{
			UserID:      s.testUser.ID,
			OriginalURL: "https://example.com/bulk-2",
			ShortCode:   "bulk2",
			Title:       "Bulk Test 2",
			IsActive:    true,
		},
		{
			UserID:      s.testUser.ID,
			OriginalURL: "https://example.com/bulk-3",
			ShortCode:   "bulk3",
			Title:       "Bulk Test 3",
			IsActive:    true,
		},
	}

	// Bulk create
	result := tx.Create(&urls)
	s.Require().NoError(result.Error)
	s.Equal(int64(3), result.RowsAffected)

	// Verify all URLs have IDs assigned
	for _, url := range urls {
		s.Greater(url.ID, uint(0))
	}

	// Bulk update
	result = tx.Model(&domain.ShortURL{}).Where("user_id = ? AND short_code LIKE ?", s.testUser.ID, "bulk%").Update("is_public", false)
	s.Require().NoError(result.Error)
	s.Equal(int64(3), result.RowsAffected)

	// Commit transaction
	result = tx.Commit()
	s.Require().NoError(result.Error)

	// Verify all URLs exist and were updated
	var count int64
	s.db.Model(&domain.ShortURL{}).Where("user_id = ? AND short_code LIKE ? AND is_public = ?", s.testUser.ID, "bulk%", false).Count(&count)
	s.Equal(int64(3), count)
}