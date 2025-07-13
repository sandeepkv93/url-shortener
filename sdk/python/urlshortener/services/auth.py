"""Authentication service for URL Shortener API."""

from typing import Optional

from ..http_client import HTTPClient, AsyncHTTPClient
from ..types.auth import (
    AuthCredentials,
    AuthTokens,
    AuthResponse,
    RegisterRequest,
    LoginRequest,
    UpdateProfileRequest,
    ChangePasswordRequest,
    TokenValidationResponse,
    RefreshTokenRequest,
    PasswordResetRequest,
    PasswordResetConfirmRequest,
    EmailVerificationRequest,
    ResendVerificationRequest,
    User
)
from ..types.common import RequestConfig


class AuthService:
    """Synchronous authentication service."""

    def __init__(self, http_client: HTTPClient):
        """Initialize auth service.
        
        Args:
            http_client: HTTP client instance
        """
        self.http_client = http_client

    def register(
        self,
        email: str,
        password: str,
        first_name: Optional[str] = None,
        last_name: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> AuthResponse:
        """Register a new user account.
        
        Args:
            email: User email address
            password: User password
            first_name: User first name
            last_name: User last name
            config: Request configuration
            
        Returns:
            AuthResponse: Authentication response with user and tokens
        """
        request = RegisterRequest(
            email=email,
            password=password,
            first_name=first_name,
            last_name=last_name
        )
        
        response = self.http_client.post(
            '/auth/register',
            data=request.dict(),
            config=config
        )
        
        return AuthResponse(**response['data'])

    def login(
        self,
        email: str,
        password: str,
        remember_me: bool = False,
        config: Optional[RequestConfig] = None
    ) -> AuthResponse:
        """Login with email and password.
        
        Args:
            email: User email address
            password: User password
            remember_me: Extended session duration
            config: Request configuration
            
        Returns:
            AuthResponse: Authentication response with user and tokens
        """
        request = LoginRequest(
            email=email,
            password=password,
            remember_me=remember_me
        )
        
        response = self.http_client.post(
            '/auth/login',
            data=request.dict(),
            config=config
        )
        
        auth_response = AuthResponse(**response['data'])
        
        # Update client with new access token
        self.http_client.set_access_token(auth_response.access_token)
        
        return auth_response

    def logout(self, config: Optional[RequestConfig] = None) -> bool:
        """Logout current user.
        
        Args:
            config: Request configuration
            
        Returns:
            bool: True if logout successful
        """
        response = self.http_client.post('/auth/logout', config=config)
        
        # Clear access token from client
        self.http_client.set_access_token(None)
        
        return response.get('success', False)

    def refresh_token(
        self,
        refresh_token: str,
        config: Optional[RequestConfig] = None
    ) -> AuthTokens:
        """Refresh access token.
        
        Args:
            refresh_token: Refresh token
            config: Request configuration
            
        Returns:
            AuthTokens: New authentication tokens
        """
        request = RefreshTokenRequest(refresh_token=refresh_token)
        
        response = self.http_client.post(
            '/auth/refresh',
            data=request.dict(),
            config=config
        )
        
        tokens = AuthTokens(**response['data'])
        
        # Update client with new access token
        self.http_client.set_access_token(tokens.access_token)
        
        return tokens

    def get_profile(self, config: Optional[RequestConfig] = None) -> User:
        """Get current user profile.
        
        Args:
            config: Request configuration
            
        Returns:
            User: User profile information
        """
        response = self.http_client.get('/auth/profile', config=config)
        return User(**response['data'])

    def update_profile(
        self,
        email: Optional[str] = None,
        first_name: Optional[str] = None,
        last_name: Optional[str] = None,
        timezone: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> User:
        """Update user profile.
        
        Args:
            email: New email address
            first_name: First name
            last_name: Last name
            timezone: User timezone
            config: Request configuration
            
        Returns:
            User: Updated user profile
        """
        request = UpdateProfileRequest(
            email=email,
            first_name=first_name,
            last_name=last_name,
            timezone=timezone
        )
        
        # Only include non-None fields
        data = {k: v for k, v in request.dict().items() if v is not None}
        
        response = self.http_client.put(
            '/auth/profile',
            data=data,
            config=config
        )
        
        return User(**response['data'])

    def change_password(
        self,
        current_password: str,
        new_password: str,
        config: Optional[RequestConfig] = None
    ) -> bool:
        """Change user password.
        
        Args:
            current_password: Current password
            new_password: New password
            config: Request configuration
            
        Returns:
            bool: True if password changed successfully
        """
        request = ChangePasswordRequest(
            current_password=current_password,
            new_password=new_password
        )
        
        response = self.http_client.post(
            '/auth/change-password',
            data=request.dict(),
            config=config
        )
        
        return response.get('success', False)

    def validate_token(
        self,
        token: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> TokenValidationResponse:
        """Validate access token.
        
        Args:
            token: Token to validate (uses current token if not provided)
            config: Request configuration
            
        Returns:
            TokenValidationResponse: Token validation result
        """
        params = {}
        if token:
            params['token'] = token
            
        response = self.http_client.get(
            '/auth/validate',
            params=params,
            config=config
        )
        
        return TokenValidationResponse(**response['data'])

    def request_password_reset(
        self,
        email: str,
        config: Optional[RequestConfig] = None
    ) -> bool:
        """Request password reset email.
        
        Args:
            email: User email address
            config: Request configuration
            
        Returns:
            bool: True if reset email sent
        """
        request = PasswordResetRequest(email=email)
        
        response = self.http_client.post(
            '/auth/password-reset',
            data=request.dict(),
            config=config
        )
        
        return response.get('success', False)

    def confirm_password_reset(
        self,
        token: str,
        new_password: str,
        config: Optional[RequestConfig] = None
    ) -> bool:
        """Confirm password reset with token.
        
        Args:
            token: Reset token
            new_password: New password
            config: Request configuration
            
        Returns:
            bool: True if password reset successful
        """
        request = PasswordResetConfirmRequest(
            token=token,
            new_password=new_password
        )
        
        response = self.http_client.post(
            '/auth/password-reset/confirm',
            data=request.dict(),
            config=config
        )
        
        return response.get('success', False)

    def verify_email(
        self,
        token: str,
        config: Optional[RequestConfig] = None
    ) -> bool:
        """Verify email address with token.
        
        Args:
            token: Verification token
            config: Request configuration
            
        Returns:
            bool: True if email verified successfully
        """
        request = EmailVerificationRequest(token=token)
        
        response = self.http_client.post(
            '/auth/verify-email',
            data=request.dict(),
            config=config
        )
        
        return response.get('success', False)

    def resend_verification(
        self,
        email: str,
        config: Optional[RequestConfig] = None
    ) -> bool:
        """Resend email verification.
        
        Args:
            email: User email address
            config: Request configuration
            
        Returns:
            bool: True if verification email sent
        """
        request = ResendVerificationRequest(email=email)
        
        response = self.http_client.post(
            '/auth/resend-verification',
            data=request.dict(),
            config=config
        )
        
        return response.get('success', False)


class AsyncAuthService:
    """Asynchronous authentication service."""

    def __init__(self, http_client: AsyncHTTPClient):
        """Initialize async auth service.
        
        Args:
            http_client: Async HTTP client instance
        """
        self.http_client = http_client

    async def register(
        self,
        email: str,
        password: str,
        first_name: Optional[str] = None,
        last_name: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> AuthResponse:
        """Register a new user account.
        
        Args:
            email: User email address
            password: User password
            first_name: User first name
            last_name: User last name
            config: Request configuration
            
        Returns:
            AuthResponse: Authentication response with user and tokens
        """
        request = RegisterRequest(
            email=email,
            password=password,
            first_name=first_name,
            last_name=last_name
        )
        
        response = await self.http_client.post(
            '/auth/register',
            data=request.dict(),
            config=config
        )
        
        return AuthResponse(**response['data'])

    async def login(
        self,
        email: str,
        password: str,
        remember_me: bool = False,
        config: Optional[RequestConfig] = None
    ) -> AuthResponse:
        """Login with email and password.
        
        Args:
            email: User email address
            password: User password
            remember_me: Extended session duration
            config: Request configuration
            
        Returns:
            AuthResponse: Authentication response with user and tokens
        """
        request = LoginRequest(
            email=email,
            password=password,
            remember_me=remember_me
        )
        
        response = await self.http_client.post(
            '/auth/login',
            data=request.dict(),
            config=config
        )
        
        auth_response = AuthResponse(**response['data'])
        
        # Update client with new access token
        self.http_client.set_access_token(auth_response.access_token)
        
        return auth_response

    async def logout(self, config: Optional[RequestConfig] = None) -> bool:
        """Logout current user.
        
        Args:
            config: Request configuration
            
        Returns:
            bool: True if logout successful
        """
        response = await self.http_client.post('/auth/logout', config=config)
        
        # Clear access token from client
        self.http_client.set_access_token(None)
        
        return response.get('success', False)

    async def refresh_token(
        self,
        refresh_token: str,
        config: Optional[RequestConfig] = None
    ) -> AuthTokens:
        """Refresh access token.
        
        Args:
            refresh_token: Refresh token
            config: Request configuration
            
        Returns:
            AuthTokens: New authentication tokens
        """
        request = RefreshTokenRequest(refresh_token=refresh_token)
        
        response = await self.http_client.post(
            '/auth/refresh',
            data=request.dict(),
            config=config
        )
        
        tokens = AuthTokens(**response['data'])
        
        # Update client with new access token
        self.http_client.set_access_token(tokens.access_token)
        
        return tokens

    async def get_profile(self, config: Optional[RequestConfig] = None) -> User:
        """Get current user profile.
        
        Args:
            config: Request configuration
            
        Returns:
            User: User profile information
        """
        response = await self.http_client.get('/auth/profile', config=config)
        return User(**response['data'])

    async def update_profile(
        self,
        email: Optional[str] = None,
        first_name: Optional[str] = None,
        last_name: Optional[str] = None,
        timezone: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> User:
        """Update user profile.
        
        Args:
            email: New email address
            first_name: First name
            last_name: Last name
            timezone: User timezone
            config: Request configuration
            
        Returns:
            User: Updated user profile
        """
        request = UpdateProfileRequest(
            email=email,
            first_name=first_name,
            last_name=last_name,
            timezone=timezone
        )
        
        # Only include non-None fields
        data = {k: v for k, v in request.dict().items() if v is not None}
        
        response = await self.http_client.put(
            '/auth/profile',
            data=data,
            config=config
        )
        
        return User(**response['data'])

    async def change_password(
        self,
        current_password: str,
        new_password: str,
        config: Optional[RequestConfig] = None
    ) -> bool:
        """Change user password.
        
        Args:
            current_password: Current password
            new_password: New password
            config: Request configuration
            
        Returns:
            bool: True if password changed successfully
        """
        request = ChangePasswordRequest(
            current_password=current_password,
            new_password=new_password
        )
        
        response = await self.http_client.post(
            '/auth/change-password',
            data=request.dict(),
            config=config
        )
        
        return response.get('success', False)

    async def validate_token(
        self,
        token: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> TokenValidationResponse:
        """Validate access token.
        
        Args:
            token: Token to validate (uses current token if not provided)
            config: Request configuration
            
        Returns:
            TokenValidationResponse: Token validation result
        """
        params = {}
        if token:
            params['token'] = token
            
        response = await self.http_client.get(
            '/auth/validate',
            params=params,
            config=config
        )
        
        return TokenValidationResponse(**response['data'])

    async def request_password_reset(
        self,
        email: str,
        config: Optional[RequestConfig] = None
    ) -> bool:
        """Request password reset email.
        
        Args:
            email: User email address
            config: Request configuration
            
        Returns:
            bool: True if reset email sent
        """
        request = PasswordResetRequest(email=email)
        
        response = await self.http_client.post(
            '/auth/password-reset',
            data=request.dict(),
            config=config
        )
        
        return response.get('success', False)

    async def confirm_password_reset(
        self,
        token: str,
        new_password: str,
        config: Optional[RequestConfig] = None
    ) -> bool:
        """Confirm password reset with token.
        
        Args:
            token: Reset token
            new_password: New password
            config: Request configuration
            
        Returns:
            bool: True if password reset successful
        """
        request = PasswordResetConfirmRequest(
            token=token,
            new_password=new_password
        )
        
        response = await self.http_client.post(
            '/auth/password-reset/confirm',
            data=request.dict(),
            config=config
        )
        
        return response.get('success', False)

    async def verify_email(
        self,
        token: str,
        config: Optional[RequestConfig] = None
    ) -> bool:
        """Verify email address with token.
        
        Args:
            token: Verification token
            config: Request configuration
            
        Returns:
            bool: True if email verified successfully
        """
        request = EmailVerificationRequest(token=token)
        
        response = await self.http_client.post(
            '/auth/verify-email',
            data=request.dict(),
            config=config
        )
        
        return response.get('success', False)

    async def resend_verification(
        self,
        email: str,
        config: Optional[RequestConfig] = None
    ) -> bool:
        """Resend email verification.
        
        Args:
            email: User email address
            config: Request configuration
            
        Returns:
            bool: True if verification email sent
        """
        request = ResendVerificationRequest(email=email)
        
        response = await self.http_client.post(
            '/auth/resend-verification',
            data=request.dict(),
            config=config
        )
        
        return response.get('success', False)