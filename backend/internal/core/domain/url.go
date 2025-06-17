package domain

import (
	"time"

	"gorm.io/gorm"
)

type ShortURL struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	ShortCode   string         `json:"short_code" gorm:"uniqueIndex;not null"`
	OriginalURL string         `json:"original_url" gorm:"type:text;not null"`
	UserID      *uint          `json:"user_id" gorm:"index"`
	CustomAlias bool           `json:"custom_alias" gorm:"default:false"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	ClickCount  int64          `json:"click_count" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Clicks []Click `json:"clicks,omitempty" gorm:"foreignKey:ShortURLID"`
}

type CreateShortURLRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
	CustomAlias string `json:"custom_alias,omitempty" validate:"omitempty,alphanum,max=50"`
	ExpiresAt   string `json:"expires_at,omitempty"`
}

type UpdateShortURLRequest struct {
	OriginalURL string `json:"original_url,omitempty" validate:"omitempty,url"`
	IsActive    *bool  `json:"is_active,omitempty"`
	ExpiresAt   string `json:"expires_at,omitempty"`
}

type ShortURLResponse struct {
	ID          uint       `json:"id"`
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	ShortURL    string     `json:"short_url"`
	CustomAlias bool       `json:"custom_alias"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IsActive    bool       `json:"is_active"`
	ClickCount  int64      `json:"click_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	User        *User      `json:"user,omitempty"`
}

type ShortURLListResponse struct {
	URLs       []ShortURLResponse `json:"urls"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

func (s *ShortURL) ToResponse(baseURL string) *ShortURLResponse {
	return &ShortURLResponse{
		ID:          s.ID,
		ShortCode:   s.ShortCode,
		OriginalURL: s.OriginalURL,
		ShortURL:    baseURL + "/" + s.ShortCode,
		CustomAlias: s.CustomAlias,
		ExpiresAt:   s.ExpiresAt,
		IsActive:    s.IsActive,
		ClickCount:  s.ClickCount,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		User:        s.User,
	}
}

type URLFilter struct {
	Status    string `json:"status,omitempty"` // active, expired, inactive
	DateFrom  string `json:"date_from,omitempty"`
	DateTo    string `json:"date_to,omitempty"`
	Search    string `json:"search,omitempty"`
	SortBy    string `json:"sort_by,omitempty"`    // created_at, click_count, expires_at
	SortOrder string `json:"sort_order,omitempty"` // asc, desc
}

type BulkUpdateRequest struct {
	IDs      []uint                 `json:"ids" validate:"required"`
	Updates  UpdateShortURLRequest  `json:"updates"`
}

type RecordClickRequest struct {
	ShortCode string `json:"short_code" validate:"required"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
}

func (s *ShortURL) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*s.ExpiresAt)
}

func (s *ShortURL) IsAccessible() bool {
	return s.IsActive && !s.IsExpired()
}