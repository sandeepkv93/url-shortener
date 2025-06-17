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

type QRHandler struct {
	qrService ports.QRService
}

func NewQRHandler(qrService ports.QRService) *QRHandler {
	return &QRHandler{
		qrService: qrService,
	}
}

// GenerateQRCode handles QR code generation requests
func (h *QRHandler) GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	var req domain.QRCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set user ID if authenticated
	if userID := middleware.GetUserIDFromContext(r.Context()); userID != 0 {
		req.UserID = userID
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate QR code
	qrResponse, err := h.qrService.GenerateQRCode(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrInvalidRequest:
			h.writeErrorResponse(w, "Invalid QR code request", http.StatusBadRequest)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		case domain.ErrURLNotFound:
			h.writeErrorResponse(w, "Short URL not found", http.StatusNotFound)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return QR code as binary data with appropriate headers
	w.Header().Set("Content-Type", qrResponse.MimeType)
	w.Header().Set("Content-Length", strconv.Itoa(len(qrResponse.Data)))
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	
	// Add custom headers with QR code info
	w.Header().Set("X-QR-Format", qrResponse.Format)
	w.Header().Set("X-QR-Size", strconv.Itoa(qrResponse.Size))
	w.Header().Set("X-QR-URL", qrResponse.URL)
	
	w.WriteHeader(http.StatusOK)
	w.Write(qrResponse.Data)
}

// GenerateQRCodeForURL handles QR code generation for a specific short URL
func (h *QRHandler) GenerateQRCodeForURL(w http.ResponseWriter, r *http.Request) {
	// Get short code from URL parameter
	shortCode := chi.URLParam(r, "shortCode")
	if shortCode == "" {
		h.writeErrorResponse(w, "Short code is required", http.StatusBadRequest)
		return
	}

	// Parse QR code options from query parameters
	options := h.parseQRCodeOptions(r)

	// Generate QR code for URL
	qrResponse, err := h.qrService.GenerateQRCodeForURL(r.Context(), shortCode, options)
	if err != nil {
		switch err {
		case domain.ErrURLNotFound:
			h.writeErrorResponse(w, "Short URL not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return QR code as binary data
	w.Header().Set("Content-Type", qrResponse.MimeType)
	w.Header().Set("Content-Length", strconv.Itoa(len(qrResponse.Data)))
	w.Header().Set("Cache-Control", "public, max-age=3600")
	
	// Add custom headers
	w.Header().Set("X-QR-Format", qrResponse.Format)
	w.Header().Set("X-QR-Size", strconv.Itoa(qrResponse.Size))
	w.Header().Set("X-QR-URL", qrResponse.URL)
	
	w.WriteHeader(http.StatusOK)
	w.Write(qrResponse.Data)
}

// GetQRCodeFormats handles getting supported QR code formats
func (h *QRHandler) GetQRCodeFormats(w http.ResponseWriter, r *http.Request) {
	formats := h.qrService.GetQRCodeFormats(r.Context())
	
	response := map[string]interface{}{
		"formats": formats,
	}
	
	h.writeJSONResponse(w, response, http.StatusOK)
}

// GetQRCodeSizes handles getting supported QR code sizes
func (h *QRHandler) GetQRCodeSizes(w http.ResponseWriter, r *http.Request) {
	sizes := h.qrService.GetQRCodeSizes(r.Context())
	
	response := map[string]interface{}{
		"sizes": sizes,
	}
	
	h.writeJSONResponse(w, response, http.StatusOK)
}

// ValidateQRCodeOptions handles validating QR code options
func (h *QRHandler) ValidateQRCodeOptions(w http.ResponseWriter, r *http.Request) {
	var options domain.QRCodeOptions
	if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate options
	if err := h.qrService.ValidateQRCodeOptions(r.Context(), options); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"valid":   true,
		"message": "QR code options are valid",
	}
	
	h.writeJSONResponse(w, response, http.StatusOK)
}

// GetQRCodePreview handles getting QR code preview information without generating the actual image
func (h *QRHandler) GetQRCodePreview(w http.ResponseWriter, r *http.Request) {
	var req domain.QRCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set user ID if authenticated
	if userID := middleware.GetUserIDFromContext(r.Context()); userID != 0 {
		req.UserID = userID
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set defaults if not provided
	if req.Size == 0 {
		req.Size = 256
	}
	if req.Format == "" {
		req.Format = "png"
	}
	if req.ErrorCorrection == "" {
		req.ErrorCorrection = "M"
	}

	// Create preview response
	preview := map[string]interface{}{
		"url":               req.URL,
		"short_code":        req.ShortCode,
		"size":              req.Size,
		"format":            req.Format,
		"foreground_color":  req.ForegroundColor,
		"background_color":  req.BackgroundColor,
		"error_correction":  req.ErrorCorrection,
		"border":            req.Border,
		"estimated_size_kb": h.estimateQRCodeSize(req.Size, req.Format),
	}

	h.writeJSONResponse(w, preview, http.StatusOK)
}

// Helper methods

func (h *QRHandler) parseQRCodeOptions(r *http.Request) domain.QRCodeOptions {
	options := domain.QRCodeOptions{}

	// Parse size
	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil {
			options.Size = size
		}
	}

	// Parse format
	if format := r.URL.Query().Get("format"); format != "" {
		options.Format = format
	}

	// Parse colors
	if fgColor := r.URL.Query().Get("fg_color"); fgColor != "" {
		options.ForegroundColor = fgColor
	}
	if bgColor := r.URL.Query().Get("bg_color"); bgColor != "" {
		options.BackgroundColor = bgColor
	}

	// Parse error correction level
	if ecLevel := r.URL.Query().Get("error_correction"); ecLevel != "" {
		options.ErrorCorrection = ecLevel
	}

	// Parse border
	if borderStr := r.URL.Query().Get("border"); borderStr != "" {
		if border, err := strconv.Atoi(borderStr); err == nil {
			options.Border = border
		}
	}

	return options
}

func (h *QRHandler) estimateQRCodeSize(size int, format string) int {
	// Simple estimation based on format and size
	// In a real implementation, you'd have more sophisticated size calculation
	baseSize := (size * size) / 8 // rough pixel to byte conversion

	switch format {
	case "png":
		return baseSize / 4 // PNG compression
	case "jpeg", "jpg":
		return baseSize / 8 // JPEG compression
	case "svg":
		return baseSize / 2 // SVG is text-based
	case "pdf":
		return baseSize + 1024 // PDF overhead
	default:
		return baseSize
	}
}

func (h *QRHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to encode response"}`))
	}
}

func (h *QRHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Write([]byte(`{"error": "Internal server error"}`))
	}
}