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