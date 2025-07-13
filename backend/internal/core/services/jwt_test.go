package services

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTService_GenerateAccessToken(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := uint(123)
	email := "test@example.com"

	token, err := service.GenerateAccessToken(userID, email)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token can be parsed
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	// Verify claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, float64(userID), claims["user_id"])
	assert.Equal(t, email, claims["email"])
	assert.Equal(t, "access", claims["type"])
	assert.Equal(t, "url-shortener", claims["iss"])
}

func TestJWTService_GenerateRefreshToken(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := uint(123)

	token, err := service.GenerateRefreshToken(userID)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token can be parsed
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	// Verify claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, float64(userID), claims["user_id"])
	assert.Equal(t, "refresh", claims["type"])
	assert.Equal(t, "url-shortener", claims["iss"])
}

func TestJWTService_ValidateAccessToken(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := uint(123)
	email := "test@example.com"

	// Generate a valid token
	token, err := service.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	// Validate the token
	claims, err := service.ValidateAccessToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, "access", claims.Type)
	assert.NotEmpty(t, claims.JTI)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestJWTService_ValidateAccessToken_InvalidToken(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	// Test with invalid token
	_, err := service.ValidateAccessToken("invalid-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestJWTService_ValidateAccessToken_WrongTokenType(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := uint(123)

	// Generate a refresh token
	refreshToken, err := service.GenerateRefreshToken(userID)
	require.NoError(t, err)

	// Try to validate it as an access token
	_, err = service.ValidateAccessToken(refreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token type")
}

func TestJWTService_ValidateRefreshToken(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := uint(123)

	// Generate a valid refresh token
	token, err := service.GenerateRefreshToken(userID)
	require.NoError(t, err)

	// Validate the token
	claims, err := service.ValidateRefreshToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "refresh", claims.Type)
	assert.NotEmpty(t, claims.JTI)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestJWTService_ValidateRefreshToken_WrongTokenType(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := uint(123)
	email := "test@example.com"

	// Generate an access token
	accessToken, err := service.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	// Try to validate it as a refresh token
	_, err = service.ValidateRefreshToken(accessToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token type")
}

func TestJWTService_ExpiredToken(t *testing.T) {
	secretKey := "test-secret-key"
	service := &jwtService{
		secretKey:       secretKey,
		accessTokenTTL:  -time.Hour, // Expired
		refreshTokenTTL: time.Hour * 24 * 7,
		issuer:          "url-shortener",
	}

	userID := uint(123)
	email := "test@example.com"

	// Generate an expired token
	token, err := service.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	// Try to validate expired token
	_, err = service.ValidateAccessToken(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestJWTService_GetTokenTTL(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	accessTTL := service.GetTokenTTL("access")
	assert.Equal(t, time.Hour, accessTTL)

	refreshTTL := service.GetTokenTTL("refresh")
	assert.Equal(t, time.Hour*24*7, refreshTTL)

	defaultTTL := service.GetTokenTTL("unknown")
	assert.Equal(t, time.Hour, defaultTTL)
}

func TestJWTService_ExtractTokenFromHeader(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	// Test valid header
	token, err := service.ExtractTokenFromHeader("Bearer abc123")
	require.NoError(t, err)
	assert.Equal(t, "abc123", token)

	// Test invalid header format
	_, err = service.ExtractTokenFromHeader("abc123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid authorization header format")

	// Test short header
	_, err = service.ExtractTokenFromHeader("Bear")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid authorization header format")
}

func TestJWTService_GetUserIDFromToken(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := uint(123)
	email := "test@example.com"

	// Generate a valid token
	token, err := service.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	// Extract user ID
	extractedUserID, err := service.GetUserIDFromToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)

	// Test with invalid token
	_, err = service.GetUserIDFromToken("invalid-token")
	assert.Error(t, err)
}

func TestJWTService_IsTokenExpired(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := uint(123)
	email := "test@example.com"

	// Generate a valid token
	token, err := service.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	// Check if token is expired (should be false)
	isExpired := service.IsTokenExpired(token)
	assert.False(t, isExpired)

	// Test with invalid token
	isExpired = service.IsTokenExpired("invalid-token")
	assert.True(t, isExpired)
}

func TestJWTService_GetTokenClaims(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := uint(123)
	email := "test@example.com"

	// Generate a valid token
	token, err := service.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	// Get token claims
	claims, err := service.GetTokenClaims(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, "access", claims.Type)
	assert.NotEmpty(t, claims.JTI)

	// Test with invalid token
	_, err = service.GetTokenClaims("invalid-token")
	assert.Error(t, err)
}

func TestJWTService_DifferentSecretKeys(t *testing.T) {
	secretKey1 := "secret1"
	secretKey2 := "secret2"
	
	service1 := NewJWTService(secretKey1)
	service2 := NewJWTService(secretKey2)

	userID := uint(123)
	email := "test@example.com"

	// Generate token with first service
	token, err := service1.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	// Try to validate with second service (different secret)
	_, err = service2.ValidateAccessToken(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestJWTService_RevokeToken(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := uint(123)
	email := "test@example.com"

	// Generate a valid token
	token, err := service.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	// Revoke token (placeholder implementation)
	err = service.RevokeToken(token)
	assert.NoError(t, err)
}

func TestJWTService_generateJTI(t *testing.T) {
	secretKey := "test-secret-key"
	service := &jwtService{
		secretKey:       secretKey,
		accessTokenTTL:  time.Hour,
		refreshTokenTTL: time.Hour * 24 * 7,
		issuer:          "url-shortener",
	}

	jti1 := service.generateJTI()
	jti2 := service.generateJTI()

	assert.NotEmpty(t, jti1)
	assert.NotEmpty(t, jti2)
	assert.NotEqual(t, jti1, jti2) // JTIs should be unique
}

func TestJWTService_TokenSigningMethod(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := uint(123)
	email := "test@example.com"

	token, err := service.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	// Parse token and check signing method
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)
}

func TestJWTService_TokenStructure(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := uint(123)
	email := "test@example.com"

	token, err := service.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	// Verify token has 3 parts (header.payload.signature)
	parts := strings.Split(token, ".")
	assert.Len(t, parts, 3)
	assert.NotEmpty(t, parts[0]) // header
	assert.NotEmpty(t, parts[1]) // payload
	assert.NotEmpty(t, parts[2]) // signature
}