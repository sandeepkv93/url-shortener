package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"url-shortener/internal/api/middleware"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type URLHandler struct {
	urlService       ports.URLService
	analyticsService ports.AnalyticsService
}

func NewURLHandler(urlService ports.URLService, analyticsService ports.AnalyticsService) *URLHandler {
	return &URLHandler{
		urlService:       urlService,
		analyticsService: analyticsService,
	}
}

// CreateShortURL handles URL shortening requests
func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req domain.ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set user ID from context
	req.UserID = userID

	// Validate request
	if err := req.Validate(); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create short URL
	shortURL, err := h.urlService.ShortenURL(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrInvalidURL:
			h.writeErrorResponse(w, "Invalid URL format", http.StatusBadRequest)
		case domain.ErrShortCodeExists:
			h.writeErrorResponse(w, "Custom alias already exists", http.StatusConflict)
		case domain.ErrInvalidShortCode:
			h.writeErrorResponse(w, "Invalid custom alias format", http.StatusBadRequest)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, shortURL, http.StatusCreated)
}

// GetUserURLs handles getting user's URLs with pagination
func (h *URLHandler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse pagination parameters
	offset, limit := h.parsePaginationParams(r)

	// Get user URLs
	urls, total, err := h.urlService.GetUserURLs(r.Context(), userID, offset, limit)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"urls":   urls,
		"total":  total,
		"offset": offset,
		"limit":  limit,
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// GetURL handles getting a specific URL by ID
func (h *URLHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse URL ID
	urlIDStr := chi.URLParam(r, "id")
	urlID, err := strconv.ParseUint(urlIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	// Get URL stats
	stats, err := h.urlService.GetURLStats(r.Context(), uint(urlID), userID)
	if err != nil {
		switch err {
		case domain.ErrURLNotFound:
			h.writeErrorResponse(w, "URL not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, stats, http.StatusOK)
}

// UpdateURL handles updating URL properties
func (h *URLHandler) UpdateURL(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse URL ID
	urlIDStr := chi.URLParam(r, "id")
	urlID, err := strconv.ParseUint(urlIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	var req domain.UpdateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update URL
	updatedURL, err := h.urlService.UpdateURL(r.Context(), uint(urlID), userID, req)
	if err != nil {
		switch err {
		case domain.ErrURLNotFound:
			h.writeErrorResponse(w, "URL not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, updatedURL, http.StatusOK)
}

// DeleteURL handles URL deletion
func (h *URLHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse URL ID
	urlIDStr := chi.URLParam(r, "id")
	urlID, err := strconv.ParseUint(urlIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	// Delete URL
	if err := h.urlService.DeleteURL(r.Context(), uint(urlID), userID); err != nil {
		switch err {
		case domain.ErrURLNotFound:
			h.writeErrorResponse(w, "URL not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, map[string]string{"message": "URL deleted successfully"}, http.StatusOK)
}

// RedirectURL handles short URL redirection
func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")
	if shortCode == "" {
		h.writeErrorResponse(w, "Short code is required", http.StatusBadRequest)
		return
	}

	// Get original URL
	shortURL, err := h.urlService.GetOriginalURL(r.Context(), shortCode)
	if err != nil {
		switch err {
		case domain.ErrURLNotFound:
			h.writeErrorResponse(w, "URL not found", http.StatusNotFound)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Check if URL is expired or inactive
	if !shortURL.IsActive {
		h.writeErrorResponse(w, "URL is inactive", http.StatusGone)
		return
	}

	// Check password protection
	if shortURL.Password != nil {
		password := r.URL.Query().Get("password")
		if password == "" {
			h.writeErrorResponse(w, "Password required", http.StatusUnauthorized)
			return
		}

		valid, err := h.urlService.ValidatePassword(r.Context(), shortCode, password)
		if err != nil || !valid {
			h.writeErrorResponse(w, "Invalid password", http.StatusUnauthorized)
			return
		}
	}

	// Record click analytics
	clickData := h.extractClickData(r)
	if err := h.urlService.RecordClick(r.Context(), shortURL, clickData); err != nil {
		// Log error but don't fail the redirect
		// In production, you might want to use a proper logger
	}

	// Redirect to original URL
	http.Redirect(w, r, shortURL.OriginalURL, http.StatusMovedPermanently)
}

// GetPopularURLs handles getting popular URLs (public endpoint)
func (h *URLHandler) GetPopularURLs(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = parsedLimit
		}
	}

	urls, err := h.urlService.GetPopularURLs(r.Context(), limit)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"urls":  urls,
		"limit": limit,
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// Helper methods

func (h *URLHandler) parsePaginationParams(r *http.Request) (offset, limit int) {
	offset = 0
	limit = 20 // default limit

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	return offset, limit
}

func (h *URLHandler) extractClickData(r *http.Request) domain.ClickData {
	// Get client IP
	clientIP := r.Header.Get("X-Real-IP")
	if clientIP == "" {
		clientIP = r.Header.Get("X-Forwarded-For")
	}
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}

	// Get user agent
	userAgent := r.Header.Get("User-Agent")

	// Get referer
	referer := r.Header.Get("Referer")

	// For a production implementation, you would:
	// 1. Parse user agent to extract device, browser, OS
	// 2. Use IP geolocation service to get country, region, city
	// 3. Implement proper device/browser detection

	return domain.ClickData{
		IPAddress: clientIP,
		UserAgent: userAgent,
		Referer:   referer,
		Country:   "Unknown", // Would be determined by IP geolocation
		Region:    "Unknown",
		City:      "Unknown",
		Device:    "Unknown", // Would be parsed from user agent
		Browser:   "Unknown", // Would be parsed from user agent
		OS:        "Unknown", // Would be parsed from user agent
	}
}

func (h *URLHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to encode response"}`))
	}
}

func (h *URLHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Write([]byte(`{"error": "Internal server error"}`))
	}
}