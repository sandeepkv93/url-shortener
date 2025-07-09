package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type geolocationService struct {
	cacheRepo   ports.CacheService
	configRepo  ports.ConfigService
	httpClient  *http.Client
	apiKey      string
	provider    string
}

// IP API response structure
type ipAPIResponse struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	ZIP         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
	Query       string  `json:"query"`
}

// Country code mapping
var countryCodeMap = map[string]string{
	"United States":      "US",
	"Canada":            "CA",
	"United Kingdom":    "GB",
	"Germany":           "DE",
	"France":            "FR",
	"Italy":             "IT",
	"Spain":             "ES",
	"Netherlands":       "NL",
	"Belgium":           "BE",
	"Switzerland":       "CH",
	"Austria":           "AT",
	"Poland":            "PL",
	"Czech Republic":    "CZ",
	"Hungary":           "HU",
	"Slovakia":          "SK",
	"Slovenia":          "SI",
	"Croatia":           "HR",
	"Serbia":            "RS",
	"Bosnia and Herzegovina": "BA",
	"Montenegro":        "ME",
	"Albania":           "AL",
	"Macedonia":         "MK",
	"Bulgaria":          "BG",
	"Romania":           "RO",
	"Greece":            "GR",
	"Turkey":            "TR",
	"Cyprus":            "CY",
	"Malta":             "MT",
	"Portugal":          "PT",
	"Ireland":           "IE",
	"Denmark":           "DK",
	"Sweden":            "SE",
	"Norway":            "NO",
	"Finland":           "FI",
	"Iceland":           "IS",
	"Estonia":           "EE",
	"Latvia":            "LV",
	"Lithuania":         "LT",
	"Russia":            "RU",
	"Ukraine":           "UA",
	"Belarus":           "BY",
	"Moldova":           "MD",
	"Georgia":           "GE",
	"Armenia":           "AM",
	"Azerbaijan":        "AZ",
	"Kazakhstan":        "KZ",
	"Kyrgyzstan":        "KG",
	"Tajikistan":        "TJ",
	"Turkmenistan":      "TM",
	"Uzbekistan":        "UZ",
	"Afghanistan":       "AF",
	"Pakistan":          "PK",
	"India":             "IN",
	"Bangladesh":        "BD",
	"Sri Lanka":         "LK",
	"Myanmar":           "MM",
	"Thailand":          "TH",
	"Laos":              "LA",
	"Vietnam":           "VN",
	"Cambodia":          "KH",
	"Malaysia":          "MY",
	"Singapore":         "SG",
	"Indonesia":         "ID",
	"Philippines":       "PH",
	"Brunei":            "BN",
	"China":             "CN",
	"Taiwan":            "TW",
	"Hong Kong":         "HK",
	"Macau":             "MO",
	"Mongolia":          "MN",
	"North Korea":       "KP",
	"South Korea":       "KR",
	"Japan":             "JP",
	"Australia":         "AU",
	"New Zealand":       "NZ",
	"Papua New Guinea":  "PG",
	"Fiji":              "FJ",
	"New Caledonia":     "NC",
	"Vanuatu":           "VU",
	"Solomon Islands":   "SB",
	"Samoa":             "WS",
	"Tonga":             "TO",
	"Tuvalu":            "TV",
	"Kiribati":          "KI",
	"Nauru":             "NR",
	"Palau":             "PW",
	"Marshall Islands":  "MH",
	"Micronesia":        "FM",
	"Egypt":             "EG",
	"Libya":             "LY",
	"Sudan":             "SD",
	"Algeria":           "DZ",
	"Morocco":           "MA",
	"Tunisia":           "TN",
	"Ethiopia":          "ET",
	"Kenya":             "KE",
	"Uganda":            "UG",
	"Tanzania":          "TZ",
	"Rwanda":            "RW",
	"Burundi":           "BI",
	"South Africa":      "ZA",
	"Namibia":           "NA",
	"Botswana":          "BW",
	"Zimbabwe":          "ZW",
	"Zambia":            "ZM",
	"Malawi":            "MW",
	"Mozambique":        "MZ",
	"Madagascar":        "MG",
	"Mauritius":         "MU",
	"Seychelles":        "SC",
	"Comoros":           "KM",
	"Mayotte":           "YT",
	"Reunion":           "RE",
	"Nigeria":           "NG",
	"Ghana":             "GH",
	"Ivory Coast":       "CI",
	"Senegal":           "SN",
	"Mali":              "ML",
	"Burkina Faso":      "BF",
	"Niger":             "NE",
	"Chad":              "TD",
	"Central African Republic": "CF",
	"Cameroon":          "CM",
	"Equatorial Guinea": "GQ",
	"Gabon":             "GA",
	"Republic of the Congo": "CG",
	"Democratic Republic of the Congo": "CD",
	"Angola":            "AO",
	"Zambia":            "ZM",
	"Brazil":            "BR",
	"Argentina":         "AR",
	"Chile":             "CL",
	"Peru":              "PE",
	"Bolivia":           "BO",
	"Paraguay":          "PY",
	"Uruguay":           "UY",
	"Ecuador":           "EC",
	"Colombia":          "CO",
	"Venezuela":         "VE",
	"Guyana":            "GY",
	"Suriname":          "SR",
	"French Guiana":     "GF",
	"Mexico":            "MX",
	"Guatemala":         "GT",
	"Belize":            "BZ",
	"El Salvador":       "SV",
	"Honduras":          "HN",
	"Nicaragua":         "NI",
	"Costa Rica":        "CR",
	"Panama":            "PA",
	"Cuba":              "CU",
	"Jamaica":           "JM",
	"Haiti":             "HT",
	"Dominican Republic": "DO",
	"Puerto Rico":       "PR",
	"Trinidad and Tobago": "TT",
	"Barbados":          "BB",
	"Saint Lucia":       "LC",
	"Grenada":           "GD",
	"Saint Vincent and the Grenadines": "VC",
	"Antigua and Barbuda": "AG",
	"Dominica":          "DM",
	"Saint Kitts and Nevis": "KN",
	"Bahamas":           "BS",
}

