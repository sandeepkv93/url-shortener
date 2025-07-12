package ssl

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/core/services"
)

// CertificateManager handles SSL certificate operations
type CertificateManager struct {
	config *config.Config
	logger *services.LoggingService
}

// CertificateRequest represents a certificate signing request
type CertificateRequest struct {
	CommonName       string
	Organization     []string
	Country          []string
	Province         []string
	Locality         []string
	DNSNames         []string
	IPAddresses      []net.IP
	ValidityDays     int
	KeySize          int
	IsCA             bool
	KeyUsage         x509.KeyUsage
	ExtKeyUsage      []x509.ExtKeyUsage
}

// CertificateInfo provides detailed information about a certificate
type CertificateInfo struct {
	Subject            string    `json:"subject"`
	Issuer             string    `json:"issuer"`
	SerialNumber       string    `json:"serial_number"`
	NotBefore          time.Time `json:"not_before"`
	NotAfter           time.Time `json:"not_after"`
	DNSNames           []string  `json:"dns_names"`
	IPAddresses        []string  `json:"ip_addresses"`
	KeyUsage           string    `json:"key_usage"`
	ExtKeyUsage        []string  `json:"ext_key_usage"`
	IsCA               bool      `json:"is_ca"`
	IsSelfSigned       bool      `json:"is_self_signed"`
	SignatureAlgorithm string    `json:"signature_algorithm"`
	PublicKeyAlgorithm string    `json:"public_key_algorithm"`
	Version            int       `json:"version"`
	ExpiresInDays      int       `json:"expires_in_days"`
	IsExpired          bool      `json:"is_expired"`
	IsValid            bool      `json:"is_valid"`
}

// NewCertificateManager creates a new certificate manager
func NewCertificateManager(cfg *config.Config, logger *services.LoggingService) *CertificateManager {
	return &CertificateManager{
		config: cfg,
		logger: logger,
	}
}

// GenerateSelfSignedCertificate generates a self-signed certificate for development
func (cm *CertificateManager) GenerateSelfSignedCertificate(req *CertificateRequest) ([]byte, []byte, error) {
	cm.logger.Info("Generating self-signed certificate", "common_name", req.CommonName)
	
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, req.KeySize)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	
	// Create certificate template
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:         req.CommonName,
			Organization:       req.Organization,
			Country:            req.Country,
			Province:           req.Province,
			Locality:           req.Locality,
			OrganizationalUnit: []string{"IT Department"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(req.ValidityDays) * 24 * time.Hour),
		KeyUsage:              req.KeyUsage,
		ExtKeyUsage:           req.ExtKeyUsage,
		BasicConstraintsValid: true,
		IsCA:                  req.IsCA,
		DNSNames:              req.DNSNames,
		IPAddresses:           req.IPAddresses,
	}
	
	// Add default key usage if not specified
	if template.KeyUsage == 0 {
		template.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
	}
	
	// Add default extended key usage if not specified
	if len(template.ExtKeyUsage) == 0 {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	}
	
	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}
	
	// Encode certificate to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})
	
	// Encode private key to PEM
	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
	}
	
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyDER,
	})
	
	cm.logger.Info("Self-signed certificate generated successfully",
		"common_name", req.CommonName,
		"validity_days", req.ValidityDays,
		"dns_names", req.DNSNames,
	)
	
	return certPEM, keyPEM, nil
}

// SaveCertificateToFile saves certificate and key to files
func (cm *CertificateManager) SaveCertificateToFile(certPEM, keyPEM []byte, certPath, keyPath string) error {
	// Create directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(certPath), 0755); err != nil {
		return fmt.Errorf("failed to create certificate directory: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(keyPath), 0755); err != nil {
		return fmt.Errorf("failed to create key directory: %w", err)
	}
	
	// Write certificate file
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return fmt.Errorf("failed to write certificate file: %w", err)
	}
	
	// Write key file with restricted permissions
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}
	
	cm.logger.Info("Certificate saved to files",
		"cert_path", certPath,
		"key_path", keyPath,
	)
	
	return nil
}

