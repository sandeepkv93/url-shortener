package services

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type notificationService struct {
	configRepo ports.ConfigService
	cacheRepo  ports.CacheService
	smtpHost   string
	smtpPort   int
	smtpUser   string
	smtpPass   string
	fromEmail  string
	fromName   string
}

func NewNotificationService(
	configRepo ports.ConfigService,
	cacheRepo ports.CacheService,
	smtpHost string,
	smtpPort int,
	smtpUser string,
	smtpPass string,
	fromEmail string,
	fromName string,
) ports.NotificationService {
	return &notificationService{
		configRepo: configRepo,
		cacheRepo:  cacheRepo,
		smtpHost:   smtpHost,
		smtpPort:   smtpPort,
		smtpUser:   smtpUser,
		smtpPass:   smtpPass,
		fromEmail:  fromEmail,
		fromName:   fromName,
	}
}

func (s *notificationService) SendWelcomeEmail(ctx context.Context, user *domain.User) error {
	subject := "Welcome to URL Shortener"
	body := s.buildWelcomeEmailBody(user)
	
	return s.sendEmail(ctx, user.Email, subject, body)
}

func (s *notificationService) SendPasswordResetEmail(ctx context.Context, user *domain.User, resetToken string) error {
	subject := "Password Reset Request"
	body := s.buildPasswordResetEmailBody(user, resetToken)
	
	return s.sendEmail(ctx, user.Email, subject, body)
}

func (s *notificationService) SendPasswordChangedNotification(ctx context.Context, user *domain.User) error {
	subject := "Password Changed Successfully"
	body := s.buildPasswordChangedEmailBody(user)
	
	return s.sendEmail(ctx, user.Email, subject, body)
}

func (s *notificationService) SendAnalyticsDigest(ctx context.Context, user *domain.User, digest *domain.AnalyticsDigest) error {
	subject := fmt.Sprintf("Your URL Analytics Digest - %s", digest.Period)
	body := s.buildAnalyticsDigestEmailBody(user, digest)
	
	return s.sendEmail(ctx, user.Email, subject, body)
}

func (s *notificationService) SendClickAlert(ctx context.Context, user *domain.User, alert *domain.ClickAlert) error {
	subject := fmt.Sprintf("Click Alert: %s", alert.ShortCode)
	body := s.buildClickAlertEmailBody(user, alert)
	
	return s.sendEmail(ctx, user.Email, subject, body)
}

func (s *notificationService) SendMaintenanceNotification(ctx context.Context, users []*domain.User, message string) error {
	subject := "Scheduled Maintenance Notification"
	body := s.buildMaintenanceEmailBody(message)
	
	var errors []string
	for _, user := range users {
		if err := s.sendEmail(ctx, user.Email, subject, body); err != nil {
			errors = append(errors, fmt.Sprintf("failed to send to %s: %v", user.Email, err))
		}
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("failed to send to some users: %s", strings.Join(errors, "; "))
	}
	
	return nil
}

func (s *notificationService) SendSecurityAlert(ctx context.Context, user *domain.User, alert *domain.SecurityAlert) error {
	subject := fmt.Sprintf("Security Alert: %s", alert.Type)
	body := s.buildSecurityAlertEmailBody(user, alert)
	
	return s.sendEmail(ctx, user.Email, subject, body)
}

func (s *notificationService) sendEmail(ctx context.Context, to, subject, body string) error {
	// Check if email notifications are enabled
	if !s.configRepo.GetFeatureEmailNotificationsEnabled() {
		return fmt.Errorf("email notifications are disabled")
	}
	
	// Rate limiting: check if we've sent too many emails to this address recently
	rateLimitKey := fmt.Sprintf("email_rate_limit:%s", to)
	if s.isRateLimited(ctx, rateLimitKey) {
		return fmt.Errorf("rate limit exceeded for email address: %s", to)
	}
	
	// Prepare email message
	msg := s.buildEmailMessage(to, subject, body)
	
	// Send email via SMTP
	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPass, s.smtpHost)
	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
	
	err := smtp.SendMail(addr, auth, s.fromEmail, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	
	// Update rate limit
	s.updateRateLimit(ctx, rateLimitKey)
	
	return nil
}

