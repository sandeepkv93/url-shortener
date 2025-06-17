import { createContext, useContext, ReactNode } from 'react'

interface AuthContextType {
  // Placeholder for auth context - will be implemented in Step 12
  isAuthenticated: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  // Placeholder implementation - will be completed in Step 12
  const value: AuthContextType = {
    isAuthenticated: false,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}