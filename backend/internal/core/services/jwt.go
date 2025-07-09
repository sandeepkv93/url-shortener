package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type jwtService struct {
	secretKey        string
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	issuer           string
}

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

func NewJWTService(secretKey string) ports.JWTService {
	return &jwtService{
		secretKey:        secretKey,
		accessTokenTTL:   time.Hour,     // 1 hour
		refreshTokenTTL:  time.Hour * 24 * 7, // 7 days
		issuer:           "url-shortener",
	}
}

func (s *jwtService) GenerateAccessToken(userID uint, email string) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			ID:        s.generateJTI(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) GenerateRefreshToken(userID uint) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		UserID: userID,
		Email:  "",
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			ID:        s.generateJTI(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ValidateAccessToken(tokenString string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Check if it's an access token
	if claims.Type != "access" {
		return nil, fmt.Errorf("invalid token type: expected access, got %s", claims.Type)
	}

	// Check expiration
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token expired")
	}

	return &domain.TokenClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		Type:      claims.Type,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
		JTI:       claims.ID,
	}, nil
}

func (s *jwtService) ValidateRefreshToken(tokenString string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Check if it's a refresh token
	if claims.Type != "refresh" {
		return nil, fmt.Errorf("invalid token type: expected refresh, got %s", claims.Type)
	}

	// Check expiration
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token expired")
	}

	return &domain.TokenClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		Type:      claims.Type,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
		JTI:       claims.ID,
	}, nil
}

func (s *jwtService) generateJTI() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

// GetTokenTTL returns the TTL for different token types
func (s *jwtService) GetTokenTTL(tokenType string) time.Duration {
	switch tokenType {
	case "access":
		return s.accessTokenTTL
	case "refresh":
		return s.refreshTokenTTL
	default:
		return s.accessTokenTTL
	}
}

// RevokeToken would be used to revoke a token (implementation depends on token blacklisting strategy)
func (s *jwtService) RevokeToken(tokenString string) error {
	// In a production system, you would implement token revocation
	// by maintaining a blacklist of revoked tokens (JTI) in Redis or database
	// For now, this is a placeholder
	return nil
}

// ExtractTokenFromHeader extracts JWT token from Authorization header
func (s *jwtService) ExtractTokenFromHeader(authHeader string) (string, error) {
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}
	
	if authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", fmt.Errorf("authorization header must start with Bearer")
	}
	
	return authHeader[len(bearerPrefix):], nil
}

// GetUserIDFromToken extracts user ID from a valid token
func (s *jwtService) GetUserIDFromToken(tokenString string) (uint, error) {
	claims, err := s.ValidateAccessToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// IsTokenExpired checks if a token is expired without validating the signature
func (s *jwtService) IsTokenExpired(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	
	if err != nil {
		return true
	}
	
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return true
	}
	
	return time.Now().After(claims.ExpiresAt.Time)
}

// GetTokenClaims extracts claims from a token without validating expiration
func (s *jwtService) GetTokenClaims(tokenString string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &domain.TokenClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		Type:      claims.Type,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
		JTI:       claims.ID,
	}, nil
}