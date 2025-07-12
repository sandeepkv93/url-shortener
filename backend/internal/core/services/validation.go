package services

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"url-shortener/internal/core/domain"
)

// ValidationService provides comprehensive input validation and sanitization
type ValidationService interface {
	ValidateAndSanitizeURL(ctx context.Context, rawURL string) (*domain.URLValidationResult, error)
	ValidateUserInput(ctx context.Context, input interface{}) error
	SanitizeString(input string) string
	ValidatePassword(password string) *domain.PasswordValidation
	ValidateEmail(email string) error
	ValidateAlias(alias string) error
	SanitizeHTML(input string) string
	ValidateJSONInput(input map[string]interface{}) error
}

type validationService struct {
	validator       *validator.Validate
	htmlPolicy      *bluemonday.Policy
	securityService SecurityService
}

// NewValidationService creates a new validation service
func NewValidationService(securityService SecurityService) ValidationService {
	v := validator.New()
	
	// Register custom validators
	v.RegisterValidation("password", validatePasswordFunc)
	v.RegisterValidation("custom_alias", validateCustomAliasFunc)
	v.RegisterValidation("safe_string", validateSafeStringFunc)
	
	// Create HTML sanitization policy
	htmlPolicy := bluemonday.UGCPolicy()
	htmlPolicy.AllowElements() // Strip all HTML tags
	
	return &validationService{
		validator:       v,
		htmlPolicy:      htmlPolicy,
		securityService: securityService,
	}
}

// ValidateAndSanitizeURL validates and sanitizes a URL with comprehensive security checks
func (s *validationService) ValidateAndSanitizeURL(ctx context.Context, rawURL string) (*domain.URLValidationResult, error) {
	result := &domain.URLValidationResult{
		OriginalURL:    rawURL,
		IsValid:       true,
		Issues:        []string{},
		Recommendations: []string{},
	}
	
	// Basic URL validation
	if rawURL == "" {
		result.IsValid = false
		result.Issues = append(result.Issues, "URL cannot be empty")
		return result, nil
	}
	
	// Sanitize URL
	sanitizedURL := s.SanitizeString(rawURL)
	result.SanitizedURL = sanitizedURL
	
	// Parse URL
	parsedURL, err := url.Parse(sanitizedURL)
	if err != nil {
		result.IsValid = false
		result.Issues = append(result.Issues, fmt.Sprintf("Invalid URL format: %v", err))
		return result, nil
	}
	
	// URL scheme validation
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		result.IsValid = false
		result.Issues = append(result.Issues, "Only HTTP and HTTPS URLs are allowed")
		return result, nil
	}
	
	// Host validation
	if parsedURL.Host == "" {
		result.IsValid = false
		result.Issues = append(result.Issues, "URL must have a valid hostname")
		return result, nil
	}
	
	// Check for suspicious patterns
	if s.hasSuspiciousPatterns(sanitizedURL) {
		result.Issues = append(result.Issues, "URL contains suspicious patterns")
		result.Recommendations = append(result.Recommendations, "Review URL for potential security risks")
	}
	
	// Security validation using SecurityService
	if s.securityService != nil {
		securityReport, err := s.securityService.ValidateURL(ctx, sanitizedURL)
		if err == nil {
			result.SecurityReport = securityReport
			if !securityReport.IsSecure {
				result.Issues = append(result.Issues, "URL failed security validation")
				result.Recommendations = append(result.Recommendations, securityReport.Recommendations...)
			}
		}
	}
	
	// Final validation
	if len(result.Issues) > 0 {
		result.IsValid = false
	}
	
	return result, nil
}

// ValidateUserInput validates user input using struct tags
func (s *validationService) ValidateUserInput(ctx context.Context, input interface{}) error {
	return s.validator.Struct(input)
}

// SanitizeString removes potentially dangerous characters and patterns
func (s *validationService) SanitizeString(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	// Remove control characters except common whitespace
	var result strings.Builder
	for _, r := range input {
		if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
			continue
		}
		result.WriteRune(r)
	}
	
	// Trim whitespace
	return strings.TrimSpace(result.String())
}

