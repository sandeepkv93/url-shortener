package middleware

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/core/services"
)

// TLSMiddleware provides SSL/TLS configuration and security enhancements
type TLSMiddleware struct {
	config       *config.Config
	logger       *services.LoggingService
	tlsConfig    *tls.Config
	certificates []tls.Certificate
}

// TLSInfo contains information about the TLS configuration
type TLSInfo struct {
	CertificateCount int                    `json:"certificate_count"`
	MinVersion       string                 `json:"min_version"`
	MaxVersion       string                 `json:"max_version"`
	CipherSuites     []string              `json:"cipher_suites"`
	Certificates     []CertificateInfo     `json:"certificates"`
	OCSP             bool                  `json:"ocsp_enabled"`
	SessionTickets   bool                  `json:"session_tickets_enabled"`
}

// CertificateInfo contains information about a specific certificate
type CertificateInfo struct {
	Subject    string    `json:"subject"`
	Issuer     string    `json:"issuer"`
	NotBefore  time.Time `json:"not_before"`
	NotAfter   time.Time `json:"not_after"`
	DNSNames   []string  `json:"dns_names"`
	ExpiresIn  int       `json:"expires_in_days"`
	IsExpired  bool      `json:"is_expired"`
	IsSelfSigned bool    `json:"is_self_signed"`
}

// NewTLSMiddleware creates a new TLS middleware
func NewTLSMiddleware(cfg *config.Config, logger *services.LoggingService) (*TLSMiddleware, error) {
	tm := &TLSMiddleware{
		config: cfg,
		logger: logger,
	}
	
	if err := tm.initializeTLSConfig(); err != nil {
		return nil, fmt.Errorf("failed to initialize TLS config: %w", err)
	}
	
	return tm, nil
}

// initializeTLSConfig sets up the TLS configuration
func (tm *TLSMiddleware) initializeTLSConfig() error {
	if !tm.config.IsTLSEnabled() {
		tm.logger.Info("TLS is not enabled, skipping TLS configuration")
		return nil
	}
	
	// Load certificates
	if err := tm.loadCertificates(); err != nil {
		return fmt.Errorf("failed to load certificates: %w", err)
	}
	
	// Create TLS configuration
	tm.tlsConfig = &tls.Config{
		Certificates: tm.certificates,
		MinVersion:   tm.getTLSVersion(tm.config.GetTLSMinVersion()),
		MaxVersion:   tls.VersionTLS13, // Always use the latest as max
		CipherSuites: tm.getCipherSuites(),
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
		PreferServerCipherSuites: false, // Modern practice
		SessionTicketsDisabled:   false, // Enable for performance
		Renegotiation:           tls.RenegotiateNever,
		InsecureSkipVerify:      false,
	}
	
	// Enable OCSP stapling if available
	for i := range tm.tlsConfig.Certificates {
		cert := &tm.tlsConfig.Certificates[i]
		if len(cert.Certificate) > 0 {
			// OCSP stapling would be implemented here
			// This is a placeholder for production implementation
		}
	}
	
	tm.logger.Info("TLS configuration initialized successfully",
		"min_version", tm.getTLSVersionString(tm.tlsConfig.MinVersion),
		"max_version", tm.getTLSVersionString(tm.tlsConfig.MaxVersion),
		"certificate_count", len(tm.certificates),
	)
	
	return nil
}

// loadCertificates loads SSL certificates from files
func (tm *TLSMiddleware) loadCertificates() error {
	certFile := tm.config.GetTLSCertFile()
	keyFile := tm.config.GetTLSKeyFile()
	
	if certFile == "" || keyFile == "" {
		return fmt.Errorf("certificate file and key file must be specified")
	}
	
	// Check if files exist
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		return fmt.Errorf("certificate file not found: %s", certFile)
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		return fmt.Errorf("key file not found: %s", keyFile)
	}
	
	// Load certificate and key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load certificate pair: %w", err)
	}
	
	tm.certificates = []tls.Certificate{cert}
	
	// Validate certificate
	if err := tm.validateCertificate(&cert); err != nil {
		tm.logger.Warn("Certificate validation warning", "error", err.Error())
	}
	
	return nil
}

