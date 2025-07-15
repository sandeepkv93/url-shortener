package domain

import "fmt"

// URLValidationResult contains the result of URL validation and sanitization
type URLValidationResult struct {
	OriginalURL     string               `json:"original_url"`
	SanitizedURL    string               `json:"sanitized_url"`
	IsValid         bool                 `json:"is_valid"`
	Issues          []string             `json:"issues"`
	Recommendations []string             `json:"recommendations"`
	SecurityReport  interface{}          `json:"security_report,omitempty"` // Will use SecurityReport type when available
}

// PasswordValidation contains password strength validation results
type PasswordValidation struct {
	IsValid     bool             `json:"is_valid"`
	Strength    PasswordStrength `json:"strength"`
	Issues      []string         `json:"issues"`
	Suggestions []string         `json:"suggestions"`
}

// PasswordStrength represents password strength levels
type PasswordStrength string

const (
	PasswordStrengthWeak   PasswordStrength = "weak"
	PasswordStrengthMedium PasswordStrength = "medium"
	PasswordStrengthStrong PasswordStrength = "strong"
)

// ValidationError represents a validation error with detailed information
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Tag)
}

// InputValidationResult contains comprehensive input validation results
type InputValidationResult struct {
	IsValid bool              `json:"is_valid"`
	Errors  []ValidationError `json:"errors"`
}

// SecurityValidationRequest represents a request for security validation
type SecurityValidationRequest struct {
	URL        string            `json:"url" validate:"required,url"`
	Alias      string            `json:"alias,omitempty" validate:"omitempty,custom_alias"`
	Title      string            `json:"title,omitempty" validate:"omitempty,max=200,safe_string"`
	Description string           `json:"description,omitempty" validate:"omitempty,max=1000,safe_string"`
	Tags       []string          `json:"tags,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// Validate validates the security validation request
func (r *SecurityValidationRequest) Validate() error {
	if r.URL == "" {
		return fmt.Errorf("URL is required")
	}
	
	// Additional business logic validation
	if len(r.Tags) > 10 {
		return fmt.Errorf("too many tags (maximum 10)")
	}
	
	for _, tag := range r.Tags {
		if len(tag) > 50 {
			return fmt.Errorf("tag too long (maximum 50 characters)")
		}
	}
	
	if len(r.Metadata) > 20 {
		return fmt.Errorf("too many metadata fields (maximum 20)")
	}
	
	return nil
}

// URLSecurityReport contains comprehensive security analysis of a URL (defined in security.go)
// This is a reference to maintain consistency