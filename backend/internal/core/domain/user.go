package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Email       string         `json:"email" gorm:"uniqueIndex;not null"`
	Password    string         `json:"-" gorm:"not null"`
	FirstName   string         `json:"first_name" gorm:"not null"`
	LastName    string         `json:"last_name" gorm:"not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	LastLoginAt *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	ShortURLs []ShortURL `json:"short_urls,omitempty" gorm:"foreignKey:UserID"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateUserRequest struct {
	Email     string `json:"email,omitempty" validate:"omitempty,email"`
	FirstName string `json:"first_name,omitempty" validate:"omitempty,max=50"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,max=50"`
}

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordResetConfirmRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,max=50"`
	LastName  string `json:"last_name" validate:"required,max=50"`
}

type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    int           `json:"expires_in"`
}

type LoginResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

type UserStats struct {
	TotalURLs       int64 `json:"total_urls"`
	TotalClicks     int64 `json:"total_clicks"`
	ActiveURLs      int64 `json:"active_urls"`
	ExpiredURLs     int64 `json:"expired_urls"`
	TopPerformingURL string `json:"top_performing_url"`
	AccountAge      int64  `json:"account_age_days"`
}

type UserAnalytics struct {
	UserID       uint              `json:"user_id"`
	TotalURLs    int64             `json:"total_urls"`
	TotalClicks  int64             `json:"total_clicks"`
	ClicksByDate map[string]int64  `json:"clicks_by_date"`
	TopURLs      []TopURLStat      `json:"top_urls"`
}

type TopURLStat struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	ClickCount  int64  `json:"click_count"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}

// Validation methods
func (r *RegisterRequest) Validate() error {
	if r.Email == "" {
		return ErrInvalidEmail
	}
	if len(r.Password) < 8 {
		return ErrInvalidPassword
	}
	if r.FirstName == "" || r.LastName == "" {
		return ErrInvalidRequest
	}
	return nil
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return ErrInvalidEmail
	}
	if r.Password == "" {
		return ErrInvalidPassword
	}
	return nil
}

func (r *UpdateUserRequest) Validate() error {
	// Basic validation for update requests
	return nil
}

func (r *ChangePasswordRequest) Validate() error {
	if r.CurrentPassword == "" {
		return ErrInvalidPassword
	}
	if len(r.NewPassword) < 8 {
		return ErrInvalidPassword
	}
	return nil
}