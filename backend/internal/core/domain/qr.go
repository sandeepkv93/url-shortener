package domain

import (
	"image/color"
	"time"
)

type QROptions struct {
	Size        int    `json:"size" validate:"min=100,max=1000"`         // QR code size in pixels
	Format      string `json:"format" validate:"oneof=png svg pdf"`     // Output format
	ErrorLevel  string `json:"error_level" validate:"oneof=L M Q H"`    // Error correction level
	BorderSize  int    `json:"border_size" validate:"min=0,max=10"`     // Border size in modules
	Foreground  string `json:"foreground,omitempty"`                    // Foreground color (hex)
	Background  string `json:"background,omitempty"`                    // Background color (hex)
	Logo        string `json:"logo,omitempty"`                          // Logo URL for branded QR codes
	LogoSize    int    `json:"logo_size" validate:"min=10,max=30"`      // Logo size as percentage
}

type CustomQRRequest struct {
	ShortURL    string     `json:"short_url" validate:"required,url"`
	Options     *QROptions `json:"options"`
	Customization *QRCustomization `json:"customization,omitempty"`
}

type QRCustomization struct {
	Style       string `json:"style" validate:"oneof=square circle rounded"`
	Pattern     string `json:"pattern" validate:"oneof=solid gradient"`
	EyeStyle    string `json:"eye_style" validate:"oneof=square circle rounded"`
	DataPattern string `json:"data_pattern" validate:"oneof=square circle diamond"`
}

type QRCodeResponse struct {
	Data     []byte `json:"-"`                // Raw QR code data
	Format   string `json:"format"`           // Image format (png, jpeg, svg, pdf)
	Size     int    `json:"size"`             // Size in pixels
	URL      string `json:"url"`              // The URL encoded in the QR code
	MimeType string `json:"mime_type"`        // MIME type for HTTP responses
}

type QRCodeListResponse struct {
	QRCodes    []QRCodeResponse `json:"qr_codes"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

type QRCodeHistory struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ShortURLID   uint      `json:"short_url_id"`
	Format       string    `json:"format"`
	Size         int       `json:"size"`
	DownloadCount int64    `json:"download_count"`
	CreatedAt    time.Time `json:"created_at"`
	LastDownloaded *time.Time `json:"last_downloaded"`
}

type QRBatch struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	URLs        []string  `json:"urls"`
	Options     *QROptions `json:"options"`
	Status      string    `json:"status"` // pending, processing, completed, failed
	Progress    int       `json:"progress"` // 0-100
	ResultURL   string    `json:"result_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Error       string    `json:"error,omitempty"`
}

// Additional models for QR service
type QRCodeRequest struct {
	URL              string `json:"url" validate:"required"`
	ShortCode        string `json:"short_code,omitempty"`
	Size             int    `json:"size,omitempty"`
	Format           string `json:"format,omitempty"`
	ForegroundColor  string `json:"foreground_color,omitempty"`
	BackgroundColor  string `json:"background_color,omitempty"`
	ErrorCorrection  string `json:"error_correction,omitempty"`
	Border           int    `json:"border,omitempty"`
	UserID           uint   `json:"user_id,omitempty"`
}

type QRCodeOptions struct {
	Size             int    `json:"size,omitempty"`
	Format           string `json:"format,omitempty"`
	ForegroundColor  string `json:"foreground_color,omitempty"`
	BackgroundColor  string `json:"background_color,omitempty"`
	ErrorCorrection  string `json:"error_correction,omitempty"`
	Border           int    `json:"border,omitempty"`
}

type QRGenerationOptions struct {
	Size             int        `json:"size"`
	Format           string     `json:"format"`
	ForegroundColor  color.RGBA `json:"foreground_color"`
	BackgroundColor  color.RGBA `json:"background_color"`
	ErrorCorrection  int        `json:"error_correction"`
	Border           int        `json:"border"`
}

// Validation methods
func (r *QRCodeRequest) Validate() error {
	if r.URL == "" && r.ShortCode == "" {
		return ErrInvalidRequest
	}
	return nil
}