"""Authentication related type definitions."""

from typing import Optional
from datetime import datetime
from pydantic import BaseModel, Field, EmailStr


class AuthCredentials(BaseModel):
    """User authentication credentials."""
    email: EmailStr = Field(description="User email address")
    password: str = Field(min_length=8, description="User password")


class AuthTokens(BaseModel):
    """Authentication tokens response."""
    access_token: str = Field(description="JWT access token")
    refresh_token: str = Field(description="JWT refresh token")
    expires_at: datetime = Field(description="Token expiration time")
    token_type: str = Field(default="Bearer", description="Token type")


class User(BaseModel):
    """User profile information."""
    id: int = Field(description="User ID")
    email: EmailStr = Field(description="User email address")
    created_at: datetime = Field(description="Account creation timestamp")
    updated_at: datetime = Field(description="Last profile update timestamp")
    
    # Optional profile fields
    first_name: Optional[str] = Field(None, description="User first name")
    last_name: Optional[str] = Field(None, description="User last name")
    avatar_url: Optional[str] = Field(None, description="User avatar URL")
    timezone: Optional[str] = Field(None, description="User timezone")
    is_verified: bool = Field(default=False, description="Email verification status")


class RegisterRequest(BaseModel):
    """User registration request."""
    email: EmailStr = Field(description="User email address")
    password: str = Field(min_length=8, description="User password")
    first_name: Optional[str] = Field(None, max_length=50, description="User first name")
    last_name: Optional[str] = Field(None, max_length=50, description="User last name")


class LoginRequest(BaseModel):
    """User login request."""
    email: EmailStr = Field(description="User email address")
    password: str = Field(description="User password")
    remember_me: bool = Field(default=False, description="Extended session duration")


class AuthResponse(BaseModel):
    """Authentication response with user and tokens."""
    user: User = Field(description="User profile")
    access_token: str = Field(description="JWT access token")
    refresh_token: str = Field(description="JWT refresh token")
    expires_at: datetime = Field(description="Token expiration time")


class UpdateProfileRequest(BaseModel):
    """User profile update request."""
    email: Optional[EmailStr] = Field(None, description="New email address")
    first_name: Optional[str] = Field(None, max_length=50, description="First name")
    last_name: Optional[str] = Field(None, max_length=50, description="Last name")
    timezone: Optional[str] = Field(None, description="User timezone")


class ChangePasswordRequest(BaseModel):
    """Password change request."""
    current_password: str = Field(description="Current password")
    new_password: str = Field(min_length=8, description="New password")


class TokenValidationResponse(BaseModel):
    """Token validation response."""
    valid: bool = Field(description="Whether token is valid")
    user: Optional[User] = Field(None, description="User profile if valid")
    expires_at: Optional[datetime] = Field(None, description="Token expiration time")


class RefreshTokenRequest(BaseModel):
    """Refresh token request."""
    refresh_token: str = Field(description="Refresh token")


class PasswordResetRequest(BaseModel):
    """Password reset request."""
    email: EmailStr = Field(description="User email address")


class PasswordResetConfirmRequest(BaseModel):
    """Password reset confirmation."""
    token: str = Field(description="Reset token")
    new_password: str = Field(min_length=8, description="New password")


class EmailVerificationRequest(BaseModel):
    """Email verification request."""
    token: str = Field(description="Verification token")


class ResendVerificationRequest(BaseModel):
    """Resend verification email request."""
    email: EmailStr = Field(description="User email address")