func NewGeolocationService(
	cacheRepo ports.CacheService,
	configRepo ports.ConfigService,
	apiKey string,
) ports.GeolocationService {
	return &geolocationService{
		cacheRepo:  cacheRepo,
		configRepo: configRepo,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		apiKey:   apiKey,
		provider: "ip-api",
	}
}

func (s *geolocationService) GetLocationFromIP(ctx context.Context, ipAddress string) (*domain.GeoLocation, error) {
	// Validate IP address
	if ipAddress == "" || ipAddress == "127.0.0.1" || ipAddress == "::1" {
		return &domain.GeoLocation{
			IPAddress: ipAddress,
			Country:   "Unknown",
			Region:    "Unknown",
			City:      "Unknown",
			Latitude:  0.0,
			Longitude: 0.0,
			Timezone:  "UTC",
			ISP:       "Unknown",
		}, nil
	}

	// Check cache first
	cacheKey := fmt.Sprintf("geo:%s", ipAddress)
	if cachedLocation, err := s.getCachedLocation(ctx, cacheKey); err == nil {
		return cachedLocation, nil
	}

	// Get location from external API
	location, err := s.getLocationFromAPI(ctx, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get location from API: %w", err)
	}

	// Cache the result for 24 hours
	if err := s.cacheLocation(ctx, cacheKey, location, 24*time.Hour); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to cache location for IP %s: %v\n", ipAddress, err)
	}

	return location, nil
}

func (s *geolocationService) GetLocationsBatch(ctx context.Context, ipAddresses []string) (map[string]*domain.GeoLocation, error) {
	locations := make(map[string]*domain.GeoLocation)
	
	for _, ip := range ipAddresses {
		location, err := s.GetLocationFromIP(ctx, ip)
		if err != nil {
			// Continue with other IPs if one fails
			locations[ip] = &domain.GeoLocation{
				IPAddress: ip,
				Country:   "Unknown",
				Region:    "Unknown",
				City:      "Unknown",
				Latitude:  0.0,
				Longitude: 0.0,
				Timezone:  "UTC",
				ISP:       "Unknown",
			}
			continue
		}
		locations[ip] = location
	}
	
	return locations, nil
}

func (s *geolocationService) ValidateLocation(ctx context.Context, location *domain.GeoLocation) error {
	if location == nil {
		return fmt.Errorf("location cannot be nil")
	}
	
	if location.IPAddress == "" {
		return fmt.Errorf("IP address is required")
	}
	
	// Validate latitude/longitude ranges
	if location.Latitude < -90 || location.Latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}
	
	if location.Longitude < -180 || location.Longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}
	
	return nil
}

func (s *geolocationService) GetCountryCode(ctx context.Context, countryName string) (string, error) {
	if countryName == "" {
		return "", fmt.Errorf("country name cannot be empty")
	}
	
	// Check direct mapping
	if code, exists := countryCodeMap[countryName]; exists {
		return code, nil
	}
	
	// Try case-insensitive lookup
	for name, code := range countryCodeMap {
		if strings.EqualFold(name, countryName) {
			return code, nil
		}
	}
	
	return "", fmt.Errorf("country code not found for: %s", countryName)
}

