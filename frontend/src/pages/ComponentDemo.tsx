import { useState } from 'react'
import LoginForm from '@/components/auth/LoginForm'
import RegisterForm from '@/components/auth/RegisterForm'
import PasswordReset from '@/components/auth/PasswordReset'
import Header from '@/components/common/Header'
import Footer from '@/components/common/Footer'
import Layout from '@/components/common/Layout'
import Loading, { 
  ButtonLoading, 
  PageLoading, 
  SkeletonLoading, 
  InlineLoading 
} from '@/components/common/Loading'

const ComponentDemo = () => {
  const [currentComponent, setCurrentComponent] = useState<'login' | 'register' | 'password-reset' | 'components'>('login')

  const components = [
    { id: 'login', name: 'Login Form', component: <LoginForm /> },
    { id: 'register', name: 'Register Form', component: <RegisterForm /> },
    { id: 'password-reset', name: 'Password Reset', component: <PasswordReset /> },
    { id: 'components', name: 'Common Components', component: null },
  ]

  const renderCurrentComponent = () => {
    switch (currentComponent) {
      case 'login':
        return <LoginForm />
      case 'register':
        return <RegisterForm />
      case 'password-reset':
        return <PasswordReset />
      case 'components':
        return (
          <div className="space-y-8">
            <div className="bg-white rounded-lg shadow-md p-6">
              <h3 className="text-lg font-semibold mb-4">Loading Components</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h4 className="font-medium mb-2">Spinner Loading</h4>
                  <Loading variant="spinner" size="md" message="Loading..." />
                </div>
                <div>
                  <h4 className="font-medium mb-2">Dots Loading</h4>
                  <Loading variant="dots" size="md" message="Processing..." />
                </div>
                <div>
                  <h4 className="font-medium mb-2">Button Loading</h4>
                  <button className="bg-primary-600 text-white px-4 py-2 rounded-md flex items-center">
                    <ButtonLoading className="mr-2" />
                    Submitting...
                  </button>
                </div>
                <div>
                  <h4 className="font-medium mb-2">Inline Loading</h4>
                  <p className="text-gray-600">
                    Fetching data <InlineLoading size="sm" />
                  </p>
                </div>
              </div>
              <div className="mt-6">
                <h4 className="font-medium mb-2">Skeleton Loading</h4>
                <SkeletonLoading />
              </div>
            </div>
            
            <div className="bg-white rounded-lg shadow-md p-6">
              <h3 className="text-lg font-semibold mb-4">Header Component</h3>
              <div className="border-2 border-dashed border-gray-300 rounded-lg p-4">
                <p className="text-gray-600 text-center">
                  Header component is shown at the top of this page
                </p>
              </div>
            </div>
            
            <div className="bg-white rounded-lg shadow-md p-6">
              <h3 className="text-lg font-semibold mb-4">Footer Component</h3>
              <div className="border-2 border-dashed border-gray-300 rounded-lg p-4">
                <p className="text-gray-600 text-center">
                  Footer component is shown at the bottom of this page
                </p>
              </div>
            </div>
          </div>
        )
      default:
        return <LoginForm />
    }
  }

  return (
    <Layout>
      <div className="min-h-screen bg-gray-50">
        {/* Demo Navigation */}
        <div className="bg-white shadow-sm border-b border-gray-200 mb-8">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex space-x-8 py-4">
              <h2 className="text-lg font-semibold text-gray-900 self-center">Component Demo:</h2>
              {components.map((comp) => (
                <button
                  key={comp.id}
                  onClick={() => setCurrentComponent(comp.id as any)}
                  className={`px-3 py-2 text-sm font-medium rounded-md transition-colors ${
                    currentComponent === comp.id
                      ? 'bg-primary-100 text-primary-700'
                      : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                  }`}
                >
                  {comp.name}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Component Display */}
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="flex justify-center">
            {renderCurrentComponent()}
          </div>
        </div>
      </div>
    </Layout>
  )
}

export default ComponentDemo