func (s *notificationService) buildEmailMessage(to, subject, body string) string {
	return fmt.Sprintf(
		"From: %s <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		s.fromName, s.fromEmail, to, subject, body,
	)
}

func (s *notificationService) buildWelcomeEmailBody(user *domain.User) string {
	baseURL := s.configRepo.GetBaseURL()
	
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to URL Shortener</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h1 style="color: #333;">Welcome to URL Shortener!</h1>
    <p>Hi %s,</p>
    <p>Thank you for signing up for URL Shortener. Your account has been successfully created.</p>
    <p>You can now start shortening URLs, tracking analytics, and managing your links.</p>
    <p><a href="%s" style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Get Started</a></p>
    <p>If you have any questions, please don't hesitate to contact our support team.</p>
    <p>Best regards,<br>The URL Shortener Team</p>
</body>
</html>
`, user.FirstName, baseURL)
}

func (s *notificationService) buildPasswordResetEmailBody(user *domain.User, resetToken string) string {
	baseURL := s.configRepo.GetBaseURL()
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, resetToken)
	
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Reset Request</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h1 style="color: #333;">Password Reset Request</h1>
    <p>Hi %s,</p>
    <p>You requested a password reset for your URL Shortener account.</p>
    <p>Click the link below to reset your password:</p>
    <p><a href="%s" style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Reset Password</a></p>
    <p>This link will expire in 1 hour for security reasons.</p>
    <p>If you didn't request this password reset, please ignore this email.</p>
    <p>Best regards,<br>The URL Shortener Team</p>
</body>
</html>
`, user.FirstName, resetURL)
}

func (s *notificationService) buildPasswordChangedEmailBody(user *domain.User) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Changed Successfully</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h1 style="color: #333;">Password Changed Successfully</h1>
    <p>Hi %s,</p>
    <p>Your password has been successfully changed.</p>
    <p>If you didn't make this change, please contact our support team immediately.</p>
    <p>Best regards,<br>The URL Shortener Team</p>
</body>
</html>
`, user.FirstName)
}

func (s *notificationService) buildAnalyticsDigestEmailBody(user *domain.User, digest *domain.AnalyticsDigest) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Analytics Digest</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h1 style="color: #333;">Your Analytics Digest</h1>
    <p>Hi %s,</p>
    <p>Here's your URL analytics summary for %s:</p>
    <div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 20px 0;">
        <h3 style="margin-top: 0;">Summary Statistics</h3>
        <ul>
            <li>Total URLs: %d</li>
            <li>Total Clicks: %d</li>
            <li>Unique Visitors: %d</li>
            <li>Top Country: %s</li>
        </ul>
    </div>
    <p>View detailed analytics in your dashboard.</p>
    <p>Best regards,<br>The URL Shortener Team</p>
</body>
</html>
`, user.FirstName, digest.Period, digest.TotalURLs, digest.TotalClicks, digest.UniqueVisitors, digest.TopCountry)
}

func (s *notificationService) buildClickAlertEmailBody(user *domain.User, alert *domain.ClickAlert) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Click Alert</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h1 style="color: #333;">Click Alert</h1>
    <p>Hi %s,</p>
    <p>Your URL <strong>%s</strong> has received %d clicks in the last %s.</p>
    <div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 20px 0;">
        <h3 style="margin-top: 0;">Alert Details</h3>
        <ul>
            <li>Short Code: %s</li>
            <li>Original URL: %s</li>
            <li>Clicks: %d</li>
            <li>Threshold: %d</li>
            <li>Time Period: %s</li>
        </ul>
    </div>
    <p>View detailed analytics in your dashboard.</p>
    <p>Best regards,<br>The URL Shortener Team</p>
