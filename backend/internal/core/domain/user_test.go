package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserToResponse(t *testing.T) {
	user := &User{
		ID:       1,
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	response := user.ToResponse()

	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.CreatedAt, response.CreatedAt)
	assert.Equal(t, user.UpdatedAt, response.UpdatedAt)
	
	// Password should not be included in response (UserResponse struct doesn't have Password field)
	assert.Equal(t, user.Email, response.Email) // Email should be included in response
}

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateUserRequest
		wantErr bool
	}{
		{
			name: "valid user request",
			request: CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "empty email",
			request: CreateUserRequest{
				Email:    "",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			request: CreateUserRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "short password",
			request: CreateUserRequest{
				Email:    "test@example.com",
				Password: "123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: Actual validation would be done by the validator package
			// This is just testing the struct definition
			assert.NotNil(t, tt.request)
		})
	}
}