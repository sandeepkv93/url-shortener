package services

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// SecurityService handles URL security validation and threat detection
type SecurityService interface {
	ValidateURL(ctx context.Context, rawURL string) (*URLSecurityReport, error)
	CheckMaliciousURL(ctx context.Context, parsedURL *url.URL) (*ThreatAssessment, error)
	IsBlacklisted(ctx context.Context, hostname string) (bool, string, error)
	ScanURLContent(ctx context.Context, targetURL string) (*ContentScanResult, error)
	CheckPhishing(ctx context.Context, parsedURL *url.URL) (*PhishingReport, error)
	UpdateBlacklist(ctx context.Context, domains []string) error
}

type securityService struct {
	blacklistedDomains map[string]string // domain -> reason
	phishingPatterns   []*regexp.Regexp
	suspiciousKeywords []string
	httpClient         *http.Client
	logger             *logrus.Logger
}

// URLSecurityReport contains comprehensive security analysis of a URL
type URLSecurityReport struct {
	URL                string             `json:"url"`
	IsSecure           bool               `json:"is_secure"`
	ThreatLevel        ThreatLevel        `json:"threat_level"`
	ThreatAssessment   *ThreatAssessment  `json:"threat_assessment,omitempty"`
	PhishingReport     *PhishingReport    `json:"phishing_report,omitempty"`
	ContentScanResult  *ContentScanResult `json:"content_scan_result,omitempty"`
	Recommendations    []string           `json:"recommendations"`
	ScannedAt          time.Time          `json:"scanned_at"`
}

// ThreatLevel represents the security threat level
type ThreatLevel string

const (
	ThreatLevelSafe     ThreatLevel = "safe"
	ThreatLevelLow      ThreatLevel = "low"
	ThreatLevelMedium   ThreatLevel = "medium"
	ThreatLevelHigh     ThreatLevel = "high"
	ThreatLevelCritical ThreatLevel = "critical"
)

// ThreatAssessment contains detailed threat analysis
type ThreatAssessment struct {
	IsBlacklisted      bool     `json:"is_blacklisted"`
	BlacklistReason    string   `json:"blacklist_reason,omitempty"`
	IsSuspiciousHost   bool     `json:"is_suspicious_host"`
	IsPrivateIP        bool     `json:"is_private_ip"`
	IsShortURL         bool     `json:"is_short_url"`
	HasSuspiciousPath  bool     `json:"has_suspicious_path"`
	RiskFactors        []string `json:"risk_factors"`
	ThreatCategories   []string `json:"threat_categories"`
}

// PhishingReport contains phishing detection results
type PhishingReport struct {
	IsPhishingSuspected bool     `json:"is_phishing_suspected"`
	PhishingIndicators  []string `json:"phishing_indicators"`
	TargetBrand         string   `json:"target_brand,omitempty"`
	SuspiciousElements  []string `json:"suspicious_elements"`
	ConfidenceScore     float64  `json:"confidence_score"`
}

// ContentScanResult contains results from scanning URL content
type ContentScanResult struct {
	IsAccessible       bool     `json:"is_accessible"`
	HTTPStatus         int      `json:"http_status"`
	ContentType        string   `json:"content_type"`
	HasJavaScript      bool     `json:"has_javascript"`
	HasForms           bool     `json:"has_forms"`
	HasExternalLinks   bool     `json:"has_external_links"`
	SuspiciousContent  []string `json:"suspicious_content"`
	SSLInfo            *SSLInfo `json:"ssl_info,omitempty"`
	ResponseTime       int64    `json:"response_time_ms"`
	ContentLength      int64    `json:"content_length"`
}

