import { ReactNode } from 'react'
import Header from './Header'
import Footer from './Footer'

interface LayoutProps {
  children: ReactNode
  className?: string
}

const Layout = ({ children, className = '' }: LayoutProps) => {
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <Header />
      
      <main className={`flex-1 max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8 w-full ${className}`}>
        {children}
      </main>
      
      <Footer />
    </div>
  )
}

export default Layout