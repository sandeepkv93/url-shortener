package domain

import (
	"time"

	"gorm.io/gorm"
)

type Click struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	ShortURLID  uint           `json:"short_url_id" gorm:"not null;index"`
	IPAddress   string         `json:"ip_address" gorm:"type:inet"`
	UserAgent   string         `json:"user_agent" gorm:"type:text"`
	Referer     string         `json:"referer" gorm:"type:text"`
	Country     string         `json:"country" gorm:"size:2"`
	Region      string         `json:"region" gorm:"size:100"`
	City        string         `json:"city" gorm:"size:100"`
	Device      string         `json:"device" gorm:"size:50"`
	Browser     string         `json:"browser" gorm:"size:50"`
	OS          string         `json:"os" gorm:"size:50"`
	ClickedAt   time.Time      `json:"clicked_at" gorm:"index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	ShortURL ShortURL `json:"short_url,omitempty" gorm:"foreignKey:ShortURLID"`
}

type ClickStats struct {
	TotalClicks  int64                    `json:"total_clicks"`
	UniqueClicks int64                    `json:"unique_clicks"`
	ClicksByDate map[string]int64         `json:"clicks_by_date"`
	ClicksByTime map[int]int64            `json:"clicks_by_time"`
	TopCountries []CountryStat           `json:"top_countries"`
	TopDevices   []DeviceStat            `json:"top_devices"`
	TopBrowsers  []BrowserStat           `json:"top_browsers"`
	TopReferers  []RefererStat           `json:"top_referers"`
	RecentClicks []RecentClickStat       `json:"recent_clicks"`
}

type CountryStat struct {
	Country string `json:"country"`
	Count   int64  `json:"count"`
}

type DeviceStat struct {
	Device string `json:"device"`
	Count  int64  `json:"count"`
}

type BrowserStat struct {
	Browser string `json:"browser"`
	Count   int64  `json:"count"`
}

type RefererStat struct {
	Referer string `json:"referer"`
	Count   int64  `json:"count"`
}

type RecentClickStat struct {
	Country   string    `json:"country"`
	City      string    `json:"city"`
	Device    string    `json:"device"`
	Browser   string    `json:"browser"`
	Referer   string    `json:"referer"`
	ClickedAt time.Time `json:"clicked_at"`
}

type GeoStats struct {
	CountryStats map[string]int64 `json:"country_stats"`
	RegionStats  map[string]int64 `json:"region_stats"`
	CityStats    map[string]int64 `json:"city_stats"`
}

type TimelineStats struct {
	Period string           `json:"period"`
	Data   map[string]int64 `json:"data"`
}

type GlobalStats struct {
	TotalURLs       int64 `json:"total_urls"`
	TotalClicks     int64 `json:"total_clicks"`
	TotalUsers      int64 `json:"total_users"`
	ActiveURLs      int64 `json:"active_urls"`
	ClicksToday     int64 `json:"clicks_today"`
	URLsCreatedToday int64 `json:"urls_created_today"`
	NewUsersToday   int64 `json:"new_users_today"`
}

type UserDashboard struct {
	UserID          uint              `json:"user_id"`
	TotalURLs       int64             `json:"total_urls"`
	TotalClicks     int64             `json:"total_clicks"`
	ActiveURLs      int64             `json:"active_urls"`
	ClicksToday     int64             `json:"clicks_today"`
	ClicksThisWeek  int64             `json:"clicks_this_week"`
	ClicksThisMonth int64             `json:"clicks_this_month"`
	TopURLs         []TopURLStat      `json:"top_urls"`
	RecentActivity  []RecentActivity  `json:"recent_activity"`
	ClickTrend      map[string]int64  `json:"click_trend"`
}

type GlobalDashboard struct {
	GlobalStats     GlobalStats      `json:"global_stats"`
	TopCountries    []CountryStat    `json:"top_countries"`
	TopDevices      []DeviceStat     `json:"top_devices"`
	TopBrowsers     []BrowserStat    `json:"top_browsers"`
	RecentURLs      []RecentURLStat  `json:"recent_urls"`
	ActivityTrend   map[string]int64 `json:"activity_trend"`
}

type RecentActivity struct {
	Type        string    `json:"type"` // url_created, click_received
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

type RecentURLStat struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	ClickCount  int64     `json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type AnalyticsReport struct {
	ShortURLID   uint               `json:"short_url_id"`
	Period       string             `json:"period"`
	GeneratedAt  time.Time          `json:"generated_at"`
	Summary      AnalyticsSummary   `json:"summary"`
	Detailed     DetailedAnalytics  `json:"detailed"`
}

type AnalyticsSummary struct {
	TotalClicks     int64   `json:"total_clicks"`
	UniqueClicks    int64   `json:"unique_clicks"`
	ClickRate       float64 `json:"click_rate"`
	PeakHour        int     `json:"peak_hour"`
	TopCountry      string  `json:"top_country"`
	TopDevice       string  `json:"top_device"`
	TopBrowser      string  `json:"top_browser"`
}

type DetailedAnalytics struct {
	ClicksByDate     map[string]int64      `json:"clicks_by_date"`
	ClicksByHour     map[int]int64         `json:"clicks_by_hour"`
	CountryStats     []CountryStat         `json:"country_stats"`
	DeviceStats      []DeviceStat          `json:"device_stats"`
	BrowserStats     []BrowserStat         `json:"browser_stats"`
	RefererStats     []RefererStat         `json:"referer_stats"`
}

type RealTimeStats struct {
	ActiveUsers     int64            `json:"active_users"`
	ClicksLastHour  int64            `json:"clicks_last_hour"`
	ClicksLastDay   int64            `json:"clicks_last_day"`
	LiveClicks      []LiveClickStat  `json:"live_clicks"`
	OnlineCountries []string         `json:"online_countries"`
}

type LiveClickStat struct {
	Country   string    `json:"country"`
	City      string    `json:"city"`
	Device    string    `json:"device"`
	Timestamp time.Time `json:"timestamp"`
}

type RealTimeUpdate struct {
	Type      string    `json:"type"` // click, view, etc.
	Data      interface{} `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}