// validateCertificate validates the loaded certificate
func (tm *TLSMiddleware) validateCertificate(cert *tls.Certificate) error {
	if len(cert.Certificate) == 0 {
		return fmt.Errorf("certificate is empty")
	}
	
	// Parse the certificate
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}
	
	// Check expiration
	now := time.Now()
	if now.After(x509Cert.NotAfter) {
		return fmt.Errorf("certificate has expired on %v", x509Cert.NotAfter)
	}
	
	// Check if expiring soon (30 days)
	expirationWarning := 30 * 24 * time.Hour
	if now.Add(expirationWarning).After(x509Cert.NotAfter) {
		tm.logger.Warn("Certificate expires soon",
			"expires_at", x509Cert.NotAfter,
			"days_remaining", int(x509Cert.NotAfter.Sub(now).Hours()/24),
		)
	}
	
	// Check if certificate is self-signed
	if x509Cert.Issuer.String() == x509Cert.Subject.String() {
		tm.logger.Warn("Using self-signed certificate - not recommended for production")
	}
	
	tm.logger.Info("Certificate validation completed",
		"subject", x509Cert.Subject.String(),
		"issuer", x509Cert.Issuer.String(),
		"expires_at", x509Cert.NotAfter,
		"dns_names", x509Cert.DNSNames,
	)
	
	return nil
}

// Handler provides HTTPS enforcement and security headers
func (tm *TLSMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add security headers for HTTPS
		tm.addSecurityHeaders(w, r)
		
		// Enforce HTTPS if enabled
		if tm.config.Security.EnableHTTPS && !tm.isHTTPS(r) {
			tm.redirectToHTTPS(w, r)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// addSecurityHeaders adds security headers for HTTPS connections
func (tm *TLSMiddleware) addSecurityHeaders(w http.ResponseWriter, r *http.Request) {
	if !tm.isHTTPS(r) {
		return
	}
	
	// HSTS Header
	if tm.config.IsProduction() {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	} else {
		w.Header().Set("Strict-Transport-Security", "max-age=3600; includeSubDomains")
	}
	
	// Additional security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	
	// Content Security Policy
	csp := tm.buildCSPHeader()
	if csp != "" {
		w.Header().Set("Content-Security-Policy", csp)
	}
	
	// Certificate Transparency
	w.Header().Set("Expect-CT", "max-age=86400, enforce")
}

// buildCSPHeader builds the Content Security Policy header
func (tm *TLSMiddleware) buildCSPHeader() string {
	directives := []string{
		"default-src 'self'",
		"script-src 'self' 'unsafe-inline'",
		"style-src 'self' 'unsafe-inline'",
		"img-src 'self' data: https:",
		"font-src 'self'",
		"connect-src 'self'",
		"frame-ancestors 'none'",
		"base-uri 'self'",
		"form-action 'self'",
	}
	
	if tm.config.IsDevelopment() {
		// More permissive CSP for development
		directives = append(directives, "unsafe-eval")
	}
	
	return strings.Join(directives, "; ")
}

// redirectToHTTPS redirects HTTP requests to HTTPS
func (tm *TLSMiddleware) redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	httpsURL := "https://" + r.Host + r.RequestURI
	
	tm.logger.Debug("Redirecting HTTP to HTTPS",
		"original_url", r.URL.String(),
		"https_url", httpsURL,
	)
	
	http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
}

// isHTTPS checks if the request is over HTTPS
func (tm *TLSMiddleware) isHTTPS(r *http.Request) bool {
	// Check direct TLS connection
	if r.TLS != nil {
		return true
	}
	
	// Check proxy headers
	if tm.config.IsTrustProxyEnabled() {
		if proto := r.Header.Get("X-Forwarded-Proto"); proto == "https" {
			return true
		}
		if r.Header.Get("X-Forwarded-SSL") == "on" {
			return true
		}
	}
	
	return false
}

// GetTLSConfig returns the configured TLS configuration
func (tm *TLSMiddleware) GetTLSConfig() *tls.Config {
	return tm.tlsConfig
}

// GetTLSInfo returns information about the TLS configuration
func (tm *TLSMiddleware) GetTLSInfo() *TLSInfo {
	if tm.tlsConfig == nil {
		return &TLSInfo{}
	}
	
	info := &TLSInfo{
		CertificateCount: len(tm.certificates),
		MinVersion:       tm.getTLSVersionString(tm.tlsConfig.MinVersion),
		MaxVersion:       tm.getTLSVersionString(tm.tlsConfig.MaxVersion),
		CipherSuites:     tm.getCipherSuiteNames(),
		OCSP:             false, // Would be true if OCSP is implemented
		SessionTickets:   !tm.tlsConfig.SessionTicketsDisabled,
	}
	
	// Get certificate information
	for _, cert := range tm.certificates {
		if len(cert.Certificate) > 0 {
			if x509Cert, err := x509.ParseCertificate(cert.Certificate[0]); err == nil {
				certInfo := CertificateInfo{
					Subject:      x509Cert.Subject.String(),
					Issuer:       x509Cert.Issuer.String(),
					NotBefore:    x509Cert.NotBefore,
					NotAfter:     x509Cert.NotAfter,
					DNSNames:     x509Cert.DNSNames,
					ExpiresIn:    int(time.Until(x509Cert.NotAfter).Hours() / 24),
					IsExpired:    time.Now().After(x509Cert.NotAfter),
					IsSelfSigned: x509Cert.Issuer.String() == x509Cert.Subject.String(),
				}
				info.Certificates = append(info.Certificates, certInfo)
			}
		}
	}
	
	return info
}

// getTLSVersion converts version string to TLS constant
func (tm *TLSMiddleware) getTLSVersion(version string) uint16 {
	switch strings.ToLower(version) {
	case "1.0":
		return tls.VersionTLS10
	case "1.1":
		return tls.VersionTLS11
	case "1.2":
		return tls.VersionTLS12
	case "1.3":
		return tls.VersionTLS13
	default:
		return tls.VersionTLS12 // Default to TLS 1.2
	}
}

// getTLSVersionString converts TLS constant to string
func (tm *TLSMiddleware) getTLSVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown"
	}
}