</body>
</html>
`, user.FirstName, alert.ShortCode, alert.ClickCount, alert.TimePeriod, alert.ShortCode, alert.OriginalURL, alert.ClickCount, alert.Threshold, alert.TimePeriod)
}

func (s *notificationService) buildMaintenanceEmailBody(message string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Scheduled Maintenance</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h1 style="color: #333;">Scheduled Maintenance Notification</h1>
    <p>Dear User,</p>
    <div style="background-color: #fff3cd; padding: 15px; border-radius: 5px; margin: 20px 0; border-left: 4px solid #ffc107;">
        <p>%s</p>
    </div>
    <p>We apologize for any inconvenience this may cause and appreciate your patience.</p>
    <p>Best regards,<br>The URL Shortener Team</p>
</body>
</html>
`, message)
}

func (s *notificationService) buildSecurityAlertEmailBody(user *domain.User, alert *domain.SecurityAlert) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Security Alert</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h1 style="color: #dc3545;">Security Alert</h1>
    <p>Hi %s,</p>
    <p>We detected suspicious activity on your account:</p>
    <div style="background-color: #f8d7da; padding: 15px; border-radius: 5px; margin: 20px 0; border-left: 4px solid #dc3545;">
        <h3 style="margin-top: 0;">Alert Details</h3>
        <ul>
            <li>Type: %s</li>
            <li>Description: %s</li>
            <li>IP Address: %s</li>
            <li>Location: %s</li>
            <li>Time: %s</li>
        </ul>
    </div>
    <p>If this was not you, please change your password immediately and contact our support team.</p>
    <p>Best regards,<br>The URL Shortener Team</p>
</body>
</html>
`, user.FirstName, alert.Type, alert.Description, alert.IPAddress, alert.Location, alert.Timestamp.Format("2006-01-02 15:04:05"))
}

func (s *notificationService) isRateLimited(ctx context.Context, rateLimitKey string) bool {
	// Check if we've sent more than 10 emails to this address in the last hour
	count, err := s.cacheRepo.GetCounter(ctx, rateLimitKey)
	if err != nil {
		return false // Allow if we can't check the rate limit
	}
	
	return count >= 10
}

func (s *notificationService) updateRateLimit(ctx context.Context, rateLimitKey string) {
	// Increment counter with 1 hour expiration
	s.cacheRepo.IncrementCounter(ctx, rateLimitKey, 1, time.Hour)
}

// Helper functions for template rendering
func (s *notificationService) renderTemplate(template string, data interface{}) string {
	// This is a simplified template rendering
	// In a real application, you'd use a proper template engine
	return template
}

func (s *notificationService) validateEmailAddress(email string) error {
	if email == "" {
		return fmt.Errorf("email address cannot be empty")
	}
	
	// Basic email validation
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email address format")
	}
	
	return nil
}

// SendBulkNotification sends notifications to multiple users
func (s *notificationService) SendBulkNotification(ctx context.Context, users []*domain.User, subject, template string, data interface{}) error {
	body := s.renderTemplate(template, data)
	
	var errors []string
	for _, user := range users {
		if err := s.sendEmail(ctx, user.Email, subject, body); err != nil {
			errors = append(errors, fmt.Sprintf("failed to send to %s: %v", user.Email, err))
		}
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("failed to send to some users: %s", strings.Join(errors, "; "))
	}
	
	return nil
}

// SendTestEmail sends a test email to verify email configuration
func (s *notificationService) SendTestEmail(ctx context.Context, toEmail string) error {
	subject := "Test Email from URL Shortener"
	body := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Test Email</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h1 style="color: #333;">Test Email</h1>
    <p>This is a test email from URL Shortener.</p>
    <p>If you received this email, the email configuration is working correctly.</p>
    <p>Best regards,<br>The URL Shortener Team</p>
</body>
</html>
`
	
	return s.sendEmail(ctx, toEmail, subject, body)
}

// GetEmailStatistics returns email sending statistics
func (s *notificationService) GetEmailStatistics(ctx context.Context) (*domain.EmailStatistics, error) {
	// This would typically come from a database or cache
	// For now, return mock statistics
	return &domain.EmailStatistics{
		TotalSent:    1000,
		TotalFailed:  10,
		TotalBounced: 5,
		DeliveryRate: 98.5,
		LastSent:     time.Now().Add(-time.Hour),
	}, nil
}