func (s *geolocationService) getLocationFromAPI(ctx context.Context, ipAddress string) (*domain.GeoLocation, error) {
	// Use ip-api.com as it's free and has good coverage
	url := fmt.Sprintf("http://ip-api.com/json/%s", ipAddress)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var apiResp ipAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	
	if apiResp.Status == "fail" {
		return nil, fmt.Errorf("API returned failure status")
	}
	
	location := &domain.GeoLocation{
		IPAddress:   ipAddress,
		Country:     apiResp.Country,
		CountryCode: apiResp.CountryCode,
		Region:      apiResp.RegionName,
		City:        apiResp.City,
		PostalCode:  apiResp.ZIP,
		Latitude:    apiResp.Lat,
		Longitude:   apiResp.Lon,
		Timezone:    apiResp.Timezone,
		ISP:         apiResp.ISP,
		Organization: apiResp.Org,
		ASN:         apiResp.AS,
		Accuracy:    "city", // ip-api provides city-level accuracy
	}
	
	return location, nil
}

func (s *geolocationService) getCachedLocation(ctx context.Context, cacheKey string) (*domain.GeoLocation, error) {
	data, err := s.cacheRepo.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	
	var location domain.GeoLocation
	if err := json.Unmarshal([]byte(data), &location); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached location: %w", err)
	}
	
	return &location, nil
}

func (s *geolocationService) cacheLocation(ctx context.Context, cacheKey string, location *domain.GeoLocation, ttl time.Duration) error {
	data, err := json.Marshal(location)
	if err != nil {
		return fmt.Errorf("failed to marshal location: %w", err)
	}
	
	return s.cacheRepo.Set(ctx, cacheKey, string(data), ttl)
}

// Helper function to check if IP is private
func (s *geolocationService) isPrivateIP(ipAddress string) bool {
	// This is a simplified check - in production you'd use proper IP parsing
	return strings.HasPrefix(ipAddress, "192.168.") ||
		strings.HasPrefix(ipAddress, "10.") ||
		strings.HasPrefix(ipAddress, "172.16.") ||
		strings.HasPrefix(ipAddress, "127.") ||
		ipAddress == "::1"
}

// Helper function to extract country from user agent (fallback)
func (s *geolocationService) extractCountryFromUserAgent(userAgent string) string {
	// This is a very basic implementation
	// In production, you'd use a proper user agent parser
	userAgent = strings.ToLower(userAgent)
	
	if strings.Contains(userAgent, "en-us") {
		return "US"
	}
	if strings.Contains(userAgent, "en-gb") {
		return "GB"
	}
	if strings.Contains(userAgent, "en-ca") {
		return "CA"
	}
	if strings.Contains(userAgent, "en-au") {
		return "AU"
	}
	if strings.Contains(userAgent, "de-de") {
		return "DE"
	}
	if strings.Contains(userAgent, "fr-fr") {
		return "FR"
	}
	if strings.Contains(userAgent, "es-es") {
		return "ES"
	}
	if strings.Contains(userAgent, "it-it") {
		return "IT"
	}
	if strings.Contains(userAgent, "pt-br") {
		return "BR"
	}
	if strings.Contains(userAgent, "ru-ru") {
		return "RU"
	}
	if strings.Contains(userAgent, "zh-cn") {
		return "CN"
	}
	if strings.Contains(userAgent, "ja-jp") {
		return "JP"
	}
	if strings.Contains(userAgent, "ko-kr") {
		return "KR"
	}
	
	return "Unknown"
}

// GetLocationWithFallback tries multiple methods to get location
func (s *geolocationService) GetLocationWithFallback(ctx context.Context, ipAddress, userAgent string) (*domain.GeoLocation, error) {
	// Try primary method
	location, err := s.GetLocationFromIP(ctx, ipAddress)
	if err == nil && location.Country != "Unknown" {
		return location, nil
	}
	
	// Fallback to user agent parsing
	country := s.extractCountryFromUserAgent(userAgent)
	countryCode, _ := s.GetCountryCode(ctx, country)
	
	return &domain.GeoLocation{
		IPAddress:   ipAddress,
		Country:     country,
		CountryCode: countryCode,
		Region:      "Unknown",
		City:        "Unknown",
		Latitude:    0.0,
		Longitude:   0.0,
		Timezone:    "UTC",
		ISP:         "Unknown",
		Accuracy:    "country",
	}, nil
}