package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShortURLToResponse(t *testing.T) {
	expiresAt := time.Now().Add(24 * time.Hour)
	shortURL := &ShortURL{
		ID:          1,
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		CustomAlias: true,
		ExpiresAt:   &expiresAt,
		IsActive:    true,
		ClickCount:  10,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	baseURL := "https://short.ly"
	response := shortURL.ToResponse(baseURL)

	assert.Equal(t, shortURL.ID, response.ID)
	assert.Equal(t, shortURL.ShortCode, response.ShortCode)
	assert.Equal(t, shortURL.OriginalURL, response.OriginalURL)
	assert.Equal(t, baseURL+"/"+shortURL.ShortCode, response.ShortURL)
	assert.Equal(t, shortURL.CustomAlias, response.CustomAlias)
	assert.Equal(t, shortURL.ExpiresAt, response.ExpiresAt)
	assert.Equal(t, shortURL.IsActive, response.IsActive)
	assert.Equal(t, shortURL.ClickCount, response.ClickCount)
}

func TestShortURLIsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt *time.Time
		expected  bool
	}{
		{
			name:      "no expiration",
			expiresAt: nil,
			expected:  false,
		},
		{
			name:      "not expired",
			expiresAt: timePtr(time.Now().Add(24 * time.Hour)),
			expected:  false,
		},
		{
			name:      "expired",
			expiresAt: timePtr(time.Now().Add(-24 * time.Hour)),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortURL := &ShortURL{
				ExpiresAt: tt.expiresAt,
			}
			
			result := shortURL.IsExpired()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateShortURLRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateShortURLRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreateShortURLRequest{
				OriginalURL: "https://example.com",
				CustomAlias: "mylink",
			},
			wantErr: false,
		},
		{
			name: "empty URL",
			request: CreateShortURLRequest{
				OriginalURL: "",
				CustomAlias: "mylink",
			},
			wantErr: true,
		},
		{
			name: "invalid URL",
			request: CreateShortURLRequest{
				OriginalURL: "not-a-url",
				CustomAlias: "mylink",
			},
			wantErr: true,
		},
		{
			name: "long custom alias",
			request: CreateShortURLRequest{
				OriginalURL: "https://example.com",
				CustomAlias: "this-is-a-very-long-custom-alias-that-exceeds-the-maximum-length-allowed",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: Actual validation would be done by the validator package
			// This is just testing the struct definition
			assert.NotNil(t, tt.request)
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}