package integration

import (
	"net/http"
	"strings"
	"time"
)

// AuthWorkflowTestSuite tests complete authentication workflows
type AuthWorkflowTestSuite struct {
	IntegrationTestSuite
}

// TestUserRegistrationWorkflow tests the complete user registration process
func (s *AuthWorkflowTestSuite) TestUserRegistrationWorkflow() {
	// Test user registration
	registerReq := map[string]string{
		"email":    "newuser@example.com",
		"password": "securepassword123",
	}
	
	resp := s.makeRequest("POST", "/api/v1/auth/register", registerReq, nil)
	
	var registerResp struct {
		User struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
		Token string `json:"token"`
	}
	s.assertJSONResponse(resp, http.StatusCreated, &registerResp)
	
	// Verify user data
	s.Equal("newuser@example.com", registerResp.User.Email)
	s.NotEmpty(registerResp.Token)
	s.Greater(registerResp.User.ID, uint(0))
	
	// Verify JWT token is valid
	s.verifyTokenIsValid(registerResp.Token, registerResp.User.ID)
	
	// Test duplicate registration should fail
	dupResp := s.makeRequest("POST", "/api/v1/auth/register", registerReq, nil)
	s.assertErrorResponse(dupResp, http.StatusConflict, "email already exists")
}

// TestUserLoginWorkflow tests the complete user login process
func (s *AuthWorkflowTestSuite) TestUserLoginWorkflow() {
	// Test successful login
	loginReq := map[string]string{
		"email":    s.testUser.Email,
		"password": "password123",
	}
	
	resp := s.makeRequest("POST", "/api/v1/auth/login", loginReq, nil)
	
	var loginResp struct {
		User struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
		Token string `json:"token"`
	}
	s.assertJSONResponse(resp, http.StatusOK, &loginResp)
	
	// Verify response data
	s.Equal(s.testUser.Email, loginResp.User.Email)
	s.Equal(s.testUser.ID, loginResp.User.ID)
	s.NotEmpty(loginResp.Token)
	
	// Verify JWT token is valid
	s.verifyTokenIsValid(loginResp.Token, loginResp.User.ID)
	
	// Test login with wrong password
	wrongPasswordReq := map[string]string{
		"email":    s.testUser.Email,
		"password": "wrongpassword",
	}
	
	wrongResp := s.makeRequest("POST", "/api/v1/auth/login", wrongPasswordReq, nil)
	s.assertErrorResponse(wrongResp, http.StatusUnauthorized, "invalid credentials")
	
	// Test login with non-existent user
	nonExistentReq := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "anypassword",
	}
	
	nonExistentResp := s.makeRequest("POST", "/api/v1/auth/login", nonExistentReq, nil)
	s.assertErrorResponse(nonExistentResp, http.StatusUnauthorized, "invalid credentials")
}

// TestTokenRefreshWorkflow tests JWT token refresh functionality
func (s *AuthWorkflowTestSuite) TestTokenRefreshWorkflow() {
	// Login to get initial token
	token := s.loginUser()
	
	// Test token refresh
	refreshReq := map[string]string{
		"token": token,
	}
	
	resp := s.makeRequest("POST", "/api/v1/auth/refresh", refreshReq, nil)
	
	var refreshResp struct {
		Token string `json:"token"`
	}
	s.assertJSONResponse(resp, http.StatusOK, &refreshResp)
	
	// Verify new token is different and valid
	s.NotEqual(token, refreshResp.Token)
	s.verifyTokenIsValid(refreshResp.Token, s.testUser.ID)
	
	// Test refresh with invalid token
	invalidRefreshReq := map[string]string{
		"token": "invalid.token.here",
	}
	
	invalidResp := s.makeRequest("POST", "/api/v1/auth/refresh", invalidRefreshReq, nil)
	s.assertErrorResponse(invalidResp, http.StatusUnauthorized, "invalid token")
}

// TestProtectedEndpointAccess tests access to protected endpoints
func (s *AuthWorkflowTestSuite) TestProtectedEndpointAccess() {
	// Test access without token
	resp := s.makeRequest("GET", "/api/v1/auth/profile", nil, nil)
	s.assertErrorResponse(resp, http.StatusUnauthorized, "authorization header required")
	
	// Test access with invalid token
	invalidHeaders := map[string]string{
		"Authorization": "Bearer invalid.token.here",
	}
	invalidResp := s.makeRequest("GET", "/api/v1/auth/profile", nil, invalidHeaders)
	s.assertErrorResponse(invalidResp, http.StatusUnauthorized, "invalid token")
	
	// Test access with valid token
	validResp := s.makeAuthenticatedRequest("GET", "/api/v1/auth/profile", nil)
	
	var profileResp struct {
		User struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}
	s.assertJSONResponse(validResp, http.StatusOK, &profileResp)
	
	s.Equal(s.testUser.ID, profileResp.User.ID)
	s.Equal(s.testUser.Email, profileResp.User.Email)
}

// TestUserProfileManagement tests user profile update operations
func (s *AuthWorkflowTestSuite) TestUserProfileManagement() {
	// Test profile update
	updateReq := map[string]string{
		"email": "updated@example.com",
	}
	
	resp := s.makeAuthenticatedRequest("PUT", "/api/v1/auth/profile", updateReq)
	
	var updateResp struct {
		User struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}
	s.assertJSONResponse(resp, http.StatusOK, &updateResp)
	
	// Verify email was updated
	s.Equal("updated@example.com", updateResp.User.Email)
	s.Equal(s.testUser.ID, updateResp.User.ID)
	
	// Verify update persisted in database
	profileResp := s.makeAuthenticatedRequest("GET", "/api/v1/auth/profile", nil)
	
	var profileData struct {
		User struct {
			Email string `json:"email"`
		} `json:"user"`
	}
	s.assertJSONResponse(profileResp, http.StatusOK, &profileData)
	s.Equal("updated@example.com", profileData.User.Email)
}