// getCipherSuites returns the configured cipher suites
func (tm *TLSMiddleware) getCipherSuites() []uint16 {
	// Use secure cipher suites
	return []uint16{
		tls.TLS_AES_256_GCM_SHA384,
		tls.TLS_AES_128_GCM_SHA256,
		tls.TLS_CHACHA20_POLY1305_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	}
}

// getCipherSuiteNames returns human-readable cipher suite names
func (tm *TLSMiddleware) getCipherSuiteNames() []string {
	suites := tm.getCipherSuites()
	names := make([]string, len(suites))
	
	suiteNames := map[uint16]string{
		tls.TLS_AES_256_GCM_SHA384:                   "TLS_AES_256_GCM_SHA384",
		tls.TLS_AES_128_GCM_SHA256:                   "TLS_AES_128_GCM_SHA256",
		tls.TLS_CHACHA20_POLY1305_SHA256:             "TLS_CHACHA20_POLY1305_SHA256",
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:  "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:    "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305:   "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305:     "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:  "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:    "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	}
	
	for i, suite := range suites {
		if name, exists := suiteNames[suite]; exists {
			names[i] = name
		} else {
			names[i] = fmt.Sprintf("Unknown(%d)", suite)
		}
	}
	
	return names
}

// TLSInfoHandler returns TLS configuration information
func (tm *TLSMiddleware) TLSInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := tm.GetTLSInfo()
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		// Simple JSON encoding (in production, use proper JSON marshaling)
		fmt.Fprintf(w, `{
			"tls_enabled": %t,
			"certificate_count": %d,
			"min_version": "%s",
			"max_version": "%s",
			"cipher_suites": %d,
			"certificates": %d
		}`, 
			tm.config.IsTLSEnabled(),
			info.CertificateCount,
			info.MinVersion,
			info.MaxVersion,
			len(info.CipherSuites),
			len(info.Certificates),
		)
	}
}

// StartCertificateMonitoring starts a goroutine to monitor certificate expiration
func (tm *TLSMiddleware) StartCertificateMonitoring() {
	if !tm.config.IsTLSEnabled() {
		return
	}
	
	go func() {
		ticker := time.NewTicker(24 * time.Hour) // Check daily
		defer ticker.Stop()
		
		for range ticker.C {
			tm.checkCertificateExpiration()
		}
	}()
}

// checkCertificateExpiration checks if certificates are expiring soon
func (tm *TLSMiddleware) checkCertificateExpiration() {
	for _, cert := range tm.certificates {
		if len(cert.Certificate) > 0 {
			if x509Cert, err := x509.ParseCertificate(cert.Certificate[0]); err == nil {
				daysUntilExpiration := int(time.Until(x509Cert.NotAfter).Hours() / 24)
				
				if daysUntilExpiration <= 30 {
					tm.logger.Warn("Certificate expiring soon",
						"subject", x509Cert.Subject.String(),
						"expires_at", x509Cert.NotAfter,
						"days_remaining", daysUntilExpiration,
					)
				}
				
				if daysUntilExpiration <= 7 {
					tm.logger.Error("Certificate expires very soon",
						"subject", x509Cert.Subject.String(),
						"expires_at", x509Cert.NotAfter,
						"days_remaining", daysUntilExpiration,
					)
				}
			}
		}
	}
}