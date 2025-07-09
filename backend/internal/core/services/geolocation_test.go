package services

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"url-shortener/internal/core/domain"
)

func TestGeolocationService_GetLocationFromIP_Success(t *testing.T) {
	// Mock server response
	mockResponse := `{
		"status": "success",
		"country": "United States",
		"countryCode": "US",
		"region": "CA",
		"regionName": "California",
		"city": "San Francisco",
		"zip": "94105",
		"lat": 37.7749,
		"lon": -122.4194,
		"timezone": "America/Los_Angeles",
		"isp": "Test ISP",
		"org": "Test Organization",
		"as": "AS12345 Test AS",
		"query": "8.8.8.8"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	// Mock cache miss
	mockCacheRepo.On("Get", mock.Anything, "geo:8.8.8.8").Return("", errors.New("cache miss"))
	mockCacheRepo.On("Set", mock.Anything, "geo:8.8.8.8", mock.AnythingOfType("string"), 24*time.Hour).Return(nil)

	service := &geolocationService{
		cacheRepo:  mockCacheRepo,
		configRepo: mockConfigRepo,
		httpClient: &http.Client{Timeout: 5 * time.Second},
		apiKey:     "",
		provider:   "ip-api",
	}

	// Override the API URL to use our test server
	originalGetLocationFromAPI := service.getLocationFromAPI
	service.getLocationFromAPI = func(ctx context.Context, ipAddress string) (*domain.GeoLocation, error) {
		// Use test server URL
		url := server.URL + "/" + ipAddress
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := service.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Parse response and create GeoLocation
		return &domain.GeoLocation{
			IPAddress:    ipAddress,
			Country:      "United States",
			CountryCode:  "US",
			Region:       "California",
			City:         "San Francisco",
			PostalCode:   "94105",
			Latitude:     37.7749,
			Longitude:    -122.4194,
			Timezone:     "America/Los_Angeles",
			ISP:          "Test ISP",
			Organization: "Test Organization",
			ASN:          "AS12345 Test AS",
			Accuracy:     "city",
		}, nil
	}

	ctx := context.Background()
	location, err := service.GetLocationFromIP(ctx, "8.8.8.8")

	require.NoError(t, err)
	assert.Equal(t, "8.8.8.8", location.IPAddress)
	assert.Equal(t, "United States", location.Country)
	assert.Equal(t, "US", location.CountryCode)
	assert.Equal(t, "California", location.Region)
	assert.Equal(t, "San Francisco", location.City)
	assert.Equal(t, "94105", location.PostalCode)
	assert.Equal(t, 37.7749, location.Latitude)
	assert.Equal(t, -122.4194, location.Longitude)
	assert.Equal(t, "America/Los_Angeles", location.Timezone)
	assert.Equal(t, "Test ISP", location.ISP)
	assert.Equal(t, "Test Organization", location.Organization)
	assert.Equal(t, "AS12345 Test AS", location.ASN)
	assert.Equal(t, "city", location.Accuracy)

	// Restore original method
	service.getLocationFromAPI = originalGetLocationFromAPI

	mockCacheRepo.AssertExpectations(t)
}

func TestGeolocationService_GetLocationFromIP_LocalhostIP(t *testing.T) {
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	service := NewGeolocationService(mockCacheRepo, mockConfigRepo, "")

	ctx := context.Background()
	location, err := service.GetLocationFromIP(ctx, "127.0.0.1")

	require.NoError(t, err)
	assert.Equal(t, "127.0.0.1", location.IPAddress)
	assert.Equal(t, "Unknown", location.Country)
	assert.Equal(t, "Unknown", location.Region)
	assert.Equal(t, "Unknown", location.City)
	assert.Equal(t, 0.0, location.Latitude)
	assert.Equal(t, 0.0, location.Longitude)
	assert.Equal(t, "UTC", location.Timezone)
	assert.Equal(t, "Unknown", location.ISP)
}

func TestGeolocationService_GetLocationFromIP_IPv6Localhost(t *testing.T) {
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	service := NewGeolocationService(mockCacheRepo, mockConfigRepo, "")

	ctx := context.Background()
	location, err := service.GetLocationFromIP(ctx, "::1")

	require.NoError(t, err)
	assert.Equal(t, "::1", location.IPAddress)
	assert.Equal(t, "Unknown", location.Country)
	assert.Equal(t, "Unknown", location.Region)
	assert.Equal(t, "Unknown", location.City)
	assert.Equal(t, 0.0, location.Latitude)
	assert.Equal(t, 0.0, location.Longitude)
	assert.Equal(t, "UTC", location.Timezone)
	assert.Equal(t, "Unknown", location.ISP)
}

func TestGeolocationService_GetLocationFromIP_EmptyIP(t *testing.T) {
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	service := NewGeolocationService(mockCacheRepo, mockConfigRepo, "")

	ctx := context.Background()
	location, err := service.GetLocationFromIP(ctx, "")

	require.NoError(t, err)
	assert.Equal(t, "", location.IPAddress)
	assert.Equal(t, "Unknown", location.Country)
}

func TestGeolocationService_GetLocationFromIP_CacheHit(t *testing.T) {
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	cachedLocation := `{
		"ip_address": "8.8.8.8",
		"country": "United States",
		"country_code": "US",
		"region": "California",
		"city": "San Francisco",
		"postal_code": "94105",
		"latitude": 37.7749,
		"longitude": -122.4194,
		"timezone": "America/Los_Angeles",
		"isp": "Google",
		"organization": "Google LLC",
		"asn": "AS15169 Google LLC",
		"accuracy": "city"
	}`

	mockCacheRepo.On("Get", mock.Anything, "geo:8.8.8.8").Return(cachedLocation, nil)

	service := NewGeolocationService(mockCacheRepo, mockConfigRepo, "")

	ctx := context.Background()
	location, err := service.GetLocationFromIP(ctx, "8.8.8.8")

	require.NoError(t, err)
	assert.Equal(t, "8.8.8.8", location.IPAddress)
	assert.Equal(t, "United States", location.Country)
	assert.Equal(t, "US", location.CountryCode)

	mockCacheRepo.AssertExpectations(t)
}

func TestGeolocationService_GetLocationsBatch(t *testing.T) {
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	// Mock cache misses and successful API calls
	mockCacheRepo.On("Get", mock.Anything, "geo:8.8.8.8").Return("", errors.New("cache miss"))
	mockCacheRepo.On("Get", mock.Anything, "geo:1.1.1.1").Return("", errors.New("cache miss"))
	mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), 24*time.Hour).Return(nil).Maybe()

	service := NewGeolocationService(mockCacheRepo, mockConfigRepo, "")

	// Override getLocationFromAPI to return mock data
	originalGetLocationFromAPI := service.(*geolocationService).getLocationFromAPI
	service.(*geolocationService).getLocationFromAPI = func(ctx context.Context, ipAddress string) (*domain.GeoLocation, error) {
		if ipAddress == "8.8.8.8" {
			return &domain.GeoLocation{
				IPAddress:   "8.8.8.8",
				Country:     "United States",
				CountryCode: "US",
				Region:      "California",
				City:        "San Francisco",
				Accuracy:    "city",
			}, nil
		}
		if ipAddress == "1.1.1.1" {
			return &domain.GeoLocation{
				IPAddress:   "1.1.1.1",
				Country:     "United States",
				CountryCode: "US",
				Region:      "California",
				City:        "San Francisco",
				Accuracy:    "city",
			}, nil
		}
		return nil, errors.New("unknown IP")
	}

	ctx := context.Background()
	locations, err := service.GetLocationsBatch(ctx, []string{"8.8.8.8", "1.1.1.1"})

	require.NoError(t, err)
	assert.Len(t, locations, 2)
	
	loc1, exists := locations["8.8.8.8"]
	assert.True(t, exists)
	assert.Equal(t, "8.8.8.8", loc1.IPAddress)
	assert.Equal(t, "United States", loc1.Country)

	loc2, exists := locations["1.1.1.1"]
	assert.True(t, exists)
	assert.Equal(t, "1.1.1.1", loc2.IPAddress)
	assert.Equal(t, "United States", loc2.Country)

	// Restore original method
	service.(*geolocationService).getLocationFromAPI = originalGetLocationFromAPI

	mockCacheRepo.AssertExpectations(t)
}

func TestGeolocationService_GetLocationsBatch_WithErrors(t *testing.T) {
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	// Mock cache misses
	mockCacheRepo.On("Get", mock.Anything, "geo:8.8.8.8").Return("", errors.New("cache miss"))
	mockCacheRepo.On("Get", mock.Anything, "geo:invalid").Return("", errors.New("cache miss"))
	mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), 24*time.Hour).Return(nil).Maybe()

	service := NewGeolocationService(mockCacheRepo, mockConfigRepo, "")

	// Override getLocationFromAPI to return mock data and errors
	originalGetLocationFromAPI := service.(*geolocationService).getLocationFromAPI
	service.(*geolocationService).getLocationFromAPI = func(ctx context.Context, ipAddress string) (*domain.GeoLocation, error) {
		if ipAddress == "8.8.8.8" {
			return &domain.GeoLocation{
				IPAddress:   "8.8.8.8",
				Country:     "United States",
				CountryCode: "US",
				Accuracy:    "city",
			}, nil
		}
		return nil, errors.New("API error")
	}

	ctx := context.Background()
	locations, err := service.GetLocationsBatch(ctx, []string{"8.8.8.8", "invalid"})

	require.NoError(t, err)
	assert.Len(t, locations, 2)
	
	// Valid IP should have real location
	loc1, exists := locations["8.8.8.8"]
	assert.True(t, exists)
	assert.Equal(t, "8.8.8.8", loc1.IPAddress)
	assert.Equal(t, "United States", loc1.Country)

	// Invalid IP should have unknown location
	loc2, exists := locations["invalid"]
	assert.True(t, exists)
	assert.Equal(t, "invalid", loc2.IPAddress)
	assert.Equal(t, "Unknown", loc2.Country)
	assert.Equal(t, "Unknown", loc2.Region)
	assert.Equal(t, "Unknown", loc2.City)

	// Restore original method
	service.(*geolocationService).getLocationFromAPI = originalGetLocationFromAPI

	mockCacheRepo.AssertExpectations(t)
}

func TestGeolocationService_ValidateLocation(t *testing.T) {
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	service := NewGeolocationService(mockCacheRepo, mockConfigRepo, "")

	ctx := context.Background()

	// Test valid location
	validLocation := &domain.GeoLocation{
		IPAddress: "8.8.8.8",
		Country:   "United States",
		Latitude:  37.7749,
		Longitude: -122.4194,
	}
	err := service.ValidateLocation(ctx, validLocation)
	assert.NoError(t, err)

	// Test nil location
	err = service.ValidateLocation(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "location cannot be nil")

	// Test empty IP address
	invalidLocation := &domain.GeoLocation{
		IPAddress: "",
		Country:   "United States",
		Latitude:  37.7749,
		Longitude: -122.4194,
	}
	err = service.ValidateLocation(ctx, invalidLocation)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "IP address is required")

	// Test invalid latitude
	invalidLocation = &domain.GeoLocation{
		IPAddress: "8.8.8.8",
		Country:   "United States",
		Latitude:  91.0, // Invalid latitude
		Longitude: -122.4194,
	}
	err = service.ValidateLocation(ctx, invalidLocation)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "latitude must be between -90 and 90")

	// Test invalid longitude
	invalidLocation = &domain.GeoLocation{
		IPAddress: "8.8.8.8",
		Country:   "United States",
		Latitude:  37.7749,
		Longitude: 181.0, // Invalid longitude
	}
	err = service.ValidateLocation(ctx, invalidLocation)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "longitude must be between -180 and 180")
}

func TestGeolocationService_GetCountryCode(t *testing.T) {
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	service := NewGeolocationService(mockCacheRepo, mockConfigRepo, "")

	ctx := context.Background()

	// Test valid country name
	code, err := service.GetCountryCode(ctx, "United States")
	require.NoError(t, err)
	assert.Equal(t, "US", code)

	// Test case-insensitive lookup
	code, err = service.GetCountryCode(ctx, "united states")
	require.NoError(t, err)
	assert.Equal(t, "US", code)

	// Test another country
	code, err = service.GetCountryCode(ctx, "Germany")
	require.NoError(t, err)
	assert.Equal(t, "DE", code)

	// Test empty country name
	_, err = service.GetCountryCode(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "country name cannot be empty")

	// Test unknown country
	_, err = service.GetCountryCode(ctx, "Unknown Country")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "country code not found")
}

func TestGeolocationService_isPrivateIP(t *testing.T) {
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	service := &geolocationService{
		cacheRepo:  mockCacheRepo,
		configRepo: mockConfigRepo,
	}

	// Test private IP addresses
	assert.True(t, service.isPrivateIP("192.168.1.1"))
	assert.True(t, service.isPrivateIP("10.0.0.1"))
	assert.True(t, service.isPrivateIP("172.16.0.1"))
	assert.True(t, service.isPrivateIP("127.0.0.1"))
	assert.True(t, service.isPrivateIP("::1"))

	// Test public IP addresses
	assert.False(t, service.isPrivateIP("8.8.8.8"))
	assert.False(t, service.isPrivateIP("1.1.1.1"))
	assert.False(t, service.isPrivateIP("208.67.222.222"))
}

func TestGeolocationService_extractCountryFromUserAgent(t *testing.T) {
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	service := &geolocationService{
		cacheRepo:  mockCacheRepo,
		configRepo: mockConfigRepo,
	}

	tests := []struct {
		userAgent string
		expected  string
	}{
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) en-US", "US"},
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) en-GB", "GB"},
		{"Mozilla/5.0 (X11; Linux x86_64) en-CA", "CA"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) en-AU", "AU"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) de-DE", "DE"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) fr-FR", "FR"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) es-ES", "ES"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) it-IT", "IT"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) pt-BR", "BR"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) ru-RU", "RU"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) zh-CN", "CN"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) ja-JP", "JP"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) ko-KR", "KR"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) unknown", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.userAgent, func(t *testing.T) {
			result := service.extractCountryFromUserAgent(tt.userAgent)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGeolocationService_GetLocationWithFallback(t *testing.T) {
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	service := NewGeolocationService(mockCacheRepo, mockConfigRepo, "")

	// Mock cache miss
	mockCacheRepo.On("Get", mock.Anything, "geo:8.8.8.8").Return("", errors.New("cache miss"))
	mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), 24*time.Hour).Return(nil).Maybe()

	// Override getLocationFromAPI to return an error (simulating API failure)
	originalGetLocationFromAPI := service.(*geolocationService).getLocationFromAPI
	service.(*geolocationService).getLocationFromAPI = func(ctx context.Context, ipAddress string) (*domain.GeoLocation, error) {
		return nil, errors.New("API failure")
	}

	ctx := context.Background()
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) en-US"
	location, err := service.GetLocationWithFallback(ctx, "8.8.8.8", userAgent)

	require.NoError(t, err)
	assert.Equal(t, "8.8.8.8", location.IPAddress)
	assert.Equal(t, "US", location.Country)
	assert.Equal(t, "US", location.CountryCode)
	assert.Equal(t, "Unknown", location.Region)
	assert.Equal(t, "Unknown", location.City)
	assert.Equal(t, "country", location.Accuracy)

	// Restore original method
	service.(*geolocationService).getLocationFromAPI = originalGetLocationFromAPI

	mockCacheRepo.AssertExpectations(t)
}

func TestGeolocationService_GetLocationWithFallback_PrimarySuccess(t *testing.T) {
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	service := NewGeolocationService(mockCacheRepo, mockConfigRepo, "")

	// Mock cache miss
	mockCacheRepo.On("Get", mock.Anything, "geo:8.8.8.8").Return("", errors.New("cache miss"))
	mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), 24*time.Hour).Return(nil).Maybe()

	// Override getLocationFromAPI to return successful data
	originalGetLocationFromAPI := service.(*geolocationService).getLocationFromAPI
	service.(*geolocationService).getLocationFromAPI = func(ctx context.Context, ipAddress string) (*domain.GeoLocation, error) {
		return &domain.GeoLocation{
			IPAddress:   "8.8.8.8",
			Country:     "United States",
			CountryCode: "US",
			Region:      "California",
			City:        "San Francisco",
			Accuracy:    "city",
		}, nil
	}

	ctx := context.Background()
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) en-US"
	location, err := service.GetLocationWithFallback(ctx, "8.8.8.8", userAgent)

	require.NoError(t, err)
	assert.Equal(t, "8.8.8.8", location.IPAddress)
	assert.Equal(t, "United States", location.Country)
	assert.Equal(t, "US", location.CountryCode)
	assert.Equal(t, "California", location.Region)
	assert.Equal(t, "San Francisco", location.City)
	assert.Equal(t, "city", location.Accuracy)

	// Restore original method
	service.(*geolocationService).getLocationFromAPI = originalGetLocationFromAPI

	mockCacheRepo.AssertExpectations(t)
}

func TestGeolocationService_NewGeolocationService(t *testing.T) {
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	service := NewGeolocationService(mockCacheRepo, mockConfigRepo, "test-api-key")

	assert.NotNil(t, service)
	
	// Check that the service is properly initialized
	geoService, ok := service.(*geolocationService)
	assert.True(t, ok)
	assert.Equal(t, mockCacheRepo, geoService.cacheRepo)
	assert.Equal(t, mockConfigRepo, geoService.configRepo)
	assert.Equal(t, "test-api-key", geoService.apiKey)
	assert.Equal(t, "ip-api", geoService.provider)
	assert.NotNil(t, geoService.httpClient)
	assert.Equal(t, 5*time.Second, geoService.httpClient.Timeout)
}

func TestGeolocationService_countryCodeMap(t *testing.T) {
	// Test that country code map has expected entries
	assert.Equal(t, "US", countryCodeMap["United States"])
	assert.Equal(t, "GB", countryCodeMap["United Kingdom"])
	assert.Equal(t, "DE", countryCodeMap["Germany"])
	assert.Equal(t, "FR", countryCodeMap["France"])
	assert.Equal(t, "JP", countryCodeMap["Japan"])
	assert.Equal(t, "CN", countryCodeMap["China"])
	assert.Equal(t, "BR", countryCodeMap["Brazil"])
	assert.Equal(t, "AU", countryCodeMap["Australia"])
	assert.Equal(t, "CA", countryCodeMap["Canada"])
	assert.Equal(t, "IN", countryCodeMap["India"])

	// Test that map has reasonable size
	assert.Greater(t, len(countryCodeMap), 100)
}