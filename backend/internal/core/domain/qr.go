package domain

import (
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
	ID          uint      `json:"id"`
	ShortURLID  uint      `json:"short_url_id"`
	Format      string    `json:"format"`
	Size        int       `json:"size"`
	URL         string    `json:"url"`          // Download URL
	DataURL     string    `json:"data_url"`     // Base64 data URL
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
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