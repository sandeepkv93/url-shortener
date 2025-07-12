import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import QRGenerator, { QRCodeOptions } from '../QRGenerator'

// Mock qrcode library
const mockToDataURL = vi.fn()
const mockToString = vi.fn()

vi.mock('qrcode', () => ({
  default: {
    toDataURL: mockToDataURL,
    toString: mockToString,
  },
}))

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

// Mock document.createElement for download functionality
const mockClick = vi.fn()
const mockAppendChild = vi.fn()
const mockRemoveChild = vi.fn()

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

describe('QRGenerator', () => {
  const mockUrl = 'https://example.com/test'
  const mockDataUrl = 'data:image/png;base64,mock-data'
  const mockSvgString = '<svg>mock svg content</svg>'

  beforeEach(() => {
    vi.clearAllMocks()
    mockToDataURL.mockResolvedValue(mockDataUrl)
    mockToString.mockResolvedValue(mockSvgString)
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('Rendering', () => {
    it('renders with default props', async () => {
      render(<QRGenerator url={mockUrl} />)
      
      expect(screen.getByText('QR Code Generator')).toBeInTheDocument()
      expect(screen.getByDisplayValue(mockUrl)).toBeInTheDocument()
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
    })

    it('renders without customization panel when showCustomization is false', () => {
      render(<QRGenerator url={mockUrl} showCustomization={false} />)
      
      expect(screen.queryByTitle('Customize QR Code')).not.toBeInTheDocument()
      expect(screen.queryByTitle('Reset to defaults')).not.toBeInTheDocument()
    })

    it('applies custom className', () => {
      const { container } = render(<QRGenerator url={mockUrl} className="custom-class" />)
      
      expect(container.firstChild).toHaveClass('custom-class')
    })

    it('shows loading state initially', () => {
      render(<QRGenerator url={mockUrl} />)
      
      expect(screen.getByRole('generic', { hidden: true })).toHaveClass('animate-spin')
    })
  })

  describe('QR Code Generation', () => {
    it('generates PNG QR code with default options', async () => {
      render(<QRGenerator url={mockUrl} />)
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledWith(mockUrl, {
          width: 256,
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

    it('generates SVG QR code when format is changed', async () => {
      render(<QRGenerator url={mockUrl} />)
      
      // Open customization panel
      fireEvent.click(screen.getByTitle('Customize QR Code'))
      
      // Change format to SVG
      const formatSelect = screen.getByLabelText('Format')
      fireEvent.change(formatSelect, { target: { value: 'image/svg+xml' } })
      
      await waitFor(() => {
        expect(mockToString).toHaveBeenCalledWith(mockUrl, {
          type: 'svg',
          width: 256,
          margin: 4,
          color: {
            dark: '#000000',
            light: '#ffffff'
          },
          errorCorrectionLevel: 'M'
        })
      })
    })

    it('uses initial options when provided', async () => {
      const initialOptions = {
        size: 400,
        backgroundColor: '#ff0000',
        foregroundColor: '#00ff00',
        margin: 2,
        errorCorrectionLevel: 'H' as const
      }
      
      render(<QRGenerator url={mockUrl} initialOptions={initialOptions} />)
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledWith(mockUrl, {
          width: 400,
          margin: 2,
          color: {
            dark: '#00ff00',
            light: '#ff0000'
          },
          errorCorrectionLevel: 'H'
        })
      })
    })

    it('calls onGenerated callback when QR code is generated', async () => {
      const mockOnGenerated = vi.fn()
      render(<QRGenerator url={mockUrl} onGenerated={mockOnGenerated} />)
      
      await waitFor(() => {
        expect(mockOnGenerated).toHaveBeenCalledWith(mockDataUrl, expect.any(Object))
      })
    })

    it('shows error when URL is empty', async () => {
      render(<QRGenerator url="" />)
      
      await waitFor(() => {
        expect(screen.getByText('URL is required')).toBeInTheDocument()
      })
    })

    it('shows error when QR code generation fails', async () => {
      mockToDataURL.mockRejectedValue(new Error('Generation failed'))
      
      render(<QRGenerator url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByText(/Failed to generate QR code/)).toBeInTheDocument()
      })
    })
  })

  describe('Customization', () => {
    beforeEach(async () => {
      render(<QRGenerator url={mockUrl} />)
      
      // Open customization panel
      fireEvent.click(screen.getByTitle('Customize QR Code'))
      
      await waitFor(() => {
        expect(screen.getByText('Customization Options')).toBeInTheDocument()
      })
    })

    it('toggles customization panel', () => {
      expect(screen.getByText('Customization Options')).toBeInTheDocument()
      
      // Close panel
      fireEvent.click(screen.getByTitle('Customize QR Code'))
      
      expect(screen.queryByText('Customization Options')).not.toBeInTheDocument()
    })

    it('updates size option', async () => {
      const sizeSlider = screen.getByLabelText('Size (px)')
      fireEvent.change(sizeSlider, { target: { value: '320' } })
      
      expect(screen.getByText('320px')).toBeInTheDocument()
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledWith(mockUrl, expect.objectContaining({
          width: 320
        }))
      })
    })

    it('updates margin option', async () => {
      const marginSlider = screen.getByLabelText('Margin')
      fireEvent.change(marginSlider, { target: { value: '6' } })
      
      expect(screen.getByText('6 units')).toBeInTheDocument()
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledWith(mockUrl, expect.objectContaining({
          margin: 6
        }))
      })
    })

    it('updates foreground color', async () => {
      const colorInputs = screen.getAllByDisplayValue('#000000')
      const textInput = colorInputs.find(input => input.getAttribute('type') === 'text')
      
      if (textInput) {
        fireEvent.change(textInput, { target: { value: '#ff0000' } })
        
        await waitFor(() => {
          expect(mockToDataURL).toHaveBeenCalledWith(mockUrl, expect.objectContaining({
            color: expect.objectContaining({
              dark: '#ff0000'
            })
          }))
        })
      }
    })

    it('updates background color', async () => {
      const colorInputs = screen.getAllByDisplayValue('#ffffff')
      const textInput = colorInputs.find(input => input.getAttribute('type') === 'text')
      
      if (textInput) {
        fireEvent.change(textInput, { target: { value: '#00ff00' } })
        
        await waitFor(() => {
          expect(mockToDataURL).toHaveBeenCalledWith(mockUrl, expect.objectContaining({
            color: expect.objectContaining({
              light: '#00ff00'
            })
          }))
        })
      }
    })

    it('updates error correction level', async () => {
      const errorLevelSelect = screen.getByLabelText('Error Correction')
      fireEvent.change(errorLevelSelect, { target: { value: 'H' } })
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledWith(mockUrl, expect.objectContaining({
          errorCorrectionLevel: 'H'
        }))
      })
    })

    it('resets to default options', async () => {
      // Change some options first
      const sizeSlider = screen.getByLabelText('Size (px)')
      fireEvent.change(sizeSlider, { target: { value: '400' } })
      
      // Reset to defaults
      fireEvent.click(screen.getByTitle('Reset to defaults'))
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledWith(mockUrl, expect.objectContaining({
          width: 256
        }))
      })
    })
  })

  describe('Download Functionality', () => {
    beforeEach(async () => {
      render(<QRGenerator url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
    })

    it('downloads PNG QR code', async () => {
      const downloadButton = screen.getByText('Download PNG')
      fireEvent.click(downloadButton)
      
      expect(document.createElement).toHaveBeenCalledWith('a')
      expect(mockAppendChild).toHaveBeenCalled()
      expect(mockClick).toHaveBeenCalled()
      expect(mockRemoveChild).toHaveBeenCalled()
    })

    it('downloads SVG QR code', async () => {
      // Change to SVG format
      fireEvent.click(screen.getByTitle('Customize QR Code'))
      const formatSelect = screen.getByLabelText('Format')
      fireEvent.change(formatSelect, { target: { value: 'image/svg+xml' } })
      
      await waitFor(() => {
        expect(screen.getByText('Download SVG')).toBeInTheDocument()
      })
      
      const downloadButton = screen.getByText('Download SVG')
      fireEvent.click(downloadButton)
      
      expect(mockCreateObjectURL).toHaveBeenCalled()
      expect(document.createElement).toHaveBeenCalledWith('a')
      expect(mockClick).toHaveBeenCalled()
    })
  })

  describe('Regeneration', () => {
    beforeEach(async () => {
      render(<QRGenerator url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
    })

    it('regenerates QR code when regenerate button is clicked', async () => {
      const regenerateButton = screen.getByText('Regenerate')
      
      vi.clearAllMocks()
      fireEvent.click(regenerateButton)
      
      await waitFor(() => {
        expect(mockToDataURL).toHaveBeenCalledWith(mockUrl, expect.any(Object))
      })
    })

    it('disables regenerate button while generating', async () => {
      // Make generation slower
      mockToDataURL.mockImplementation(() => new Promise(resolve => setTimeout(() => resolve(mockDataUrl), 100)))
      
      const regenerateButton = screen.getByText('Regenerate')
      fireEvent.click(regenerateButton)
      
      expect(regenerateButton).toBeDisabled()
      
      await waitFor(() => {
        expect(regenerateButton).not.toBeDisabled()
      })
    })
  })

  describe('URL Display', () => {
    it('displays the target URL', () => {
      render(<QRGenerator url={mockUrl} />)
      
      expect(screen.getByText('Target URL:')).toBeInTheDocument()
      expect(screen.getByText(mockUrl)).toBeInTheDocument()
    })

    it('handles long URLs properly', () => {
      const longUrl = 'https://example.com/very/long/url/that/should/break/properly/in/the/display'
      render(<QRGenerator url={longUrl} />)
      
      expect(screen.getByText(longUrl)).toBeInTheDocument()
    })
  })

  describe('Accessibility', () => {
    it('has proper labels for form controls', async () => {
      render(<QRGenerator url={mockUrl} />)
      
      // Open customization panel
      fireEvent.click(screen.getByTitle('Customize QR Code'))
      
      await waitFor(() => {
        expect(screen.getByLabelText('Size (px)')).toBeInTheDocument()
        expect(screen.getByLabelText('Margin')).toBeInTheDocument()
        expect(screen.getByLabelText('Foreground Color')).toBeInTheDocument()
        expect(screen.getByLabelText('Background Color')).toBeInTheDocument()
        expect(screen.getByLabelText('Error Correction')).toBeInTheDocument()
        expect(screen.getByLabelText('Format')).toBeInTheDocument()
      })
    })

    it('has proper alt text for QR code image', async () => {
      render(<QRGenerator url={mockUrl} />)
      
      await waitFor(() => {
        expect(screen.getByAltText('QR Code')).toBeInTheDocument()
      })
    })

    it('has proper button titles', async () => {
      render(<QRGenerator url={mockUrl} />)
      
      expect(screen.getByTitle('Customize QR Code')).toBeInTheDocument()
      expect(screen.getByTitle('Reset to defaults')).toBeInTheDocument()
    })
  })
})