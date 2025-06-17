package middleware

import (
	"context"
	"net/http"
	"strings"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type AuthMiddleware struct {
	jwtService ports.JWTService
	userRepo   ports.UserRepository
}

func NewAuthMiddleware(jwtService ports.JWTService, userRepo ports.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

// RequireAuth middleware validates JWT token and sets user context
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token == "" {
			m.writeErrorResponse(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			m.writeErrorResponse(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Get user from database to ensure they still exist and are active
		user, err := m.userRepo.GetByID(r.Context(), claims.UserID)
		if err != nil {
			m.writeErrorResponse(w, "User not found", http.StatusUnauthorized)
			return
		}

		if !user.IsActive {
			m.writeErrorResponse(w, "User account is inactive", http.StatusUnauthorized)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "user_id", user.ID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth middleware validates JWT token if present but doesn't require it
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Get user from database
		user, err := m.userRepo.GetByID(r.Context(), claims.UserID)
		if err != nil || !user.IsActive {
			next.ServeHTTP(w, r)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "user_id", user.ID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnly middleware ensures the user is an admin
func (m *AuthMiddleware) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil {
			m.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Check if user has admin privileges (assuming we add this field later)
		// For now, we'll skip this check since admin functionality isn't in the domain model yet
		
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) extractToken(r *http.Request) string {
	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// Check query parameter as fallback
	return r.URL.Query().Get("token")
}

func (m *AuthMiddleware) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error": "` + message + `"}`))
}

// Helper functions to get user from context
func GetUserFromContext(ctx context.Context) *domain.User {
	if user, ok := ctx.Value("user").(*domain.User); ok {
		return user
	}
	return nil
}

func GetUserIDFromContext(ctx context.Context) uint {
	if userID, ok := ctx.Value("user_id").(uint); ok {
		return userID
	}
	return 0
}

func IsAuthenticated(ctx context.Context) bool {
	return GetUserFromContext(ctx) != nil
}