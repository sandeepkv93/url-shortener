import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { createElement } from 'react'
import App from './App'

// Mock the router with complete set of components and hooks
vi.mock('react-router-dom', () => ({
  BrowserRouter: ({ children }: { children: React.ReactNode }) => createElement('div', null, children),
  Routes: ({ children }: { children: React.ReactNode }) => createElement('div', null, children),
  Route: () => createElement('div', null, 'Route'),
  Link: ({ children, to, ...props }: any) => createElement('a', { href: to, ...props }, children),
  NavLink: ({ children, to, ...props }: any) => createElement('a', { href: to, ...props }, children),
  useLocation: () => ({
    pathname: '/',
    search: '',
    hash: '',
    state: null,
    key: 'default',
  }),
  useNavigate: () => vi.fn(),
  useParams: () => ({}),
  useSearchParams: () => [new URLSearchParams(), vi.fn()],
}))

describe('App', () => {
  it('renders without crashing', () => {
    render(<App />)
    // Check for the header brand specifically
    expect(screen.getByRole('heading', { level: 1, name: 'URL Shortener' })).toBeInTheDocument()
  })
})