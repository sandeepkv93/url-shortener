package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type urlService struct {
	urlRepo     ports.URLRepository
	clickRepo   ports.ClickRepository
	cacheRepo   ports.CacheService
	configRepo  ports.ConfigService
}

const (
	shortCodeChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortCodeLength = 6
	maxRetries = 10
)

func NewURLService(
	urlRepo ports.URLRepository,
	clickRepo ports.ClickRepository,
	cacheRepo ports.CacheService,
	configRepo ports.ConfigService,
) ports.URLService {
	return &urlService{
		urlRepo:    urlRepo,
		clickRepo:  clickRepo,
		cacheRepo:  cacheRepo,
		configRepo: configRepo,
	}
}

func (s *urlService) ShortenURL(ctx context.Context, req domain.ShortenURLRequest) (*domain.ShortURL, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Validate URL format
	if !s.isValidURL(req.OriginalURL) {
		return nil, domain.ErrInvalidURL
	}

	// Generate unique short code
	shortCode, err := s.generateUniqueShortCode(ctx, req.CustomAlias)
	if err != nil {
		return nil, fmt.Errorf("failed to generate short code: %w", err)
	}

	// Create short URL
	shortURL := &domain.ShortURL{
		ShortCode:   shortCode,
		OriginalURL: req.OriginalURL,
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		IsActive:    true,
		ClickCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set expiration if provided
	if req.ExpiresAt != nil {
		shortURL.ExpiresAt = req.ExpiresAt
	}

	// Set password if provided
	if req.Password != "" {
		hashedPassword, err := s.hashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		shortURL.Password = &hashedPassword
	}

	// Save to database
	if err := s.urlRepo.Create(ctx, shortURL); err != nil {
		return nil, fmt.Errorf("failed to create short URL: %w", err)
	}

	// Cache the URL for faster lookup
	if err := s.cacheRepo.CacheURL(ctx, shortCode, shortURL.OriginalURL, shortURL.UserID, time.Hour*24); err != nil {
		// Log error but don't fail the creation
		fmt.Printf("Failed to cache URL: %v", err)
	}

	return shortURL, nil
}

func (s *urlService) GetOriginalURL(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	// Try cache first
	cachedURL, _, err := s.cacheRepo.GetCachedURL(ctx, shortCode)
	if err == nil && cachedURL != "" {
		// Get full URL details from database for complete response
		shortURL, dbErr := s.urlRepo.GetActiveByShortCode(ctx, shortCode)
		if dbErr == nil {
			return shortURL, nil
		}
	}

	// Get from database
	shortURL, err := s.urlRepo.GetActiveByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	// Cache for future requests
	if err := s.cacheRepo.CacheURL(ctx, shortCode, shortURL.OriginalURL, shortURL.UserID, time.Hour*24); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to cache URL: %v", err)
	}

	return shortURL, nil
}

func (s *urlService) GetUserURLs(ctx context.Context, userID uint, offset, limit int) ([]*domain.ShortURL, int64, error) {
	return s.urlRepo.GetByUserID(ctx, userID, offset, limit)
}

func (s *urlService) UpdateURL(ctx context.Context, id uint, userID uint, req domain.UpdateURLRequest) (*domain.ShortURL, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get existing URL
	shortURL, err := s.urlRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if shortURL.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	// Update fields
	if req.Title != nil {
		shortURL.Title = *req.Title
	}
	if req.Description != nil {
		shortURL.Description = *req.Description
	}
	if req.IsActive != nil {
		shortURL.IsActive = *req.IsActive
	}
	if req.ExpiresAt != nil {
		shortURL.ExpiresAt = req.ExpiresAt
	}

	shortURL.UpdatedAt = time.Now()

	// Update in database
	if err := s.urlRepo.Update(ctx, shortURL); err != nil {
		return nil, fmt.Errorf("failed to update URL: %w", err)
	}

	// Update cache
	if shortURL.IsActive {
		if err := s.cacheRepo.CacheURL(ctx, shortURL.ShortCode, shortURL.OriginalURL, shortURL.UserID, time.Hour*24); err != nil {
			fmt.Printf("Failed to update cache: %v", err)
		}
	} else {
		// Remove from cache if deactivated
		if err := s.cacheRepo.InvalidateURL(ctx, shortURL.ShortCode); err != nil {
			fmt.Printf("Failed to remove from cache: %v", err)
		}
	}

	return shortURL, nil
}

func (s *urlService) DeleteURL(ctx context.Context, id uint, userID uint) error {
	// Get existing URL
	shortURL, err := s.urlRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check ownership
	if shortURL.UserID != userID {
		return domain.ErrUnauthorized
	}

	// Delete from database
	if err := s.urlRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete URL: %w", err)
	}

	// Remove from cache
	if err := s.cacheRepo.InvalidateURL(ctx, shortURL.ShortCode); err != nil {
		fmt.Printf("Failed to remove from cache: %v", err)
	}

	return nil
}