// LoadCertificateInfo loads and analyzes a certificate file
func (cm *CertificateManager) LoadCertificateInfo(certPath string) (*CertificateInfo, error) {
	// Read certificate file
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	
	// Decode PEM block
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	
	// Parse certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}
	
	// Extract IP addresses
	ipAddresses := make([]string, len(cert.IPAddresses))
	for i, ip := range cert.IPAddresses {
		ipAddresses[i] = ip.String()
	}
	
	// Extract extended key usage
	extKeyUsage := make([]string, len(cert.ExtKeyUsage))
	for i, eku := range cert.ExtKeyUsage {
		extKeyUsage[i] = cm.extKeyUsageToString(eku)
	}
	
	// Check if certificate is valid
	now := time.Now()
	isValid := now.After(cert.NotBefore) && now.Before(cert.NotAfter)
	isExpired := now.After(cert.NotAfter)
	expiresInDays := int(time.Until(cert.NotAfter).Hours() / 24)
	
	info := &CertificateInfo{
		Subject:            cert.Subject.String(),
		Issuer:             cert.Issuer.String(),
		SerialNumber:       cert.SerialNumber.String(),
		NotBefore:          cert.NotBefore,
		NotAfter:           cert.NotAfter,
		DNSNames:           cert.DNSNames,
		IPAddresses:        ipAddresses,
		KeyUsage:           cm.keyUsageToString(cert.KeyUsage),
		ExtKeyUsage:        extKeyUsage,
		IsCA:               cert.IsCA,
		IsSelfSigned:       cert.Subject.String() == cert.Issuer.String(),
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
		Version:            cert.Version,
		ExpiresInDays:      expiresInDays,
		IsExpired:          isExpired,
		IsValid:            isValid,
	}
	
	return info, nil
}

// ValidateCertificateChain validates a certificate chain
func (cm *CertificateManager) ValidateCertificateChain(certPath string, intermediatesPath ...string) error {
	// Load main certificate
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("failed to read certificate: %w", err)
	}
	
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return fmt.Errorf("failed to decode certificate PEM")
	}
	
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}
	
	// Create certificate pool for intermediates
	intermediates := x509.NewCertPool()
	
	// Load intermediate certificates
	for _, intermediatePath := range intermediatesPath {
		intermediatePEM, err := os.ReadFile(intermediatePath)
		if err != nil {
			cm.logger.Warn("Failed to read intermediate certificate", "path", intermediatePath, "error", err)
			continue
		}
		
		if !intermediates.AppendCertsFromPEM(intermediatePEM) {
			cm.logger.Warn("Failed to parse intermediate certificate", "path", intermediatePath)
		}
	}
	
	// Create verification options
	opts := x509.VerifyOptions{
		Intermediates: intermediates,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	
	// Verify certificate chain
	chains, err := cert.Verify(opts)
	if err != nil {
		return fmt.Errorf("certificate chain validation failed: %w", err)
	}
	
	cm.logger.Info("Certificate chain validation successful",
		"chains_found", len(chains),
		"subject", cert.Subject.String(),
	)
	
	return nil
}

// CheckCertificateExpiration checks if certificates are expiring soon
func (cm *CertificateManager) CheckCertificateExpiration(certPaths []string, warningDays int) []string {
	var expiringSoon []string
	
	for _, certPath := range certPaths {
		info, err := cm.LoadCertificateInfo(certPath)
		if err != nil {
			cm.logger.Error("Failed to load certificate info", "path", certPath, "error", err)
			continue
		}
		
		if info.ExpiresInDays <= warningDays {
			expiringSoon = append(expiringSoon, certPath)
			cm.logger.Warn("Certificate expiring soon",
				"path", certPath,
				"subject", info.Subject,
				"expires_in_days", info.ExpiresInDays,
				"expires_at", info.NotAfter,
			)
		}
	}
	
	return expiringSoon
}

