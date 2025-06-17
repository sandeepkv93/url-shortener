package domain

import (
	"time"
)

type GeoLocation struct {
	IPAddress     string  `json:"ip_address"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	City          string  `json:"city"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Timezone      string  `json:"timezone"`
	ISP           string  `json:"isp"`
	Organization  string  `json:"organization"`
	ASN           string  `json:"asn"`
	Accuracy      string  `json:"accuracy"` // city, region, country
	Source        string  `json:"source"`   // ipapi, maxmind, ipgeolocation
	CachedAt      time.Time `json:"cached_at"`
}

type GeoBounds struct {
	NorthEast GeoPoint `json:"northeast"`
	SouthWest GeoPoint `json:"southwest"`
}

type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CountryInfo struct {
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Capital     string   `json:"capital"`
	Currency    string   `json:"currency"`
	Languages   []string `json:"languages"`
	Timezone    string   `json:"timezone"`
	Population  int64    `json:"population"`
	Area        float64  `json:"area"`
	Region      string   `json:"region"`
	Subregion   string   `json:"subregion"`
}

type GeoAnalytics struct {
	TopCountries    []GeoCountryStat `json:"top_countries"`
	TopRegions      []GeoRegionStat  `json:"top_regions"`
	TopCities       []GeoCityStat    `json:"top_cities"`
	ClicksOverTime  map[string]GeoTimeStats `json:"clicks_over_time"`
	HeatmapData     []GeoHeatmapPoint `json:"heatmap_data"`
	Bounds          GeoBounds        `json:"bounds"`
}

type GeoCountryStat struct {
	Country      string  `json:"country"`
	CountryCode  string  `json:"country_code"`
	ClickCount   int64   `json:"click_count"`
	Percentage   float64 `json:"percentage"`
	UniqueVisitors int64 `json:"unique_visitors"`
}

type GeoRegionStat struct {
	Region       string  `json:"region"`
	Country      string  `json:"country"`
	ClickCount   int64   `json:"click_count"`
	Percentage   float64 `json:"percentage"`
}

type GeoCityStat struct {
	City         string  `json:"city"`
	Region       string  `json:"region"`
	Country      string  `json:"country"`
	ClickCount   int64   `json:"click_count"`
	Percentage   float64 `json:"percentage"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}

type GeoTimeStats struct {
	Date        string `json:"date"`
	ClickCount  int64  `json:"click_count"`
	Countries   int    `json:"countries"`
	TopCountry  string `json:"top_country"`
}

type GeoHeatmapPoint struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	ClickCount int64   `json:"click_count"`
	Weight     float64 `json:"weight"`
}