func (s *urlService) RecordClick(ctx context.Context, shortURL *domain.ShortURL, clickData domain.ClickData) error {
	// Create click record
	click := &domain.Click{
		ShortURLID: shortURL.ID,
		IPAddress:  clickData.IPAddress,
		UserAgent:  clickData.UserAgent,
		Referer:    clickData.Referer,
		Country:    clickData.Country,
		Region:     clickData.Region,
		City:       clickData.City,
		Device:     clickData.Device,
		Browser:    clickData.Browser,
		OS:         clickData.OS,
		ClickedAt:  time.Now(),
	}

	// Save click record
	if err := s.clickRepo.Create(ctx, click); err != nil {
		return fmt.Errorf("failed to record click: %w", err)
	}

	// Increment click count
	if err := s.urlRepo.IncrementClickCount(ctx, shortURL.ID); err != nil {
		return fmt.Errorf("failed to increment click count: %w", err)
	}

	// Cache unique click for analytics
	shortCodeStr := fmt.Sprintf("%d", shortURL.ID)
	if _, err := s.cacheRepo.CacheUniqueClick(ctx, shortCodeStr, clickData.IPAddress); err != nil {
		fmt.Printf("Failed to cache unique click: %v", err)
	}

	return nil
}

func (s *urlService) GetURLStats(ctx context.Context, id uint, userID uint) (*domain.URLStats, error) {
	// Get URL
	shortURL, err := s.urlRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if shortURL.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	// Get click stats
	clickStats, err := s.clickRepo.GetClickStats(ctx, shortURL.ID, "month")
	if err != nil {
		return nil, fmt.Errorf("failed to get click stats: %w", err)
	}

	// Get geo stats
	geoStats, err := s.clickRepo.GetGeoStats(ctx, shortURL.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get geo stats: %w", err)
	}

	// Get timeline stats
	timelineStats, err := s.clickRepo.GetTimelineStats(ctx, shortURL.ID, "week")
	if err != nil {
		return nil, fmt.Errorf("failed to get timeline stats: %w", err)
	}

	return &domain.URLStats{
		ShortURL:      shortURL,
		ClickStats:    clickStats,
		GeoStats:      geoStats,
		TimelineStats: timelineStats,
	}, nil
}

func (s *urlService) GetPopularURLs(ctx context.Context, limit int) ([]*domain.ShortURL, error) {
	return s.urlRepo.GetPopularURLs(ctx, limit)
}

func (s *urlService) ValidatePassword(ctx context.Context, shortCode, password string) (bool, error) {
	shortURL, err := s.urlRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return false, err
	}

	if shortURL.Password == nil {
		return true, nil // No password required
	}

	return s.checkPassword(password, *shortURL.Password), nil
}

func (s *urlService) CleanupExpiredURLs(ctx context.Context) error {
	expiredURLs, err := s.urlRepo.GetExpiredURLs(ctx, 100)
	if err != nil {
		return fmt.Errorf("failed to get expired URLs: %w", err)
	}

	for _, url := range expiredURLs {
		// Deactivate expired URL
		url.IsActive = false
		url.UpdatedAt = time.Now()
		
		if err := s.urlRepo.Update(ctx, url); err != nil {
			fmt.Printf("Failed to deactivate expired URL %s: %v", url.ShortCode, err)
			continue
		}

		// Remove from cache
		if err := s.cacheRepo.InvalidateURL(ctx, url.ShortCode); err != nil {
			fmt.Printf("Failed to remove expired URL from cache %s: %v", url.ShortCode, err)
		}
	}

	return nil
}

func (s *urlService) generateUniqueShortCode(ctx context.Context, customAlias string) (string, error) {
	// If custom alias is provided, validate and use it
	if customAlias != "" {
		if !s.isValidShortCode(customAlias) {
			return "", domain.ErrInvalidShortCode
		}

		exists, err := s.urlRepo.ExistsByShortCode(ctx, customAlias)
		if err != nil {
			return "", fmt.Errorf("failed to check short code existence: %w", err)
		}
		if exists {
			return "", domain.ErrShortCodeExists
		}

		return customAlias, nil
	}

	// Generate random short code
	for i := 0; i < maxRetries; i++ {
		shortCode, err := s.generateRandomShortCode()
		if err != nil {
			return "", fmt.Errorf("failed to generate random short code: %w", err)
		}

		exists, err := s.urlRepo.ExistsByShortCode(ctx, shortCode)
		if err != nil {
			return "", fmt.Errorf("failed to check short code existence: %w", err)
		}
		if !exists {
			return shortCode, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique short code after %d retries", maxRetries)
}

func (s *urlService) generateRandomShortCode() (string, error) {
	result := make([]byte, shortCodeLength)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(shortCodeChars))))
		if err != nil {
			return "", err
		}
		result[i] = shortCodeChars[num.Int64()]
	}
	return string(result), nil
}

func (s *urlService) isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	return u.Scheme == "http" || u.Scheme == "https"
}

func (s *urlService) isValidShortCode(shortCode string) bool {
	if len(shortCode) < 3 || len(shortCode) > 20 {
		return false
	}

	for _, char := range shortCode {
		if !strings.ContainsRune(shortCodeChars, char) {
			return false
		}
	}

	return true
}

func (s *urlService) hashPassword(password string) (string, error) {
	// Simple hash for URL passwords (not as secure as bcrypt for user passwords)
	// In production, you might want to use bcrypt here too
	return fmt.Sprintf("%x", password), nil
}

func (s *urlService) checkPassword(password, hashedPassword string) bool {
	expectedHash := fmt.Sprintf("%x", password)
	return expectedHash == hashedPassword
}