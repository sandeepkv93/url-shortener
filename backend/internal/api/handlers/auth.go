package handlers

import (
	"encoding/json"
	"net/http"

	"url-shortener/internal/api/middleware"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Register user
	response, err := h.authService.Register(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrUserAlreadyExists:
			h.writeErrorResponse(w, "User already exists", http.StatusConflict)
		case domain.ErrInvalidCredentials:
			h.writeErrorResponse(w, "Invalid credentials", http.StatusBadRequest)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, response, http.StatusCreated)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Authenticate user
	response, err := h.authService.Login(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			h.writeErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
		case domain.ErrUserNotFound:
			h.writeErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		h.writeErrorResponse(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	// Refresh token
	response, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			h.writeErrorResponse(w, "Invalid refresh token", http.StatusUnauthorized)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Logout user
	if err := h.authService.Logout(r.Context(), userID); err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, map[string]string{"message": "Logged out successfully"}, http.StatusOK)
}

// GetProfile handles getting user profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Get user profile
	profile, err := h.authService.GetProfile(r.Context(), user.ID)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, profile, http.StatusOK)
}

// UpdateProfile handles updating user profile
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req domain.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update profile
	profile, err := h.authService.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			h.writeErrorResponse(w, "User not found", http.StatusNotFound)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, profile, http.StatusOK)
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req domain.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Change password
	if err := h.authService.ChangePassword(r.Context(), userID, req); err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			h.writeErrorResponse(w, "Current password is incorrect", http.StatusBadRequest)
		case domain.ErrUserNotFound:
			h.writeErrorResponse(w, "User not found", http.StatusNotFound)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, map[string]string{"message": "Password changed successfully"}, http.StatusOK)
}

// ValidateToken handles token validation for clients
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.writeErrorResponse(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"valid":   true,
		"user_id": user.ID,
		"email":   user.Email,
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// Helper methods

func (h *AuthHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, write a simple error response
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to encode response"}`))
	}
}

func (h *AuthHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Fallback to simple string response
		w.Write([]byte(`{"error": "Internal server error"}`))
	}
}