// ValidatePassword validates password strength and requirements
func (s *validationService) ValidatePassword(password string) *domain.PasswordValidation {
	validation := &domain.PasswordValidation{
		IsValid:   true,
		Strength:  domain.PasswordStrengthWeak,
		Issues:    []string{},
		Suggestions: []string{},
	}
	
	// Minimum length check
	if len(password) < 8 {
		validation.IsValid = false
		validation.Issues = append(validation.Issues, "Password must be at least 8 characters long")
		validation.Suggestions = append(validation.Suggestions, "Use at least 8 characters")
	}
	
	// Maximum length check
	if len(password) > 128 {
		validation.IsValid = false
		validation.Issues = append(validation.Issues, "Password must not exceed 128 characters")
	}
	
	// Character requirements
	var (
		hasUpper   = false
		hasLower   = false
		hasDigit   = false
		hasSpecial = false
	)
	
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		validation.Issues = append(validation.Issues, "Password must contain at least one uppercase letter")
		validation.Suggestions = append(validation.Suggestions, "Add uppercase letters (A-Z)")
	}
	
	if !hasLower {
		validation.Issues = append(validation.Issues, "Password must contain at least one lowercase letter")
		validation.Suggestions = append(validation.Suggestions, "Add lowercase letters (a-z)")
	}
	
	if !hasDigit {
		validation.Issues = append(validation.Issues, "Password must contain at least one digit")
		validation.Suggestions = append(validation.Suggestions, "Add numbers (0-9)")
	}
	
	if !hasSpecial {
		validation.Issues = append(validation.Issues, "Password must contain at least one special character")
		validation.Suggestions = append(validation.Suggestions, "Add special characters (!@#$%^&*)")
	}
	
	// Common password patterns
	commonPatterns := []string{
		"password", "123456", "qwerty", "abc123", "admin", "user",
		"pass", "login", "welcome", "monkey", "dragon",
	}
	
	lowerPassword := strings.ToLower(password)
	for _, pattern := range commonPatterns {
		if strings.Contains(lowerPassword, pattern) {
			validation.Issues = append(validation.Issues, "Password contains common patterns")
			validation.Suggestions = append(validation.Suggestions, "Avoid common words and patterns")
			break
		}
	}
	
	// Calculate strength
	if len(validation.Issues) == 0 {
		strengthScore := 0
		if len(password) >= 12 {
			strengthScore++
		}
		if hasUpper && hasLower && hasDigit && hasSpecial {
			strengthScore++
		}
		if len(password) >= 16 {
			strengthScore++
		}
		
		switch strengthScore {
		case 0, 1:
			validation.Strength = domain.PasswordStrengthWeak
		case 2:
			validation.Strength = domain.PasswordStrengthMedium
		case 3:
			validation.Strength = domain.PasswordStrengthStrong
		}
	}
	
	return validation
}

// ValidateEmail validates email format and checks for suspicious patterns
func (s *validationService) ValidateEmail(email string) error {
	email = s.SanitizeString(email)
	
	// Basic format validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	
	// Length validation
	if len(email) > 254 {
		return fmt.Errorf("email address too long")
	}
	
	// Check for suspicious patterns
	suspiciousPatterns := []string{
		"script", "javascript", "vbscript", "onload", "onerror",
		"<", ">", "\"", "'", "&",
	}
	
	lowerEmail := strings.ToLower(email)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerEmail, pattern) {
			return fmt.Errorf("email contains suspicious characters")
		}
	}
	
	return nil
}

