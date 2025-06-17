package services

import (
	"bytes"
	"context"
	"fmt"
	"image/color"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type qrService struct {
	urlRepo    ports.URLRepository
	configRepo ports.ConfigService
	qrProvider ports.QRCodeProvider
}

func NewQRService(
	urlRepo ports.URLRepository,
	configRepo ports.ConfigService,
	qrProvider ports.QRCodeProvider,
) ports.QRService {
	return &qrService{
		urlRepo:    urlRepo,
		configRepo: configRepo,
		qrProvider: qrProvider,
	}
}

func (s *qrService) GenerateQRCode(ctx context.Context, req domain.QRCodeRequest) (*domain.QRCodeResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// If short code is provided, verify it exists and get the full URL
	var targetURL string
	if req.ShortCode != "" {
		shortURL, err := s.urlRepo.GetByShortCode(ctx, req.ShortCode)
		if err != nil {
			return nil, fmt.Errorf("failed to get short URL: %w", err)
		}

		// Check ownership if user ID is provided
		if req.UserID != 0 && shortURL.UserID != req.UserID {
			return nil, domain.ErrUnauthorized
		}

		// Build the full short URL
		baseURL := s.configRepo.GetBaseURL()
		targetURL = fmt.Sprintf("%s/%s", baseURL, req.ShortCode)
	} else if req.URL != "" {
		targetURL = req.URL
	} else {
		return nil, domain.ErrInvalidRequest
	}

	// Set default values
	if req.Size == 0 {
		req.Size = 256
	}
	if req.Format == "" {
		req.Format = "png"
	}

	// Generate QR code
	qrData, err := s.generateQRCodeData(targetURL, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	response := &domain.QRCodeResponse{
		Data:     qrData,
		Format:   req.Format,
		Size:     req.Size,
		URL:      targetURL,
		MimeType: s.getMimeType(req.Format),
	}

	return response, nil
}

func (s *qrService) GenerateQRCodeForURL(ctx context.Context, shortCode string, options domain.QRCodeOptions) (*domain.QRCodeResponse, error) {
	// Get the short URL
	shortURL, err := s.urlRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get short URL: %w", err)
	}

	// Build the full short URL
	baseURL := s.configRepo.GetBaseURL()
	targetURL := fmt.Sprintf("%s/%s", shortCode, baseURL)

	// Create QR code request
	req := domain.QRCodeRequest{
		URL:              targetURL,
		ShortCode:        shortCode,
		Size:             options.Size,
		Format:           options.Format,
		ForegroundColor:  options.ForegroundColor,
		BackgroundColor:  options.BackgroundColor,
		ErrorCorrection:  options.ErrorCorrection,
		Border:           options.Border,
		UserID:           shortURL.UserID,
	}

	return s.GenerateQRCode(ctx, req)
}

func (s *qrService) GetQRCodeFormats(ctx context.Context) []string {
	return []string{"png", "jpeg", "svg", "pdf"}
}

func (s *qrService) GetQRCodeSizes(ctx context.Context) []int {
	return []int{128, 256, 512, 1024}
}

func (s *qrService) ValidateQRCodeOptions(ctx context.Context, options domain.QRCodeOptions) error {
	// Validate size
	validSizes := s.GetQRCodeSizes(ctx)
	sizeValid := false
	for _, size := range validSizes {
		if options.Size == size {
			sizeValid = true
			break
		}
	}
	if !sizeValid && options.Size != 0 {
		return fmt.Errorf("invalid size: %d, valid sizes are %v", options.Size, validSizes)
	}

	// Validate format
	validFormats := s.GetQRCodeFormats(ctx)
	formatValid := false
	for _, format := range validFormats {
		if options.Format == format {
			formatValid = true
			break
		}
	}
	if !formatValid && options.Format != "" {
		return fmt.Errorf("invalid format: %s, valid formats are %v", options.Format, validFormats)
	}

	// Validate error correction level
	if options.ErrorCorrection != "" {
		validLevels := []string{"L", "M", "Q", "H"}
		levelValid := false
		for _, level := range validLevels {
			if options.ErrorCorrection == level {
				levelValid = true
				break
			}
		}
		if !levelValid {
			return fmt.Errorf("invalid error correction level: %s, valid levels are %v", options.ErrorCorrection, validLevels)
		}
	}

	// Validate colors (basic hex color validation)
	if options.ForegroundColor != "" && !s.isValidHexColor(options.ForegroundColor) {
		return fmt.Errorf("invalid foreground color: %s", options.ForegroundColor)
	}
	if options.BackgroundColor != "" && !s.isValidHexColor(options.BackgroundColor) {
		return fmt.Errorf("invalid background color: %s", options.BackgroundColor)
	}

	return nil
}

func (s *qrService) generateQRCodeData(url string, req domain.QRCodeRequest) ([]byte, error) {
	// Use the QR code provider to generate the actual QR code
	options := domain.QRGenerationOptions{
		Size:             req.Size,
		Format:           req.Format,
		ForegroundColor:  s.parseColor(req.ForegroundColor, color.RGBA{0, 0, 0, 255}),
		BackgroundColor:  s.parseColor(req.BackgroundColor, color.RGBA{255, 255, 255, 255}),
		ErrorCorrection:  s.parseErrorCorrectionLevel(req.ErrorCorrection),
		Border:           req.Border,
	}

	return s.qrProvider.GenerateQRCode(url, options)
}

func (s *qrService) parseColor(hexColor string, defaultColor color.RGBA) color.RGBA {
	if hexColor == "" {
		return defaultColor
	}

	// Simple hex color parsing (without # prefix validation for brevity)
	if len(hexColor) == 7 && hexColor[0] == '#' {
		var r, g, b uint8
		fmt.Sscanf(hexColor[1:], "%02x%02x%02x", &r, &g, &b)
		return color.RGBA{r, g, b, 255}
	}

	return defaultColor
}

func (s *qrService) parseErrorCorrectionLevel(level string) int {
	switch level {
	case "L":
		return 0 // Low
	case "M":
		return 1 // Medium
	case "Q":
		return 2 // Quartile
	case "H":
		return 3 // High
	default:
		return 1 // Medium as default
	}
}

func (s *qrService) isValidHexColor(hex string) bool {
	if len(hex) != 7 || hex[0] != '#' {
		return false
	}

	for i := 1; i < 7; i++ {
		c := hex[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}

	return true
}

func (s *qrService) getMimeType(format string) string {
	switch format {
	case "png":
		return "image/png"
	case "jpeg", "jpg":
		return "image/jpeg"
	case "svg":
		return "image/svg+xml"
	case "pdf":
		return "application/pdf"
	default:
		return "image/png"
	}
}

// Simple QR code provider implementation for basic functionality
type simpleQRProvider struct{}

func NewSimpleQRProvider() ports.QRCodeProvider {
	return &simpleQRProvider{}
}

func (p *simpleQRProvider) GenerateQRCode(url string, options domain.QRGenerationOptions) ([]byte, error) {
	// This is a placeholder implementation
	// In a real implementation, you would use a proper QR code library like go-qrcode
	
	// Create a simple placeholder image
	img := make([][]bool, options.Size/8)
	for i := range img {
		img[i] = make([]bool, options.Size/8)
		for j := range img[i] {
			// Simple pattern for demonstration
			img[i][j] = (i+j)%2 == 0
		}
	}

	// Convert to PNG bytes (simplified)
	var buf bytes.Buffer
	
	// This is a very basic placeholder - in reality you'd use proper QR code generation
	// and image encoding libraries
	buf.WriteString("QR-CODE-PLACEHOLDER")
	
	return buf.Bytes(), nil
}