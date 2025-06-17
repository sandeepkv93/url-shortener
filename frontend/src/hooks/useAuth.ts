import { useContext } from 'react'
import { AuthContext } from '@/context/AuthContext'
import { AuthContextType } from '@/types/auth'

/**
 * Hook to access the authentication context
 * 
 * @returns The authentication context including user state and auth methods
 * @throws Error if used outside of AuthProvider
 */
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext)
  
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  
  return context
}

export default useAuth