// User types
export interface User {
  id: string
  email: string
  name: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

// Authentication request types
export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  name: string
  email: string
  password: string
  confirmPassword: string
}

export interface RefreshTokenRequest {
  refreshToken: string
}

export interface ChangePasswordRequest {
  currentPassword: string
  newPassword: string
  confirmPassword: string
}

export interface UpdateProfileRequest {
  name: string
  email: string
}

// Authentication response types
export interface AuthResponse {
  user: User
  accessToken: string
  refreshToken: string
  expiresAt: string
}

export interface TokenResponse {
  accessToken: string
  refreshToken: string
  expiresAt: string
}

// Authentication state types
export interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
}

// Authentication context types
export interface AuthContextType {
  // State
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null

  // Actions
  login: (credentials: LoginRequest) => Promise<void>
  register: (userData: RegisterRequest) => Promise<void>
  logout: () => Promise<void>
  refreshAuth: () => Promise<void>
  updateProfile: (data: UpdateProfileRequest) => Promise<void>
  changePassword: (data: ChangePasswordRequest) => Promise<void>
  clearError: () => void
  validateToken: () => Promise<boolean>
}

// Token management types
export interface TokenData {
  accessToken: string
  refreshToken: string
  expiresAt: string
}

// API Error types
export interface ApiError {
  message: string
  status: number
  code?: string
}

// Validation types
export interface ValidationError {
  field: string
  message: string
}

export interface AuthValidationErrors {
  email?: string
  password?: string
  name?: string
  confirmPassword?: string
  currentPassword?: string
  newPassword?: string
}