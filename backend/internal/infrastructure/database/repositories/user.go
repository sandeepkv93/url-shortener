package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Exists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}

func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users with pagination
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

func (r *userRepository) GetUserStats(ctx context.Context, userID uint) (*domain.UserStats, error) {
	var stats domain.UserStats
	
	// Get total URLs count
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalURLs).Error; err != nil {
		return nil, fmt.Errorf("failed to count user URLs: %w", err)
	}

	// Get active URLs count
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Where("user_id = ? AND is_active = ? AND (expires_at IS NULL OR expires_at > ?)", userID, true, now).
		Count(&stats.ActiveURLs).Error; err != nil {
		return nil, fmt.Errorf("failed to count active URLs: %w", err)
	}

	// Get expired URLs count
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Where("user_id = ? AND expires_at IS NOT NULL AND expires_at <= ?", userID, now).
		Count(&stats.ExpiredURLs).Error; err != nil {
		return nil, fmt.Errorf("failed to count expired URLs: %w", err)
	}

	// Get total clicks count
	if err := r.db.WithContext(ctx).
		Table("clicks").
		Joins("JOIN short_urls ON clicks.short_url_id = short_urls.id").
		Where("short_urls.user_id = ?", userID).
		Count(&stats.TotalClicks).Error; err != nil {
		return nil, fmt.Errorf("failed to count user clicks: %w", err)
	}

	// Get top performing URL
	var topURL struct {
		ShortCode string
	}
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Select("short_code").
		Where("user_id = ?", userID).
		Order("click_count DESC").
		First(&topURL).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get top performing URL: %w", err)
	}
	stats.TopPerformingURL = topURL.ShortCode

	// Calculate account age in days
	var user domain.User
	if err := r.db.WithContext(ctx).Select("created_at").First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("failed to get user creation date: %w", err)
	}
	
	// Use raw SQL for date calculation to be database-agnostic
	var accountAge int64
	if err := r.db.WithContext(ctx).
		Raw("SELECT EXTRACT(DAY FROM NOW() - ?) as age", user.CreatedAt).
		Scan(&accountAge).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate account age: %w", err)
	}
	stats.AccountAge = accountAge

	return &stats, nil
}

// Helper function to check for duplicate key errors
func isDuplicateKeyError(err error) bool {
	// This is a simplified check. In a production environment,
	// you'd want to check for specific database error codes
	return err != nil && (
		containsString(err.Error(), "duplicate") ||
		containsString(err.Error(), "unique") ||
		containsString(err.Error(), "UNIQUE"))
}

func containsString(str, substr string) bool {
	return len(str) >= len(substr) && 
		   (str == substr || 
		    (len(str) > len(substr) && 
		     (hasSubstring(str, substr))))
}

func hasSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}