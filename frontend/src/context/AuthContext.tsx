import React, { createContext, useReducer, useEffect, ReactNode, useCallback } from 'react'
import {
  AuthContextType,
  AuthState,
  User,
  LoginRequest,
  RegisterRequest,
  ChangePasswordRequest,
  UpdateProfileRequest,
  ApiError,
} from '@/types/auth'
import { authService, tokenManager, setupTokenRefresh } from '@/services/auth'

// Auth actions
type AuthAction =
  | { type: 'AUTH_START' }
  | { type: 'AUTH_SUCCESS'; payload: { user: User } }
  | { type: 'AUTH_FAILURE'; payload: { error: string } }
  | { type: 'AUTH_LOGOUT' }
  | { type: 'UPDATE_USER'; payload: { user: User } }
  | { type: 'CLEAR_ERROR' }
  | { type: 'SET_LOADING'; payload: { isLoading: boolean } }

// Initial state
const initialState: AuthState = {
  user: null,
  accessToken: null,
  refreshToken: null,
  isAuthenticated: false,
  isLoading: true, // Start with loading true to check existing auth
  error: null,
}

// Auth reducer
const authReducer = (state: AuthState, action: AuthAction): AuthState => {
  switch (action.type) {
    case 'AUTH_START':
      return {
        ...state,
        isLoading: true,
        error: null,
      }

    case 'AUTH_SUCCESS':
      return {
        ...state,
        user: action.payload.user,
        isAuthenticated: true,
        isLoading: false,
        error: null,
      }

    case 'AUTH_FAILURE':
      return {
        ...state,
        user: null,
        accessToken: null,
        refreshToken: null,
        isAuthenticated: false,
        isLoading: false,
        error: action.payload.error,
      }

    case 'AUTH_LOGOUT':
      return {
        ...initialState,
        isLoading: false,
      }

    case 'UPDATE_USER':
      return {
        ...state,
        user: action.payload.user,
      }

    case 'CLEAR_ERROR':
      return {
        ...state,
        error: null,
      }

    case 'SET_LOADING':
      return {
        ...state,
        isLoading: action.payload.isLoading,
      }

    default:
      return state
  }
}

// Create context
export const AuthContext = createContext<AuthContextType | undefined>(undefined)

// AuthProvider component
interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [state, dispatch] = useReducer(authReducer, initialState)

  // Handle API errors
  const handleError = useCallback((error: unknown): string => {
    if (error && typeof error === 'object' && 'message' in error) {
      return (error as ApiError).message
    }
    return 'An unexpected error occurred'
  }, [])

  // Initialize auth state from stored tokens
  const initializeAuth = useCallback(async () => {
    try {
      const tokens = tokenManager.getTokens()
      const user = tokenManager.getUser()

      if (tokens && user && !tokenManager.isTokenExpired(tokens.expiresAt)) {
        // Validate token with server
        const isValid = await authService.validateToken()
        if (isValid) {
          dispatch({ type: 'AUTH_SUCCESS', payload: { user } })
        } else {
          // Token is invalid, clear storage
          tokenManager.clearTokens()
          tokenManager.clearUser()
          dispatch({ type: 'SET_LOADING', payload: { isLoading: false } })
        }
      } else {
        // No valid tokens found
        tokenManager.clearTokens()
        tokenManager.clearUser()
        dispatch({ type: 'SET_LOADING', payload: { isLoading: false } })
      }
    } catch (error) {
      console.error('Auth initialization error:', error)
      tokenManager.clearTokens()
      tokenManager.clearUser()
      dispatch({ type: 'SET_LOADING', payload: { isLoading: false } })
    }
  }, [])

  // Login function
  const login = useCallback(async (credentials: LoginRequest): Promise<void> => {
    dispatch({ type: 'AUTH_START' })
    try {
      const authData = await authService.login(credentials)
      dispatch({ type: 'AUTH_SUCCESS', payload: { user: authData.user } })
    } catch (error) {
      const errorMessage = handleError(error)
      dispatch({ type: 'AUTH_FAILURE', payload: { error: errorMessage } })
      throw error
    }
  }, [handleError])

  // Register function
  const register = useCallback(async (userData: RegisterRequest): Promise<void> => {
    dispatch({ type: 'AUTH_START' })
    try {
      const authData = await authService.register(userData)
      dispatch({ type: 'AUTH_SUCCESS', payload: { user: authData.user } })
    } catch (error) {
      const errorMessage = handleError(error)
      dispatch({ type: 'AUTH_FAILURE', payload: { error: errorMessage } })
      throw error
    }
  }, [handleError])

  // Logout function
  const logout = useCallback(async (): Promise<void> => {
    try {
      await authService.logout()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      dispatch({ type: 'AUTH_LOGOUT' })
    }
  }, [])

  // Refresh auth function
  const refreshAuth = useCallback(async (): Promise<void> => {
    try {
      await authService.refreshToken()
      // Get updated user profile
      const user = await authService.getProfile()
      dispatch({ type: 'UPDATE_USER', payload: { user } })
    } catch (error) {
      console.error('Auth refresh error:', error)
      dispatch({ type: 'AUTH_LOGOUT' })
      throw error
    }
  }, [])

  // Update profile function
  const updateProfile = useCallback(async (data: UpdateProfileRequest): Promise<void> => {
    try {
      const updatedUser = await authService.updateProfile(data)
      dispatch({ type: 'UPDATE_USER', payload: { user: updatedUser } })
    } catch (error) {
      const errorMessage = handleError(error)
      dispatch({ type: 'AUTH_FAILURE', payload: { error: errorMessage } })
      throw error
    }
  }, [handleError])

  // Change password function
  const changePassword = useCallback(async (data: ChangePasswordRequest): Promise<void> => {
    try {
      await authService.changePassword(data)
    } catch (error) {
      const errorMessage = handleError(error)
      dispatch({ type: 'AUTH_FAILURE', payload: { error: errorMessage } })
      throw error
    }
  }, [handleError])

  // Clear error function
  const clearError = useCallback((): void => {
    dispatch({ type: 'CLEAR_ERROR' })
  }, [])

  // Validate token function
  const validateToken = useCallback(async (): Promise<boolean> => {
    try {
      return await authService.validateToken()
    } catch {
      return false
    }
  }, [])

  // Set up automatic token refresh and auth state management
  useEffect(() => {
    let cleanup: (() => void) | undefined

    const onTokenRefreshed = () => {
      // Optionally refresh user data when token is refreshed
      console.log('Token refreshed successfully')
    }

    const onAuthExpired = () => {
      dispatch({ type: 'AUTH_LOGOUT' })
    }

    // Initialize auth state
    initializeAuth().then(() => {
      // Set up token refresh only after initialization
      if (state.isAuthenticated) {
        cleanup = setupTokenRefresh(onTokenRefreshed, onAuthExpired)
      }
    })

    return () => {
      if (cleanup) {
        cleanup()
      }
    }
  }, [initializeAuth, state.isAuthenticated])

  // Context value
  const contextValue: AuthContextType = {
    // State
    user: state.user,
    isAuthenticated: state.isAuthenticated,
    isLoading: state.isLoading,
    error: state.error,

    // Actions
    login,
    register,
    logout,
    refreshAuth,
    updateProfile,
    changePassword,
    clearError,
    validateToken,
  }

  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  )
}