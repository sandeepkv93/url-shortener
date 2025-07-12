import { Suspense, lazy } from 'react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { AuthProvider } from '@/context/AuthContext'
import Layout from '@/components/common/Layout'
import Loading from '@/components/common/Loading'
import ErrorBoundary from '@/components/common/ErrorBoundary'

// Lazy load pages for better performance
const Home = lazy(() => import('@/pages/Home'))
const Dashboard = lazy(() => import('@/pages/Dashboard'))
const Analytics = lazy(() => import('@/pages/Analytics'))
const Profile = lazy(() => import('@/pages/Profile'))
const NotFound = lazy(() => import('@/pages/NotFound'))
const ComponentDemo = lazy(() => import('@/pages/ComponentDemo'))

// Wrapper component for lazy-loaded routes with error boundary
const LazyRouteWrapper = ({ children }: { children: React.ReactNode }) => (
  <ErrorBoundary>
    <Layout>
      <Suspense fallback={<Loading />}>
        {children}
      </Suspense>
    </Layout>
  </ErrorBoundary>
)

// Standalone wrapper for demo page (no layout)
const DemoWrapper = ({ children }: { children: React.ReactNode }) => (
  <ErrorBoundary>
    <Suspense fallback={<Loading />}>
      {children}
    </Suspense>
  </ErrorBoundary>
)

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route 
            path="/demo" 
            element={
              <DemoWrapper>
                <ComponentDemo />
              </DemoWrapper>
            } 
          />
          <Route 
            path="/" 
            element={
              <LazyRouteWrapper>
                <Home />
              </LazyRouteWrapper>
            } 
          />
          <Route 
            path="/dashboard" 
            element={
              <LazyRouteWrapper>
                <Dashboard />
              </LazyRouteWrapper>
            } 
          />
          <Route 
            path="/analytics" 
            element={
              <LazyRouteWrapper>
                <Analytics />
              </LazyRouteWrapper>
            } 
          />
          <Route 
            path="/profile" 
            element={
              <LazyRouteWrapper>
                <Profile />
              </LazyRouteWrapper>
            } 
          />
          <Route 
            path="*" 
            element={
              <LazyRouteWrapper>
                <NotFound />
              </LazyRouteWrapper>
            } 
          />
        </Routes>
      </Router>
    </AuthProvider>
  )
}

export default App