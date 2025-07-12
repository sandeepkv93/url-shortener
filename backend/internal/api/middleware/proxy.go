package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"url-shortener/internal/config"
	"url-shortener/internal/core/services"
)

// ProxyMiddleware handles reverse proxy configurations and header processing
type ProxyMiddleware struct {
	config *config.Config
	logger *services.LoggingService
	trustedProxies []net.IPNet
}

// NewProxyMiddleware creates a new proxy middleware
func NewProxyMiddleware(cfg *config.Config, logger *services.LoggingService) *ProxyMiddleware {
	pm := &ProxyMiddleware{
		config: cfg,
		logger: logger,
	}
	
	// Initialize trusted proxy networks
	if cfg.IsTrustProxyEnabled() {
		pm.initializeTrustedProxies()
	}
	
	return pm
}

// Handler provides the middleware handler for proxy processing
func (pm *ProxyMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if pm.config.IsTrustProxyEnabled() {
			pm.processProxyHeaders(r)
		}
		
		next.ServeHTTP(w, r)
	})
}

// processProxyHeaders processes and validates proxy headers
func (pm *ProxyMiddleware) processProxyHeaders(r *http.Request) {
	// Get the immediate client IP
	clientIP := pm.getClientIP(r)
	
	// Check if the client IP is from a trusted proxy
	if !pm.isTrustedProxy(clientIP) {
		pm.logger.Debug("Request from untrusted proxy, ignoring proxy headers", 
			"client_ip", clientIP,
			"remote_addr", r.RemoteAddr,
		)
		return
	}
	
	// Process X-Forwarded-For header
	pm.processForwardedFor(r)
	
	// Process X-Real-IP header
	pm.processRealIP(r)
	
	// Process X-Forwarded-Proto header
	pm.processForwardedProto(r)
	
	// Process X-Forwarded-Host header
	pm.processForwardedHost(r)
	
	// Process X-Forwarded-Port header
	pm.processForwardedPort(r)
	
	// Log processed headers in debug mode
	if pm.logger != nil {
		pm.logger.Debug("Processed proxy headers",
			"original_remote_addr", r.RemoteAddr,
			"x_forwarded_for", r.Header.Get("X-Forwarded-For"),
			"x_real_ip", r.Header.Get("X-Real-IP"),
			"x_forwarded_proto", r.Header.Get("X-Forwarded-Proto"),
			"x_forwarded_host", r.Header.Get("X-Forwarded-Host"),
		)
	}
}

// processForwardedFor processes the X-Forwarded-For header
func (pm *ProxyMiddleware) processForwardedFor(r *http.Request) {
	forwardedFor := r.Header.Get(pm.config.GetProxyHeader())
	if forwardedFor == "" {
		return
	}
	
	// Parse the forwarded-for chain
	ips := strings.Split(forwardedFor, ",")
	if len(ips) == 0 {
		return
	}
	
	// Get the first (original client) IP
	originalIP := strings.TrimSpace(ips[0])
	if originalIP != "" && pm.isValidIP(originalIP) {
		// Update RemoteAddr to reflect the original client IP
		host, port, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// If no port, just use the original IP
			r.RemoteAddr = originalIP
		} else {
			// Preserve the port from the original connection
			r.RemoteAddr = net.JoinHostPort(originalIP, port)
		}
		
		pm.logger.Debug("Updated RemoteAddr from X-Forwarded-For", 
			"original_ip", originalIP,
			"new_remote_addr", r.RemoteAddr,
		)
	}
}

// processRealIP processes the X-Real-IP header
func (pm *ProxyMiddleware) processRealIP(r *http.Request) {
	realIP := r.Header.Get(pm.config.GetRealIPHeader())
	if realIP == "" {
		return
	}
	
	if pm.isValidIP(realIP) {
		// X-Real-IP takes precedence over X-Forwarded-For
		host, port, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			r.RemoteAddr = realIP
		} else {
			r.RemoteAddr = net.JoinHostPort(realIP, port)
		}
		
		pm.logger.Debug("Updated RemoteAddr from X-Real-IP", 
			"real_ip", realIP,
			"new_remote_addr", r.RemoteAddr,
		)
	}
}

// processForwardedProto processes the X-Forwarded-Proto header
func (pm *ProxyMiddleware) processForwardedProto(r *http.Request) {
	proto := r.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		return
	}
	
	// Normalize protocol
	proto = strings.ToLower(strings.TrimSpace(proto))
	
	switch proto {
	case "https", "wss":
		r.URL.Scheme = "https"
		r.Header.Set("X-Forwarded-Proto", "https")
		pm.logger.Debug("Updated request scheme to HTTPS from proxy header")
	case "http", "ws":
		r.URL.Scheme = "http"
		r.Header.Set("X-Forwarded-Proto", "http")
		pm.logger.Debug("Updated request scheme to HTTP from proxy header")
	default:
		pm.logger.Warn("Unknown protocol in X-Forwarded-Proto header", "proto", proto)
	}
}

// processForwardedHost processes the X-Forwarded-Host header
func (pm *ProxyMiddleware) processForwardedHost(r *http.Request) {
	forwardedHost := r.Header.Get("X-Forwarded-Host")
	if forwardedHost == "" {
		return
	}
	
	// Update the Host header
	r.Host = forwardedHost
	r.URL.Host = forwardedHost
	
	pm.logger.Debug("Updated Host from X-Forwarded-Host", "host", forwardedHost)
}