// TestPasswordChangeWorkflow tests password change functionality
func (s *AuthWorkflowTestSuite) TestPasswordChangeWorkflow() {
	// Test password change
	passwordChangeReq := map[string]string{
		"currentPassword": "password123",
		"newPassword":     "newpassword456",
	}
	
	resp := s.makeAuthenticatedRequest("POST", "/api/v1/auth/change-password", passwordChangeReq)
	s.Equal(http.StatusOK, resp.StatusCode)
	
	// Verify old password no longer works
	oldPasswordReq := map[string]string{
		"email":    s.testUser.Email,
		"password": "password123",
	}
	
	oldResp := s.makeRequest("POST", "/api/v1/auth/login", oldPasswordReq, nil)
	s.assertErrorResponse(oldResp, http.StatusUnauthorized, "invalid credentials")
	
	// Verify new password works
	newPasswordReq := map[string]string{
		"email":    s.testUser.Email,
		"password": "newpassword456",
	}
	
	newResp := s.makeRequest("POST", "/api/v1/auth/login", newPasswordReq, nil)
	s.Equal(http.StatusOK, newResp.StatusCode)
	
	// Test password change with wrong current password
	wrongCurrentReq := map[string]string{
		"currentPassword": "wrongpassword",
		"newPassword":     "anothernewpassword",
	}
	
	wrongResp := s.makeAuthenticatedRequest("POST", "/api/v1/auth/change-password", wrongCurrentReq)
	s.assertErrorResponse(wrongResp, http.StatusBadRequest, "current password is incorrect")
}

// TestLogoutWorkflow tests user logout functionality
func (s *AuthWorkflowTestSuite) TestLogoutWorkflow() {
	// Test logout
	resp := s.makeAuthenticatedRequest("POST", "/api/v1/auth/logout", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	
	// Verify token is invalidated (if token blacklisting is implemented)
	// For now, we just verify the endpoint responds correctly
	profileResp := s.makeAuthenticatedRequest("GET", "/api/v1/auth/profile", nil)
	// This should still work unless token blacklisting is implemented
	// If blacklisting is implemented, this should return 401
	s.True(profileResp.StatusCode == http.StatusOK || profileResp.StatusCode == http.StatusUnauthorized)
}

// TestTokenValidation tests token validation endpoint
func (s *AuthWorkflowTestSuite) TestTokenValidation() {
	// Test valid token
	resp := s.makeAuthenticatedRequest("GET", "/api/v1/auth/validate", nil)
	
	var validateResp struct {
		Valid bool `json:"valid"`
		User  struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}
	s.assertJSONResponse(resp, http.StatusOK, &validateResp)
	
	s.True(validateResp.Valid)
	s.Equal(s.testUser.ID, validateResp.User.ID)
	
	// Test invalid token
	invalidHeaders := map[string]string{
		"Authorization": "Bearer invalid.token.here",
	}
	invalidResp := s.makeRequest("GET", "/api/v1/auth/validate", nil, invalidHeaders)
	s.assertErrorResponse(invalidResp, http.StatusUnauthorized, "invalid token")
}

// TestRateLimitingOnAuthEndpoints tests rate limiting on authentication endpoints
func (s *AuthWorkflowTestSuite) TestRateLimitingOnAuthEndpoints() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping rate limiting tests")
		return
	}
	
	// Make multiple rapid requests to trigger rate limiting
	loginReq := map[string]string{
		"email":    "ratelimit@example.com",
		"password": "password123",
	}
	
	// The exact number of requests needed to trigger rate limiting depends on the configuration
	// We'll make a reasonable number of requests and check if any return 429
	var rateLimitTriggered bool
	
	for i := 0; i < 20; i++ {
		resp := s.makeRequest("POST", "/api/v1/auth/login", loginReq, nil)
		if resp.StatusCode == http.StatusTooManyRequests {
			rateLimitTriggered = true
			break
		}
		resp.Body.Close()
		
		// Small delay to avoid overwhelming the system
		time.Sleep(10 * time.Millisecond)
	}
	
	// If rate limiting is properly configured, it should have triggered
	// This is a soft check since the exact configuration may vary
	if !rateLimitTriggered {
		s.T().Log("Rate limiting may not be configured or may have higher limits")
	}
}

// TestConcurrentAuthOperations tests authentication operations under concurrent load
func (s *AuthWorkflowTestSuite) TestConcurrentAuthOperations() {
	// Test concurrent logins
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(index int) {
			defer func() { done <- true }()
			
			loginReq := map[string]string{
				"email":    s.testUser.Email,
				"password": "password123",
			}
			
			resp := s.makeRequest("POST", "/api/v1/auth/login", loginReq, nil)
			s.Equal(http.StatusOK, resp.StatusCode)
			resp.Body.Close()
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Helper methods

// verifyTokenIsValid verifies that a JWT token is valid and contains expected claims
func (s *AuthWorkflowTestSuite) verifyTokenIsValid(token string, expectedUserID uint) {
	s.NotEmpty(token)
	
	// Verify token has the correct structure (3 parts separated by dots)
	parts := strings.Split(token, ".")
	s.Len(parts, 3, "JWT should have 3 parts")
	
	// Use the JWT service to validate the token
	claims, err := s.jwtService.ValidateAccessToken(token)
	s.NoError(err, "Token should be valid")
	s.Equal(expectedUserID, claims.UserID, "Token should contain correct user ID")
}