// ValidateAlias validates custom alias format and security
func (s *validationService) ValidateAlias(alias string) error {
	alias = s.SanitizeString(alias)
	
	// Length validation
	if len(alias) < 3 {
		return fmt.Errorf("alias must be at least 3 characters long")
	}
	
	if len(alias) > 50 {
		return fmt.Errorf("alias must not exceed 50 characters")
	}
	
	// Character validation - only alphanumeric, hyphens, and underscores
	aliasRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !aliasRegex.MatchString(alias) {
		return fmt.Errorf("alias can only contain letters, numbers, hyphens, and underscores")
	}
	
	// Must start and end with alphanumeric
	if !regexp.MustCompile(`^[a-zA-Z0-9]`).MatchString(alias) ||
		!regexp.MustCompile(`[a-zA-Z0-9]$`).MatchString(alias) {
		return fmt.Errorf("alias must start and end with a letter or number")
	}
	
	// Reserved words
	reservedWords := []string{
		"api", "admin", "www", "ftp", "mail", "email", "test",
		"app", "application", "service", "system", "user", "auth",
		"login", "register", "dashboard", "profile", "settings",
		"help", "support", "contact", "about", "terms", "privacy",
		"health", "status", "ping", "robots", "sitemap",
	}
	
	lowerAlias := strings.ToLower(alias)
	for _, word := range reservedWords {
		if lowerAlias == word {
			return fmt.Errorf("alias '%s' is reserved", alias)
		}
	}
	
	return nil
}

// SanitizeHTML removes or escapes HTML content
func (s *validationService) SanitizeHTML(input string) string {
	return s.htmlPolicy.Sanitize(input)
}

// ValidateJSONInput validates JSON input for security issues
func (s *validationService) ValidateJSONInput(input map[string]interface{}) error {
	return s.validateMapRecursively(input, 0, 10) // Max depth of 10
}

// Helper functions

func (s *validationService) hasSuspiciousPatterns(input string) bool {
	suspiciousPatterns := []string{
		"javascript:", "vbscript:", "data:", "blob:",
		"<script", "</script>", "onload=", "onerror=",
		"eval(", "setTimeout(", "setInterval(",
		"document.cookie", "window.location",
		"innerHTML", "outerHTML",
	}
	
	lowerInput := strings.ToLower(input)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}
	
	return false
}

func (s *validationService) validateMapRecursively(data map[string]interface{}, depth, maxDepth int) error {
	if depth > maxDepth {
		return fmt.Errorf("JSON depth exceeds maximum allowed depth of %d", maxDepth)
	}
	
	for key, value := range data {
		// Validate key
		if len(key) > 100 {
			return fmt.Errorf("JSON key too long: %s", key)
		}
		
		// Check for suspicious key patterns
		if s.hasSuspiciousPatterns(key) {
			return fmt.Errorf("suspicious pattern in JSON key: %s", key)
		}
		
		// Validate value based on type
		switch v := value.(type) {
		case string:
			if len(v) > 10000 { // 10KB limit for string values
				return fmt.Errorf("string value too long for key: %s", key)
			}
			if s.hasSuspiciousPatterns(v) {
				return fmt.Errorf("suspicious pattern in value for key: %s", key)
			}
		case map[string]interface{}:
			if err := s.validateMapRecursively(v, depth+1, maxDepth); err != nil {
				return err
			}
		case []interface{}:
			if len(v) > 1000 { // Limit array size
				return fmt.Errorf("array too large for key: %s", key)
			}
			for _, item := range v {
				if nestedMap, ok := item.(map[string]interface{}); ok {
					if err := s.validateMapRecursively(nestedMap, depth+1, maxDepth); err != nil {
						return err
					}
				}
			}
		}
	}
	
	return nil
}

// Custom validator functions

func validatePasswordFunc(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	// Basic requirements
	if len(password) < 8 {
		return false
	}
	
	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}
	
	return hasUpper && hasLower && hasDigit
}

func validateCustomAliasFunc(fl validator.FieldLevel) bool {
	alias := fl.Field().String()
	
	if len(alias) < 3 || len(alias) > 50 {
		return false
	}
	
	// Only alphanumeric, hyphens, and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, alias)
	return matched
}

func validateSafeStringFunc(fl validator.FieldLevel) bool {
	input := fl.Field().String()
	
	// Check for suspicious patterns
	suspiciousPatterns := []string{
		"<script", "</script>", "javascript:", "vbscript:",
		"onload=", "onerror=", "eval(", "innerHTML",
	}
	
	lowerInput := strings.ToLower(input)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerInput, pattern) {
			return false
		}
	}
	
	return true
}