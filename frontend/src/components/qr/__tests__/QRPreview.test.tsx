import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import QRPreview from '../QRPreview'

// Mock qrcode library
const mockToDataURL = vi.fn()
const mockToString = vi.fn()

vi.mock('qrcode', () => ({
  default: {
    toDataURL: mockToDataURL,
    toString: mockToString,
  },
}))

// Mock clipboard API
const mockWriteText = vi.fn()
const mockWrite = vi.fn()

Object.defineProperty(navigator, 'clipboard', {
  value: {
    writeText: mockWriteText,
    write: mockWrite,
  },
  writable: true,
})

// Mock ClipboardItem
window.ClipboardItem = vi.fn().mockImplementation((data) => ({ data }))

// Mock URL.createObjectURL and URL.revokeObjectURL
const mockCreateObjectURL = vi.fn(() => 'mock-blob-url')
const mockRevokeObjectURL = vi.fn()

Object.defineProperty(URL, 'createObjectURL', {
  value: mockCreateObjectURL,
  writable: true,
})

Object.defineProperty(URL, 'revokeObjectURL', {
  value: mockRevokeObjectURL,
  writable: true,
})

// Mock fetch for clipboard functionality
global.fetch = vi.fn()

// Mock document.createElement and window.open for download/open functionality
const mockClick = vi.fn()
const mockAppendChild = vi.fn()
const mockRemoveChild = vi.fn()
const mockWindowOpen = vi.fn()

Object.defineProperty(document, 'createElement', {
  value: vi.fn(() => ({
    href: '',
    download: '',
    click: mockClick,
  })),
  writable: true,
})

Object.defineProperty(document.body, 'appendChild', {
  value: mockAppendChild,
  writable: true,
})

Object.defineProperty(document.body, 'removeChild', {
  value: mockRemoveChild,
  writable: true,
})

Object.defineProperty(window, 'open', {
  value: mockWindowOpen,
  writable: true,
})

