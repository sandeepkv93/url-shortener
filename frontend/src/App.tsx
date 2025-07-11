import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { AuthProvider } from '@/context/AuthContext'
import Layout from '@/components/common/Layout'
import Home from '@/pages/Home'
import Dashboard from '@/pages/Dashboard'
import Analytics from '@/pages/Analytics'
import Profile from '@/pages/Profile'
import NotFound from '@/pages/NotFound'
import ComponentDemo from '@/pages/ComponentDemo'

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/demo" element={<ComponentDemo />} />
          <Route path="/" element={<Layout><Home /></Layout>} />
          <Route path="/dashboard" element={<Layout><Dashboard /></Layout>} />
          <Route path="/analytics" element={<Layout><Analytics /></Layout>} />
          <Route path="/profile" element={<Layout><Profile /></Layout>} />
          <Route path="*" element={<Layout><NotFound /></Layout>} />
        </Routes>
      </Router>
    </AuthProvider>
  )
}

export default App