// SSLInfo contains SSL certificate information
type SSLInfo struct {
	IsSecure     bool      `json:"is_secure"`
	Issuer       string    `json:"issuer"`
	Subject      string    `json:"subject"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsExpired    bool      `json:"is_expired"`
	IsSelfSigned bool      `json:"is_self_signed"`
}

// NewSecurityService creates a new security service instance
func NewSecurityService(logger *logrus.Logger) SecurityService {
	// Create HTTP client with security settings
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
		},
	}

	return &securityService{
		blacklistedDomains: initializeBlacklist(),
		phishingPatterns:   initializePhishingPatterns(),
		suspiciousKeywords: initializeSuspiciousKeywords(),
		httpClient:         client,
		logger:             logger,
	}
}

// ValidateURL performs comprehensive security validation of a URL
func (s *securityService) ValidateURL(ctx context.Context, rawURL string) (*URLSecurityReport, error) {
	report := &URLSecurityReport{
		URL:             rawURL,
		IsSecure:        true,
		ThreatLevel:     ThreatLevelSafe,
		Recommendations: []string{},
		ScannedAt:       time.Now(),
	}

	// Parse URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format: %w", err)
	}

	// Basic URL structure validation
	if err := s.validateURLStructure(parsedURL); err != nil {
		report.IsSecure = false
		report.ThreatLevel = ThreatLevelHigh
		report.Recommendations = append(report.Recommendations, err.Error())
		return report, nil
	}

	// Threat assessment
	threatAssessment, err := s.CheckMaliciousURL(ctx, parsedURL)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to perform threat assessment")
	} else {
		report.ThreatAssessment = threatAssessment
		if len(threatAssessment.RiskFactors) > 0 {
			report.IsSecure = false
			report.ThreatLevel = s.calculateThreatLevel(threatAssessment)
		}
	}

	// Phishing detection
	phishingReport, err := s.CheckPhishing(ctx, parsedURL)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to perform phishing check")
	} else {
		report.PhishingReport = phishingReport
		if phishingReport.IsPhishingSuspected {
			report.IsSecure = false
			if report.ThreatLevel < ThreatLevelHigh {
				report.ThreatLevel = ThreatLevelHigh
			}
		}
	}

	// Content scanning (for HTTP/HTTPS URLs)
	if parsedURL.Scheme == "http" || parsedURL.Scheme == "https" {
		contentScan, err := s.ScanURLContent(ctx, rawURL)
		if err != nil {
			s.logger.WithError(err).Warn("Failed to scan URL content")
		} else {
			report.ContentScanResult = contentScan
			if len(contentScan.SuspiciousContent) > 0 {
				report.IsSecure = false
				if report.ThreatLevel < ThreatLevelMedium {
					report.ThreatLevel = ThreatLevelMedium
				}
			}
		}
	}

	// Generate recommendations
	report.Recommendations = s.generateRecommendations(report)

	return report, nil
}

// CheckMaliciousURL performs malicious URL detection
func (s *securityService) CheckMaliciousURL(ctx context.Context, parsedURL *url.URL) (*ThreatAssessment, error) {
	assessment := &ThreatAssessment{
		RiskFactors:      []string{},
		ThreatCategories: []string{},
	}

	hostname := strings.ToLower(parsedURL.Hostname())

	// Check blacklist
	isBlacklisted, reason, err := s.IsBlacklisted(ctx, hostname)
	if err != nil {
		return nil, err
	}
	assessment.IsBlacklisted = isBlacklisted
	assessment.BlacklistReason = reason

	if isBlacklisted {
		assessment.RiskFactors = append(assessment.RiskFactors, "Domain is blacklisted")
		assessment.ThreatCategories = append(assessment.ThreatCategories, "blacklisted")
	}

	// Check for private/localhost IPs
	if isPrivateOrLocalhost(hostname) {
		assessment.IsPrivateIP = true
		assessment.RiskFactors = append(assessment.RiskFactors, "Private IP or localhost detected")
		assessment.ThreatCategories = append(assessment.ThreatCategories, "private_ip")
	}

	// Check for suspicious hosts
	if s.isSuspiciousHost(hostname) {
		assessment.IsSuspiciousHost = true
		assessment.RiskFactors = append(assessment.RiskFactors, "Suspicious hostname pattern")
		assessment.ThreatCategories = append(assessment.ThreatCategories, "suspicious_host")
	}

	// Check for known short URL services
	if s.isShortURLService(hostname) {
		assessment.IsShortURL = true
		assessment.RiskFactors = append(assessment.RiskFactors, "URL redirection service detected")
		assessment.ThreatCategories = append(assessment.ThreatCategories, "url_shortener")
	}

	// Check for suspicious path patterns
	if s.hasSuspiciousPath(parsedURL.Path) {
		assessment.HasSuspiciousPath = true
		assessment.RiskFactors = append(assessment.RiskFactors, "Suspicious URL path detected")
		assessment.ThreatCategories = append(assessment.ThreatCategories, "suspicious_path")
	}

	return assessment, nil
}

// IsBlacklisted checks if a domain is blacklisted
func (s *securityService) IsBlacklisted(ctx context.Context, hostname string) (bool, string, error) {
	// Check exact match
	if reason, exists := s.blacklistedDomains[hostname]; exists {
		return true, reason, nil
	}

	// Check subdomain patterns
	parts := strings.Split(hostname, ".")
	for i := 1; i < len(parts); i++ {
		domain := strings.Join(parts[i:], ".")
		if reason, exists := s.blacklistedDomains[domain]; exists {
			return true, fmt.Sprintf("Parent domain %s: %s", domain, reason), nil
		}
	}

	return false, "", nil
}

// ScanURLContent scans the content of a URL for suspicious elements
func (s *securityService) ScanURLContent(ctx context.Context, targetURL string) (*ContentScanResult, error) {
	result := &ContentScanResult{
		SuspiciousContent: []string{},
	}

	start := time.Now()

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set user agent
	req.Header.Set("User-Agent", "URLShortener-SecurityScanner/1.0")

	// Make request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		result.IsAccessible = false
		return result, nil
	}
	defer resp.Body.Close()

	result.ResponseTime = time.Since(start).Milliseconds()
	result.IsAccessible = true
	result.HTTPStatus = resp.StatusCode
	result.ContentType = resp.Header.Get("Content-Type")
	result.ContentLength = resp.ContentLength

	// Analyze SSL info for HTTPS
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		result.SSLInfo = &SSLInfo{
			IsSecure:     true,
			Issuer:       cert.Issuer.String(),
			Subject:      cert.Subject.String(),
			ExpiresAt:    cert.NotAfter,
			IsExpired:    time.Now().After(cert.NotAfter),
			IsSelfSigned: cert.Issuer.String() == cert.Subject.String(),
		}

		if result.SSLInfo.IsExpired {
			result.SuspiciousContent = append(result.SuspiciousContent, "Expired SSL certificate")
		}
		if result.SSLInfo.IsSelfSigned {
			result.SuspiciousContent = append(result.SuspiciousContent, "Self-signed SSL certificate")
		}
	}

	// Basic content analysis (limited to avoid heavy processing)
	if strings.Contains(result.ContentType, "text/html") {
		// Read limited content for analysis
		buffer := make([]byte, 8192) // Read first 8KB
		n, _ := resp.Body.Read(buffer)
		content := string(buffer[:n])

		// Check for suspicious elements
		result.HasJavaScript = strings.Contains(content, "<script")
		result.HasForms = strings.Contains(content, "<form")
		result.HasExternalLinks = s.hasExternalLinks(content, targetURL)

		// Check for suspicious keywords
		for _, keyword := range s.suspiciousKeywords {
			if strings.Contains(strings.ToLower(content), keyword) {
				result.SuspiciousContent = append(result.SuspiciousContent, fmt.Sprintf("Suspicious keyword: %s", keyword))
			}
		}
	}

	return result, nil
}

// CheckPhishing performs phishing detection
func (s *securityService) CheckPhishing(ctx context.Context, parsedURL *url.URL) (*PhishingReport, error) {
	report := &PhishingReport{
		PhishingIndicators: []string{},
		SuspiciousElements: []string{},
		ConfidenceScore:    0.0,
	}

	hostname := strings.ToLower(parsedURL.Hostname())
	fullURL := parsedURL.String()

	// Check against phishing patterns
	for _, pattern := range s.phishingPatterns {
		if pattern.MatchString(fullURL) || pattern.MatchString(hostname) {
			report.PhishingIndicators = append(report.PhishingIndicators, fmt.Sprintf("Matches phishing pattern: %s", pattern.String()))
			report.ConfidenceScore += 0.3
		}
	}

	// Check for homograph attacks (internationalized domain names)
	if s.hasHomographAttack(hostname) {
		report.PhishingIndicators = append(report.PhishingIndicators, "Potential homograph attack (suspicious unicode characters)")
		report.ConfidenceScore += 0.4
	}

	// Check for suspicious TLD combinations
	if s.hasSuspiciousTLD(hostname) {
		report.PhishingIndicators = append(report.PhishingIndicators, "Suspicious top-level domain")
		report.ConfidenceScore += 0.2
	}

	// Check for brand impersonation
	targetBrand := s.detectBrandImpersonation(hostname)
	if targetBrand != "" {
		report.TargetBrand = targetBrand
		report.PhishingIndicators = append(report.PhishingIndicators, fmt.Sprintf("Potential %s brand impersonation", targetBrand))
		report.ConfidenceScore += 0.5
	}

	// Check URL structure indicators
	if s.hasPhishingURLStructure(parsedURL) {
		report.PhishingIndicators = append(report.PhishingIndicators, "Suspicious URL structure")
		report.ConfidenceScore += 0.3
	}

	// Determine if phishing is suspected
	report.IsPhishingSuspected = report.ConfidenceScore >= 0.5

	return report, nil
}

// UpdateBlacklist updates the blacklisted domains
func (s *securityService) UpdateBlacklist(ctx context.Context, domains []string) error {
	for _, domain := range domains {
		s.blacklistedDomains[strings.ToLower(domain)] = "User reported"
	}
	
	s.logger.WithField("count", len(domains)).Info("Updated blacklist with new domains")
	return nil
}

// Helper functions

func (s *securityService) validateURLStructure(parsedURL *url.URL) error {
	// Check scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	// Check host
	if parsedURL.Host == "" {
		return fmt.Errorf("missing hostname")
	}

	// Check for malformed URLs
	if strings.Contains(parsedURL.Host, "..") {
		return fmt.Errorf("malformed hostname")
	}

	return nil
}

func (s *securityService) calculateThreatLevel(assessment *ThreatAssessment) ThreatLevel {
	if assessment.IsBlacklisted {
		return ThreatLevelCritical
	}

	riskCount := len(assessment.RiskFactors)
	switch {
	case riskCount >= 4:
		return ThreatLevelHigh
	case riskCount >= 2:
		return ThreatLevelMedium
	case riskCount >= 1:
		return ThreatLevelLow
	default:
		return ThreatLevelSafe
	}
}

func (s *securityService) isSuspiciousHost(hostname string) bool {
	suspiciousPatterns := []string{
		"bit\\.ly", "tinyurl", "goo\\.gl", "t\\.co", "short",
		"[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}", // IP addresses
		"[a-z]{10,}", // Very long subdomains
		"[0-9]{5,}",  // Long numeric strings
	}

	for _, pattern := range suspiciousPatterns {
		matched, _ := regexp.MatchString(pattern, hostname)
		if matched {
			return true
		}
	}

	return false
}

func (s *securityService) isShortURLService(hostname string) bool {
	shortURLServices := []string{
		"bit.ly", "tinyurl.com", "goo.gl", "t.co", "short.link",
		"ow.ly", "buff.ly", "rebrand.ly", "tiny.cc", "is.gd",
	}

	for _, service := range shortURLServices {
		if hostname == service || strings.HasSuffix(hostname, "."+service) {
			return true
		}
	}

	return false
}

func (s *securityService) hasSuspiciousPath(path string) bool {
	suspiciousPatterns := []string{
		"\\.\\.[\\/\\\\]", // Directory traversal
		"<script",        // XSS attempts
		"javascript:",    // JavaScript URLs
		"vbscript:",      // VBScript URLs
		"%3cscript",      // URL encoded script tags
	}

	lowerPath := strings.ToLower(path)
	for _, pattern := range suspiciousPatterns {
		matched, _ := regexp.MatchString(pattern, lowerPath)
		if matched {
			return true
		}
	}

	return false
}

func (s *securityService) hasExternalLinks(content, baseURL string) bool {
	parsedBase, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	// Simple check for external links
	linkPattern := regexp.MustCompile(`href\s*=\s*["']([^"']+)["']`)
	matches := linkPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			linkURL, err := url.Parse(match[1])
			if err == nil && linkURL.Host != "" && linkURL.Host != parsedBase.Host {
				return true
			}
		}
	}

	return false
}

func (s *securityService) hasHomographAttack(hostname string) bool {
	// Check for common homograph attack patterns
	suspiciousChars := []string{
		"а", "е", "о", "р", "с", "х", "у", // Cyrillic letters that look like Latin
		"ο", "ρ", "α", "ε",                 // Greek letters
	}

	for _, char := range suspiciousChars {
		if strings.Contains(hostname, char) {
			return true
		}
	}

	return false
}

func (s *securityService) hasSuspiciousTLD(hostname string) bool {
	suspiciousTLDs := []string{
		".tk", ".ml", ".ga", ".cf", // Free domains often used in phishing
		".click", ".download", ".review",
	}

	for _, tld := range suspiciousTLDs {
		if strings.HasSuffix(hostname, tld) {
			return true
		}
	}

	return false
}

func (s *securityService) detectBrandImpersonation(hostname string) string {
	brands := map[string]string{
		"paypal":   "PayPal",
		"amazon":   "Amazon",
		"google":   "Google",
		"facebook": "Facebook",
		"apple":    "Apple",
		"microsoft": "Microsoft",
		"netflix":  "Netflix",
		"spotify":  "Spotify",
		"instagram": "Instagram",
		"twitter":  "Twitter",
	}

	for brand, displayName := range brands {
		if strings.Contains(hostname, brand) && !strings.HasSuffix(hostname, "."+brand+".com") {
			return displayName
		}
	}

	return ""
}

func (s *securityService) hasPhishingURLStructure(parsedURL *url.URL) bool {
	// Check for multiple subdomains (common in phishing)
	subdomains := strings.Split(parsedURL.Hostname(), ".")
	if len(subdomains) > 4 {
		return true
	}

	// Check for suspicious URL patterns
	fullURL := parsedURL.String()
	if strings.Count(fullURL, "-") > 3 || strings.Count(fullURL, ".") > 5 {
		return true
	}

	return false
}

func (s *securityService) generateRecommendations(report *URLSecurityReport) []string {
	recommendations := []string{}

	if report.ThreatAssessment != nil {
		if report.ThreatAssessment.IsBlacklisted {
			recommendations = append(recommendations, "This domain is blacklisted and should not be accessed")
		}
		if report.ThreatAssessment.IsPrivateIP {
			recommendations = append(recommendations, "Avoid using private IP addresses in public URLs")
		}
		if report.ThreatAssessment.IsShortURL {
			recommendations = append(recommendations, "Be cautious with shortened URLs - verify the destination")
		}
	}

	if report.PhishingReport != nil && report.PhishingReport.IsPhishingSuspected {
		recommendations = append(recommendations, "This URL shows signs of phishing - verify authenticity before clicking")
	}

	if report.ContentScanResult != nil {
		if len(report.ContentScanResult.SuspiciousContent) > 0 {
			recommendations = append(recommendations, "URL contains suspicious content - exercise caution")
		}
		if report.ContentScanResult.SSLInfo != nil && report.ContentScanResult.SSLInfo.IsExpired {
			recommendations = append(recommendations, "SSL certificate is expired - connection may not be secure")
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "URL appears to be safe for general use")
	}

	return recommendations
}

// Helper function for private IP detection (reused from middleware)
func isPrivateOrLocalhost(hostname string) bool {
	if hostname == "localhost" || hostname == "127.0.0.1" || hostname == "::1" {
		return true
	}

	privatePatterns := []string{
		"^10\\.",
		"^172\\.(1[6-9]|2[0-9]|3[01])\\.",
		"^192\\.168\\.",
		"^169\\.254\\.",
	}

	for _, pattern := range privatePatterns {
		matched, _ := regexp.MatchString(pattern, hostname)
		if matched {
			return true
		}
	}

	return false
}

// Initialize blacklist with known malicious domains
func initializeBlacklist() map[string]string {
	return map[string]string{
		"malware.com":     "Known malware distribution",
		"phishing.test":   "Phishing test domain",
		"spam.example":    "Spam source",
		"malicious.site":  "Malicious content",
		"fake-bank.com":   "Banking phishing",
		"virus.download":  "Malware download site",
		"scam.online":     "Online scam",
		"exploit.kit":     "Exploit kit distribution",
	}
}

// Initialize phishing detection patterns
func initializePhishingPatterns() []*regexp.Regexp {
	patterns := []string{
		`(?i)(paypal|amazon|google|facebook|apple)[\w\-]*\.(tk|ml|ga|cf)`,
		`(?i)(secure|login|account|verify)[\w\-]*\.(tk|ml|ga|cf)`,
		`(?i)(bank|visa|mastercard)[\w\-]*\.(tk|ml|ga|cf)`,
		`(?i)[\w\-]*(login|signin|account)[\w\-]*\.(tk|ml|ga|cf)`,
	}

	compiledPatterns := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		if re, err := regexp.Compile(pattern); err == nil {
			compiledPatterns = append(compiledPatterns, re)
		}
	}

	return compiledPatterns
}

// Initialize suspicious keywords for content analysis
func initializeSuspiciousKeywords() []string {
	return []string{
		"urgent action required",
		"verify your account",
		"suspended account",
		"click here now",
		"limited time offer",
		"congratulations you won",
		"claim your prize",
		"free money",
		"nigerian prince",
		"inheritance",
		"lottery winner",
		"tax refund",
		"phishing",
		"malware",
		"virus",
		"trojan",
	}
}