describe('QRPreview', () => {
  const mockUrl = 'https://example.com/test'
  const mockDataUrl = 'data:image/png;base64,mock-data'
  const mockSvgString = '<svg>mock svg content</svg>'

  beforeEach(() => {
    vi.clearAllMocks()
    mockToDataURL.mockResolvedValue(mockDataUrl)
    mockToString.mockResolvedValue(mockSvgString)
    mockWindowOpen.mockReturnValue({
      document: {
        write: vi.fn(),
      },
    })
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('Rendering', () => {
    it('renders with default props', async () => {
      render(<QRPreview url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
      
      expect(screen.getByText('Target URL:')).toBeInTheDocument()
      expect(screen.getByText(mockUrl)).toBeInTheDocument()
    })

    it('renders with custom size', async () => {
      render(<QRPreview url={mockUrl} size={300} />)
      
      await waitFor(() => {
        const qrImage = screen.getByAltText('QR Code')
        expect(qrImage).toHaveAttribute('style', expect.stringContaining('width: 300px'))
      })
    })

    it('hides actions when showActions is false', async () => {
      render(<QRPreview url={mockUrl} showActions={false} />)
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
      
      expect(screen.queryByTitle('Download QR Code')).not.toBeInTheDocument()
      expect(screen.queryByTitle('Copy QR Code')).not.toBeInTheDocument()
    })

    it('hides URL when showUrl is false', async () => {
      render(<QRPreview url={mockUrl} showUrl={false} />)
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
      
      expect(screen.queryByText('Target URL:')).not.toBeInTheDocument()
    })

    it('applies custom className', () => {
      const { container } = render(<QRPreview url={mockUrl} className="custom-class" />)
      
      expect(container.firstChild).toHaveClass('custom-class')
    })

    it('shows loading state initially', () => {
      render(<QRPreview url={mockUrl} />)
      
      expect(screen.getByRole('generic', { hidden: true })).toHaveClass('animate-spin')
    })
  })

  describe('QR Code Generation', () => {
    it('generates PNG QR code by default', async () => {
      render(<QRPreview url={mockUrl} />)
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledWith(mockUrl, {
          width: 200,
          margin: 4,
          color: {
            dark: '#000000',
            light: '#ffffff'
          },
          errorCorrectionLevel: 'M'
        })
      })

      expect(screen.getByAltText('QR Code')).toHaveAttribute('src', mockDataUrl)
    })

    it('generates SVG QR code when format is svg', async () => {
      render(<QRPreview url={mockUrl} format="svg" />)
      
      await waitFor(() => {
        expect(mockToString).toHaveBeenCalledWith(mockUrl, {
          type: 'svg',
          width: 200,
          margin: 4,
          color: {
            dark: '#000000',
            light: '#ffffff'
          },
          errorCorrectionLevel: 'M'
        })
      })
    })

    it('uses custom colors', async () => {
      render(
        <QRPreview 
          url={mockUrl} 
          backgroundColor="#ff0000" 
          foregroundColor="#00ff00" 
        />
      )
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledWith(mockUrl, expect.objectContaining({
          color: {
            dark: '#00ff00',
            light: '#ff0000'
          }
        }))
      })
    })

    it('uses custom margin and error correction level', async () => {
      render(
        <QRPreview 
          url={mockUrl} 
          margin={6} 
          errorCorrectionLevel="H" 
        />
      )
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledWith(mockUrl, expect.objectContaining({
          margin: 6,
          errorCorrectionLevel: 'H'
        }))
      })
    })

    it('shows error when URL is empty', async () => {
      render(<QRPreview url="" />)
      
      await waitFor(() => {
        expect(screen.getByText('Failed to generate QR code')).toBeInTheDocument()
      })
    })

    it('shows error when QR code generation fails', async () => {
      mockToDataURL.mockRejectedValue(new Error('Generation failed'))
      
      render(<QRPreview url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByText('Failed to generate QR code')).toBeInTheDocument()
      })
    })

    it('regenerates QR code when props change', async () => {
      const { rerender } = render(<QRPreview url={mockUrl} size={200} />)
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledTimes(1)
      })
      
      vi.clearAllMocks()
      
      rerender(<QRPreview url={mockUrl} size={300} />)
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledWith(mockUrl, expect.objectContaining({
          width: 300
        }))
      })
    })
  })

  describe('Download Functionality', () => {
    beforeEach(async () => {
      render(<QRPreview url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
    })

    it('downloads PNG QR code', async () => {
      const downloadButton = screen.getByTitle('Download QR Code')
      fireEvent.click(downloadButton)
      
      expect(document.createElement).toHaveBeenCalledWith('a')
      expect(mockAppendChild).toHaveBeenCalled()
      expect(mockClick).toHaveBeenCalled()
      expect(mockRemoveChild).toHaveBeenCalled()
    })

    it('downloads SVG QR code', async () => {
      const { rerender } = render(<QRPreview url={mockUrl} format="svg" />)
      
      await waitFor(() => {
        expect(mockToString).toHaveBeenCalled()
      })
      
      const downloadButton = screen.getByTitle('Download QR Code')
      fireEvent.click(downloadButton)
      
      expect(mockCreateObjectURL).toHaveBeenCalled()
      expect(document.createElement).toHaveBeenCalledWith('a')
      expect(mockClick).toHaveBeenCalled()
    })
  })

  describe('Copy Functionality', () => {
    beforeEach(async () => {
      render(<QRPreview url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
    })

    it('copies QR code to clipboard for PNG format', async () => {
      const mockBlob = new Blob(['mock-data'], { type: 'image/png' })
      
      global.fetch = vi.fn().mockResolvedValue({
        blob: () => Promise.resolve(mockBlob),
      } as Response)
      
      mockWrite.mockResolvedValue(undefined)
      
      const copyButton = screen.getByTitle('Copy QR Code')
      fireEvent.click(copyButton)
      
      await waitFor(() => {
        expect(global.fetch).toHaveBeenCalledWith(mockDataUrl)
        expect(mockWrite).toHaveBeenCalled()
      })
      
      expect(screen.getByText('Copied!')).toBeInTheDocument()
    })

    it('falls back to download when clipboard API fails', async () => {
      global.fetch = vi.fn().mockRejectedValue(new Error('Fetch failed'))
      
      const copyButton = screen.getByTitle('Copy QR Code')
      fireEvent.click(copyButton)
      
      await waitFor(() => {
        expect(document.createElement).toHaveBeenCalledWith('a')
        expect(mockClick).toHaveBeenCalled()
      })
    })

    it('does not show copy button for SVG format', async () => {
      const { rerender } = render(<QRPreview url={mockUrl} format="svg" />)
      
      await waitFor(() => {
        expect(mockToString).toHaveBeenCalled()
      })
      
      expect(screen.queryByTitle('Copy QR Code')).not.toBeInTheDocument()
    })
  })

  describe('Open in New Tab Functionality', () => {
    beforeEach(async () => {
      render(<QRPreview url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
    })

    it('opens PNG QR code in new tab', () => {
      const openButton = screen.getByTitle('Open in new tab')
      fireEvent.click(openButton)
      
      expect(mockWindowOpen).toHaveBeenCalled()
    })

    it('opens SVG QR code in new tab', async () => {
      const { rerender } = render(<QRPreview url={mockUrl} format="svg" />)
      
      await waitFor(() => {
        expect(mockToString).toHaveBeenCalled()
      })
      
      const openButton = screen.getByTitle('Open in new tab')
      fireEvent.click(openButton)
      
      expect(mockWindowOpen).toHaveBeenCalled()
    })
  })

  describe('Modal Functionality', () => {
    beforeEach(async () => {
      render(<QRPreview url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
    })

    it('opens modal when QR code is clicked', async () => {
      const qrImage = screen.getByAltText('QR Code')
      fireEvent.click(qrImage.closest('[style]') || qrImage)
      
      await waitFor(() => {
        expect(screen.getByText('QR Code')).toBeInTheDocument() // Modal title
        expect(screen.getByRole('button', { name: /close/i })).toBeInTheDocument()
      })
    })

    it('closes modal when close button is clicked', async () => {
      // Open modal
      const qrImage = screen.getByAltText('QR Code')
      fireEvent.click(qrImage.closest('[style]') || qrImage)
      
      await waitFor(() => {
        expect(screen.getByText('QR Code')).toBeInTheDocument()
      })
      
      // Close modal
      const closeButton = screen.getByRole('button', { name: /close/i })
      fireEvent.click(closeButton)
      
      await waitFor(() => {
        expect(screen.queryByText('QR Code')).not.toBeInTheDocument()
      })
    })

    it('shows download and copy actions in modal', async () => {
      // Open modal
      const qrImage = screen.getByAltText('QR Code')
      fireEvent.click(qrImage.closest('[style]') || qrImage)
      
      await waitFor(() => {
        expect(screen.getByText('Download PNG')).toBeInTheDocument()
        expect(screen.getByText('Copy')).toBeInTheDocument()
      })
    })
  })

  describe('Hover Effects', () => {
    beforeEach(async () => {
      render(<QRPreview url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
    })

    it('shows hover overlay with maximize icon', async () => {
      const qrContainer = screen.getByAltText('QR Code').closest('.group')
      
      expect(qrContainer).toBeInTheDocument()
      
      if (qrContainer) {
        fireEvent.mouseEnter(qrContainer)
        
        // Check if Maximize2 icon becomes visible (through group-hover class)
        const overlay = qrContainer.querySelector('.group-hover\\:bg-opacity-10')
        expect(overlay).toBeInTheDocument()
      }
    })
  })

  describe('Accessibility', () => {
    it('has proper alt text for QR code image', async () => {
      render(<QRPreview url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
    })

    it('has proper button titles', async () => {
      render(<QRPreview url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByTitle('Download QR Code')).toBeInTheDocument()
        expect(screen.getByTitle('Copy QR Code')).toBeInTheDocument()
        expect(screen.getByTitle('Open in new tab')).toBeInTheDocument()
      })
    })

    it('is keyboard accessible for modal', async () => {
      render(<QRPreview url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
      
      const user = userEvent.setup()
      
      // Use tab to navigate to the QR code container
      await user.tab()
      
      // Press Enter to open modal
      await user.keyboard('{Enter}')
      
      // Modal should open
      await waitFor(() => {
        expect(screen.getByText('QR Code')).toBeInTheDocument()
      })
    })
  })

  describe('Edge Cases', () => {
    it('handles very long URLs', async () => {
      const longUrl = 'https://example.com/' + 'very-long-path/'.repeat(50)
      render(<QRPreview url={longUrl} />)
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledWith(longUrl, expect.any(Object))
      })
      
      expect(screen.getByText(longUrl)).toBeInTheDocument()
    })

    it('handles special characters in URL', async () => {
      const specialUrl = 'https://example.com/path?query=test&param=value#anchor'
      render(<QRPreview url={specialUrl} />)
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledWith(specialUrl, expect.any(Object))
      })
    })

    it('handles no actions gracefully when QR generation fails', async () => {
      mockToDataURL.mockRejectedValue(new Error('Generation failed'))
      
      render(<QRPreview url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByText('Failed to generate QR code')).toBeInTheDocument()
      })
      
      expect(screen.queryByTitle('Download QR Code')).not.toBeInTheDocument()
    })
  })
})