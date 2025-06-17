package domain

import (
	"time"

	"gorm.io/gorm"
)

type ShortURL struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	ShortCode   string         `json:"short_code" gorm:"uniqueIndex;not null"`
	OriginalURL string         `json:"original_url" gorm:"type:text;not null"`
	UserID      uint           `json:"user_id" gorm:"index;not null"`
	Title       string         `json:"title" gorm:"size:255"`
	Description string         `json:"description" gorm:"type:text"`
	Password    *string        `json:"-" gorm:"size:255"`
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

// Request/Response models for services
type ShortenURLRequest struct {
	OriginalURL string     `json:"original_url" validate:"required,url"`
	UserID      uint       `json:"user_id" validate:"required"`
	Title       string     `json:"title" validate:"omitempty,max=255"`
	Description string     `json:"description" validate:"omitempty,max=1000"`
	CustomAlias string     `json:"custom_alias" validate:"omitempty,alphanum,max=50"`
	Password    string     `json:"password" validate:"omitempty,min=4"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

type UpdateURLRequest struct {
	Title       *string    `json:"title,omitempty" validate:"omitempty,max=255"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=1000"`
	IsActive    *bool      `json:"is_active,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type ClickData struct {
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
	Country   string `json:"country"`
	Region    string `json:"region"`
	City      string `json:"city"`
	Device    string `json:"device"`
	Browser   string `json:"browser"`
	OS        string `json:"os"`
}

type URLStats struct {
	ShortURL      *ShortURL      `json:"short_url"`
	ClickStats    *ClickStats    `json:"click_stats"`
	GeoStats      *GeoStats      `json:"geo_stats"`
	TimelineStats *TimelineStats `json:"timeline_stats"`
}

type DashboardStats struct {
	TotalURLs       int64                `json:"total_urls"`
	ActiveURLs      int64                `json:"active_urls"`
	TotalClicks     int64                `json:"total_clicks"`
	ClickGrowthRate float64              `json:"click_growth_rate"`
	URLGrowthRate   float64              `json:"url_growth_rate"`
	ClicksByDate    map[string]int64     `json:"clicks_by_date"`
	TopURLs         []TopURLStat         `json:"top_urls"`
	RecentActivity  []ActivityItem       `json:"recent_activity"`
}

type URLAnalytics struct {
	ShortURLID   uint                  `json:"short_url_id"`
	ShortCode    string                `json:"short_code"`
	TotalClicks  int64                 `json:"total_clicks"`
	UniqueClicks int64                 `json:"unique_clicks"`
	ClicksByDate map[string]int64      `json:"clicks_by_date"`
	ClicksByTime map[int]int64         `json:"clicks_by_time"`
	CountryStats map[string]int64      `json:"country_stats"`
	RegionStats  map[string]int64      `json:"region_stats"`
	CityStats    map[string]int64      `json:"city_stats"`
	TopDevices   []DeviceStat          `json:"top_devices"`
	TopBrowsers  []BrowserStat         `json:"top_browsers"`
	TopReferers  []RefererStat         `json:"top_referers"`
	RecentClicks []RecentClickStat     `json:"recent_clicks"`
}

type URLPerformance struct {
	ShortURL     *ShortURL `json:"short_url"`
	TotalClicks  int64     `json:"total_clicks"`
	UniqueClicks int64     `json:"unique_clicks"`
	ClickRate    float64   `json:"click_rate"`
}

type DeviceStats struct {
	TopDevices  []DeviceStat  `json:"top_devices"`
	TopBrowsers []BrowserStat `json:"top_browsers"`
}

type ActivityItem struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

type DateRange struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// Validation methods
func (r *ShortenURLRequest) Validate() error {
	if r.OriginalURL == "" {
		return ErrInvalidURL
	}
	if r.UserID == 0 {
		return ErrInvalidRequest
	}
	return nil
}

func (r *UpdateURLRequest) Validate() error {
	return nil
}