// processForwardedPort processes the X-Forwarded-Port header
func (pm *ProxyMiddleware) processForwardedPort(r *http.Request) {
	forwardedPort := r.Header.Get("X-Forwarded-Port")
	if forwardedPort == "" {
		return
	}
	
	// If we have a host without port, add the forwarded port
	if r.Host != "" && !strings.Contains(r.Host, ":") {
		r.Host = net.JoinHostPort(r.Host, forwardedPort)
		r.URL.Host = r.Host
		
		pm.logger.Debug("Added port to Host from X-Forwarded-Port", 
			"port", forwardedPort,
			"new_host", r.Host,
		)
	}
}

// getClientIP extracts the immediate client IP from the request
func (pm *ProxyMiddleware) getClientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If there's no port, RemoteAddr is just the IP
		return r.RemoteAddr
	}
	return ip
}

// isTrustedProxy checks if the given IP is from a trusted proxy
func (pm *ProxyMiddleware) isTrustedProxy(ip string) bool {
	if len(pm.trustedProxies) == 0 {
		// If no trusted proxies configured, trust all (not recommended for production)
		return true
	}
	
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false
	}
	
	for _, trustedNet := range pm.trustedProxies {
		if trustedNet.Contains(clientIP) {
			return true
		}
	}
	
	return false
}

// isValidIP validates if a string is a valid IP address
func (pm *ProxyMiddleware) isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// initializeTrustedProxies initializes the trusted proxy networks
func (pm *ProxyMiddleware) initializeTrustedProxies() {
	// Default trusted proxy networks (common cloud provider ranges)
	defaultTrustedProxies := []string{
		"10.0.0.0/8",     // Private network
		"172.16.0.0/12",  // Private network
		"192.168.0.0/16", // Private network
		"127.0.0.0/8",    // Loopback
		"::1/128",        // IPv6 loopback
		"fc00::/7",       // IPv6 private
	}
	
	// Add common cloud provider ranges
	cloudProviderRanges := []string{
		// AWS ALB ranges (these would need to be updated based on actual deployment)
		"52.0.0.0/8",
		"54.0.0.0/8",
		// Google Cloud Load Balancer ranges
		"35.191.0.0/16",
		"130.211.0.0/22",
		// Azure Load Balancer ranges
		"168.63.129.16/32",
		// Cloudflare ranges (sample)
		"103.21.244.0/22",
		"103.22.200.0/22",
	}
	
	allRanges := append(defaultTrustedProxies, cloudProviderRanges...)
	
	for _, cidr := range allRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			pm.logger.Warn("Invalid CIDR in trusted proxies", "cidr", cidr, "error", err.Error())
			continue
		}
		pm.trustedProxies = append(pm.trustedProxies, *network)
	}
	
	pm.logger.Info("Initialized trusted proxy networks", "count", len(pm.trustedProxies))
}

// GetRealClientIP returns the real client IP considering proxy headers
func (pm *ProxyMiddleware) GetRealClientIP(r *http.Request) string {
	if !pm.config.IsTrustProxyEnabled() {
		return pm.getClientIP(r)
	}
	
	// Check X-Real-IP first (most specific)
	if realIP := r.Header.Get(pm.config.GetRealIPHeader()); realIP != "" && pm.isValidIP(realIP) {
		return realIP
	}
	
	// Check X-Forwarded-For (may contain multiple IPs)
	if forwardedFor := r.Header.Get(pm.config.GetProxyHeader()); forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			originalIP := strings.TrimSpace(ips[0])
			if pm.isValidIP(originalIP) {
				return originalIP
			}
		}
	}
	
	// Fall back to direct connection IP
	return pm.getClientIP(r)
}

// IsHTTPS checks if the request is HTTPS considering proxy headers
func (pm *ProxyMiddleware) IsHTTPS(r *http.Request) bool {
	// Check if proxy indicates HTTPS
	if pm.config.IsTrustProxyEnabled() {
		if proto := r.Header.Get("X-Forwarded-Proto"); proto == "https" {
			return true
		}
		if r.Header.Get("X-Forwarded-SSL") == "on" {
			return true
		}
	}
	
	// Check direct TLS connection
	return r.TLS != nil
}

// SetSecurityHeaders sets appropriate security headers based on proxy configuration
func (pm *ProxyMiddleware) SetSecurityHeaders(w http.ResponseWriter, r *http.Request) {
	// Set HSTS header if HTTPS
	if pm.IsHTTPS(r) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}
	
	// Set other security headers
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
}

// ProxyInfoMiddleware adds proxy information to the request context
func (pm *ProxyMiddleware) ProxyInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add proxy information to response headers for debugging (in development)
		if pm.config.IsDevelopment() {
			w.Header().Set("X-Real-IP-Detected", pm.GetRealClientIP(r))
			w.Header().Set("X-HTTPS-Detected", fmt.Sprintf("%t", pm.IsHTTPS(r)))
			w.Header().Set("X-Proxy-Trust-Enabled", fmt.Sprintf("%t", pm.config.IsTrustProxyEnabled()))
		}
		
		// Set security headers
		if pm.config.Security.EnableSecurityHeaders {
			pm.SetSecurityHeaders(w, r)
		}
		
		next.ServeHTTP(w, r)
	})
}