// CreateCSR creates a Certificate Signing Request
func (cm *CertificateManager) CreateCSR(req *CertificateRequest) ([]byte, []byte, error) {
	cm.logger.Info("Creating Certificate Signing Request", "common_name", req.CommonName)
	
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, req.KeySize)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	
	// Create CSR template
	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:         req.CommonName,
			Organization:       req.Organization,
			Country:            req.Country,
			Province:           req.Province,
			Locality:           req.Locality,
			OrganizationalUnit: []string{"IT Department"},
		},
		DNSNames:    req.DNSNames,
		IPAddresses: req.IPAddresses,
	}
	
	// Create CSR
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, template, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CSR: %w", err)
	}
	
	// Encode CSR to PEM
	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDER,
	})
	
	// Encode private key to PEM
	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
	}
	
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyDER,
	})
	
	cm.logger.Info("CSR created successfully", "common_name", req.CommonName)
	
	return csrPEM, keyPEM, nil
}

// GetDefaultCertificateRequest returns a default certificate request for development
func (cm *CertificateManager) GetDefaultCertificateRequest() *CertificateRequest {
	return &CertificateRequest{
		CommonName:   "localhost",
		Organization: []string{"URL Shortener Development"},
		Country:      []string{"US"},
		Province:     []string{""},
		Locality:     []string{""},
		DNSNames:     []string{"localhost", "127.0.0.1", "url-shortener.local"},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		ValidityDays: 365,
		KeySize:      2048,
		IsCA:         false,
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
}

// Helper methods

func (cm *CertificateManager) keyUsageToString(usage x509.KeyUsage) string {
	var usages []string
	
	if usage&x509.KeyUsageDigitalSignature != 0 {
		usages = append(usages, "DigitalSignature")
	}
	if usage&x509.KeyUsageContentCommitment != 0 {
		usages = append(usages, "ContentCommitment")
	}
	if usage&x509.KeyUsageKeyEncipherment != 0 {
		usages = append(usages, "KeyEncipherment")
	}
	if usage&x509.KeyUsageDataEncipherment != 0 {
		usages = append(usages, "DataEncipherment")
	}
	if usage&x509.KeyUsageKeyAgreement != 0 {
		usages = append(usages, "KeyAgreement")
	}
	if usage&x509.KeyUsageCertSign != 0 {
		usages = append(usages, "CertSign")
	}
	if usage&x509.KeyUsageCRLSign != 0 {
		usages = append(usages, "CRLSign")
	}
	if usage&x509.KeyUsageEncipherOnly != 0 {
		usages = append(usages, "EncipherOnly")
	}
	if usage&x509.KeyUsageDecipherOnly != 0 {
		usages = append(usages, "DecipherOnly")
	}
	
	if len(usages) == 0 {
		return "None"
	}
	
	return fmt.Sprintf("[%s]", fmt.Sprintf("%v", usages))
}

func (cm *CertificateManager) extKeyUsageToString(usage x509.ExtKeyUsage) string {
	switch usage {
	case x509.ExtKeyUsageServerAuth:
		return "ServerAuth"
	case x509.ExtKeyUsageClientAuth:
		return "ClientAuth"
	case x509.ExtKeyUsageCodeSigning:
		return "CodeSigning"
	case x509.ExtKeyUsageEmailProtection:
		return "EmailProtection"
	case x509.ExtKeyUsageTimeStamping:
		return "TimeStamping"
	case x509.ExtKeyUsageOCSPSigning:
		return "OCSPSigning"
	default:
		return fmt.Sprintf("Unknown(%d)", int(usage))
	}
}

// LoadTLSConfig loads TLS configuration from certificate files
func (cm *CertificateManager) LoadTLSConfig(certPath, keyPath string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate pair: %w", err)
	}
	
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
		PreferServerCipherSuites: false,
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
	}
	
	return config, nil
}